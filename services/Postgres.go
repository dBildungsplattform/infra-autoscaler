package services

import (
	"fmt"
	"math"
	s "scaler/shared"
)

type PostgresService struct {
	Config PostgresServiceConfig `yaml:"postgres_config"`
}

type PostgresServiceConfig struct {
	CycleTimeSeconds int         `yaml:"cycle_time_seconds"`
	Resources        s.Resources `yaml:"resources"`
}

func (postgres PostgresService) Init() error {
	return initMetricsExporter("postgres")
}

func (postgres *PostgresService) GetConfig() PostgresServiceConfig {
	return postgres.Config
}

func (postgres PostgresService) GetResources() s.Resources {
	return postgres.Config.Resources
}

func (postgres PostgresService) GetCycleTimeSeconds() int {
	return postgres.Config.CycleTimeSeconds
}

func (postgres PostgresService) ComputeScalingProposal(object s.ScaledObject) (s.ResourceScalingProposal, error) {
	var cluster *s.Cluster
	switch objectType := object.(type) {
	case *s.Cluster:
		cluster = objectType
	default:
		return s.ResourceScalingProposal{}, fmt.Errorf("unsupported scaled object type: %s", object.GetType())
	}
	if !cluster.Ready {
		return s.ResourceScalingProposal{}, fmt.Errorf("cluster %s (%s) is not ready", cluster.ClusterName, cluster.ClusterId)
	}
	return postgres.applyRules(*cluster), nil
}

func (postgres PostgresService) applyRules(cluster s.Cluster) s.ResourceScalingProposal {
	targetResource := &s.ResourceScalingProposal{
		Cpu: s.ScaleOp{
			Direction: s.ScaleNone,
			Reason:    "Default",
			Amount:    0,
		},
		Mem: s.ScaleOp{
			Direction: s.ScaleNone,
			Reason:    "Default",
			Amount:    0,
		},
	}
	// Scaling rules:
	// 1. Scale up if current resource is below configured minimum
	// 2. Scale up if current resource usage exceeds maximum usage
	// Add enough resources to either reach usage below the maximum usage or the maximum amount of resources
	// 3. Scale down if current resource is above configured maximum
	// 4. Scale down if current resource usage is below minimum usage
	// Remove enough resources to either reach usage above the minimum usage or the minimum amount of resources

	currentCores := cluster.ResourceState.Cpu.CurrentCores
	minCores := postgres.Config.Resources.Cpu.MinCores
	maxCores := postgres.Config.Resources.Cpu.MaxCores
	currentCpuUsage := cluster.ResourceState.Cpu.CurrentUsage
	minCpuUsage := postgres.Config.Resources.Cpu.MinUsage
	maxCpuUsage := postgres.Config.Resources.Cpu.MaxUsage

	// Rule 1 CPU
	if currentCores < int32(minCores) {
		targetResource.Cpu.Direction = s.ScaleUp
		targetResource.Cpu.Reason = targetResource.Cpu.Reason + ",Rule 1: resource below minimum"
		targetResource.Cpu.Amount = int32(minCores) - currentCores
	}

	// Rule 2 CPU
	if cpuMaxUsageDelta := currentCpuUsage - maxCpuUsage; cpuMaxUsageDelta > 0 && currentCores < int32(maxCores) {
		targetResource.Cpu.Direction = s.ScaleUp
		targetResource.Cpu.Reason = targetResource.Cpu.Reason + ",Rule 2: usage above maximum"
		cpuInc := cpuMaxUsageDelta * float32(currentCores) / currentCpuUsage
		targetHeuristic := currentCores + int32(math.Ceil(float64(cpuInc)))
		targetResource.Cpu.Amount = int32(math.Min(float64(targetHeuristic), float64((maxCores)))) - currentCores
	}

	// Rule 3 CPU
	if currentCores > int32(maxCores) {
		targetResource.Cpu.Direction = s.ScaleDown
		targetResource.Cpu.Reason = targetResource.Cpu.Reason + ",Rule 3: resource above maximum"
		targetResource.Cpu.Amount = currentCores - int32(maxCores)
	}

	// Rule 4 CPU
	if cpuMinUsageDelta := minCpuUsage - currentCpuUsage; cpuMinUsageDelta > 0 && currentCores > int32(minCores) {
		targetResource.Cpu.Direction = s.ScaleDown
		targetResource.Cpu.Reason = targetResource.Cpu.Reason + ",Rule 4: usage below minimum"
		cpuDec := cpuMinUsageDelta * float32(currentCores) / currentCpuUsage
		targetHeuristic := currentCores - int32(math.Ceil(float64(cpuDec)))
		targetResource.Cpu.Amount = currentCores - int32(math.Max(float64(targetHeuristic), float64((minCores))))
	}

	currentMemory := cluster.ResourceState.Memory.CurrentBytes
	minMemory := postgres.Config.Resources.Memory.MinBytes
	maxMemory := postgres.Config.Resources.Memory.MaxBytes
	currentMemoryUsage := cluster.ResourceState.Memory.CurrentUsage
	minMemoryUsage := postgres.Config.Resources.Memory.MinUsage
	maxMemoryUsage := postgres.Config.Resources.Memory.MaxUsage

	// Rule 1 memory
	if currentMemory < int32(minMemory) {
		targetResource.Mem.Direction = s.ScaleUp
		targetResource.Mem.Reason = targetResource.Mem.Reason + ",Rule 1: resource below minimum"
		targetResource.Mem.Amount = int32(minMemory) - currentMemory
	}

	// Rule 2 memory
	if memMaxUsageDelta := currentMemoryUsage - maxMemoryUsage; memMaxUsageDelta > 0 && currentMemory < int32(maxMemory) {
		targetResource.Mem.Direction = s.ScaleUp
		targetResource.Mem.Reason = targetResource.Mem.Reason + ",Rule 2: usage above maximum"
		memInc := memMaxUsageDelta * float32(currentMemory) / currentMemoryUsage
		targetHeuristic := currentMemory + int32(math.Ceil(float64(memInc)))
		targetResource.Mem.Amount = int32(math.Min(float64(targetHeuristic), float64((maxMemory)))) - currentMemory
	}

	// Rule 3 memory
	if currentMemory > int32(maxMemory) {
		targetResource.Mem.Direction = s.ScaleDown
		targetResource.Mem.Reason = targetResource.Mem.Reason + ",Rule 3: resource above maximum"
		targetResource.Mem.Amount = currentMemory - int32(maxMemory)
	}

	// Rule 4 memory
	if memMinUsageDelta := minMemoryUsage - currentMemoryUsage; memMinUsageDelta > 0 && currentMemory > int32(minMemory) {
		targetResource.Mem.Direction = s.ScaleDown
		targetResource.Mem.Reason = targetResource.Mem.Reason + ",Rule 4: usage below minimum"
		memDec := memMinUsageDelta * float32(currentMemory) / currentMemoryUsage
		targetHeuristic := currentMemory - int32(math.Ceil(float64(memDec)))
		targetResource.Mem.Amount = currentMemory - int32(math.Max(float64(targetHeuristic), float64((minMemory))))
	}

	return *targetResource
}

func (service PostgresService) Validate() error {
	if err := service.Config.Validate(); err != nil {
		return err
	}
	return nil
}

func (config PostgresServiceConfig) Validate() error {
	if config.CycleTimeSeconds <= 0 {
		return fmt.Errorf("cycle time seconds must be greater than 0")
	}
	if err := config.Resources.Validate(); err != nil {
		return err
	}
	return nil
}
