package services

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io"
	"math"
	"net/http"
	s "scaler/shared"
)

type BBBService struct {
	State  BBBServiceState  `yaml:"-"`
	Config BBBServiceConfig `yaml:"bbb_config"`
}

type BBBServiceState struct {
	Name string
}

func (bbb BBBServiceState) GetName() string {
	return bbb.Name
}

type BBBServiceConfig struct {
	CycleTimeSeconds int             `yaml:"cycle_time_seconds"`
	Resources        s.Resources     `yaml:"resources"`
	ApiToken         s.StringFromEnv `yaml:"api_token"`
}

// BBBGetMeetingsResponseXML is the XML response from the BBB API when calling getMeetings
// We only keep the fields we need from the response
type BBBGetMeetingsResponseXML struct {
	XMLName    xml.Name `xml:"response"`
	Returncode string   `xml:"returncode"`
	MessageKey string   `xml:"messageKey"`
	Message    string   `xml:"message"`
	Meetings   struct {
		Meeting []struct {
			ParticipantCount int `xml:"participantCount"`
		} `xml:"meeting"`
	} `xml:"meetings"`
}

func (bbb BBBService) Init() error {
	if err := initMetricsExporter("bbb"); err != nil {
		return fmt.Errorf("error while registering metrics: %s", err)
	}
	return nil
}

func (bbb *BBBService) GetState() s.ServiceState {
	return bbb.State
}

func (bbb *BBBService) GetConfig() BBBServiceConfig {
	return bbb.Config
}

// See https://docs.bigbluebutton.org/development/api/#usage
func signedBBBAPIRequest(serverUrl, endpoint, parameters, apiToken string) string {
	queryString := endpoint + parameters + apiToken
	checksumRaw := sha1.Sum([]byte(queryString))
	checksumHex := hex.EncodeToString(checksumRaw[:])
	return fmt.Sprintf("https://%s/bigbluebutton/api/%s?%s&checksum=%s", serverUrl, endpoint, parameters, checksumHex)
}

