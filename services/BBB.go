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

func (bbb BBBService) ShouldScale(server s.Server) (s.ScaleResource, error) {
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

	participantsCount, err := bbb.GetParticipantsCount(server.ServerName)
	if err != nil {
		return targetResource, fmt.Errorf("error while getting participants count: %s", err)
	}

	// Scaling rules

	// Rule 1: Scale up cpu if current cpu is below configured minimum for the service
	if server.ServerCpu < int32(bbb.Config.Resources.Cpu.MinCores) {
		targetResource.Cpu.Direction = s.ScaleUp
		targetResource.Cpu.Reason = targetResource.Cpu.Reason + ",Rule 1"
		targetResource.Cpu.Amount = int32(bbb.Config.Resources.Cpu.MinCores)
	}

	// Rule 2: Scale up cpu if current cpu usage is above configured maximum usage for the service
	// Scale up to either reach cpu usage below the configured maximum usage or to the configured maximum cpu
	if cpuMaxUsageDelta := server.ServerCpuUsage - bbb.Config.Resources.Cpu.MaxUsage; cpuMaxUsageDelta > 0 {
		targetResource.Cpu.Direction = s.ScaleUp
		targetResource.Cpu.Reason = targetResource.Cpu.Reason + ",Rule 2"
		cpuInc := cpuMaxUsageDelta * float32(server.ServerCpu) / server.ServerCpuUsage
		targetHeuristic := server.ServerCpu + int32(math.Ceil(float64(cpuInc)))
		targetResource.Cpu.Amount = int32(math.Min(float64(targetHeuristic), float64((bbb.Config.Resources.Cpu.MaxCores))))
	}

	// Rule 3: RuleScale down cpu if current cpu usage is below configured minimum usage for the service
	if cpuMinUsageDelta := server.ServerCpuUsage - bbb.Config.Resources.Cpu.MinUsage; cpuMinUsageDelta < 0 && server.ServerCpu > int32(bbb.Config.Resources.Cpu.MinCores) {
		targetResource.Cpu.Direction = s.ScaleDown
		targetResource.Cpu.Reason = targetResource.Cpu.Reason + ",Rule 3"
		cpuDec := cpuMinUsageDelta * float32(server.ServerCpu) / server.ServerCpuUsage
		targetHeuristic := server.ServerCpu + int32(math.Floor(float64(cpuDec)))
		targetResource.Cpu.Amount = int32(math.Max(float64(targetHeuristic), float64((bbb.Config.Resources.Cpu.MinCores))))
	}

	// Rule 4: Scale up memory if current ram is below configured minimum for the service
	if server.ServerRam < int32(bbb.Config.Resources.Memory.MinBytes) {
		targetResource.Mem.Direction = s.ScaleUp
		targetResource.Mem.Reason = targetResource.Mem.Reason + ",Rule 4"
		targetResource.Mem.Amount = int32(bbb.Config.Resources.Memory.MinBytes)
	}

	// Rule 5: Scale up memory if current ram usage is above configured maximum usage for the service
	// Scale up to either reach ram usage below the configured maximum usage or to the configured maximum memory
	if memMaxUsageDelta := server.ServerRamUsage - bbb.Config.Resources.Memory.MaxUsage; memMaxUsageDelta > 0 {
		targetResource.Mem.Direction = s.ScaleUp
		targetResource.Mem.Reason = targetResource.Mem.Reason + ",Rule 5"
		memInc := memMaxUsageDelta * float32(server.ServerRam) / server.ServerRamUsage
		targetHeuristic := server.ServerRam + int32(math.Ceil(float64(memInc)))
		targetResource.Mem.Amount = int32(math.Min(float64(targetHeuristic), float64((bbb.Config.Resources.Memory.MaxBytes))))
	}

	// Rule 6: Scale down memory if current ram usage is below configured minimum usage for the service
	// Scale down to either reach ram usage above the configured minimum usage or to the configured minimum memory
	if memMinUsageDelta := server.ServerRamUsage - bbb.Config.Resources.Memory.MinUsage; memMinUsageDelta < 0 && server.ServerRam > int32(bbb.Config.Resources.Memory.MinBytes) {
		targetResource.Mem.Direction = s.ScaleDown
		targetResource.Mem.Reason = targetResource.Mem.Reason + ",Rule 6"
		memDec := memMinUsageDelta * float32(server.ServerRam) / server.ServerRamUsage
		targetHeuristic := server.ServerRam + int32(math.Floor(float64(memDec)))
		targetResource.Mem.Amount = int32(math.Max(float64(targetHeuristic), float64((bbb.Config.Resources.Memory.MinBytes))))
	}

	// Rule 7: Scale down resources to the configured minimum if there are no participants
	if participantsCount == 0 {
		if server.ServerRam > int32(bbb.Config.Resources.Memory.MinBytes) {
			targetResource.Mem.Direction = s.ScaleDown
			targetResource.Mem.Reason = targetResource.Mem.Reason + ",Rule 7"
			targetResource.Mem.Amount = int32(bbb.Config.Resources.Memory.MinBytes)
		}
		if server.ServerCpu > int32(bbb.Config.Resources.Cpu.MinCores) {
			targetResource.Cpu.Direction = s.ScaleDown
			targetResource.Cpu.Reason = targetResource.Cpu.Reason + ",Rule 7"
			targetResource.Cpu.Amount = int32(bbb.Config.Resources.Cpu.MinCores)
		}
	}

	return targetResource, nil
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
