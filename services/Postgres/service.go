package Postgres

import (
	"fmt"
	c "scaler/common"
)

type PostgresService struct {
	name   c.ServiceName
	state  PostgresServiceState
	config PostgresServiceConfig
}

type PostgresServiceState struct {
	Name c.ServiceName
}

func (postgres PostgresServiceState) Get_name() c.ServiceName {
	return postgres.Name
}

type PostgresServiceConfig struct {
	Name         string
	Type         c.InfrastructureType
	ProviderName c.InfrastructureProviderName
	Provider     c.IonosProvider
}

func (postgres PostgresServiceConfig) Get_name() string {
	return postgres.Name
}

func (postgres PostgresServiceConfig) Get_type() c.InfrastructureType {
	return postgres.Type
}

func (postgres PostgresServiceConfig) Get_provider_name() c.InfrastructureProviderName {
	return postgres.ProviderName
}

func (postgres PostgresService) Init(name c.ServiceName) {
	fmt.Println("Initializing Postgres service")
	postgres.name = name
	postgres.state = PostgresServiceState{Name: name}
	postgres.config = PostgresServiceConfig(*load_config())
	fmt.Println("Postgres service initialized")
	fmt.Printf("Config: \n %+v \n", postgres.config)
}

func (postgres *PostgresService) Get_state() c.ServiceState {
	return postgres.state
}

func (postgres *PostgresService) Get_config() PostgresServiceConfig {
	return postgres.config
}

func load_config() *PostgresServiceConfig {
	return &PostgresServiceConfig{
		Name:     "Postgres",
		Type:     c.Server,
		Provider: c.Load_provider(c.Ionos),
	}
}
