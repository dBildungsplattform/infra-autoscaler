package Ionos

import (
	c "scaler/common"
)

type Provider struct {
	ProviderName  string
	Username      string
	Password      string
	DatacenterIds []string
	ProviderType  c.ProviderType
}

func (Ionos Provider) Get_login_id() string {
	return Ionos.Username
}

func (Ionos Provider) Get_login_secret() string {
	return Ionos.Password
}

func (Ionos Provider) Get_type() c.ProviderType {
	return Ionos.ProviderType
}

func (Ionos Provider) Get_name() string {
	return Ionos.ProviderName
}
