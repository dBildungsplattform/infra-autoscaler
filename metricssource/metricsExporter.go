package metricssource

import "github.com/prometheus/client_golang/prometheus"

var (
	errorsTotalCounter prometheus.Counter
)

func initMetricsExporter(serviceName string) error {
	constLabels := map[string]string{
		"metrics": serviceName,
	}
	errorsTotalCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name:        "autoscaler_metrics_errors_total",
		Help:        "The total number of errors when communicating with the metrics API",
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
