package common

type Service interface {
	Init(ServiceName)
}

type ServiceState interface {
	Get_name() ServiceName
}

type ServiceConfig interface {
	Get_name() ServiceName
	Get_type() InfrastructureType
	Get_provider_name() InfrastructureProviderName
}

type IonosProvider struct {
	ProviderName  InfrastructureProviderName
	Username      string
	Password      string
	DatacenterIds []string
}

type InfrastructureType int

const (
	_ InfrastructureType = iota
	Server
	Kubernetes
)

type InfrastructureProvider interface {
	IonosProvider
}

type InfrastructureProviderName string

const (
	Ionos = "Ionos"
)

type ServiceDefinition struct {
	Name ServiceName
	Type ServiceType
}

type ServiceName string

type ServiceType string

const (
	BBB      = "BBB"
	Postgres = "Postgres"
)
