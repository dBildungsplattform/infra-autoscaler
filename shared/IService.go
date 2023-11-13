package shared

/*** Service definition ***/
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

type ServiceDefinition struct {
	Name string
	Type ServiceType
}

type ServiceType string

const (
	BBB      = "BBB"
	Postgres = "Postgres"
)
