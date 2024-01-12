package services

import (
	"fmt"
	s "scaler/shared"
)

type PostgresService struct {
	State  PostgresServiceState  `yaml:"-"`
	Config PostgresServiceConfig `yaml:"postgres_config"`
}

type PostgresServiceState struct {
	Name string
}

func (postgres PostgresServiceState) GetName() string {
	return postgres.Name
}

type PostgresServiceConfig struct {
	CycleTimeSeconds int         `yaml:"cycle_time_seconds"`
	Resources        s.Resources `yaml:"resources"`
}

func (postgres PostgresService) Init() error {
	return initMetricsExporter("postgres")
}

func (postgres *PostgresService) GetState() s.ServiceState {
	return postgres.State
}

func (postgres *PostgresService) GetConfig() PostgresServiceConfig {
	return postgres.Config
}

func (postgres PostgresService) GetResources() s.Resources {
	return postgres.Config.Resources
}

func (postgres PostgresService) GetCycleTimeSeconds() int {
	return postgres.Config.CycleTimeSeconds
}

func (postgres PostgresService) ComputeScalingProposal(s.ScaledObject) (s.ResourceScalingProposal, error) {
	return s.ResourceScalingProposal{}, nil // TODO: implement
}

func (service PostgresService) Validate() error {
	if err := service.Config.Validate(); err != nil {
		return err
	}
	if err := service.State.Validate(); err != nil {
		return err
	}
	return nil
}

func (state PostgresServiceState) Validate() error {
	return nil
}

func (config PostgresServiceConfig) Validate() error {
	if config.CycleTimeSeconds <= 0 {
		return fmt.Errorf("cycle time seconds must be greater than 0")
	}
	if err := config.Resources.Validate(); err != nil {
		return err
	}
	return nil
}
