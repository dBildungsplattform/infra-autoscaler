package services

import (
	"os"
	s "scaler/shared"
	"testing"
	"time"
)

var validPostgresConfig = &PostgresServiceConfig{
	CycleTimeSeconds: 60,
	Resources: s.Resources{
		Cpu: &s.CpuResources{
			MinCores: 2,
			MaxCores: 4,
			MaxUsage: 0.5,
		},
		Memory: &s.MemoryResources{
			MinBytes: 2048,
			MaxBytes: 4096,
			MinUsage: 0.1,
			MaxUsage: 0.5,
		},
	},
}
var samplePostgresCluster = s.Cluster{
	ClusterId:   "5678",
	ClusterName: "postgres",
	ResourceState: s.ResourceState{
		Cpu: &s.CpuResourceState{
			CurrentCores: 1,
			CurrentUsage: 0.5,
		},
		Memory: &s.MemoryResourceState{
			CurrentBytes: 1024,
			CurrentUsage: 0.5,
		},
	},
	LastUpdated: time.Now(),
	Ready:       true,
}

func TestValidateConfigOk(t *testing.T) {
	postgresConfig := validPostgresConfig
	s.ValidatePass(t, postgresConfig)
}

func TestValidateConfigNotOk(t *testing.T) {
	postgresConfig := &PostgresServiceConfig{}
	s.ValidateFail(t, postgresConfig)
}

func TestParseConfigOk(t *testing.T) {
	os.Setenv("IONOS_TOKEN", "1234567890")
	defer os.Unsetenv("IONOS_TOKEN")

	config, ok := s.OpenConfig("test_files/postgres_config_ok.yml")
	if ok != nil {
		t.Fatalf("Failed to open config: %v", ok)
	}

	c, ok := s.LoadConfig[PostgresService](config)
	if ok != nil {
		t.Fatalf("Failed to parse config: %v", ok)
	}

	s.ValidatePass(t, c)
}

func TestParseConfigNotOk(t *testing.T) {
	config, ok := s.OpenConfig("test_files/postgres_config_not_ok.yml")
	if ok != nil {
		t.Fatalf("Failed to open config: %v", ok)
	}

	_, ok = s.LoadConfig[PostgresService](config)
	if ok == nil {
		t.Fatalf("Expected error but got nil")
	}
}

func testPostgresApplyRulesCPU(t *testing.T, resourceState s.CpuResourceState, resources s.CpuResources, expected s.ScaleDirection) {
	postgresConfig := validPostgresConfig
	postgresConfig.Resources.Cpu = &resources
	postgresService := PostgresService{
		Config: *postgresConfig,
	}

	cluster := samplePostgresCluster
	cluster.ResourceState.Cpu = &resourceState

	proposal := postgresService.applyRules(cluster)
	if proposal.Cpu.Direction != expected {
		t.Fatalf("Expected CPU scale direction to be %s but got %s", expected, proposal.Cpu.Direction)
	}
}

// Check that a cluster with below minimum resources is scaled up
func TestPostgresApplyRulesRule1(t *testing.T) {
	testPostgresApplyRulesCPU(t, s.CpuResourceState{
		CurrentCores: 1,
		CurrentUsage: 0.5,
	}, s.CpuResources{
		MinCores: 2,
		MaxCores: 4,
		MinUsage: 0.1,
		MaxUsage: 0.5,
	}, s.ScaleUp)
}

// Check that a cluster with above maximum usage is scaled up
func TestPostgresApplyRulesRule2(t *testing.T) {
	testPostgresApplyRulesCPU(t, s.CpuResourceState{
		CurrentCores: 2,
		CurrentUsage: 0.6,
	}, s.CpuResources{
		MinCores: 2,
		MaxCores: 4,
		MinUsage: 0.1,
		MaxUsage: 0.5,
	}, s.ScaleUp)
}

// Check that a cluster with above maximum resources is scaled down
func TestPostgresApplyRulesRule3(t *testing.T) {
	testPostgresApplyRulesCPU(t, s.CpuResourceState{
		CurrentCores: 5,
		CurrentUsage: 0,
	}, s.CpuResources{
		MinCores: 2,
		MaxCores: 4,
		MinUsage: 0.1,
		MaxUsage: 0.5,
	}, s.ScaleDown)
}

// Check that a cluster with below minimum usage is scaled down
func TestPostgresApplyRulesRule4(t *testing.T) {
	testPostgresApplyRulesCPU(t, s.CpuResourceState{
		CurrentCores: 3,
		CurrentUsage: 0.01,
	}, s.CpuResources{
		MinCores: 2,
		MaxCores: 4,
		MinUsage: 0.1,
		MaxUsage: 0.5,
	}, s.ScaleDown)
}
