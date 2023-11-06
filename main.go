package main

import (
	"scaler/BBB"
	"scaler/scaler"
)

func main() {
	var bbb scaler.Service = &BBB.BBBService{}
	bbb.Init()
}
