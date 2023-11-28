package core

import (
	"net/http"
	s "scaler/shared"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	cyclesCounter           prometheus.Counter
	capacityTotalGauge      *prometheus.GaugeVec
	capacityUsedGauge       *prometheus.GaugeVec
	maxScaledInstancesGauge *prometheus.GaugeVec
	lastScaleTimeGauge      prometheus.Gauge
)

func initMetricsExporter() error {
	cyclesCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "autoscaler_cycles_total",
		Help: "The total number of cycles the autoscaler has run",
	})
	capacityTotalGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "autoscaler_capacity_total",
		Help: "The maximum amount of resource that can be used",
	}, []string{"resource_type"})
	capacityUsedGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "autoscaler_capacity_used",
		Help: "The amount of resource that is currently used",
	}, []string{"resource_type"})
	maxScaledInstancesGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "autoscaler_max_scaled_instances",
		Help: "The amount of instances that are scaled to the maximum for a given resource type",
	}, []string{"resource_type"})
	lastScaleTimeGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "autoscaler_last_scale_time",
		Help: "The time of the last scale operation",
	})
	metrics := []prometheus.Collector{cyclesCounter, capacityTotalGauge, capacityUsedGauge, maxScaledInstancesGauge, lastScaleTimeGauge}
	for _, metric := range metrics {
		if err := prometheus.Register(metric); err != nil {
			return err
		}
	}
	cyclesCounter.Add(0)
	return nil
}

func (sc ScalerApp) calculateMetrics(servers []s.Server) {
	resources := (*sc.service).GetResources()

	cpuUsedCapacity := 0
	memoryUsedCapacity := 0
	cpuMaxScaledInstances := 0
	memoryMaxScaledInstances := 0

	capacityTotalGauge.Reset()
	capacityUsedGauge.Reset()

	for _, server := range servers {
		cpuUsedCapacity += int(server.ServerCpu)
		memoryUsedCapacity += int(server.ServerRam)

		if int(server.ServerCpu) == resources.Cpu.MaxCores {
			cpuMaxScaledInstances++
		}
		if int(server.ServerRam) == resources.Memory.MaxBytes {
			memoryMaxScaledInstances++
		}
	}
	if resources.Cpu != nil {
		totalCpuCapacity := len(servers) * resources.Cpu.MaxCores
		capacityTotalGauge.WithLabelValues("cpu").Set(float64(totalCpuCapacity))
		capacityUsedGauge.WithLabelValues("cpu").Set(float64(cpuUsedCapacity))
		maxScaledInstancesGauge.WithLabelValues("cpu").Set(float64(cpuMaxScaledInstances))
	}
	if resources.Memory != nil {
		totalMemoryCapacity := len(servers) * resources.Memory.MaxBytes
		capacityTotalGauge.WithLabelValues("memory").Set(float64(totalMemoryCapacity))
		capacityUsedGauge.WithLabelValues("memory").Set(float64(memoryUsedCapacity))
		maxScaledInstancesGauge.WithLabelValues("memory").Set(float64(memoryMaxScaledInstances))
	}
}

func ServeMetrics() error {
	http.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(":8080", nil)
}
