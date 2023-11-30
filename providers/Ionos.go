package providers

import (
	"context"
	"fmt"
	"regexp"
	s "scaler/shared"
	"time"

	ic "github.com/ionos-cloud/sdk-go/v6"
)

type ProviderConfig struct {
	Username     s.StringFromEnv `yaml:"username"`
	Password     s.StringFromEnv `yaml:"password"`
	ContractId   s.IntFromEnv    `yaml:"contract_id"`
	ServerSource *s.ServerSource `yaml:"server_source"`
}

type Ionos struct {
	Config   ProviderConfig `yaml:"ionos_config"`
	Contract *ic.Contract   `yaml:"-"`
	Stage    s.Stage        `yaml:"-"`
	Api      ic.APIClient   `yaml:"-"`
	Servers  []s.Server     `yaml:"-"`
}

func (i *Ionos) Init() error {
	i.Api = *ic.NewAPIClient(ic.NewConfiguration(
		string(i.Config.Username),
		string(i.Config.Password),
		"",
		""))
	if err := validateAndLoadContract(i); err != nil {
		return fmt.Errorf("error while validating contract: %s", err)
	}
	if err := initMetricsExporter("ionos"); err != nil {
		return fmt.Errorf("error while registering metrics: %s", err)
	}
	return nil
}

func (i Ionos) GetServers(depth int) ([]s.Server, error) {
	var servers []s.Server
	var err error

	if i.Config.ServerSource.Static != nil {
		err = getServersStatic(&servers, i)
	} else if i.Config.ServerSource.Dynamic != nil {
		err = getServersDynamic(&servers, i, depth)
	}
	if err != nil {
		errorsTotalCounter.Inc()
		return nil, fmt.Errorf("error while getting servers: %s", err)
	}
	return servers, nil
}

func getServersStatic(servers *[]s.Server, i Ionos) error {
	for _, serverSource := range *i.Config.ServerSource.Static {
		dcServer, _, err := i.Api.ServersApi.DatacentersServersFindById(
			context.TODO(),
			serverSource.DatacenterId,
			serverSource.ServerId).XContractNumber(int32(i.Config.ContractId)).Execute()
		if err != nil {
			return fmt.Errorf("error while getting servers: %s", err)
		}
		addServer(servers, dcServer, serverSource.DatacenterId)
	}
	return nil
}

func getServersDynamic(servers *[]s.Server, i Ionos, depth int) error {
	for _, datacenterId := range i.Config.ServerSource.Dynamic.DatacenterIds {
		fmt.Println("Getting servers from datacenter: ", datacenterId)
		dcServers, _, err := i.Api.ServersApi.DatacentersServersGet(context.TODO(), datacenterId).Depth(int32(depth)).XContractNumber(int32(i.Config.ContractId)).Execute()
		if err != nil {
			return fmt.Errorf("error while getting servers: %s", err)
		}
		for _, dcServer := range *dcServers.Items {
			if match, _ := regexp.MatchString(i.Config.ServerSource.Dynamic.ServerNameRegex, *dcServer.Properties.Name); match {
				addServer(servers, dcServer, datacenterId)
			}
		}
	}
	return nil
}

func addServer(servers *[]s.Server, dcServer ic.Server, datacenterId string) {
	*servers = append(*servers, s.Server{
		DatacenterId:    datacenterId,
		ServerId:        *dcServer.Id,
		ServerName:      *dcServer.Properties.Name,
		CpuArchitecture: *dcServer.Properties.CpuFamily,
		ServerCpu:       *dcServer.Properties.Cores,
		ServerRam:       *dcServer.Properties.Ram,
		ServerCpuUsage:  0,
		ServerRamUsage:  0,
		LastUpdated:     time.Now(),
	})
}

func (i Ionos) SetServerResources(server s.Server, targetRes s.ScaleResource) error {
	if targetRes.Cpu.Direction == s.ScaleUp && targetRes.Mem.Direction == s.ScaleDown || targetRes.Cpu.Direction == s.ScaleDown && targetRes.Mem.Direction == s.ScaleUp {
		return fmt.Errorf("cannot scale cpu and memory in opposite directions")
	}

	// Validate and scale up server
	targetServer := *ic.NewServer(ic.ServerProperties{
		Cores: &targetRes.Cpu.Amount,
		Ram:   &targetRes.Mem.Amount,
	})
	validServer := validateServer(targetServer, *i.Contract)
	if !validServer {
		errorsTotalCounter.Inc()
		return fmt.Errorf("server is not valid")
	}

	fmt.Printf("targetServer: %+v \n", targetServer.Properties) // Check mode
	_, _, err := i.Api.ServersApi.DatacentersServersPut(context.TODO(), server.DatacenterId, server.ServerId).Server(targetServer).XContractNumber(int32(i.Config.ContractId)).Execute()
	if err != nil {
		errorsTotalCounter.Inc()
		return fmt.Errorf("error while setting server resources: %s", err)
	}
	return nil
}

func validateServer(server ic.Server, contract ic.Contract) bool {
	if *server.Properties.Cores > *contract.Properties.ResourceLimits.CoresPerServer || *server.Properties.Ram > *contract.Properties.ResourceLimits.RamPerServer {
		return false
	}

	return true
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

func validateAndLoadContract(i *Ionos) error {
	if i.Stage == s.DevStage { // Assume API is not initialized in dev stage
		return nil
	}
	contracts, _, err := i.Api.ContractResourcesApi.ContractsGet(context.TODO()).Execute()
	if err != nil {
		return fmt.Errorf("error while retrieving contract: %s", err)
	}
	for _, contract := range *contracts.Items {
		if *contract.Properties.ContractNumber == int64(i.Config.ContractId) {
			i.Contract = &contract
			return nil
		}
	}
	return fmt.Errorf("contract_id %d not found", i.Config.ContractId)
}
