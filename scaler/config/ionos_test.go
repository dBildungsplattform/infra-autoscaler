package config

import (
	"testing"
)

func TestValidateServerSourceOK(t *testing.T) {
	serverSource := &IonosServerInstancesSource{
		DatacenterIds:   []string{"123"},
		ServerNameRegex: ".*",
	}
	ValidatePass(t, serverSource)
}

func TestValidateServerSourceEmptyDatacenterIds(t *testing.T) {
	serverSource := &IonosServerInstancesSource{
		DatacenterIds:   []string{},
		ServerNameRegex: ".*",
	}
	ValidateFail(t, serverSource)
}

func TestValidateServerSourceBadRegex(t *testing.T) {
	serverSource := &IonosServerInstancesSource{
		DatacenterIds:   []string{"123"},
		ServerNameRegex: "*",
	}
	ValidateFail(t, serverSource)
}

func TestValidateInlineServerSourceOK(t *testing.T) {
	serverSource := InlineIonosServerInstancesSource{
		{
			DatacenterId: "123",
			ServerId:     "456",
		},
	}
	ValidatePass(t, serverSource)
}

func TestValidateInlineServerSourceEmpty(t *testing.T) {
	serverSource := InlineIonosServerInstancesSource{}
	ValidateFail(t, serverSource)
}

func TestValidateInlineServerSourceEmptyDatacenterId(t *testing.T) {
	serverSource := InlineIonosServerInstancesSource{
		{
			DatacenterId: "",
			ServerId:     "456",
		},
	}
	ValidateFail(t, serverSource)
}

func TestValidateInlineServerSourceEmptyServerId(t *testing.T) {
	serverSource := InlineIonosServerInstancesSource{
		{
			DatacenterId: "123",
			ServerId:     "",
		},
	}
	ValidateFail(t, serverSource)
}
