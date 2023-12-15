package shared

import "fmt"

/*** Metrics definition ***/
type MetricsSource interface {
	Validate() error
	GetServerCpuUsage(string) (float32, error)
	GetServerMemoryUsage(string) (float32, error)
	GetClusterCpuUsage(string) (float32, error)
	GetClusterMemoryUsage(string) (float32, error)
}

type MetricsSourceType string

const (
	Prometheus = "Prometheus"
)

func (m MetricsSourceType) Validate() error {
	switch m {
	case Prometheus:
		return nil
	default:
		return fmt.Errorf("unknown metrics type: %s", m)
	}
}
