package Postgres

import (
	"fmt"
	c "scaler/common"
)

type PostgresService struct {
	name   string
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
	ServiceDef   c.ServiceDefinition
	ProviderType c.ProviderType
	InfraType    c.InfrastructureType
}

func (postgres PostgresServiceConfig) Get_name() string {
	return postgres.ServiceDef.Name
}

func (postgres PostgresServiceConfig) Get_provider_type() c.ProviderType {
	return postgres.ProviderType
}

func (postgres PostgresServiceConfig) Get_infrastructure_type() c.InfrastructureType {
	return postgres.InfraType
}

func (postgres PostgresService) Init(sd *c.ServiceDefinition) {
	fmt.Println("Initializing Postgres service")
	postgres.name = sd.Name
	postgres.state = PostgresServiceState{Name: sd.Name}
	postgres.config = PostgresServiceConfig(*load_config(sd))
	fmt.Println("Postgres service initialized")
	fmt.Printf("Config: \n %+v \n", postgres.config)
}

func (postgres *PostgresService) Get_state() c.ServiceState {
	return postgres.state
}

func (postgres *PostgresService) Get_config() PostgresServiceConfig {
	return postgres.config
}

func load_config(sd *c.ServiceDefinition) *PostgresServiceConfig {
	return &PostgresServiceConfig{
		ServiceDef:   *sd,
		ProviderType: c.Ionos,
		InfraType:    c.Server, // TODO: Additional 'managed service' type?
	}
}
