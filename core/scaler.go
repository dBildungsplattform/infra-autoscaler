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
	service       s.Service
	provider      s.Provider
	metricsSource s.MetricsSource
}

// TODO: make these configurable
var memIncrease int32 = 1024
var memDecrease int32 = -1024
var cpuIncrease int32 = 1
var cpuDecrease int32 = -1

func InitApp(configPath string) (*ScalerApp, error) {
	configFile, err := s.OpenConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("error while opening config file: %s", err)
	}

	app, err := s.LoadConfig[s.AppDefinition](configFile)
	if app.ScalingMode == "" {
		app.ScalingMode = s.DirectScaling
	}
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
		service:       *service,
		provider:      *provider,
		metricsSource: *metricsSource,
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
	cyclesCounter.Inc()

	servers, err := sc.provider.GetServers(1)
	if err != nil {
		panic(err)
	}

	go sc.calculateMetrics(servers)

	for _, server := range servers {

		// Get current resource usage
		server.ServerCpuUsage, err = sc.metricsSource.GetServerCpuUsage(server.ServerName)
		if err != nil {
			panic(err)
		}
		server.ServerRamUsage, err = sc.metricsSource.GetServerMemoryUsage(server.ServerName)
		if err != nil {
			// TODO: Should be possible to recover from this, at least for a couple of cycles
			panic(err)
		}

		// Get scaling proposal from service
		scalingProposal, err := sc.service.ShouldScale(server)
		fmt.Printf("Scaling proposal for %+v: %+v\n", server.ServerName, scalingProposal)
		if err != nil {
			panic(err)
		}

		// Scale
		// Override heuristic target resource
		if sc.appDefinition.ScalingMode == s.DirectScaling {
			// Scale up CPU
			if scalingProposal.Cpu.Direction == s.ScaleUp {
				scalingProposal.Cpu.Amount = cpuIncrease
			}
			// Scale up RAM
			if scalingProposal.Mem.Direction == s.ScaleUp {
				scalingProposal.Mem.Amount = memIncrease
			}
			// Scale down CPU
			if scalingProposal.Cpu.Direction == s.ScaleDown {
				scalingProposal.Cpu.Amount = cpuDecrease
			}
			// Scale down RAM
			if scalingProposal.Mem.Direction == s.ScaleDown {
				scalingProposal.Mem.Amount = memDecrease
			}
		}

		err = sc.provider.SetServerResources(server, scalingProposal)
		if err != nil {
			panic(err)
		}
		// Placeholder to have metrics
		lastScaleTimeGauge.SetToCurrentTime()
	}
}
