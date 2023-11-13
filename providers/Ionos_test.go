package providers

import (
	s "scaler/shared"
	"testing"
)

func TestValidateIonosProviderOK(t *testing.T) {
	cloudProvider := &Ionos{
		Username: "username",
		Password: "password",
	}
	s.ValidatePass(t, cloudProvider)
}

func TestValidateIonosProviderNotOK(t *testing.T) {
	cloudProvider := &Ionos{}
	s.ValidateFail(t, cloudProvider)
}
