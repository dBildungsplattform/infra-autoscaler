package core

import (
	"scaler/providers"
	"scaler/services/BBB"
	"scaler/services/Postgres"
	s "scaler/shared"
)

type ScalerApp struct {
	serviceDefinitions  *[]s.ServiceDefinition
	providerDefinitions *[]s.ProviderDefinition
	services            []*s.Service
	providers           []*s.Provider
}

func Init_app(sd *[]s.ServiceDefinition, pd *[]s.ProviderDefinition) *ScalerApp {
	return &ScalerApp{
		serviceDefinitions:  sd,
		providerDefinitions: pd,
		services:            init_services(sd),
		providers:           init_providers(pd),
	}
}

func init_services(sd *[]s.ServiceDefinition) []*s.Service {
	service := []*s.Service{}
	for _, serviceDef := range *sd {
		switch t := serviceDef.Type; t {
		case s.BBB:
			var bbb s.Service = BBB.BBBService{}
			bbb.Init(&serviceDef)
			service = append(service, &bbb)
		case s.Postgres:
			var postgres s.Service = Postgres.PostgresService{}
			postgres.Init(&serviceDef)
			service = append(service, &postgres)
		}
	}
	return service
}

func init_providers(pd *[]s.ProviderDefinition) []*s.Provider {
	p := []*s.Provider{}
	for _, providerDef := range *pd {
		switch t := providerDef.Type; t {
		case s.Ionos:
			var ionos s.Provider = s.Provider(load_provider[providers.Ionos](providerDef))
			p = append(p, &ionos)
		}
	}
	return p
}

func (sc *ScalerApp) Scale() {
	panic("not implemented")
}
