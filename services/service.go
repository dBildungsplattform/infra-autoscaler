package services

import "github.com/prometheus/client_golang/prometheus"

var (
	errorsMetric = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "autoscaler_service_errors_total",
		Help: "The total number of errors when communicating with the service API",
	}, []string{"service"})
)

func registerMetrics(serviceName string) error {
	metrics := []prometheus.Collector{errorsMetric}
	for _, metric := range metrics {
		if err := prometheus.Register(metric); err != nil {
			return err
		}
	}
	errorsMetric.WithLabelValues(serviceName).Add(0)
	return nil
}
