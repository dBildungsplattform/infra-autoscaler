package services

import (
	"fmt"
	"math"
	s "scaler/shared"
)

type PostgresService struct {
	State  PostgresServiceState  `yaml:"-"`
	Config PostgresServiceConfig `yaml:"postgres_config"`
}

type PostgresServiceState struct {
	Name string
}

func (postgres PostgresServiceState) GetName() string {
	return postgres.Name
}

type PostgresServiceConfig struct {
	CycleTimeSeconds int         `yaml:"cycle_time_seconds"`
	Resources        s.Resources `yaml:"resources"`
}

func (postgres PostgresService) Init() error {
	return initMetricsExporter("postgres")
}

func (postgres *PostgresService) GetState() s.ServiceState {
	return postgres.State
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

func (postgres PostgresService) ShouldScale(obj s.ScaledObject) (s.ScaleResource, error) {
	if obj.GetType() != s.ClusterType {
		return s.ScaleResource{}, fmt.Errorf("scaled object %s is not a cluster resource", obj.GetName())
	}
	cluster := obj.(s.Cluster)

	targetResource := s.ScaleResource{
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

	postgres.applyRules(&targetResource, cluster)

	return targetResource, nil
}

func (postgres PostgresService) applyRules(targetResource *s.ScaleResource, cluster s.Cluster) {

	// Scaling rules:
	// 1. Scale up if current resource is below configured minimum
	// 2. Scale up if current resource usage exceeds maximum usage
	// Add enough resources to either reach usage below the maximum usage or the maximum amount of resources
	// 3. Scale down if current resource is above configured maximum
	// 4. Scale down if current resource usage is below minimum usage
	// Remove enough resources to either reach usage above the minimum usage or the minimum amount of resources

	// Rule 1 CPU
	if cluster.ClusterCpu < int32(postgres.Config.Resources.Cpu.MinCores) {
		targetResource.Cpu.Direction = s.ScaleUp
		targetResource.Cpu.Reason = targetResource.Cpu.Reason + ",Rule 1: resource below minimum"
		targetResource.Cpu.Amount = int32(postgres.Config.Resources.Cpu.MinCores) - cluster.ClusterCpu
	}

	// Rule 2 CPU
	if cpuMaxUsageDelta := cluster.ClusterCpuUsage - postgres.Config.Resources.Cpu.MaxUsage; cpuMaxUsageDelta > 0 && cluster.ClusterCpu < int32(postgres.Config.Resources.Cpu.MaxCores) {
		targetResource.Cpu.Direction = s.ScaleUp
		targetResource.Cpu.Reason = targetResource.Cpu.Reason + ",Rule 2: usage above maximum"
		cpuInc := cpuMaxUsageDelta * float32(cluster.ClusterCpu) / cluster.ClusterCpuUsage
		targetHeuristic := cluster.ClusterCpu + int32(math.Ceil(float64(cpuInc)))
		targetResource.Cpu.Amount = int32(math.Min(float64(targetHeuristic), float64((postgres.Config.Resources.Cpu.MaxCores)))) - cluster.ClusterCpu
	}

	// Rule 1 memory
	if cluster.ClusterRam < int32(postgres.Config.Resources.Memory.MinBytes) {
		targetResource.Mem.Direction = s.ScaleUp
		targetResource.Mem.Reason = targetResource.Mem.Reason + ",Rule 1: resource below minimum"
		targetResource.Mem.Amount = int32(postgres.Config.Resources.Memory.MinBytes) - cluster.ClusterRam
	}

	// Rule 2 memory
	if memMaxUsageDelta := cluster.ClusterRamUsage - postgres.Config.Resources.Memory.MaxUsage; memMaxUsageDelta > 0 && cluster.ClusterRam < int32(postgres.Config.Resources.Memory.MaxBytes) {
		targetResource.Mem.Direction = s.ScaleUp
		targetResource.Mem.Reason = targetResource.Mem.Reason + ",Rule 2: usage above maximum"
		memInc := memMaxUsageDelta * float32(cluster.ClusterRam) / cluster.ClusterRamUsage
		targetHeuristic := cluster.ClusterRam + int32(math.Ceil(float64(memInc)))
		targetResource.Mem.Amount = int32(math.Min(float64(targetHeuristic), float64((postgres.Config.Resources.Memory.MaxBytes)))) - cluster.ClusterRam
	}

	// Rule 3 CPU
	if cluster.ClusterCpu > int32(postgres.Config.Resources.Cpu.MaxCores) {
		targetResource.Cpu.Direction = s.ScaleDown
		targetResource.Cpu.Reason = targetResource.Cpu.Reason + ",Rule 3: resource above maximum"
		targetResource.Cpu.Amount = cluster.ClusterCpu - int32(postgres.Config.Resources.Cpu.MaxCores)
	}

	// Rule 4 CPU
	if cpuMinUsageDelta := postgres.Config.Resources.Cpu.MinUsage - cluster.ClusterCpuUsage; cpuMinUsageDelta > 0 && cluster.ClusterCpu > int32(postgres.Config.Resources.Cpu.MinCores) {
		targetResource.Cpu.Direction = s.ScaleDown
		targetResource.Cpu.Reason = targetResource.Cpu.Reason + ",Rule 4: usage below minimum"
		cpuDec := cpuMinUsageDelta * float32(cluster.ClusterCpu) / cluster.ClusterCpuUsage
		targetHeuristic := cluster.ClusterCpu - int32(math.Ceil(float64(cpuDec)))
		targetResource.Cpu.Amount = cluster.ClusterCpu - int32(math.Max(float64(targetHeuristic), float64((postgres.Config.Resources.Cpu.MinCores))))
	}

	// Rule 3 memory
	if cluster.ClusterRam > int32(postgres.Config.Resources.Memory.MaxBytes) {
		targetResource.Mem.Direction = s.ScaleDown
		targetResource.Mem.Reason = targetResource.Mem.Reason + ",Rule 3: resource above maximum"
		targetResource.Mem.Amount = cluster.ClusterRam - int32(postgres.Config.Resources.Memory.MaxBytes)
	}

	// Rule 4 memory
	if memMinUsageDelta := postgres.Config.Resources.Memory.MinUsage - cluster.ClusterRamUsage; memMinUsageDelta > 0 && cluster.ClusterRam > int32(postgres.Config.Resources.Memory.MinBytes) {
		targetResource.Mem.Direction = s.ScaleDown
		targetResource.Mem.Reason = targetResource.Mem.Reason + ",Rule 4: usage below minimum"
		memDec := memMinUsageDelta * float32(cluster.ClusterRam) / cluster.ClusterRamUsage
		targetHeuristic := cluster.ClusterRam - int32(math.Ceil(float64(memDec)))
		targetResource.Mem.Amount = cluster.ClusterRam - int32(math.Max(float64(targetHeuristic), float64((postgres.Config.Resources.Memory.MinBytes))))
	}
}

func (service PostgresService) Validate() error {
	if err := service.Config.Validate(); err != nil {
		return err
	}
	if err := service.State.Validate(); err != nil {
		return err
	}
	return nil
}

func (state PostgresServiceState) Validate() error {
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
