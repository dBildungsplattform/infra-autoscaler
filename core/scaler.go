package core

import (
	"fmt"
	"scaler/metricssource"
	"scaler/providers"
	"scaler/services"
	s "scaler/shared"
)

type ScalerApp struct {
	appDefinition *s.AppDefinition
	service       *s.Service
	provider      *s.Provider
	metricsSource *s.MetricsSource
}

// TODO: make these configurable
var memIncrease int32 = 1
var memDecrease int32 = 1
var cpuIncrease int32 = 1
var cpuDecrease int32 = 1

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

	metricsSource, err := initMetricsSource(&app.MetricsSourceType, configFile)
	if err != nil {
		return nil, fmt.Errorf("error while initializing metrics: %s", err)
	}

	initMetricsExporter()

	return &ScalerApp{
		appDefinition: app,
		service:       service,
		provider:      provider,
		metricsSource: metricsSource,
	}, nil
}

func initService(t *s.ServiceType, configFile []byte) (*s.Service, error) {
	switch *t {
	case s.BBB:
		bbb, err := s.LoadConfig[services.BBBService](configFile)
		if err != nil {
			return nil, fmt.Errorf("error while loading BBB config: %s", err)
		}
		init_err := bbb.Init()
		if init_err != nil {
			return nil, fmt.Errorf("error while initializing BBB: %s", init_err)
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

func initMetricsSource(t *s.MetricsSourceType, configFile []byte) (*s.MetricsSource, error) {
	switch *t {
	case s.Prometheus:
		prometheus, err := s.LoadConfig[metricssource.Prometheus](configFile)
		if err != nil {
			return nil, fmt.Errorf("error while loading prometheus config: %s", err)
		}
		if err := prometheus.Init(); err != nil {
			return nil, fmt.Errorf("error while initializing prometheus: %s", err)
		}
		metrics := s.MetricsSource(prometheus)
		return &metrics, nil
	}
	return nil, fmt.Errorf("unknown metrics type: %s", *t)
}

func (sc *ScalerApp) Scale() {
	var service s.Service = *sc.service

	cyclesCounter.Inc()

	servers, err := (*sc.provider).GetServers(1)
	if err != nil {
		panic(err)
	}

	go sc.calculateMetrics(servers)

	metricsSource := *sc.metricsSource
	for _, server := range servers {

		// Get current resource usage
		currCpuUsage, err := metricsSource.GetServerCpuUsage(server.ServerName)
		if err != nil {
			panic(err)
		}
		server.ServerCpuUsage = currCpuUsage
		currMemUsage, err := metricsSource.GetServerMemoryUsage(server.ServerName)
		if err != nil {
			panic(err)
		}
		server.ServerRamUsage = currMemUsage

		// Get scaling proposal from service
		targetResource, err := service.ShouldScale(server)
		if err != nil {
			panic(err)
		}

		// Scale
		if targetResource.Cpu.Direction == s.ScaleUp {
			if sc.appDefinition.ScalingMode == s.DirectScaling {
				targetResource.Cpu.Amount = server.ServerCpu + cpuIncrease
			}
			// If heuristic scaling, get the increase from the service
		}
		if targetResource.Mem.Direction == s.ScaleUp {
			if sc.appDefinition.ScalingMode == s.DirectScaling {
				targetResource.Mem.Amount = server.ServerRam + memIncrease
			}
			// If heuristic scaling, get the increase from the service
		}

		// Scale down
		if targetResource.Cpu.Direction == s.ScaleDown {
			if sc.appDefinition.ScalingMode == s.DirectScaling {
				targetResource.Cpu.Amount = server.ServerCpu - cpuDecrease
			}
			// If heuristic scaling, get the decrease from the service
		}
		if targetResource.Mem.Direction == s.ScaleDown {
			if sc.appDefinition.ScalingMode == s.DirectScaling {
				targetResource.Mem.Amount = server.ServerRam - memDecrease
			}
			// If heuristic scaling, get the decrease from the service
		}

		provider := *sc.provider
		provider.SetServerResources(server, targetResource)
		// Placeholder to have metrics
		lastScaleTimeGauge.SetToCurrentTime()
	}
}
