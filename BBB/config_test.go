package BBB

import (
	"scaler/scaler/config"
	"testing"
)

func TestValidateConfigOK(t *testing.T) {
	bbbConfig := &Config{
		CycleTimeSeconds: 60,
		Resources: Resources{
			Cpu: &config.CpuResources{
				MinCores: 1,
				MaxCores: 2,
				MaxUsage: 0.5,
			},
			Memory: &config.MemoryResources{
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
		InstancesSource: InstanceSources{
			Ionos: &config.IonosServerInstancesSource{
				DatacenterIds:   []string{"1234567890"},
				ServerNameRegex: ".*",
			},
		},
	}
	config.ValidatePass(t, bbbConfig)
}

func TestValidateConfigNotOK(t *testing.T) {
	bbbConfig := &Config{}
	config.ValidateFail(t, bbbConfig)
}

func TestParseConfigOK(t *testing.T) {
	c, ok := config.LoadConfig[Config]("test_files/bbb_config_ok.yml")
	if ok != nil {
		t.Fatalf("Failed to parse config: %v", ok)
	}
	config.ValidatePass(t, c)
}

func TestParseConfigNotOK(t *testing.T) {
	_, ok := config.LoadConfig[Config]("test_files/bbb_config_not_ok.yml")
	if ok == nil {
		t.Fatalf("Expected error but got nil")
	}
}
