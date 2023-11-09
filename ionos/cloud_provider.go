package ionos

import "fmt"

type CloudProvider struct {
	Username     string
	Password     string
	ServerSource *ServerSource `yaml:"server_source"`
	// PostgresSource *PostgresSource `yaml:"postgres_source"`
}

func (c CloudProvider) Validate() error {
	if c.Username == "" {
		return fmt.Errorf("username is empty")
	}
	if c.Password == "" {
		return fmt.Errorf("password is empty")
	}
	ss := c.ServerSource
	if ss == nil {
		return fmt.Errorf("instances_source is nil")
	}
	if err := ss.Validate(); err != nil {
		return err
	}
	return nil
}
