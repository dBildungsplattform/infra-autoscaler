package shared

// Interface that implements the scaling logic for a service and communicates with it if needed
type Service interface {
	Validate() error
	Init() error
	GetResources() Resources
	GetCycleTimeSeconds() int
	ComputeScalingProposal(ScaledObject) (ResourceScalingProposal, error)
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
