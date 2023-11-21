package metrics

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/config"
	"github.com/prometheus/common/model"
)

type PrometheusConfig struct {
	Url string
	// TODO: Use StringFromEnv once DBP-363 is merged
	Token string
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
	if _, err := url.ParseRequestURI(p.PrometheusConfig.Url); err != nil {
		return fmt.Errorf("url is invalid: %v", err)
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
	return nil
}

func (p *Prometheus) QueryServerCPUUsage(serverLabels string) string {
	// TODO: Move the queries to config ?
	cpuUsageQuery := fmt.Sprintf("avg without (mode,cpu) (1 - rate(node_cpu_seconds_total{mode=\"idle\",%s}[30s]))", serverLabels)
	return cpuUsageQuery
}

func (p *Prometheus) QueryServerMemoryUsage(serverLabels string) string {
	memoryUsageQuery := fmt.Sprintf("(node_memory_MemFree_bytes + node_memory_Cached_bytes + node_memory_Buffers_bytes) / node_memory_MemTotal_bytes{%s}", serverLabels)
	return memoryUsageQuery
}

func (p *Prometheus) Query(query string) (float64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*timeout)
	defer cancel()
	result, warnings, err := p.API.Query(ctx, query, time.Now(), v1.WithTimeout(timeout))
	if err != nil {
		return 0, err
	}
	if len(warnings) > 0 {
		fmt.Printf("Warnings: %v\n", warnings)
	}
	if result.Type() == model.ValScalar {
		vector := result.(model.Vector)
		if len(vector) == 0 {
			return 0, fmt.Errorf("no data found")
		}
		if len(vector) != 1 {
			// TODO: Should duplicate metrics trigger an error?
			fmt.Printf("Unexpected vector length: %v\n", len(vector))
		}
		return float64(vector[0].Value), nil
	} else {
		return 0, fmt.Errorf("unexpected type: %v", result.Type())
	}
}
