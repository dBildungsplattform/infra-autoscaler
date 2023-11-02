package types

type Service interface {
	Load_config() *ServiceConfig
	Get_state() *ServiceState
	Set_state(*ServiceState)
	Start_service()
}

type ServiceConfig struct {
	name string
}

type ServiceState struct {
	name string
}