func doBBBAPICall(serverUrl, endpoint, parameters, apiToken string) ([]byte, error) {
	url := signedBBBAPIRequest(serverUrl, endpoint, parameters, apiToken)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func countParticipants(meetingsResponse *BBBGetMeetingsResponseXML) int {
	count := 0
	for _, meeting := range meetingsResponse.Meetings.Meeting {
		count += meeting.ParticipantCount
	}
	return count
}

func parseBBBGetMeetingsResponseXML(xmlRaw []byte) (*BBBGetMeetingsResponseXML, error) {
	xmlParsed, err := s.ParseXML[BBBGetMeetingsResponseXML](xmlRaw)
	if err != nil {
		return nil, err
	}
	if xmlParsed.Returncode != "SUCCESS" {
		return nil, fmt.Errorf("BBB API returned error: %s - %s", xmlParsed.MessageKey, xmlParsed.Message)
	}
	return xmlParsed, nil
}

func getMeetings(serverUrl, apiToken string) (*BBBGetMeetingsResponseXML, error) {
	body, err := doBBBAPICall(serverUrl, "getMeetings", "", apiToken)
	if err != nil {
		return nil, err
	}
	return parseBBBGetMeetingsResponseXML(body)
}

func (bbb BBBService) GetParticipantsCount(serverUrl string) (int, error) {
	meetingsResponse, err := getMeetings(serverUrl, string(bbb.Config.ApiToken))
	if err != nil {
		errorsTotalCounter.Inc()
		return 0, err
	}
	return countParticipants(meetingsResponse), nil
}

func (bbb BBBService) GetResources() s.Resources {
	return bbb.Config.Resources
}

func (bbb BBBService) GetCycleTimeSeconds() int {
	return bbb.Config.CycleTimeSeconds
}

func (bbb BBBService) ShouldScale(object s.ScaledObject) (s.ScaleResource, error) {
	var server *s.Server
	switch object.(type) {
	case *s.Server:
		server = object.(*s.Server)
	default:
		return s.ScaleResource{}, fmt.Errorf("unsupported scaled object type: %s", object.GetType())
	}

	if !server.Ready {
		return s.ScaleResource{}, fmt.Errorf("server %s is not ready", server.ServerName)
	}

	participantsCount, err := bbb.GetParticipantsCount(server.ServerName)
	if err != nil {
		return s.ScaleResource{}, fmt.Errorf("error while getting participants count: %s", err)
	}

	return applyRules(*server, participantsCount, bbb), nil
}

func applyRules(server s.Server, participantsCount int, bbb BBBService) s.ScaleResource {
	targetResource := s.ScaleResource{
		Cpu: s.ScaleOp{
			Direction: s.ScaleNone,
			Reason:    "Default",
			Amount:    0,
		},
		Mem: s.ScaleOp{
			Direction: s.ScaleNone,
			Reason:    "Default",
			Amount:    0,
		},
	}

	// Scaling rules:
	// 1. Scale up if current resource is below configured minimum
	// 2. Scale up if current resource usage exceeds maximum usage
	// Add enough resources to either reach usage below the maximum usage or the maximum amount of resources
	// 3. Scale down to the configured minimum if there are no participants

	// Rule 1 CPU
	if server.ResourceState.Cpu.CurrentCores < int32(bbb.Config.Resources.Cpu.MinCores) {
		targetResource.Cpu.Direction = s.ScaleUp
		targetResource.Cpu.Reason = targetResource.Cpu.Reason + ",Rule 1: resource below minimum"
		targetResource.Cpu.Amount = int32(bbb.Config.Resources.Cpu.MinCores) - server.ResourceState.Cpu.CurrentCores
	}

	// Rule 2 CPU
	if cpuMaxUsageDelta := server.ResourceState.Cpu.CurrentUsage - bbb.Config.Resources.Cpu.MaxUsage; cpuMaxUsageDelta > 0 && server.ResourceState.Cpu.CurrentCores < int32(bbb.Config.Resources.Cpu.MaxCores) {
		targetResource.Cpu.Direction = s.ScaleUp
		targetResource.Cpu.Reason = targetResource.Cpu.Reason + ",Rule 2: usage above maximum"
		cpuInc := cpuMaxUsageDelta * float32(server.ResourceState.Cpu.CurrentCores) / server.ResourceState.Cpu.CurrentUsage
		targetHeuristic := server.ResourceState.Cpu.CurrentCores + int32(math.Ceil(float64(cpuInc)))
		targetResource.Cpu.Amount = int32(math.Min(float64(targetHeuristic), float64((bbb.Config.Resources.Cpu.MaxCores)))) - server.ResourceState.Cpu.CurrentCores
	}

	// Rule 1 memory
	if server.ResourceState.Memory.CurrentBytes < int32(bbb.Config.Resources.Memory.MinBytes) {
		targetResource.Mem.Direction = s.ScaleUp
		targetResource.Mem.Reason = targetResource.Mem.Reason + ",Rule 1: resource below minimum"
		targetResource.Mem.Amount = int32(bbb.Config.Resources.Memory.MinBytes) - server.ResourceState.Memory.CurrentBytes
	}

	// Rule 2 memory
	if memMaxUsageDelta := server.ResourceState.Memory.CurrentUsage - bbb.Config.Resources.Memory.MaxUsage; memMaxUsageDelta > 0 && server.ResourceState.Memory.CurrentBytes < int32(bbb.Config.Resources.Memory.MaxBytes) {
		targetResource.Mem.Direction = s.ScaleUp
		targetResource.Mem.Reason = targetResource.Mem.Reason + ",Rule 2: usage above maximum"
		memInc := memMaxUsageDelta * float32(server.ResourceState.Memory.CurrentBytes) / server.ResourceState.Memory.CurrentUsage
		targetHeuristic := server.ResourceState.Memory.CurrentBytes + int32(math.Ceil(float64(memInc)))
		targetResource.Mem.Amount = int32(math.Min(float64(targetHeuristic), float64((bbb.Config.Resources.Memory.MaxBytes)))) - server.ResourceState.Memory.CurrentBytes
	}

	// Rule 3 CPU and memory
	if participantsCount == 0 {
		if server.ResourceState.Memory.CurrentBytes > int32(bbb.Config.Resources.Memory.MinBytes) {
			targetResource.Mem.Direction = s.ScaleDown
			targetResource.Mem.Reason = targetResource.Mem.Reason + ",Rule 3: no participants"
			targetResource.Mem.Amount = int32(bbb.Config.Resources.Memory.MinBytes) - server.ResourceState.Memory.CurrentBytes
		}
		if server.ResourceState.Cpu.CurrentCores > int32(bbb.Config.Resources.Cpu.MinCores) {
			targetResource.Cpu.Direction = s.ScaleDown
			targetResource.Cpu.Reason = targetResource.Cpu.Reason + ",Rule 3: no participants"
			targetResource.Cpu.Amount = int32(bbb.Config.Resources.Cpu.MinCores) - server.ResourceState.Cpu.CurrentCores
		}
	}

	return targetResource
}

func (service BBBService) Validate() error {
	if err := service.Config.Validate(); err != nil {
		return err
	}
	return nil
}

func (config BBBServiceConfig) Validate() error {
	if err := config.Resources.Validate(); err != nil {
		return err
	}
	if config.ApiToken == "" {
		return fmt.Errorf("bbb.api_token is empty")
	}
	return nil
}
