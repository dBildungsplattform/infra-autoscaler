package BBB

import (
	"fmt"
	"scaler/scaler/config"
)

type Resources struct {
	Cpu    *config.CpuResources
	Memory *config.MemoryResources
}

type Config struct {
	CycleTimeSeconds int `yaml:"cycle_time_seconds"`
	Resources        Resources
	BBB              struct {
		ApiToken string `yaml:"api_token"`
	}
	CloudProvider config.CloudProvider `yaml:"cloud_provider"`
}

func (r Resources) Validate() error {
	if cpu := r.Cpu; cpu != nil {
		if err := cpu.Validate(); err != nil {
			return err
		}
	}
	if memory := r.Memory; memory != nil {
		if err := memory.Validate(); err != nil {
			return err
		}
	}
	if r.Cpu == nil && r.Memory == nil {
		return fmt.Errorf("resources.cpu and resources.memory are nil, at least one must be set")
	}
	return nil
}

func (config Config) Validate() error {
	if err := config.Resources.Validate(); err != nil {
		return err
	}
	if config.BBB.ApiToken == "" {
		return fmt.Errorf("bbb.api_token is empty")
	}
	if err := config.CloudProvider.Validate(); err != nil {
		return err
	}
	return nil
}
