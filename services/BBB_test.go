package services

import (
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
