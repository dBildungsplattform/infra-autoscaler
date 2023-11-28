package services

import "github.com/prometheus/client_golang/prometheus"

var (
	errorsTotalCounter prometheus.Counter
)

func initMetricsExporter(serviceName string) error {
	constLabels := map[string]string{
		"service": serviceName,
	}
	errorsTotalCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name:        "autoscaler_service_errors_total",
		Help:        "The total number of errors when communicating with the service API",
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
