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
		return nil, fmt.Errorf("error while opening config file: %s", err)
	}

	app, err := s.LoadConfig[s.AppDefinition](configFile)
	if err != nil {
		return nil, fmt.Errorf("error while loading app config: %s", err)
	}

	service, err := initService(&app.ServiceType, configFile)
	if err != nil {
		return nil, fmt.Errorf("error while initializing service: %s", err)
	}

	provider, err := initProvider(&app.ProviderType, configFile)
	if err != nil {
		return nil, fmt.Errorf("error while initializing provider: %s", err)
	}

	metrics, err := initMetrics(&app.MetricsType, configFile)
	if err != nil {
		return nil, fmt.Errorf("error while initializing metrics: %s", err)
	}

	return &ScalerApp{
		appDefinition: app,
		service:       service,
		provider:      provider,
		metrics:       metrics,
	}, nil
}

func initService(t *s.ServiceType, configFile []byte) (*s.Service, error) {
	switch *t {
	case s.BBB:
		bbb, err := s.LoadConfig[services.BBBService](configFile)
		if err != nil {
			return nil, fmt.Errorf("error while loading BBB config: %s", err)
		}
		service := s.Service(bbb)
		return &service, nil
	case s.Postgres:
		postgres, err := s.LoadConfig[services.PostgresService](configFile)
		if err != nil {
			return nil, fmt.Errorf("error while loading postgres config: %s", err)
		}
		service := s.Service(postgres)
		return &service, nil
	}
	return nil, fmt.Errorf("unknown service type: %s", *t)
}

func initProvider(t *s.ProviderType, configFile []byte) (*s.Provider, error) {
	switch *t {
	case s.Ionos:
		ionos, load_err := s.LoadConfig[providers.Ionos](configFile)
		if load_err != nil {
			return nil, fmt.Errorf("error while loading ionos config: %s", load_err)
		}

		init_err := ionos.Init()
		if init_err != nil {
			return nil, fmt.Errorf("error while initializing ionos: %s", init_err)
		}

		provider := s.Provider(ionos)
		return &provider, nil
	}
	return nil, fmt.Errorf("unknown provider type: %s", *t)
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
	servers, err := (*sc.provider).GetServers(1)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Servers: %+v \n", servers)
}
