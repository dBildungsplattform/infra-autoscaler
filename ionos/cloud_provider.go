package ionos

import "fmt"

type CloudProvider struct {
	Username        string
	Password        string
	InstancesSource *ServerSource `yaml:"instances_source"`
}

func (c CloudProvider) Validate() error {
	if c.Username == "" {
		return fmt.Errorf("username is empty")
	}
	if c.Password == "" {
		return fmt.Errorf("password is empty")
	}
	if c.InstancesSource == nil {
		return fmt.Errorf("instances_source is nil")
	}
	if err := c.InstancesSource.Validate(); err != nil {
		return err
	}
	return nil
}
