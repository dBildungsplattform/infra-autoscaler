package BBB

import (
	"fmt"
	s "scaler/shared"
)

type BBBService struct {
	Name   string
	state  BBBServiceState
	config BBBServiceConfig
}

type BBBServiceState struct {
	Name string
}

func (bbb BBBServiceState) Get_name() string {
	return bbb.Name
}

type BBBServiceConfig struct {
	ServiceDef       s.ServiceDefinition
	ProviderType     s.ProviderType `yaml:"provider_type"`
	InfraType        s.InfrastructureType
	CycleTimeSeconds int             `yaml:"cycle_time_seconds"`
	ServerSource     *s.ServerSource `yaml:"server_source"`
	Resources        s.Resources
	BBB              struct {
		ApiToken string `yaml:"api_token"`
	}
}

func (bbb BBBServiceConfig) Get_name() string {
	return bbb.ServiceDef.Name
}

func (bbb BBBServiceConfig) Get_provider_type() s.ProviderType {
	return bbb.ProviderType
}

func (bbb BBBServiceConfig) Get_infrastructure_type() s.InfrastructureType {
	return bbb.InfraType
}

func (bbb BBBService) Init(sd *s.ServiceDefinition) {
	fmt.Println("Initializing BBB service")
	bbb.Name = sd.Name
	bbb.state = BBBServiceState{Name: sd.Name} // TODO: Load proper state from provider API
	bbb.config = BBBServiceConfig(*load_config(sd))
	fmt.Println("BBB service initialized")
	fmt.Printf("Config: \n %+v \n", bbb.config)
}

func (bbb *BBBService) Get_state() s.ServiceState {
	return bbb.state
}

func (bbb *BBBService) Get_config() BBBServiceConfig {
	return bbb.config
}

func load_config(sd *s.ServiceDefinition) *BBBServiceConfig {
	return &BBBServiceConfig{
		ServiceDef:   *sd,
		ProviderType: s.Ionos,
		InfraType:    s.Server,
	}
}

func (config BBBServiceConfig) Validate() error {
	if err := config.Resources.Validate(); err != nil {
		return err
	}
	if config.BBB.ApiToken == "" {
		return fmt.Errorf("bbb.api_token is empty")
	}
	if err := config.ProviderType.Validate(); err != nil {
		return err
	}
	ss := config.ServerSource
	if ss == nil {
		return fmt.Errorf("instances_source is nil")
	}
	if err := ss.Validate(); err != nil {
		return err
	}
	return nil
}
