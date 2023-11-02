package bbb

import (
	types "scaler/Interface"
)

type Service struct {
	config *types.ServiceConfig
	state  *types.ServiceState
}

// get_state implements types.Service.
func (Service) Get_state() *types.ServiceState {
	panic("unimplemented")
}

// load_config implements types.Service.
func (Service) Load_config() *types.ServiceConfig {
	panic("unimplemented")
}

// set_state implements types.Service.
func (Service) Set_state(*types.ServiceState) {
	panic("unimplemented")
}

// start_service implements types.Service.
func (Service) Start_service() {
	panic("unimplemented")
}
