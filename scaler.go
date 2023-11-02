package scaler

import (
	bbb "scaler/BBB"
	types "scaler/Interface"
)

func main() {
	var BBB types.Service = bbb.Service{}
	BBB.Load_config()
	BBB.Start_service()
}
