package shared

type ScaledObject interface {
	GetType() ScaledObjectType
	GetName() string
	GetResourceState() ResourceState
	SetResourceState(resourceState ResourceState)
}

type ScaledObjectType string

const (
	ServerType  = "Server"
	ClusterType = "Cluster"
)
