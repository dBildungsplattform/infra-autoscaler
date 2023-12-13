package main

import (
	"flag"
	c "scaler/core"

	"golang.org/x/exp/slog"
)

func main() {
	configPath := flag.String("config", "config/scaler_config.yml", "path to config file")
	flag.Parse()

	app, err := c.InitApp(*configPath)
	if err != nil {
		panic(err)
	}
	go app.Scale()
	slog.Error(app.ServeMetrics().Error())
}
