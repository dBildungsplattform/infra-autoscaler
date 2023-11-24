package shared

/*** Service definition ***/
type Service interface {
	Validate() error
	GetResources() Resources
	ShouldScale(cores int, memory int) (ScaleResource, error)
}

type ServiceState interface {
	GetName() string
}

type ServiceConfig interface {
	GetProviderType() ProviderType
}

type ServiceType string

const (
	BBB      = "BBB"
	Postgres = "Postgres"
)
