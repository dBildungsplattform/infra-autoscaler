package metricssource

import (
	s "scaler/shared"
	"testing"
)

func TestValidatePrometheusOK(t *testing.T) {
	prometheus := &Prometheus{
		PrometheusConfig: PrometheusConfig{
			Url: "https://prometheus.example.com",
		},
	}
	s.ValidatePass(t, prometheus)
}

func TestValidatePrometheusConfigNotOK(t *testing.T) {
	prometheus := &Prometheus{}
	s.ValidateFail(t, prometheus)
}

func TestValidatePrometheusConfigBadUrl(t *testing.T) {
	prometheus := &Prometheus{
		PrometheusConfig: PrometheusConfig{
			Url: "not a url",
		},
	}
	s.ValidateFail(t, prometheus)
	// Missing scheme
	prometheus.PrometheusConfig.Url = "prometheus.example.com"
	s.ValidateFail(t, prometheus)
	prometheus.PrometheusConfig.Url = "/prometheus.example.com"
	s.ValidateFail(t, prometheus)
}

func TestInitPrometheusOK(t *testing.T) {
	config := PrometheusConfig{
		Url: "https://prometheus.example.com",
	}
	prometheus := &Prometheus{
		PrometheusConfig: config,
	}
	err := prometheus.Init()
	if err != nil {
		t.Errorf("Error: %s", err)
	}
}
