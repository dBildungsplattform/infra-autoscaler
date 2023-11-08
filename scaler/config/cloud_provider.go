package config

import "scaler/ionos"

type CloudProvider struct {
	Ionos *ionos.CloudProvider
}

func (s CloudProvider) Validate() error {
	if s.Ionos != nil {
		return s.Ionos.Validate()
	}
	return nil
}
