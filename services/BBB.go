package services

import (
	"fmt"
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
	CycleTimeSeconds int         `yaml:"cycle_time_seconds"`
	Resources        s.Resources `yaml:"resources"`
	ApiToken         string      `yaml:"api_token"`
}

func (bbb *BBBService) GetState() s.ServiceState {
	return bbb.State
}

func (bbb *BBBService) GetConfig() BBBServiceConfig {
	return bbb.Config
}

func (bbb BBBService) GetResources() s.Resources {
	return bbb.Config.Resources
}

func (bbb BBBService) ShouldScale(cores int, memory int) (s.ScaleResource, error) {
	targetResource := s.ScaleResource{
		Cpu: s.ScaleOp{
			Direction: s.ScaleNone,
			Amount:    0,
		},
		Mem: s.ScaleOp{
			Direction: s.ScaleNone,
			Amount:    0,
		},
	}

	// Scaling cores
	coresMaxThreshold := int(float32(bbb.Config.Resources.Cpu.MaxCores) * bbb.Config.Resources.Cpu.MaxUsage)
	coresMinThreshold := int(float32(bbb.Config.Resources.Cpu.MinCores) * bbb.Config.Resources.Cpu.MinUsage)

	if cores >= bbb.Config.Resources.Cpu.MinCores && cores <= coresMaxThreshold {
		targetResource.Cpu.Direction = s.ScaleNone
	}
	if cores < bbb.Config.Resources.Cpu.MinCores || cores > coresMaxThreshold {
		targetResource.Cpu.Direction = s.ScaleUp
	}
	if cores < coresMinThreshold {
		targetResource.Cpu.Direction = s.ScaleDown
	}

	// Scaling memory
	memoryMaxThreshold := int(float32(bbb.Config.Resources.Memory.MaxBytes) * bbb.Config.Resources.Memory.MaxUsage)
	memoryMinThreshold := int(float32(bbb.Config.Resources.Memory.MinBytes) * bbb.Config.Resources.Memory.MinUsage)

	if memory >= bbb.Config.Resources.Memory.MinBytes && memory <= memoryMaxThreshold {
		targetResource.Mem.Direction = s.ScaleNone
	}
	if memory < bbb.Config.Resources.Memory.MinBytes || memory > memoryMaxThreshold {
		targetResource.Mem.Direction = s.ScaleUp
	}
	if memory < memoryMinThreshold {
		targetResource.Mem.Direction = s.ScaleDown
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
