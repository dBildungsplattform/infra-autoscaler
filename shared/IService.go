package shared

/*** Service definition ***/
type Service interface {
	Validate() error
}

type ServiceState interface {
	Get_name() string
}

type ServiceConfig interface {
	Get_provider_type() ProviderType
}

type ServiceType string

const (
	BBB      = "BBB"
	Postgres = "Postgres"
)
