package shared

type ScaledObject interface {
	GetType() ScaledObjectType
	GetName() string
}

type ScaledObjectType string

const (
	ServerType  = "Server"
	ClusterType = "Cluster"
)
