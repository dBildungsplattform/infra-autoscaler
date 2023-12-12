package shared

import "fmt"

type AppDefinition struct {
	Name                string              `yaml:"app_name"`
	Stage               Stage               `yaml:"stage"`
	ScalingMode         ScalingMode         `yaml:"scaling_mode"`
	DirectScalingConfig DirectScalingConfig `yaml:"direct_scaling_config"`
	ServiceType         ServiceType         `yaml:"service_type"`
	ProviderType        ProviderType        `yaml:"provider_type"`
	MetricsSourceType   MetricsSourceType   `yaml:"metrics_source_type"`
	MetricsExporterPort IntFromEnv          `yaml:"metrics_exporter_port"`
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

type DirectScalingConfig struct {
	CpuIncrease int32 `yaml:"cpu_increase"`
	CpuDecrease int32 `yaml:"cpu_decrease"`
	MemIncrease int32 `yaml:"mem_increase"`
	MemDecrease int32 `yaml:"mem_decrease"`
}

func (a AppDefinition) Validate() error {
	if a.Name == "" {
		return fmt.Errorf("AppDefinition.Name is empty")
	}
	if a.Stage == "" {
		return fmt.Errorf("AppDefinition.Stage is empty")
	} else if !(a.Stage == DevStage || a.Stage == ProdStage) {
		return fmt.Errorf("AppDefinition.Stage is invalid")
	}
	if !(a.ScalingMode == DirectScaling || a.ScalingMode == HeuristicScaling || a.ScalingMode == "") {
		return fmt.Errorf("AppDefinition.ScalingMode is invalid")
	}
	if a.ScalingMode == DirectScaling {
		if a.DirectScalingConfig.CpuIncrease < 0 || a.DirectScalingConfig.CpuDecrease > 0 {
			return fmt.Errorf("AppDefinition.DirectScalingConfig.CpuIncrease %d is invalid", a.DirectScalingConfig.CpuIncrease)
		}
		if a.DirectScalingConfig.MemIncrease < 0 || a.DirectScalingConfig.MemDecrease > 0 {
			return fmt.Errorf("AppDefinition.DirectScalingConfig.MemIncrease %d is invalid", a.DirectScalingConfig.MemIncrease)
		}
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
