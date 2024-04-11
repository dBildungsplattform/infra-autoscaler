package shared

import "fmt"

/*** Resource definition ***/
type Resources struct {
	Cpu    *CpuResources    `yaml:"cpu"`
	Memory *MemoryResources `yaml:"memory"`
	Replica *ReplicaResources `yaml:"replicas"`
}

// TODO: replace this with a generic resource interface
type CpuResources struct {
	MinCores int     `yaml:"min_cores"`
	MaxCores int     `yaml:"max_cores"`
	MinUsage float32 `yaml:"min_usage"`
	MaxUsage float32 `yaml:"max_usage"`
}

type MemoryResources struct {
	MinBytes int     `yaml:"min_bytes"`
	MaxBytes int     `yaml:"max_bytes"`
	MinUsage float32 `yaml:"min_usage"`
	MaxUsage float32 `yaml:"max_usage"`
}

type ReplicaResources struct {
	MinReplicas int 	`yaml:"min_replicas"`
	MaxReplicas	int 	`yaml:"max_replicas"`
	MinUsage	float32	`yaml:"min_usage"`
	MaxUsage	float32	`yaml:"max_usage"`
}

type ResourceState struct {
	Cpu    *CpuResourceState
	Memory *MemoryResourceState
	Replica *ReplicaResourceState
}

type CpuResourceState struct {
	CurrentCores int32
	CurrentUsage float32
}

type MemoryResourceState struct {
	CurrentBytes int32
	CurrentUsage float32
}

type ReplicaResourceState struct {
	CurrentReplicas int
	CurrentUsage 	float32
}

type ResourceScalingProposal struct {
	Cpu ScaleOp
	Mem ScaleOp
	Replica ScaleOp
}

type ScaleOp struct {
	Direction ScaleDirection
	Reason    string
	Amount    int32
}

type ScaleDirection string

const (
	ScaleUp   = "up"
	ScaleDown = "down"
	ScaleNone = "none"
)

func (r Resources) Validate() error {
	if cpu := r.Cpu; cpu != nil {
		if err := cpu.Validate(); err != nil {
			return err
		}
	}
	if memory := r.Memory; memory != nil {
		if err := memory.Validate(); err != nil {
			return err
		}
	}
	if replica := r.Replica; replica != nil {
		if err := replica.Validate(); err != nil {
			return err
		}
	}
	if r.Cpu == nil && r.Memory == nil && r.Replica == nil{
		return fmt.Errorf("resources.cpu and resources.memory and resources.replica are nil, at least one must be set")
	}
	return nil
}

func (c CpuResources) Validate() error {
	if c.MinCores <= 0 {
		return fmt.Errorf("cpu.min_cores must be greater than 0 but got %d", c.MinCores)
	}
	if c.MaxCores < c.MinCores {
		return fmt.Errorf("cpu.max_cores must be greater than or equal to min_cores (%d) but got %d", c.MinCores, c.MaxCores)
	}
	if c.MinUsage < 0 || c.MinUsage > 1 {
		return fmt.Errorf("cpu.min_usage must be greater than 0 and less than or equal to 1 but got %f", c.MinUsage)
	}
	if c.MaxUsage <= c.MinUsage || c.MaxUsage > 1 {
		return fmt.Errorf("cpu.max_usage must be greater than min_usage (%f) and less than or equal to 1 but got %f", c.MinUsage, c.MaxUsage)
	}
	return nil
}

func (m MemoryResources) Validate() error {
	if m.MinBytes < 1024 {
		return fmt.Errorf("memory.min_bytes must be greater than or equal to 1024 but got %d", m.MinBytes)
	}
	if m.MaxBytes < m.MinBytes {
		return fmt.Errorf("memory.max_bytes must be greater than or equal to min_bytes (%d) but got %d", m.MinBytes, m.MaxBytes)
	}
	if m.MinUsage < 0 || m.MinUsage > 1 {
		return fmt.Errorf("memory.min_usage must be greater than 0 and less than or equal to 1 but got %f", m.MinUsage)
	}
	if m.MaxUsage <= m.MinUsage || m.MaxUsage > 1 {
		return fmt.Errorf("memory.max_usage must be greater than min_usage (%f) and less than or equal to 1 but got %f", m.MinUsage, m.MaxUsage)
	}
	return nil
}

func (r ReplicaResources) Validate() error {
	if r.MinReplicas < 1 {
		return fmt.Errorf("replicas.min_replicas must be greater than or equal to 1 but got %d", r.MinReplicas)
	}
	if r.MaxReplicas < r.MinReplicas {
		return fmt.Errorf("r.MaxReplicas must be greater than or equal to min_replicas (%d) but got %d", r.MinReplicas, r.MaxReplicas)
	}
	if r.MinUsage < 0 || r.MinUsage > 1 {
		return fmt.Errorf("replicas.min_usage must be greater than 0 and less than or equal to 1 but got %f", r.MinUsage)
	}
	if r.MaxUsage <= r.MinUsage || r.MaxUsage > 1 {
		return fmt.Errorf("replicas.max_usage must be greater than min_usage (%f) and less than or equal to 1 but got %f", r.MinUsage, r.MaxUsage)
	}
	return nil
}
