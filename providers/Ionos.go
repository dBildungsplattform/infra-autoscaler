package providers

import (
	"context"
	"fmt"
	s "scaler/shared"

	ic "github.com/ionos-cloud/sdk-go/v6"
)

type ProviderConfig struct {
	Username     s.StringFromEnv `yaml:"username"`
	Password     s.StringFromEnv `yaml:"password"`
	ServerSource *s.ServerSource `yaml:"server_source"`
}

type Ionos struct {
	Config ProviderConfig `yaml:"ionos_config"`
	Api    ic.APIClient   `yaml:"-"`
}

func (i *Ionos) Init() error {
	i.Api = *ic.NewAPIClient(ic.NewConfiguration(
		string(i.Config.Username),
		string(i.Config.Password),
		"",
		""))
	return nil
}

func (i Ionos) GetServers(depth int) ([]s.Server, error) {
	var servers []s.Server
	for _, datacenterId := range i.Config.ServerSource.Dynamic.DatacenterIds {
		fmt.Println(datacenterId)
		dc_servers, _, err := i.Api.ServersApi.DatacentersServersGet(context.Background(), datacenterId).Depth(int32(depth)).Execute()
		if err != nil {
			return nil, fmt.Errorf("error while getting servers: %s", err)
		}
		for _, dc_server := range *dc_servers.Items {
			server := s.Server{
				DatacenterId: datacenterId,
				ServerId:     *dc_server.Id,
				ServerCpu:    float64(*dc_server.Properties.Cores),
				ServerRam:    float64(*dc_server.Properties.Ram),
			}
			servers = append(servers, server)
		}
	}
	return servers, nil
}

func (i Ionos) Validate() error {
	if i.Config.Username == "" {
		return fmt.Errorf("username is empty")
	}
	if i.Config.Password == "" {
		return fmt.Errorf("password is empty")
	}
	if i.Config.ServerSource == nil {
		return fmt.Errorf("server_source is nil")
	} else {
		if err := i.Config.ServerSource.Validate(); err != nil {
			return err
		}
	}
	return nil
}
