package services

import (
	"os"
	s "scaler/shared"
	"testing"
)

func TestValidateConfigOk(t *testing.T) {
	postgresConfig := &PostgresServiceConfig{
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
	}
	s.ValidatePass(t, postgresConfig)
}

func TestValidateConfigNotOk(t *testing.T) {
	postgresConfig := &PostgresServiceConfig{}
	s.ValidateFail(t, postgresConfig)
}

func TestParseConfigOk(t *testing.T) {
	os.Setenv("IONOS_TOKEN", "1234567890")
	defer os.Unsetenv("IONOS_TOKEN")

	config, ok := s.OpenConfig("test_files/postgres_config_ok.yml")
	if ok != nil {
		t.Fatalf("Failed to open config: %v", ok)
	}

	c, ok := s.LoadConfig[PostgresService](config)
	if ok != nil {
		t.Fatalf("Failed to parse config: %v", ok)
	}

	s.ValidatePass(t, c)
}

func TestParseConfigNotOk(t *testing.T) {
	config, ok := s.OpenConfig("test_files/postgres_config_not_ok.yml")
	if ok != nil {
		t.Fatalf("Failed to open config: %v", ok)
	}

	_, ok = s.LoadConfig[PostgresService](config)
	if ok == nil {
		t.Fatalf("Expected error but got nil")
	}
}
