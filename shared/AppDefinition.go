package shared

import "fmt"

type AppDefinition struct {
	Name         string       `yaml:"app_name"`
	ServiceType  ServiceType  `yaml:"service_type"`
	ProviderType ProviderType `yaml:"provider_type"`
	MetricsType  MetricsType  `yaml:"metrics_type"`
}

func (a AppDefinition) Validate() error {
	if a.Name == "" {
		return fmt.Errorf("AppDefinition.Name is empty")
	}
	if a.ServiceType == "" {
		return fmt.Errorf("AppDefinition.ServiceType is empty")
	}
	if a.ProviderType == "" {
		return fmt.Errorf("AppDefinition.ProviderType is empty")
	}
	if a.MetricsType == "" {
		return fmt.Errorf("AppDefinition.MetricsType is empty")
	}
	return nil
}
