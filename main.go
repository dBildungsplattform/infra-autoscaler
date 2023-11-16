package main

import (
	c "scaler/core"
)

func main() {
	configPath := "config/scaler_config.yml" // TODO make this a command line arg
	app := c.InitApp(configPath)
	app.Scale()
}
