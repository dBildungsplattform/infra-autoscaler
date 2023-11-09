package config

import (
	"scaler/ionos"
	v "scaler/validater"
	"testing"
)

func TestValidateCloudProviderOK(t *testing.T) {
	cloudProvider := &CloudProvider{
		Ionos: &ionos.CloudProvider{
			Username: "username",
			Password: "password",
			ServerSource: &ionos.ServerSource{
				Dynamic: &ionos.ServerDynamicSource{
					DatacenterIds:   []string{"123"},
					ServerNameRegex: ".*",
				},
			},
		},
	}
	v.ValidatePass(t, cloudProvider)
}

func TestValidateCloudProviderNotOK(t *testing.T) {
	cloudProvider := &CloudProvider{}
	v.ValidateFail(t, cloudProvider)
}
