package main

import (
	"scaler/common"
	"scaler/core"
)

func main() {
	ServiceDefinitions := []common.ServiceDefinition{
		{Name: "bbb1", Type: common.BBB},
		{Name: "bbb2", Type: common.BBB},
		{Name: "postgres1", Type: common.Postgres},
	}
	ProviderDefinition := []common.ProviderDefinition{
		{Name: "Ionos", Type: common.Ionos},
	}
	app := core.Init_app(&ServiceDefinitions, &ProviderDefinition)
	app.Scale()
}
