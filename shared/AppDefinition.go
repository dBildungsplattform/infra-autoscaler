package shared

import "fmt"

type AppDefinition struct {
	Name              string            `yaml:"app_name"`
	Stage             Stage             `yaml:"stage"`
	ServiceType       ServiceType       `yaml:"service_type"`
	ProviderType      ProviderType      `yaml:"provider_type"`
	MetricsSourceType MetricsSourceType `yaml:"metrics_source_type"`
}

type Stage string

const (
	DevStage  = "dev"
	ProdStage = "prod"
)

func (a AppDefinition) Validate() error {
	if a.Name == "" {
		return fmt.Errorf("AppDefinition.Name is empty")
	}
	if a.Stage == "" {
		return fmt.Errorf("AppDefinition.Stage is empty")
	} else if !(a.Stage == DevStage || a.Stage == ProdStage) {
		return fmt.Errorf("AppDefinition.Stage is invalid")
	}
	if a.ServiceType == "" {
		return fmt.Errorf("AppDefinition.Type is empty")
	}
	if a.ProviderType == "" {
		return fmt.Errorf("AppDefinition.Type is empty")
	}
	if a.MetricsSourceType == "" {
		return fmt.Errorf("AppDefinition.MetricsSourceType is empty")
	}
	return nil
}
