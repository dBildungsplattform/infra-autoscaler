package shared

import (
	"fmt"
	"regexp"
	"time"
)

type Cluster struct {
	ClusterId          string    `yaml:"cluster_id"`
	ClusterName        string    `yaml:"cluster_name"`
	ClusterCpu         int32     `yaml:"cluster_cpu"`
	ClusterRam         int32     `yaml:"cluster_ram"`
	ClusterCpuUsage    float32   `yaml:"cluster_cpu_usage"`
	ClusterRamUsage    float32   `yaml:"cluster_ram_usage"`
	ClusterStorageSize int32     `yaml:"cluster_storage_size"`
	ClusterStorageType string    `yaml:"cluster_storage_type"`
	LastUpdated        time.Time `yaml:"last_updated"`
	Ready              bool      `yaml:"ready"`
}

func (c Cluster) GetType() ScaledObjectType {
	return ClusterType
}

func (c Cluster) GetName() string {
	return c.ClusterName
}

type ClusterSource struct {
	ClusterFilterName string `yaml:"cluster_filter_name"`
	ClusterNameRegex  string `yaml:"cluster_name_regex"`
}

func (ionos ClusterSource) Validate() error {
	if ionos.ClusterNameRegex == "" {
		return fmt.Errorf("cluster_name_regex is empty")
	}
	if _, err := regexp.Compile(ionos.ClusterNameRegex); err != nil {
		return fmt.Errorf("cluster_name_regex is not a valid regex: %s", err)
	}
	return nil
}