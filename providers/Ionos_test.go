package providers

import (
	s "scaler/shared"
	"testing"

	ic "github.com/ionos-cloud/sdk-go/v6"
)

func TestValidateIonosProviderOK(t *testing.T) {
	cloudProvider := &Ionos{
		Config: ProviderConfig{
			Username:   "username",
			Password:   "password",
			ContractId: 1234,
			ServerSource: &s.ServerSource{
				Dynamic: &s.ServerDynamicSource{
					DatacenterIds:   []string{"datacenter-id-1", "datacenter-id-2"},
					ServerNameRegex: "server-name-regex",
				},
			},
		},
	}
	s.ValidatePass(t, cloudProvider)
}

func TestValidateIonosProviderNotOK(t *testing.T) {
	cloudProvider := &Ionos{}
	s.ValidateFail(t, cloudProvider)
}

func TestValidateServer(t *testing.T) {
	serverName := "server-name-1"
	var serverCore, serverRam int32 = 1, 1024
	server := ic.Server{
		Properties: &ic.ServerProperties{
			Name:  &serverName,
			Cores: &serverCore,
			Ram:   &serverRam,
		},
	}

	var coresLimit, ramLimit int32 = 2, 2048
	contract := &ic.Contract{
		Properties: &ic.ContractProperties{
			ResourceLimits: &ic.ResourceLimits{
				CoresPerServer: &coresLimit,
				RamPerServer:   &ramLimit,
			},
		},
	}

	err := validateServer(server, *contract)
	if err != nil {
		t.Errorf("validateServer() failed: %v", err)
	}

	*server.Properties.Cores = 3
	err = validateServer(server, *contract)
	if err == nil {
		t.Errorf("validateServer() should fail")
	}
}
