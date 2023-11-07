package BBB

import (
	"fmt"
	c "scaler/common"
)

type BBBService struct {
	name   c.ServiceName
	state  BBBServiceState
	config BBBServiceConfig
}

type BBBServiceState struct {
	Name c.ServiceName
}

func (bbb BBBServiceState) Get_name() c.ServiceName {
	return bbb.Name
}

type BBBServiceConfig struct {
	Name         string
	Type         c.InfrastructureType
	ProviderName c.InfrastructureProviderName
	Provider     c.IonosProvider
}

func (bbb BBBServiceConfig) Get_name() string {
	return bbb.Name
}

func (bbb BBBServiceConfig) Get_type() c.InfrastructureType {
	return bbb.Type
}

func (bbb BBBServiceConfig) Get_provider_name() c.InfrastructureProviderName {
	return bbb.ProviderName
}

func (bbb BBBService) Init(name c.ServiceName) {
	fmt.Println("Initializing BBB service")
	bbb.name = name
	bbb.state = BBBServiceState{Name: name}
	bbb.config = BBBServiceConfig(*load_config())
	fmt.Println("BBB service initialized")
	fmt.Printf("Config: \n %+v \n", bbb.config)
}

func (bbb *BBBService) Get_state() c.ServiceState {
	return bbb.state
}

func (bbb *BBBService) Get_config() BBBServiceConfig {
	return bbb.config
}

func load_config() *BBBServiceConfig {
	return &BBBServiceConfig{
		Name:     "BBB",
		Type:     c.Server,
		Provider: c.Load_provider(c.Ionos),
	}
}
