package shared

type ScaledObject interface {
	GetType() ScaledObjectType
	GetName() string
	GetResourceState() ResourceState
	SetResourceState(resourceState ResourceState)
	IsReady() bool
}

type ScaledObjectType string

const (
	ServerType  = "Server"
	ClusterType = "Cluster"
)
