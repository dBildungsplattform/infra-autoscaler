package ionos

import (
	v "scaler/validater"
	"testing"
)

func TestValidateCloudProviderOK(t *testing.T) {
	cloudProvider := &CloudProvider{
		Username: "username",
		Password: "password",
		ServerSource: &ServerSource{
			Dynamic: &ServerDynamicSource{
				DatacenterIds:   []string{"123"},
				ServerNameRegex: ".*",
			},
		},
	}
	v.ValidatePass(t, cloudProvider)
}

func TestValidateCloudProviderNotOK(t *testing.T) {
	cloudProvider := &CloudProvider{}
	v.ValidateFail(t, cloudProvider)
}
