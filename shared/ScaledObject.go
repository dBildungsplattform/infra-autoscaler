package shared

type ScaledObject interface {
	GetType() ScaledObjectType
}

type ScaledObjectType string

const (
	ServerType  = "Server"
	ClusterType = "Cluster"
)
