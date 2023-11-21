package main

import (
	"flag"
	c "scaler/core"
)

func main() {
	configPath := flag.String("config", "config/scaler_config.yml", "path to config file")
	flag.Parse()

	app, err := c.InitApp(*configPath)
	if err != nil {
		panic(err)
	}
	app.Scale()
}
