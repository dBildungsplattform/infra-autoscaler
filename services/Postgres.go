package services

import (
	"fmt"
	s "scaler/shared"
)

type PostgresService struct {
	state  PostgresServiceState
	config PostgresServiceConfig
}

type PostgresServiceState struct {
	Name string
}

func (postgres PostgresServiceState) Get_name() string {
	return postgres.Name
}

type PostgresServiceConfig struct {
}

func (postgres *PostgresService) Get_state() s.ServiceState {
	return postgres.state
}

func (postgres *PostgresService) Get_config() PostgresServiceConfig {
	return postgres.config
}

func (service PostgresService) Validate() error {
	if err := service.config.Validate(); err != nil {
		return err
	}
	if err := service.state.Validate(); err != nil {
		return err
	}
	return nil
}

func (state PostgresServiceState) Validate() error {
	if state.Name == "" {
		return fmt.Errorf("name is empty")
	}
	return nil
}

func (config PostgresServiceConfig) Validate() error {
	return nil
}
