package services

import (
	"os"
	s "scaler/shared"
	"testing"
	"time"
)

var validBBBConfig = &BBBServiceConfig{
	CycleTimeSeconds: 60,
	Resources: s.Resources{
		Cpu: &s.CpuResources{
			MinCores: 2,
			MaxCores: 4,
			MaxUsage: 0.5,
		},
		Memory: &s.MemoryResources{
			MinBytes: 2048,
			MaxBytes: 4096,
			MaxUsage: 0.5,
		},
	},
	ApiToken: "1234567890",
}

var sampleBBBServer = s.Server{
	DatacenterId:    "1234",
	ServerId:        "5678",
	ServerName:      "bbb.example.com",
	CpuArchitecture: "x86",
	ResourceState: s.ResourceState{
		Cpu: &s.CpuResourceState{
			CurrentCores: 1,
			CurrentUsage: 0.5,
		},
		Memory: &s.MemoryResourceState{
			CurrentBytes: 1024,
			CurrentUsage: 0.5,
		},
	},
	LastUpdated: time.Now(),
	Ready:       true,
}

func TestValidateConfigOK(t *testing.T) {
	bbbConfig := validBBBConfig
	s.ValidatePass(t, bbbConfig)
}

func TestValidateConfigNotOK(t *testing.T) {
	bbbConfig := &BBBServiceConfig{}
	s.ValidateFail(t, bbbConfig)
}

func TestParseConfigOK(t *testing.T) {
	os.Setenv("BBB_API_TOKEN", "1234567890")
	defer os.Unsetenv("BBB_API_TOKEN")

	config, ok := s.OpenConfig("test_files/bbb_config_ok.yml")
	if ok != nil {
		t.Fatalf("Failed to open config: %v", ok)
	}

	c, ok := s.LoadConfig[BBBService](config)
	if ok != nil {
		t.Fatalf("Failed to parse config: %v", ok)
	}
	s.ValidatePass(t, c)
}

func TestParseConfigNotOK(t *testing.T) {
	config, ok := s.OpenConfig("test_files/bbb_config_not_ok.yml")
	if ok != nil {
		t.Fatalf("Failed to open config: %v", ok)
	}

	_, ok = s.LoadConfig[BBBService](config)
	if ok == nil {
		t.Fatalf("Expected error but got nil")
	}
}

func TestParseGetParticipantsCountResponseOK(t *testing.T) {
	responseBytes, err := os.ReadFile("test_files/bbb_get_participants_ok.xml")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}
	_, err = parseBBBGetMeetingsResponseXML(responseBytes)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}
}

func TestParseGetParticipantsCountResponseNotOK(t *testing.T) {
	responseBytes, err := os.ReadFile("test_files/bbb_get_participants_not_ok.xml")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}
	_, err = parseBBBGetMeetingsResponseXML(responseBytes)
	if err == nil {
		t.Fatalf("Expected error but got nil")
	}
}

func TestCountParticipants(t *testing.T) {
	responseBytes, err := os.ReadFile("test_files/bbb_get_participants_ok.xml")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}
	response, err := parseBBBGetMeetingsResponseXML(responseBytes)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}
	count := countParticipants(response)
	if count != 2 {
		t.Fatalf("Expected count to be 2 but got %d", count)
	}
}

// Example from https://docs.bigbluebutton.org/development/api/#usage
func TestSignedBBBAPIRequest(t *testing.T) {
	serverUrl := "bbb.example.com"
	endpoint := "create"
	parameters := "name=Test+Meeting&meetingID=abc123&attendeePW=111222&moderatorPW=333444"
	apiToken := "639259d4-9dd8-4b25-bf01-95f9567eaf4b"
	expected := "https://bbb.example.com/bigbluebutton/api/create?name=Test+Meeting&meetingID=abc123&attendeePW=111222&moderatorPW=333444&checksum=1fcbb0c4fc1f039f73aa6d697d2db9ba7f803f17"
	got := signedBBBAPIRequest(serverUrl, endpoint, parameters, apiToken)
	if got != expected {
		t.Fatalf("Expected %s but got %s", expected, got)
	}
}

func testApplyRulesCPU(t *testing.T, bbbParticipants int, resourceState s.CpuResourceState, resources s.CpuResources, expected s.ScaleDirection) {
	bbbConfig := validBBBConfig
	bbbConfig.Resources.Cpu = &resources
	bbbState := &BBBServiceState{
		Name: "test-meeting",
	}
	bbbService := BBBService{
		Config: *bbbConfig,
		State:  *bbbState,
	}

	server := sampleBBBServer
	server.ResourceState.Cpu = &resourceState

	proposal := applyRules(server, bbbParticipants, bbbService)
	if proposal.Cpu.Direction != expected {
		t.Fatalf("Expected CPU scale direction to be %s but got %s", expected, proposal.Cpu.Direction)
	}
}

// Check that a server with below minimum resources is scaled up even if there are no participants
func TestApplyRulesRule1(t *testing.T) {
	resources := s.CpuResources{
		MinCores: 2,
		MaxCores: 4,
		MaxUsage: 0.5,
	}
	resourceState := s.CpuResourceState{
		CurrentCores: 1,
		CurrentUsage: 0,
	}
	testApplyRulesCPU(t, 0, resourceState, resources, s.ScaleUp)
}

// Check that a server with participants and above maximum usage is scaled up
func TestApplyRulesRule2ScaleUp(t *testing.T) {
	resources := s.CpuResources{
		MinCores: 2,
		MaxCores: 4,
		MaxUsage: 0.5,
	}
	resourceState := s.CpuResourceState{
		CurrentCores: 2,
		CurrentUsage: 0.6,
	}
	testApplyRulesCPU(t, 2, resourceState, resources, s.ScaleUp)
}

// Check that a server with 0 participants and minimum resources is not modified
func TestApplyRulesRule3NoChanges(t *testing.T) {
	resources := s.CpuResources{
		MinCores: 2,
		MaxCores: 4,
		MaxUsage: 0.5,
	}
	resourceState := s.CpuResourceState{
		CurrentCores: 2,
		CurrentUsage: 0,
	}
	testApplyRulesCPU(t, 0, resourceState, resources, s.ScaleNone)
}

// Check that a server with 0 participants and above minimum resources is scaled down
func TestApplyRulesRule3ScaleDown(t *testing.T) {
	resources := s.CpuResources{
		MinCores: 2,
		MaxCores: 4,
		MaxUsage: 0.5,
	}
	resourceState := s.CpuResourceState{
		CurrentCores: 3,
		CurrentUsage: 0,
	}
	testApplyRulesCPU(t, 0, resourceState, resources, s.ScaleDown)
}

// Check that a server with participants and below maximum usage is not modified
func TestApplyRulesRuleDefault(t *testing.T) {
	resources := s.CpuResources{
		MinCores: 2,
		MaxCores: 4,
		MaxUsage: 0.5,
	}
	resourceState := s.CpuResourceState{
		CurrentCores: 3,
		CurrentUsage: 0.4,
	}
	testApplyRulesCPU(t, 2, resourceState, resources, s.ScaleNone)
}
