package providers

import (
	s "scaler/shared"
	"testing"
)

func TestValidateIonosProviderOK(t *testing.T) {
	cloudProvider := &Ionos{
		IonosConfig: ProviderConfig{
			Username: "username",
			Password: "password",
			ServerSource: &s.ServerSource{
				Dynamic: &s.ServerDynamicSource{
					DatacenterIds:   []string{"datacenter-id-1", "datacenter-id-2"},
					ServerNameRegex: "server-name-regex",
				},
			},
			PrometheusConfig: PrometheusConfig{
				Url: "url",
			},
		},
	}
	s.ValidatePass(t, cloudProvider)
}

func TestValidateIonosProviderNotOK(t *testing.T) {
	cloudProvider := &Ionos{}
	s.ValidateFail(t, cloudProvider)
}
