package providers

import (
	"os"

	s "scaler/shared"
	"testing"
)

func TestValidatePrometheusOK(t *testing.T) {
	prometheus := &Prometheus{
		PrometheusConfig: PrometheusConfig{
			Url: "url",
		},
	}
	s.ValidatePass(t, prometheus)
}

func TestValidatePrometheusNotOK(t *testing.T) {
	prometheus := &Prometheus{}
	s.ValidateFail(t, prometheus)
}

func TestInitPrometheusClientOK(t *testing.T) {
	config := PrometheusConfig{
		Url:   "https://grafana.dbildungscloud.org/api/datasources/proxy/1/",
		Token: os.Getenv("PROMETHEUS_TOKEN"),
	}
	client := &PrometheusClient{}
	err := client.Init(config)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
}
