package shared

import "fmt"

type AppDefinition struct {
	Name                string            `yaml:"app_name"`
	Stage               Stage             `yaml:"stage"`
	ScalingMode         ScalingMode       `yaml:"scaling_mode"`
	ServiceType         ServiceType       `yaml:"service_type"`
	ProviderType        ProviderType      `yaml:"provider_type"`
	MetricsSourceType   MetricsSourceType `yaml:"metrics_source_type"`
	MetricsExporterPort IntFromEnv        `yaml:"metrics_exporter_port"`
}

type Stage string

const (
	DevStage  = "dev"
	ProdStage = "prod"
)

type ScalingMode string

const (
	DirectScaling    = "direct"
	HeuristicScaling = "heuristic"
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
	if a.ScalingMode == "" {
		return fmt.Errorf("AppDefinition.ScalingMode is empty")
	} else if !(a.ScalingMode == DirectScaling || a.ScalingMode == HeuristicScaling) {
		return fmt.Errorf("AppDefinition.ScalingMode is invalid")
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
	if a.MetricsExporterPort < 0 || a.MetricsExporterPort > 65535 {
		return fmt.Errorf("AppDefinition.MetricsExporterPort %d is invalid", a.MetricsExporterPort)
	}
	return nil
}
