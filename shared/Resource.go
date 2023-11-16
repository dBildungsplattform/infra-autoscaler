package shared

import "fmt"

/*** Resource definition ***/
type Resources struct {
	Cpu    *CpuResources    `yaml:"cpu"`
	Memory *MemoryResources `yaml:"memory"`
}

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
	if r.Cpu == nil && r.Memory == nil {
		return fmt.Errorf("resources.cpu and resources.memory are nil, at least one must be set")
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
