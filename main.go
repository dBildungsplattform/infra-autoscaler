package main

import (
	"flag"
	c "scaler/core"
)

func main() {
	configPath := flag.String("config", "config/scaler_config.yml", "path to config file")
	flag.Parse()

	app := c.InitApp(*configPath)
	app.Scale()
}
