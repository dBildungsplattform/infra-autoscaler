package shared

import "fmt"

/*** Metrics definition ***/
type Metrics interface {
	Validate() error
}

type MetricsType string

const (
	Prometheus = "Prometheus"
)

func (m MetricsType) Validate() error {
	switch m {
	case Prometheus:
		return nil
	default:
		return fmt.Errorf("unknown metrics type: %s", m)
	}
}
