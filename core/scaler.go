package core

import (
	c "scaler/common"
	"scaler/providers/Ionos"
	"scaler/services/BBB"
	"scaler/services/Postgres"
)

type ScalerApp struct {
	serviceDefinitions  *[]c.ServiceDefinition
	providerDefinitions *[]c.ProviderDefinition
	services            []*c.Service
	providers           []*c.Provider
}

func Init_app(sd *[]c.ServiceDefinition, pd *[]c.ProviderDefinition) *ScalerApp {
	return &ScalerApp{
		serviceDefinitions:  sd,
		providerDefinitions: pd,
		services:            init_services(sd),
		providers:           init_providers(pd),
	}
}

func init_services(sd *[]c.ServiceDefinition) []*c.Service {
	s := []*c.Service{}
	for _, serviceDef := range *sd {
		switch t := serviceDef.Type; t {
		case c.BBB:
			var bbb c.Service = BBB.BBBService{}
			bbb.Init(&serviceDef)
			s = append(s, &bbb)
		case c.Postgres:
			var postgres c.Service = Postgres.PostgresService{}
			postgres.Init(&serviceDef)
			s = append(s, &postgres)
		}
	}
	return s
}

func init_providers(pd *[]c.ProviderDefinition) []*c.Provider {
	p := []*c.Provider{}
	for _, providerDef := range *pd {
		switch t := providerDef.Type; t {
		case c.Ionos:
			var ionos c.Provider = Ionos.Provider(load_provider[Ionos.Provider](providerDef))
			p = append(p, &ionos)
		}
	}
	return p
}

func (sc *ScalerApp) Scale() {
	panic("not implemented")
}
