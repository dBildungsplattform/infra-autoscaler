package BBB

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
		BBB: struct {
			ApiToken string `yaml:"api_token"`
		}{
			ApiToken: "1234567890",
		},
		ServiceDef: s.ServiceDefinition{
			Name: "bbb",
			Type: s.BBB,
		},
		ProviderType: s.Ionos,
		InfraType:    s.Server,
		ServerSource: &s.ServerSource{
			Dynamic: &s.ServerDynamicSource{
				DatacenterIds:   []string{"1234567890"},
				ServerNameRegex: ".*",
			},
		},
	}
	s.ValidatePass(t, bbbConfig)
}

func TestValidateConfigNotOK(t *testing.T) {
	bbbConfig := &BBBServiceConfig{}
	s.ValidateFail(t, bbbConfig)
}

func TestParseConfigOK(t *testing.T) {
	c, ok := s.LoadConfig[BBBServiceConfig]("test_files/bbb_config_ok.yml")
	if ok != nil {
		t.Fatalf("Failed to parse config: %v", ok)
	}
	s.ValidatePass(t, c)
}

func TestParseConfigNotOK(t *testing.T) {
	_, ok := s.LoadConfig[BBBServiceConfig]("test_files/bbb_config_not_ok.yml")
	if ok == nil {
		t.Fatalf("Expected error but got nil")
	}
}
