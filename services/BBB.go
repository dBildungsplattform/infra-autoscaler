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

func (bbb BBBServiceState) Get_name() string {
	return bbb.Name
}

type BBBServiceConfig struct {
	CycleTimeSeconds int         `yaml:"cycle_time_seconds"`
	Resources        s.Resources `yaml:"resources"`
	ApiToken         string      `yaml:"api_token"`
}

func (bbb *BBBService) Get_state() s.ServiceState {
	return bbb.State
}

func (bbb *BBBService) Get_config() BBBServiceConfig {
	return bbb.Config
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
