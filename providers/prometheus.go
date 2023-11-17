package providers

import (
	"context"
	"fmt"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/config"
	"github.com/prometheus/common/model"
)

type PrometheusConfig struct {
	Url   string
	Token string
}

type Prometheus struct {
	PrometheusConfig PrometheusConfig `yaml:"prometheus_config"`
}

func (p Prometheus) Validate() error {
	if p.PrometheusConfig.Url == "" {
		return fmt.Errorf("url is empty")
	}
	return nil
}

type PrometheusClient struct {
	Client api.Client
	API    v1.API
}

func (c *PrometheusClient) Init(prometheusConfig PrometheusConfig) error {
	var err error = nil
	c.Client, err = api.NewClient(api.Config{
		Address:      prometheusConfig.Url,
		RoundTripper: config.NewAuthorizationCredentialsRoundTripper("Bearer", config.Secret(prometheusConfig.Token), api.DefaultRoundTripper),
	})
	if err != nil {
		return err
	}
	c.API = v1.NewAPI(c.Client)
	return nil
}

func (c *PrometheusClient) QueryCPUUsage(serverLabels string) (float64, error) {
	cpuUsageQuery := fmt.Sprintf("avg without (mode,cpu) (1 - rate(node_cpu_seconds_total{mode=\"idle\",%s}[30s]))", serverLabels)
	return c.Query(cpuUsageQuery)
}

func (c *PrometheusClient) QueryMemoryUsage(serverLabels string) (float64, error) {
	memoryUsageQuery := fmt.Sprintf("(node_memory_MemFree_bytes + node_memory_Cached_bytes + node_memory_Buffers_bytes) / node_memory_MemTotal_bytes{%s}", serverLabels)
	return c.Query(memoryUsageQuery)
}

func (c *PrometheusClient) Query(query string) (float64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, warnings, err := c.API.Query(ctx, query, time.Now(), v1.WithTimeout(5*time.Second))
	if err != nil {
		return 0, err
	}
	if len(warnings) > 0 {
		fmt.Printf("Warnings: %v\n", warnings)
	}
	switch result.Type() {
	case model.ValVector:
		vector := result.(model.Vector)
		if len(vector) == 0 {
			return 0, fmt.Errorf("no data found")
		}
		if len(vector) != 1 {
			// TODO: Should duplicate metrics trigger an error?
			fmt.Printf("Unexpected vector length: %v\n", len(vector))
		}
		return float64(vector[0].Value), nil
	default:
		return 0, fmt.Errorf("enexpected type: %v", result.Type())
	}
}
