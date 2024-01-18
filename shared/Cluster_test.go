package shared

import (
	"testing"
)

func TestValidateClusterDynamicSourceOK(t *testing.T) {
	clusterSource := &ClusterDynamicSource{
		ClusterNameRegex: ".*",
	}
	ValidatePass(t, clusterSource)
}

func TestValidateClusterDynamicSourceBadRegex(t *testing.T) {
	clusterSource := &ClusterDynamicSource{
		ClusterNameRegex: "*",
	}
	ValidateFail(t, clusterSource)
}

func TestValidateClusterStaticSourceOK(t *testing.T) {
	clusterSource := &ClusterStaticSource{
		ClusterIds: []string{"123"},
	}
	ValidatePass(t, clusterSource)
}

func TestValidateClusterStaticSourceEmpty(t *testing.T) {
	clusterSource := &ClusterStaticSource{}
	ValidateFail(t, clusterSource)
}

func TestValidateClusterSourceOK(t *testing.T) {
	clusterSource := &ClusterSource{
		Dynamic: &ClusterDynamicSource{
			ClusterNameRegex: ".*",
		},
	}
	ValidatePass(t, clusterSource)
}
