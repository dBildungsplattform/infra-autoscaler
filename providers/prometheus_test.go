package providers

import (
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
		Url: "https://prometheus.example.com",
	}
	client := &PrometheusClient{}
	err := client.Init(config)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
}
