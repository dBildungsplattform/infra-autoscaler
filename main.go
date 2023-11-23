package main

import (
	"flag"
	"net/http"
	c "scaler/core"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	configPath := flag.String("config", "config/scaler_config.yml", "path to config file")
	flag.Parse()

	app, err := c.InitApp(*configPath)
	if err != nil {
		panic(err)
	}
	go app.Scale()

	http.Handle("/metrics", promhttp.Handler())
	// TODO make port configurable
	http.ListenAndServe(":8080", nil)
}
