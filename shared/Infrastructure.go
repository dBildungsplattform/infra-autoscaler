package shared

/*** Infrastructure definition ***/
type InfrastructureType int

const (
	_ InfrastructureType = iota
	Server
	Kubernetes
)
