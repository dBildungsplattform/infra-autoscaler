package core

import (
	"fmt"
	"scaler/providers"
	"scaler/services/BBB"
	"scaler/services/Postgres"
	s "scaler/shared"
)

type ScalerApp struct {
	appDefinition *s.AppDefinition
	service       *s.Service
	provider      *s.Provider
}

func InitApp(configPath string) *ScalerApp {
	configFile, err := s.OpenConfig(configPath)
	if err != nil {
		panic(err)
	}
	app, err := s.LoadConfig[s.AppDefinition](configFile)
	if err != nil {
		panic(err)
	}

	return &ScalerApp{
		appDefinition: app,
		service:       initService(&app.ServiceType, configFile),
		provider:      initProvider(&app.ProviderType, configFile),
	}
}

func initService(t *s.ServiceType, configFile []byte) *s.Service {
	switch *t {
	case s.BBB:
		bbb, err := s.LoadConfig[BBB.BBBService](configFile)
		if err != nil {
			panic(err)
		}
		service := s.Service(bbb)
		return &service
	case s.Postgres:
		postgres, err := s.LoadConfig[Postgres.PostgresService](configFile)
		if err != nil {
			panic(err)
		}
		service := s.Service(postgres)
		return &service
	}
	return nil // TODO: return error
}

func initProvider(t *s.ProviderType, configFile []byte) *s.Provider {
	switch *t {
	case s.Ionos:
		ionos, err := s.LoadConfig[providers.Ionos](configFile)
		if err != nil {
			panic(err)
		}
		provider := s.Provider(ionos)
		return &provider
	}
	return nil // TODO: return error
}

func (sc *ScalerApp) Scale() {
	fmt.Printf("App: %+v \n", sc.appDefinition)
	fmt.Printf("Service: %+v \n", *sc.service)
	fmt.Printf("Provider: %+v \n", *sc.provider)
	panic("not implemented")
}
