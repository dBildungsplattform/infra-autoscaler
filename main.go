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
	app := core.Load_app(&ServiceDefinitions)
	app.Scale()
}
