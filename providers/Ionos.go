package providers

import (
	"context"
	"fmt"
	"regexp"
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

	if i.Config.ServerSource.Static != nil {
		getServersStatic(&servers, i)
	} else if i.Config.ServerSource.Dynamic != nil {
		getServersDynamic(&servers, i, depth)
	}
	return servers, nil
}

func getServersStatic(servers *[]s.Server, i Ionos) {
	for _, serverSource := range *i.Config.ServerSource.Static {
		dcServer, _, err := i.Api.ServersApi.DatacentersServersFindById(context.Background(), serverSource.DatacenterId, serverSource.ServerId).Execute()
		if err != nil {
			fmt.Printf("error while getting servers: %s", err)
		}
		*servers = append(*servers, s.Server{
			DatacenterId: serverSource.DatacenterId,
			ServerId:     *dcServer.Id,
			ServerCpu:    float64(*dcServer.Properties.Cores),
			ServerRam:    float64(*dcServer.Properties.Ram),
		})
	}
}

func getServersDynamic(servers *[]s.Server, i Ionos, depth int) {
	for _, datacenterId := range i.Config.ServerSource.Dynamic.DatacenterIds {
		dcServers, _, err := i.Api.ServersApi.DatacentersServersGet(context.Background(), datacenterId).Depth(int32(depth)).Execute()
		if err != nil {
			fmt.Printf("error while getting servers: %s", err)
		}
		for _, dcServer := range *dcServers.Items {
			if match, _ := regexp.MatchString(i.Config.ServerSource.Dynamic.ServerNameRegex, *dcServer.Properties.Name); match {
				*servers = append(*servers, s.Server{
					DatacenterId: datacenterId,
					ServerId:     *dcServer.Id,
					ServerCpu:    float64(*dcServer.Properties.Cores),
					ServerRam:    float64(*dcServer.Properties.Ram),
				})
			}
		}
	}
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
