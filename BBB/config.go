package BBB

import (
	"fmt"
	"scaler/scaler/config"
)

type Resources struct {
	Cpu    *config.CpuResources
	Memory *config.MemoryResources
}

type InstanceSources struct {
	Ionos       *config.IonosServerInstancesSource       `yaml:"ionos"`
	InlineIonos *config.InlineIonosServerInstancesSource `yaml:"inline_ionos"`
}

type Config struct {
	CycleTimeSeconds int `yaml:"cycle_time_seconds"`
	Resources        Resources
	BBB              struct {
		ApiToken string `yaml:"api_token"`
	}
	InstancesSource InstanceSources `yaml:"instances_source"`
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

func (is InstanceSources) Validate() error {
	ionos, inline := is.Ionos, is.InlineIonos
	if ionos == nil && inline == nil {
		return fmt.Errorf("instances_source.ionos and instances_source.inline_ionos are nil, one must be set")
	}
	if ionos != nil && inline != nil {
		return fmt.Errorf("instances_source.ionos and instances_source.inline_ionos are both set, only one must be set")
	}
	if ionos != nil {
		err := ionos.Validate()
		if err != nil {
			return err
		}
	}
	if inline != nil {
		err := inline.Validate()
		if err != nil {
			return err
		}
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
	if err := config.InstancesSource.Validate(); err != nil {
		return err
	}
	return nil
}
