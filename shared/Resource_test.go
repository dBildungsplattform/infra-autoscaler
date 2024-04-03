package shared

import (
	"testing"
)

func TestValidateCpuResourcesOK(t *testing.T) {
	cpuResources := &CpuResources{
		MinCores: 1,
		MaxCores: 2,
		MinUsage: 0.1,
		MaxUsage: 0.5,
	}
	ValidatePass(t, cpuResources)
}

func TestValidateCpuResourcesMinCoresZero(t *testing.T) {
	cpuResources := &CpuResources{
		MinCores: 0,
		MaxCores: 2,
		MinUsage: 0.1,
		MaxUsage: 0.5,
	}
	ValidateFail(t, cpuResources)
}

func TestValidateCpuResourcesMaxCoresLessThanMinCores(t *testing.T) {
	cpuResources := &CpuResources{
		MinCores: 2,
		MaxCores: 1,
		MinUsage: 0.1,
		MaxUsage: 0.5,
	}
	ValidateFail(t, cpuResources)
}

func TestValidateCpuResourcesMaxUsageZero(t *testing.T) {
	cpuResources := &CpuResources{
		MinCores: 2,
		MaxCores: 2,
		MinUsage: 0.1,
		MaxUsage: 0,
	}
	ValidateFail(t, cpuResources)
}

func TestValidateCpuResourcesMaxUsageGreaterThanOne(t *testing.T) {
	cpuResources := &CpuResources{
		MinCores: 2,
		MaxCores: 2,
		MinUsage: 0.1,
		MaxUsage: 1.1,
	}
	ValidateFail(t, cpuResources)
}

func TestValidateMemoryResourcesOK(t *testing.T) {
	memoryResources := &MemoryResources{
		MinBytes: 1024,
		MaxBytes: 2048,
		MinUsage: 0.1,
		MaxUsage: 0.5,
	}
	ValidatePass(t, memoryResources)
}

func TestValidateMemoryResourcesMinBytesLessThan1024(t *testing.T) {
	memoryResources := &MemoryResources{
		MinBytes: 1023,
		MaxBytes: 2048,
		MinUsage: 0.1,
		MaxUsage: 0.5,
	}
	ValidateFail(t, memoryResources)
}

func TestValidateMemoryResourcesMaxBytesLessThanMinBytes(t *testing.T) {
	memoryResources := &MemoryResources{
		MinBytes: 2048,
		MaxBytes: 1024,
		MinUsage: 0.1,
		MaxUsage: 0.5,
	}
	ValidateFail(t, memoryResources)
}

func TestValidateMemoryResourcesMaxUsageZero(t *testing.T) {
	memoryResources := &MemoryResources{
		MinBytes: 2048,
		MaxBytes: 2048,
		MinUsage: 0.1,
		MaxUsage: 0,
	}
	ValidateFail(t, memoryResources)
}

func TestValidateMemoryResourcesMaxUsageGreaterThanOne(t *testing.T) {
	memoryResources := &MemoryResources{
		MinBytes: 2048,
		MaxBytes: 2048,
		MinUsage: 0.1,
		MaxUsage: 1.1,
	}
	ValidateFail(t, memoryResources)
}

func TestValidateReplicaResourcesOK(t *testing.T) {
	ReplicaResources := &ReplicaResources{
		MinReplicas: 2,
		MaxReplicas: 4,
		MinUsage: 0.3,
		MaxUsage: 0.7,
	}
	ValidatePass(t, ReplicaResources)
}

func TestValidateReplicaResourcesMinReplicasZero(t *testing.T) {
	ReplicaResources := &ReplicaResources{
		MinReplicas: 0,
		MaxReplicas: 4,
		MinUsage: 0.3,
		MaxUsage: 0.7,
	}
	ValidateFail(t, ReplicaResources)
}

func TestValidateReplicaResourcesMaxLessThanMin(t *testing.T) {
	ReplicaResources := &ReplicaResources{
		MinReplicas: 2,
		MaxReplicas: 1,
		MinUsage: 0.3,
		MaxUsage: 0.7,
	}
	ValidateFail(t, ReplicaResources)
}

func TestValidateReplicaResourcesMaxUsageGreaterThanOne(t *testing.T) {
	ReplicaResources := &ReplicaResources{
		MinReplicas: 2,
		MaxReplicas: 4,
		MinUsage: 0.3,
		MaxUsage: 1.1,
	}
	ValidateFail(t, ReplicaResources)
}

func TestValidateReplicaResourcesMaxUsageZero(t *testing.T) {
	ReplicaResources := &ReplicaResources{
		MinReplicas: 2,
		MaxReplicas: 4,
		MinUsage: 0.3,
		MaxUsage: 0,
	}
	ValidateFail(t, ReplicaResources)
}