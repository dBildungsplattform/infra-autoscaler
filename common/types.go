package common

type Service interface {
	Init(*ServiceDefinition)
}

type ServiceState interface {
	Get_name() string
}

type ServiceConfig interface {
	Get_name() string
	Get_provider_type() ProviderType
	Get_infrastructure_type() InfrastructureType
}

type InfrastructureType int

const (
	_ InfrastructureType = iota
	Server
	Kubernetes
)

type Provider interface {
	Get_login_id() string
	Get_login_secret() string
	Get_type() ProviderType
	Get_name() string
}

type ProviderDefinition struct {
	Name string
	Type ProviderType
}

type ProviderType string

const (
	Ionos = "Ionos"
)

type ServiceDefinition struct {
	Name string
	Type ServiceType
}

type ServiceType string

const (
	BBB      = "BBB"
	Postgres = "Postgres"
)
