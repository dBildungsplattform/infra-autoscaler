package services

import (
	"os"
	s "scaler/shared"
	"testing"
)

func TestValidateConfigOK(t *testing.T) {
	bbbConfig := &BBBServiceConfig{
		CycleTimeSeconds: 60,
		Resources: s.Resources{
			Cpu: &s.CpuResources{
				MinCores: 1,
				MaxCores: 2,
				MaxUsage: 0.5,
			},
			Memory: &s.MemoryResources{
				MinBytes: 1024,
				MaxBytes: 2048,
				MaxUsage: 0.5,
			},
		},
		ApiToken: "1234567890",
	}
	s.ValidatePass(t, bbbConfig)
}

func TestValidateConfigNotOK(t *testing.T) {
	bbbConfig := &BBBServiceConfig{}
	s.ValidateFail(t, bbbConfig)
}

func TestParseConfigOK(t *testing.T) {
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
