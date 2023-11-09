package ionos

import (
	v "scaler/validater"
	"testing"
)

func TestValidateServerDynamicSourceOK(t *testing.T) {
	serverSource := &ServerDynamicSource{
		DatacenterIds:   []string{"123"},
		ServerNameRegex: ".*",
	}
	v.ValidatePass(t, serverSource)
}

func TestValidateServerDynamicSourceEmptyDatacenterIds(t *testing.T) {
	serverSource := &ServerDynamicSource{
		DatacenterIds:   []string{},
		ServerNameRegex: ".*",
	}
	v.ValidateFail(t, serverSource)
}

func TestValidateServerDynamicSourceBadRegex(t *testing.T) {
	serverSource := &ServerDynamicSource{
		DatacenterIds:   []string{"123"},
		ServerNameRegex: "*",
	}
	v.ValidateFail(t, serverSource)
}

func TestValidateServerStaticSourceOK(t *testing.T) {
	serverSource := ServerStaticSource{
		{
			DatacenterId: "123",
			ServerId:     "456",
		},
	}
	v.ValidatePass(t, serverSource)
}

func TestValidateServerStaticSourceEmpty(t *testing.T) {
	serverSource := ServerStaticSource{}
	v.ValidateFail(t, serverSource)
}

func TestValidateServerStaticSourceEmptyDatacenterId(t *testing.T) {
	serverSource := ServerStaticSource{
		{
			DatacenterId: "",
			ServerId:     "456",
		},
	}
	v.ValidateFail(t, serverSource)
}

func TestValidateServerStaticSourceEmptyServerId(t *testing.T) {
	serverSource := ServerStaticSource{
		{
			DatacenterId: "123",
			ServerId:     "",
		},
	}
	v.ValidateFail(t, serverSource)
}

func TestValidateServerSourceOK(t *testing.T) {
	serverSource := &ServerSource{
		Dynamic: &ServerDynamicSource{
			DatacenterIds:   []string{"123"},
			ServerNameRegex: ".*",
		},
	}
	v.ValidatePass(t, serverSource)
}

func TestValidateServerSourceNotOk(t *testing.T) {
	serverSource := &ServerSource{}
	v.ValidateFail(t, serverSource)
}
