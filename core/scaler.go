package core

import (
	"fmt"
	"scaler/metricssource"
	"scaler/providers"
	"scaler/services"
	s "scaler/shared"
	"time"

	"golang.org/x/exp/slog"
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

func (sc ScalerApp) scaleObject(object s.ScaledObject) error {
	var err error
	resourceState := object.GetResourceState()
	resourceState.Cpu.CurrentUsage, err = sc.metricsSource.GetCpuUsage(object)
	if err != nil {
		return fmt.Errorf("error while getting cpu usage for %s %s: %s", object.GetType(), object.GetName(), err)
	}
	slog.Info(fmt.Sprintf("CPU usage for %s %s: %f\n", object.GetType(), object.GetName(), resourceState.Cpu.CurrentUsage))
	resourceState.Memory.CurrentUsage, err = sc.metricsSource.GetMemoryUsage(object)
	if err != nil {
		return fmt.Errorf("error while getting memory usage for %s %s: %s", object.GetType(), object.GetName(), err)
	}
	slog.Info(fmt.Sprintf("Memory usage for %s %s: %f\n", object.GetType(), object.GetName(), resourceState.Memory.CurrentUsage))
	object.SetResourceState(resourceState)

	// Get scaling proposal from service
	scalingProposal, err := sc.service.ComputeScalingProposal(object)
	if err != nil {
		return fmt.Errorf("error while getting scaling proposal for %s %s: %s", object.GetType(), object.GetName(), err)
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
	slog.Info(fmt.Sprintf("Scaling proposal for %s: %+v\n", object.GetName(), scalingProposal))

	err = sc.provider.UpdateScaledObject(object, scalingProposal)
	if err != nil {
		return fmt.Errorf("error while setting resources for %s %s: %s", object.GetType(), object.GetName(), err)
	}
	if scalingProposal.Cpu.Direction != s.ScaleNone || scalingProposal.Mem.Direction != s.ScaleNone {
		lastScaleTimeGauge.SetToCurrentTime()
	}
	return nil
}

func (sc *ScalerApp) Scale() {
	for {
		cyclesCounter.Inc()

		scaledObjects, err := sc.provider.GetScaledObjects()
		if err != nil {
			slog.Error(fmt.Sprint("Error while getting scaled objects: ", err))
		}

		go sc.calculateMetrics(scaledObjects)

		for _, scaledObject := range scaledObjects {
			err := sc.scaleObject(scaledObject)
			if err != nil {
				slog.Error(err.Error())
			}
		}
		time.Sleep(time.Duration(sc.service.GetCycleTimeSeconds()) * time.Second)
	}
}
