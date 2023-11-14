package providers

import (
	"fmt"
	s "scaler/shared"
)

type Ionos struct {
	ProviderName  string
	Username      string
	Password      string
	DatacenterIds []string
	ProviderType  s.ProviderType
	// PostgresSource *PostgresSource `yaml:"postgres_source"`
}

func (i Ionos) Get_login_id() string {
	return i.Username
}

func (i Ionos) Get_login_secret() string {
	return i.Password
}

func (i Ionos) Get_type() s.ProviderType {
	return i.ProviderType
}

func (i Ionos) Get_name() string {
	return i.ProviderName
}

func (p Ionos) Validate() error {
	if p.Username == "" {
		return fmt.Errorf("username is empty")
	}
	if p.Password == "" {
		return fmt.Errorf("password is empty")
	}
	return nil
}
