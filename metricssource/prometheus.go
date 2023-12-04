package metricssource

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	s "scaler/shared"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/config"
	"github.com/prometheus/common/model"
	"golang.org/x/exp/slog"
)

type PrometheusConfig struct {
	Url   string
	Token s.StringFromEnv `yaml:"token"`
}

// TODO: Move the timeout to config ?
var timeout = 5 * time.Second

type Prometheus struct {
	PrometheusConfig PrometheusConfig `yaml:"prometheus_config"`
	API              v1.API           `yaml:"-"`
}

func (p Prometheus) Validate() error {
	if p.PrometheusConfig.Url == "" {
		return fmt.Errorf("url is empty")
	}
	urlParsed, err := url.Parse(p.PrometheusConfig.Url)
	if err != nil {
		return fmt.Errorf("url is invalid: %v", err)
	}
	if urlParsed.Scheme != "http" && urlParsed.Scheme != "https" {
		return fmt.Errorf("url scheme is invalid: %s", urlParsed.Scheme)
	}
	if urlParsed.Host == "" {
		return fmt.Errorf("url host is empty")
	}
	return nil
}

func (p *Prometheus) Init() error {
	var err error = nil
	var rt http.RoundTripper = api.DefaultRoundTripper
	if p.PrometheusConfig.Token != "" {
		rt = config.NewAuthorizationCredentialsRoundTripper("Bearer", config.Secret(p.PrometheusConfig.Token), api.DefaultRoundTripper)
	}
	client, err := api.NewClient(api.Config{
		Address:      p.PrometheusConfig.Url,
		RoundTripper: rt,
	})
	if err != nil {
		return err
	}
	p.API = v1.NewAPI(client)
	if err := initMetricsExporter("prometheus"); err != nil {
		return fmt.Errorf("error while registering metrics: %s", err)
	}
	return nil
}

func (p *Prometheus) QueryServerCPUUsage(serverLabels string) string {
	// TODO: Move the queries to config ?
	cpuUsageQuery := fmt.Sprintf("avg without (mode,cpu) (1 - rate(node_cpu_seconds_total{mode=\"idle\",%s}[30s]))", serverLabels)
	return cpuUsageQuery
}

func (p *Prometheus) QueryServerMemoryUsage(serverLabels string) string {
	memoryUsageQuery := fmt.Sprintf("1 - (node_memory_MemFree_bytes + node_memory_Cached_bytes + node_memory_Buffers_bytes) / node_memory_MemTotal_bytes{%s}", serverLabels)
	return memoryUsageQuery
}

func (p *Prometheus) Query(query string) (float32, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*timeout)
	defer cancel()
	result, warnings, err := p.API.Query(ctx, query, time.Now(), v1.WithTimeout(timeout))
	if err != nil {
		errorsTotalCounter.Inc()
		return 0, err
	}
	if len(warnings) > 0 {
		slog.Warn(fmt.Sprintf("Warnings: %v\n", warnings))
	}
	if result.Type() == model.ValVector {
		vector := result.(model.Vector)
		if len(vector) == 0 {
			errorsTotalCounter.Inc()
			return 0, fmt.Errorf("no data found")
		}
		if len(vector) != 1 {
			// Duplicate metrics can occur if Prometheus has multiple jobs with the same targets
			// This is not a scaler error but we should log it
			slog.Warn(fmt.Sprintf("Unexpected vector length: %v\n", len(vector)))
		}
		return float32(vector[0].Value), nil
	} else {
		errorsTotalCounter.Inc()
		return 0, fmt.Errorf("unexpected type: %v", result.Type())
	}
}

func (p Prometheus) GetServerCpuUsage(serverName string) (float32, error) {
	serverLabels := fmt.Sprintf("instance=~\"%s\"", serverName)
	query := p.QueryServerCPUUsage(serverLabels)
	return p.Query(query)
}

func (p Prometheus) GetServerMemoryUsage(serverName string) (float32, error) {
	serverLabels := fmt.Sprintf("instance=~\"%s\"", serverName)
	query := p.QueryServerMemoryUsage(serverLabels)
	return p.Query(query)
}
