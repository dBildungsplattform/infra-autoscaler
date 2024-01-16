package shared

import (
	"fmt"
	"regexp"
	"time"
)

type Cluster struct {
	ClusterId     string `yaml:"cluster_id"`
	ClusterName   string `yaml:"cluster_name"`
	ResourceState ResourceState
	LastUpdated   time.Time `yaml:"last_updated"`
	Ready         bool      `yaml:"ready"`
}

func (c Cluster) GetType() ScaledObjectType {
	return ClusterType
}

func (c Cluster) GetName() string {
	return c.ClusterName
}

func (c Cluster) GetResourceState() ResourceState {
	return c.ResourceState
}

func (c *Cluster) SetResourceState(resourceState ResourceState) {
	c.ResourceState = resourceState
}

type ClusterSource struct {
	Dynamic *ClusterDynamicSource `yaml:"dynamic"`
	Static  *ClusterStaticSource  `yaml:"static"`
}

type ClusterDynamicSource struct {
	ClusterNameRegex string `yaml:"cluster_name_regex"`
}

type ClusterStaticSource struct {
	ClusterIds []string `yaml:"cluster_ids"`
}

func (ionos ClusterSource) Validate() error {
	static, dynamic := ionos.Static, ionos.Dynamic
	if static == nil && dynamic == nil {
		return fmt.Errorf("ionos.static and ionos.dynamic are nil, one must be set")
	}
	if static != nil && dynamic != nil {
		return fmt.Errorf("ionos.static and ionos.dynamic are both set, only one must be set")
	}
	if static != nil {
		err := static.Validate()
		if err != nil {
			return err
		}
	}
	if dynamic != nil {
		err := dynamic.Validate()
		if err != nil {
			return err
		}
	}
	return nil
}

func (ionos ClusterDynamicSource) Validate() error {
	if ionos.ClusterNameRegex == "" {
		return fmt.Errorf("ionos.cluster_name_regex is empty")
	}
	if _, err := regexp.Compile(ionos.ClusterNameRegex); err != nil {
		return fmt.Errorf("ionos.cluster_name_regex is invalid: %s", err)
	}
	return nil
}

func (ionos ClusterStaticSource) Validate() error {
	if len(ionos.ClusterIds) == 0 {
		return fmt.Errorf("ionos.cluster_ids is empty")
	}
	return nil
}
