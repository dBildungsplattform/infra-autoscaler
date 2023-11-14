package main

import (
	c "scaler/core"
	s "scaler/shared"
)

func main() {
	ServiceDefinitions := []s.ServiceDefinition{
		{Name: "bbb1", Type: s.BBB},
		{Name: "bbb2", Type: s.BBB},
		{Name: "postgres1", Type: s.Postgres},
	}
	ProviderDefinition := []s.ProviderDefinition{
		{Name: "Ionos", Type: s.Ionos},
	}
	app := c.Init_app(&ServiceDefinitions, &ProviderDefinition)
	app.Scale()
}
