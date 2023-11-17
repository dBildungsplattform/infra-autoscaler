package providers

import (
	"fmt"
	s "scaler/shared"
)

type ProviderConfig struct {
	Username     string
	Password     string
	ServerSource *s.ServerSource `yaml:"server_source"`
	// PostgresSource *PostgresSource `yaml:"postgres_source"`
	PrometheusConfig PrometheusConfig `yaml:"prometheus_config"`
}

type Ionos struct {
	IonosConfig ProviderConfig `yaml:"ionos_config"`
}

func (i Ionos) Get_login_id() string {
	return i.IonosConfig.Username
}

func (i Ionos) Get_login_secret() string {
	return i.IonosConfig.Password
}

func (p Ionos) Validate() error {
	if p.IonosConfig.Username == "" {
		return fmt.Errorf("username is empty")
	}
	if p.IonosConfig.Password == "" {
		return fmt.Errorf("password is empty")
	}
	if p.IonosConfig.ServerSource == nil {
		return fmt.Errorf("server_source is nil")
	} else {
		if err := p.IonosConfig.ServerSource.Validate(); err != nil {
			return err
		}
	}
	if err := p.IonosConfig.PrometheusConfig.Validate(); err != nil {
		return err
	}
	return nil
}
