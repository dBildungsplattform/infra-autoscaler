package Postgres

import (
	"fmt"
	s "scaler/shared"
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
	ServiceDef   s.ServiceDefinition
	ProviderType s.ProviderType
	InfraType    s.InfrastructureType
}

func (postgres PostgresServiceConfig) Get_name() string {
	return postgres.ServiceDef.Name
}

func (postgres PostgresServiceConfig) Get_provider_type() s.ProviderType {
	return postgres.ProviderType
}

func (postgres PostgresServiceConfig) Get_infrastructure_type() s.InfrastructureType {
	return postgres.InfraType
}

func (postgres PostgresService) Init(sd *s.ServiceDefinition) {
	fmt.Println("Initializing Postgres service")
	postgres.name = sd.Name
	postgres.state = PostgresServiceState{Name: sd.Name}
	postgres.config = PostgresServiceConfig(*load_config(sd))
	fmt.Println("Postgres service initialized")
	fmt.Printf("Config: \n %+v \n", postgres.config)
}

func (postgres *PostgresService) Get_state() s.ServiceState {
	return postgres.state
}

func (postgres *PostgresService) Get_config() PostgresServiceConfig {
	return postgres.config
}

func load_config(sd *s.ServiceDefinition) *PostgresServiceConfig {
	return &PostgresServiceConfig{
		ServiceDef:   *sd,
		ProviderType: s.Ionos,
		InfraType:    s.Server, // TODO: Additional 'managed service' type?
	}
}
