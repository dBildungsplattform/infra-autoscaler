package core

import (
	c "scaler/common"
	"scaler/services/BBB"
	"scaler/services/Postgres"
)

type ScalerConfig struct {
	ServiceDefinitions *[]c.ServiceDefinition
	ProviderNames      []c.InfrastructureProviderName
	//Providers     []InfrastructureProvider
	Services []*c.Service
}

func Load_app(sd *[]c.ServiceDefinition) *ScalerConfig {
	// TODO load config from file, env, CLI args, etc.
	var sc *ScalerConfig = &ScalerConfig{
		ServiceDefinitions: sd,
		ProviderNames:      []c.InfrastructureProviderName{c.Ionos},
		//Providers:     []InfrastructureProvider{},
		Services: []*c.Service{},
	}
	for _, serviceDef := range *sc.ServiceDefinitions {
		switch t := serviceDef.Type; t {
		case c.BBB:
			var bbb c.Service = BBB.BBBService{}
			bbb.Init(serviceDef.Name)
			sc.Services = append(sc.Services, &bbb)
		case c.Postgres:
			var postgres c.Service = Postgres.PostgresService{}
			postgres.Init(serviceDef.Name)
			sc.Services = append(sc.Services, &postgres)
		}
	}
	return sc
}

func (sc *ScalerConfig) Scale() {
	panic("not implemented")
}
