package core

import (
	"fmt"
	"scaler/metrics"
	"scaler/providers"
	"scaler/services"
	s "scaler/shared"
)

type ScalerApp struct {
	appDefinition *s.AppDefinition
	service       *s.Service
	provider      *s.Provider
	metrics       *s.Metrics
}

func InitApp(configPath string) (*ScalerApp, error) {
	configFile, err := s.OpenConfig(configPath)
	if err != nil {
		panic(err)
	}
	app, err := s.LoadConfig[s.AppDefinition](configFile)
	if err != nil {
		panic(err)
	}

	metrics, err := initMetrics(&app.MetricsType, configFile)
	if err != nil {
		return nil, fmt.Errorf("error while initializing metrics: %s", err)
	}
	return &ScalerApp{
		appDefinition: app,
		service:       initService(&app.ServiceType, configFile),
		provider:      initProvider(&app.ProviderType, configFile),
		metrics:       metrics,
	}, nil
}

func initService(t *s.ServiceType, configFile []byte) *s.Service {
	switch *t {
	case s.BBB:
		bbb, err := s.LoadConfig[services.BBBService](configFile)
		if err != nil {
			panic(err)
		}
		service := s.Service(bbb)
		return &service
	case s.Postgres:
		postgres, err := s.LoadConfig[services.PostgresService](configFile)
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

func initMetrics(t *s.MetricsType, configFile []byte) (*s.Metrics, error) {
	switch *t {
	case s.Prometheus:
		prometheus, err := s.LoadConfig[metrics.Prometheus](configFile)
		if err != nil {
			return nil, fmt.Errorf("error while loading prometheus config: %s", err)
		}
		if err := prometheus.Init(); err != nil {
			return nil, fmt.Errorf("error while initializing prometheus: %s", err)
		}
		metrics := s.Metrics(prometheus)
		return &metrics, nil
	}
	return nil, fmt.Errorf("unknown metrics type: %s", *t)
}

func (sc *ScalerApp) Scale() {
	fmt.Printf("App: %+v \n", sc.appDefinition)
	fmt.Printf("Service: %+v \n", *sc.service)
	fmt.Printf("Provider: %+v \n", *sc.provider)
	fmt.Printf("Metrics: %+v \n", *sc.metrics)
	panic("not implemented")
}
