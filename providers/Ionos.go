package providers

import (
	"context"
	"fmt"
	"regexp"
	s "scaler/shared"
	"time"

	ic "github.com/ionos-cloud/sdk-go/v6"
	"golang.org/x/exp/slog"
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
		return nil, err
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
			return fmt.Errorf("error while getting server %s in datacenter %s: %s", serverSource.ServerId, serverSource.DatacenterId, err)
		}
		slog.Info(fmt.Sprintf("Found server %s (%s) in datacenter %s\n", *dcServer.Properties.Name, serverSource.ServerId, serverSource.DatacenterId))
		addServer(servers, dcServer, serverSource.DatacenterId)
	}
	return nil
}

func getServersDynamic(servers *[]s.Server, i Ionos, depth int) error {
	for _, datacenterId := range i.Config.ServerSource.Dynamic.DatacenterIds {
		slog.Info(fmt.Sprint("Getting servers from datacenter: ", datacenterId))
		dcServers, _, err := i.Api.ServersApi.DatacentersServersGet(context.TODO(), datacenterId).Depth(int32(depth)).XContractNumber(int32(i.Config.ContractId)).Execute()
		if err != nil {
			return fmt.Errorf("error while getting servers in datacenter %s: %s", datacenterId, err)
		}
		slog.Info(fmt.Sprintf("Found %d servers in datacenter %s\n", len(*dcServers.Items), datacenterId))
		matchCount := 0
		for _, dcServer := range *dcServers.Items {
			if match, _ := regexp.MatchString(i.Config.ServerSource.Dynamic.ServerNameRegex, *dcServer.Properties.Name); match {
				matchCount++
				addServer(servers, dcServer, datacenterId)
			}
		}
		slog.Info(fmt.Sprintf("Matched %d servers in datacenter %s\n", matchCount, datacenterId))
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

func (i Ionos) SetServerResources(server s.Server, scalingProposal s.ScaleResource) error {
	if scalingProposal.Cpu.Direction == s.ScaleUp && scalingProposal.Mem.Direction == s.ScaleDown || scalingProposal.Cpu.Direction == s.ScaleDown && scalingProposal.Mem.Direction == s.ScaleUp {
		return fmt.Errorf("cannot scale cpu and memory in opposite directions")
	}

	if scalingProposal.Cpu.Direction == s.ScaleNone && scalingProposal.Mem.Direction == s.ScaleNone {
		return nil
	}

	targetCpu := server.ServerCpu + scalingProposal.Cpu.Amount
	targetMem := server.ServerRam + scalingProposal.Mem.Amount

	// Validate and scale server
	targetServer := *ic.NewServer(ic.ServerProperties{
		Cores: &targetCpu,
		Ram:   &targetMem,
	})
	err := validateServer(targetServer, *i.Contract)
	if err != nil {
		errorsTotalCounter.Inc()
		return fmt.Errorf("target server for %s is not valid: %s", *targetServer.Properties.Name, err)
	}

	slog.Info(fmt.Sprintf("Target for server %s: %d cores, %d bytes\n", server.ServerName, *targetServer.Properties.Cores, *targetServer.Properties.Ram))
	//_, _, err := i.Api.ServersApi.DatacentersServersPut(context.TODO(), server.DatacenterId, server.ServerId).Server(targetServer).XContractNumber(int32(i.Config.ContractId)).Execute()
	//if err != nil {
	//	errorsTotalCounter.Inc()
	//	return fmt.Errorf("error while setting server resources: %s", err)
	//}
	return nil
}

func validateServer(server ic.Server, contract ic.Contract) error {
	coresOk := *server.Properties.Cores <= *contract.Properties.ResourceLimits.CoresPerServer
	ramOk := *server.Properties.Ram <= *contract.Properties.ResourceLimits.RamPerServer
	errorMessage := ""
	if !coresOk {
		errorMessage += fmt.Sprintf("cores %d are above contract limit %d", *server.Properties.Cores, *contract.Properties.ResourceLimits.CoresPerServer)
	}
	if !ramOk {
		if errorMessage != "" {
			errorMessage += ", "
		}
		errorMessage += fmt.Sprintf("memory %d is above contract limit %d", *server.Properties.Ram, *contract.Properties.ResourceLimits.RamPerServer)
	}
	if errorMessage != "" {
		return fmt.Errorf(errorMessage)
	}
	return nil
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
