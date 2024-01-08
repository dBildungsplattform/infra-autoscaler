package services

import "github.com/prometheus/client_golang/prometheus"

var (
	errorsTotalCounter prometheus.Counter
)

func initMetricsExporter(serviceName string) error {
	constLabels := map[string]string{
		"component_type": serviceName,
		"component":      "service",
	}
	errorsTotalCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name:        "autoscaler_component_errors_total",
		Help:        "The total number of errors encountered by a component of the autoscaler",
		ConstLabels: constLabels,
	})
	metrics := []prometheus.Collector{errorsTotalCounter}
	for _, metric := range metrics {
		if err := prometheus.Register(metric); err != nil {
			return err
		}
	}
	errorsTotalCounter.Add(0)
	return nil
}
