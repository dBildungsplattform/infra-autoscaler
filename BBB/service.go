package BBB

import (
	"fmt"
	"scaler/scaler"
)

type BBBService struct {
	name   string
	state  scaler.ServiceState
	config scaler.ServiceConfig[any]
}

func (bbb *BBBService) Init() {
	fmt.Println("Initializing BBB service")
	bbb.name = "BBB"
	bbb.state = scaler.ServiceState{}
	bbb.config = scaler.ServiceConfig[any](*scaler.Load_config())
	fmt.Println("BBB service initialized")
	fmt.Printf("Config: \n %+v \n", bbb.config)
}

func (bbb *BBBService) Get_state() scaler.ServiceState {
	return bbb.state
}

func (bbb *BBBService) Get_config() scaler.ServiceConfig[any] {
	return bbb.config
}
