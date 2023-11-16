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
}

type Ionos struct {
	ProviderConfig ProviderConfig `yaml:"provider_config"`
}

func (i Ionos) Get_login_id() string {
	return i.ProviderConfig.Username
}

func (i Ionos) Get_login_secret() string {
	return i.ProviderConfig.Password
}

func (p Ionos) Validate() error {
	if p.ProviderConfig.Username == "" {
		return fmt.Errorf("username is empty")
	}
	if p.ProviderConfig.Password == "" {
		return fmt.Errorf("password is empty")
	}
	if p.ProviderConfig.ServerSource == nil {
		return fmt.Errorf("server_source is nil")
	} else {
		if err := p.ProviderConfig.ServerSource.Validate(); err != nil {
			return err
		}
	}
	return nil
}
