package config

import (
	"fmt"
	"scaler/ionos"
)

type CloudProvider struct {
	Ionos *ionos.CloudProvider
}

func (s CloudProvider) Validate() error {
	if s.Ionos != nil {
		if err := s.Ionos.Validate(); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("cloud_provider is empty")
}
