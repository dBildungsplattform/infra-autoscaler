package core

import (
	"fmt"
	"net/http"
	s "scaler/shared"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	cyclesCounter           prometheus.Counter
	cycleTimeGauge          prometheus.Gauge
	capacityTotalGauge      *prometheus.GaugeVec
	capacityUsedGauge       *prometheus.GaugeVec
	instancesGauge          *prometheus.GaugeVec
	maxScaledInstancesGauge *prometheus.GaugeVec
	lastScaleTimeGauge      prometheus.Gauge
)

func initMetricsExporter() error {
	cyclesCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "autoscaler_cycle_count",
		Help: "The total number of cycles the autoscaler has run",
	})
	cycleTimeGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "autoscaler_cycle_time_seconds",
		Help: "Autoscaler cycle time in seconds",
	})
	capacityTotalGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "autoscaler_capacity_total",
		Help: "The maximum amount of resource that can be used",
	}, []string{"resource_type"})
	capacityUsedGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "autoscaler_capacity_used",
		Help: "The amount of resource that is currently used",
	}, []string{"resource_type"})
	instancesGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "autoscaler_instances_count",
		Help: "The amount of instances currently loaded by the autoscaler",
	}, []string{"ready"})
	maxScaledInstancesGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "autoscaler_max_scaled_instances",
		Help: "The amount of instances that are scaled to the maximum for a given resource type",
	}, []string{"resource_type"})
	lastScaleTimeGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "autoscaler_last_scale_time",
		Help: "The time of the last scale operation",
	})
	metrics := []prometheus.Collector{cyclesCounter, cycleTimeGauge, capacityTotalGauge, capacityUsedGauge, instancesGauge, maxScaledInstancesGauge, lastScaleTimeGauge}
	for _, metric := range metrics {
		if err := prometheus.Register(metric); err != nil {
			return err
		}
	}
	cyclesCounter.Add(0)
	return nil
}

func (sc ScalerApp) calculateMetrics(scaledObjects []s.ScaledObject) {
	resources := sc.service.GetResources()

	cpuUsedCapacity := 0
	memoryUsedCapacity := 0
	readyInstances := 0
	notReadyInstances := 0
	cpuMaxScaledInstances := 0
	memoryMaxScaledInstances := 0

	capacityTotalGauge.Reset()
	capacityUsedGauge.Reset()

	for _, object := range scaledObjects {
		currentCpuCores := int(object.GetResourceState().Cpu.CurrentCores)
		currentMemoryBytes := int(object.GetResourceState().Memory.CurrentBytes)
		cpuUsedCapacity += currentCpuCores
		memoryUsedCapacity += currentMemoryBytes

		if object.IsReady() {
			readyInstances++
		} else {
			notReadyInstances++
		}

		if currentCpuCores == resources.Cpu.MaxCores {
			cpuMaxScaledInstances++
		}
		if currentMemoryBytes == resources.Memory.MaxBytes {
			memoryMaxScaledInstances++
		}
	}
	instancesGauge.WithLabelValues("true").Set(float64(readyInstances))
	instancesGauge.WithLabelValues("false").Set(float64(notReadyInstances))

	if resources.Cpu != nil {
		totalCpuCapacity := len(scaledObjects) * resources.Cpu.MaxCores
		capacityTotalGauge.WithLabelValues("cpu").Set(float64(totalCpuCapacity))
		capacityUsedGauge.WithLabelValues("cpu").Set(float64(cpuUsedCapacity))
		maxScaledInstancesGauge.WithLabelValues("cpu").Set(float64(cpuMaxScaledInstances))
	}
	if resources.Memory != nil {
		totalMemoryCapacity := len(scaledObjects) * resources.Memory.MaxBytes
		capacityTotalGauge.WithLabelValues("memory").Set(float64(totalMemoryCapacity))
		capacityUsedGauge.WithLabelValues("memory").Set(float64(memoryUsedCapacity))
		maxScaledInstancesGauge.WithLabelValues("memory").Set(float64(memoryMaxScaledInstances))
	}
}

func (sc ScalerApp) ServeMetrics() error {
	http.Handle("/metrics", promhttp.Handler())
	port := sc.appDefinition.MetricsExporterPort
	if port == 0 {
		port = 8080
	}
	portString := fmt.Sprint(":", port)
	return http.ListenAndServe(portString, nil)
}
