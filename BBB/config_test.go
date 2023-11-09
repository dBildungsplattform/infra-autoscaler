package BBB

import (
	"scaler/ionos"
	"scaler/scaler/config"
	v "scaler/validater"
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
		CloudProvider: config.CloudProvider{
			Ionos: &ionos.CloudProvider{
				Username: "username",
				Password: "password",
				ServerSource: &ionos.ServerSource{
					Dynamic: &ionos.ServerDynamicSource{
						DatacenterIds:   []string{"1234567890"},
						ServerNameRegex: ".*",
					},
				},
			},
		},
	}
	v.ValidatePass(t, bbbConfig)
}

func TestValidateConfigNotOK(t *testing.T) {
	bbbConfig := &Config{}
	v.ValidateFail(t, bbbConfig)
}

func TestParseConfigOK(t *testing.T) {
	c, ok := config.LoadConfig[Config]("test_files/bbb_config_ok.yml")
	if ok != nil {
		t.Fatalf("Failed to parse config: %v", ok)
	}
	v.ValidatePass(t, c)
}

func TestParseConfigNotOK(t *testing.T) {
	_, ok := config.LoadConfig[Config]("test_files/bbb_config_not_ok.yml")
	if ok == nil {
		t.Fatalf("Expected error but got nil")
	}
}
