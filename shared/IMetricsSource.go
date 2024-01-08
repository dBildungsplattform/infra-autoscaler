package shared

import "fmt"

/*** Metrics definition ***/
type MetricsSource interface {
	Validate() error
	GetCpuUsage(ScaledObject) (float32, error)
	GetMemoryUsage(ScaledObject) (float32, error)
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
