package shared

import (
	"testing"
)

func TestValidateServerDynamicSourceOK(t *testing.T) {
	serverSource := &ServerDynamicSource{
		DatacenterIds:   []string{"123"},
		ServerNameRegex: ".*",
	}
	ValidatePass(t, serverSource)
}

func TestValidateServerDynamicSourceEmptyDatacenterIds(t *testing.T) {
	serverSource := &ServerDynamicSource{
		DatacenterIds:   []string{},
		ServerNameRegex: ".*",
	}
	ValidateFail(t, serverSource)
}

func TestValidateServerDynamicSourceBadRegex(t *testing.T) {
	serverSource := &ServerDynamicSource{
		DatacenterIds:   []string{"123"},
		ServerNameRegex: "*",
	}
	ValidateFail(t, serverSource)
}

func TestValidateServerStaticSourceOK(t *testing.T) {
	serverSource := ServerStaticSource{
		{
			DatacenterId: "123",
			ServerId:     "456",
		},
	}
	ValidatePass(t, serverSource)
}

func TestValidateServerStaticSourceEmpty(t *testing.T) {
	serverSource := ServerStaticSource{}
	ValidateFail(t, serverSource)
}

func TestValidateServerStaticSourceEmptyDatacenterId(t *testing.T) {
	serverSource := ServerStaticSource{
		{
			DatacenterId: "",
			ServerId:     "456",
		},
	}
	ValidateFail(t, serverSource)
}

func TestValidateServerStaticSourceEmptyServerId(t *testing.T) {
	serverSource := ServerStaticSource{
		{
			DatacenterId: "123",
			ServerId:     "",
		},
	}
	ValidateFail(t, serverSource)
}

func TestValidateServerSourceOK(t *testing.T) {
	serverSource := &ServerSource{
		Dynamic: &ServerDynamicSource{
			DatacenterIds:   []string{"123"},
			ServerNameRegex: ".*",
		},
	}
	ValidatePass(t, serverSource)
}

func TestValidateServerSourceNotOk(t *testing.T) {
	serverSource := &ServerSource{}
	ValidateFail(t, serverSource)
}
