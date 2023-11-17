package providers

import (
	s "scaler/shared"
	"testing"
)

func TestValidatePrometheusOK(t *testing.T) {
	prometheus := &PrometheusConfig{
		Url: "url",
	}
	s.ValidatePass(t, prometheus)
}

func TestValidatePrometheusConfigNotOK(t *testing.T) {
	prometheus := &PrometheusConfig{}
	s.ValidateFail(t, prometheus)
}

func TestInitPrometheusClientOK(t *testing.T) {
	config := PrometheusConfig{
		Url: "https://prometheus.example.com",
	}
	client := &PrometheusClient{}
	err := client.Init(config)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
}
