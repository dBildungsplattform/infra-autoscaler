package BBB

import (
	"fmt"
	c "scaler/common"
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
	ServiceDef   c.ServiceDefinition
	ProviderType c.ProviderType
	InfraType    c.InfrastructureType
}

func (bbb BBBServiceConfig) Get_name() string {
	return bbb.ServiceDef.Name
}

func (bbb BBBServiceConfig) Get_provider_type() c.ProviderType {
	return bbb.ProviderType
}

func (bbb BBBServiceConfig) Get_infrastructure_type() c.InfrastructureType {
	return bbb.InfraType
}

func (bbb BBBService) Init(sd *c.ServiceDefinition) {
	fmt.Println("Initializing BBB service")
	bbb.Name = sd.Name
	bbb.state = BBBServiceState{Name: sd.Name} // TODO: Load proper state from provider API
	bbb.config = BBBServiceConfig(*load_config(sd))
	fmt.Println("BBB service initialized")
	fmt.Printf("Config: \n %+v \n", bbb.config)
}

func (bbb *BBBService) Get_state() c.ServiceState {
	return bbb.state
}

func (bbb *BBBService) Get_config() BBBServiceConfig {
	return bbb.config
}

func load_config(sd *c.ServiceDefinition) *BBBServiceConfig {
	return &BBBServiceConfig{
		ServiceDef:   *sd,
		ProviderType: c.Ionos,
		InfraType:    c.Server,
	}
}
