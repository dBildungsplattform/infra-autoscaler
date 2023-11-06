package scaler

type Service interface {
	Init()
	Get_state() ServiceState
	Get_config() ServiceConfig[any]
}

type ServiceState struct {
	Name string
}

type ServiceConfig[provider any] struct {
	Name     string
	Type     InfrastructureType
	Provider provider
}

type IonosProvider struct {
	ProviderName  InfrastructureProvider
	Username      string
	Password      string
	DatacenterIds []string
}

type InfrastructureType int

const (
	UndefinedInfra InfrastructureType = iota
	Server
	Kubernetes
)

type InfrastructureProvider int

const (
	UndefinedProvider InfrastructureProvider = iota
	Ionos
)

func Load_config() *ServiceConfig[any] {
	// TODO load config from file, env, CLI args, etc.
	Ionos := IonosProvider{
		ProviderName:  Ionos,
		Username:      "username",
		Password:      "password",
		DatacenterIds: []string{"datacenter1", "datacenter2"},
	}
	return &ServiceConfig[any]{
		Name:     "BBB",
		Type:     Server,
		Provider: Ionos,
	}
}
