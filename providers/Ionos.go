package providers

import (
	"context"
	"fmt"
	"regexp"
	s "scaler/shared"
	"time"

	icDbaas "github.com/ionos-cloud/sdk-go-dbaas-postgres"
	ic "github.com/ionos-cloud/sdk-go/v6"
	"golang.org/x/exp/slog"
)

type ProviderConfig struct {
	Token         s.StringFromEnv  `yaml:"token"`
	Username      s.StringFromEnv  `yaml:"username"`
	Password      s.StringFromEnv  `yaml:"password"`
	ContractId    s.IntFromEnv     `yaml:"contract_id"`
	ServerSource  *s.ServerSource  `yaml:"server_source"`
	ClusterSource *s.ClusterSource `yaml:"cluster_source"`
}

type Ionos struct {
	Config   ProviderConfig    `yaml:"ionos_config"`
	Contract *ic.Contract      `yaml:"-"`
	Stage    s.Stage           `yaml:"-"`
	Api      ic.APIClient      `yaml:"-"`
	DbaasApi icDbaas.APIClient `yaml:"-"`
}

func (i *Ionos) Init() error {
	i.Api = *ic.NewAPIClient(ic.NewConfiguration(
		string(i.Config.Username),
		string(i.Config.Password),
		string(i.Config.Token),
		""))
	i.DbaasApi = *icDbaas.NewAPIClient(icDbaas.NewConfiguration(
		string(i.Config.Username),
		string(i.Config.Password),
		string(i.Config.Token), 
		""))
	if err := validateAndLoadContract(i); err != nil {
		return fmt.Errorf("error while validating contract: %s", err)
	}
	if err := initMetricsExporter("ionos"); err != nil {
		return fmt.Errorf("error while registering metrics: %s", err)
	}
	return nil
}

func (i Ionos) getServers(depth int) ([]*s.Server, error) {
	var servers []*s.Server
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

func getServersStatic(servers *[]*s.Server, i Ionos) error {
	for _, serverSource := range *i.Config.ServerSource.Static {
		dcServer, _, err := i.Api.ServersApi.DatacentersServersFindById(
			context.TODO(),
			serverSource.DatacenterId,
			serverSource.ServerId).XContractNumber(int32(i.Config.ContractId)).Execute()
		if err != nil {
			return fmt.Errorf("error while getting server %s in datacenter %s: %s", serverSource.ServerId, serverSource.DatacenterId, err)
		}
		slog.Info(fmt.Sprintf("Found server %s (%s) in datacenter %s\n", *dcServer.Properties.Name, serverSource.ServerId, serverSource.DatacenterId))
		server := responseToServer(dcServer, serverSource.DatacenterId)
		*servers = append(*servers, &server)
	}
	return nil
}

func getServersDynamic(servers *[]*s.Server, i Ionos, depth int) error {
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
				server := responseToServer(dcServer, datacenterId)
				*servers = append(*servers, &server)
			}
		}
		slog.Info(fmt.Sprintf("Matched %d servers in datacenter %s\n", matchCount, datacenterId))
	}
	return nil
}

func responseToServer(response ic.Server, datacenterId string) s.Server {
	return s.Server{
		DatacenterId:    datacenterId,
		ServerId:        *response.Id,
		ServerName:      *response.Properties.Name,
		CpuArchitecture: *response.Properties.CpuFamily,
		ResourceState: s.ResourceState{
			Cpu: &s.CpuResourceState{
				CurrentCores: *response.Properties.Cores,
				CurrentUsage: 0,
			},
			Memory: &s.MemoryResourceState{
				CurrentBytes: *response.Properties.Ram,
				CurrentUsage: 0,
			},
		},
		LastUpdated: time.Now(),
		Ready:       *response.Properties.VmState == "RUNNING" && *response.Metadata.State == "AVAILABLE",
	}
}

func (i Ionos) updateServer(server s.Server, scalingProposal s.ResourceScalingProposal) error {
	// When scaling in different directions, scaling up overrides scaling down
	if scalingProposal.Cpu.Direction == s.ScaleUp && scalingProposal.Mem.Direction == s.ScaleDown {
		scalingProposal.Mem.Direction = s.ScaleNone
		scalingProposal.Mem.Amount = 0
	}
	if scalingProposal.Cpu.Direction == s.ScaleDown && scalingProposal.Mem.Direction == s.ScaleUp {
		scalingProposal.Cpu.Direction = s.ScaleNone
		scalingProposal.Cpu.Amount = 0
	}

	if scalingProposal.Cpu.Direction == s.ScaleNone && scalingProposal.Mem.Direction == s.ScaleNone {
		return nil
	}

	targetCpu := server.ResourceState.Cpu.CurrentCores + scalingProposal.Cpu.Amount
	targetMem := server.ResourceState.Memory.CurrentBytes + scalingProposal.Mem.Amount

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
	_, _, err = i.Api.ServersApi.DatacentersServersPut(context.TODO(), server.DatacenterId, server.ServerId).Server(targetServer).XContractNumber(int32(i.Config.ContractId)).Execute()
	if err != nil {
		errorsTotalCounter.Inc()
		return fmt.Errorf("error while setting server resources: %s", err)
	}
	return nil
}

func (i Ionos) getClusters() ([]*s.Cluster, error) {
	var clusters []*s.Cluster
	var err error

	if i.Config.ClusterSource.Static != nil {
		err = getClustersStatic(&clusters, i)
	} else if i.Config.ClusterSource.Dynamic != nil {
		err = getClustersDynamic(&clusters, i)
	}
	if err != nil {
		errorsTotalCounter.Inc()
		return nil, err
	}
	return clusters, nil
}

func getClustersStatic(clusters *[]*s.Cluster, i Ionos) error {
	for _, clusterId := range i.Config.ClusterSource.Static.ClusterIds {
		response, _, err := i.DbaasApi.ClustersApi.ClustersFindById(context.TODO(), clusterId).Execute()
		if err != nil {
			return fmt.Errorf("error while getting cluster %s: %s", clusterId, err)
		}

		slog.Info(fmt.Sprintf("Found cluster %s (%s)\n", *response.Properties.DisplayName, clusterId))
		cluster := responseToCluster(response)
		*clusters = append(*clusters, &cluster)
	}
	return nil
}

func getClustersDynamic(clusters *[]*s.Cluster, i Ionos) error {
	clustersResponse, _, err := i.DbaasApi.ClustersApi.ClustersGet(context.TODO()).Execute()
	if err != nil {
		return fmt.Errorf("error while getting clusters: %s", err)
	}
	for _, response := range *clustersResponse.Items {
		if match, _ := regexp.MatchString(i.Config.ClusterSource.Dynamic.ClusterNameRegex, *response.Properties.DisplayName); match {
			slog.Info(fmt.Sprintf("Matched cluster %s (%s)\n", *response.Properties.DisplayName, *response.Id))
			cluster := responseToCluster(response)
			*clusters = append(*clusters, &cluster)
		}
	}
	return nil
}

func responseToCluster(response icDbaas.ClusterResponse) s.Cluster {
	return s.Cluster{
		ClusterId:   *response.Id,
		ClusterName: *response.Properties.DisplayName,
		ResourceState: s.ResourceState{
			Cpu: &s.CpuResourceState{
				CurrentCores: *response.Properties.Cores,
				CurrentUsage: 0,
			},
			Memory: &s.MemoryResourceState{
				CurrentBytes: *response.Properties.Ram,
				CurrentUsage: 0,
			},
		},
		LastUpdated: time.Now(),
		Ready:       *response.Metadata.State == "AVAILABLE",
	}
}

func (i Ionos) updateCluster(cluster s.Cluster, scalingProposal s.ResourceScalingProposal) error {
	if scalingProposal.Cpu.Direction == s.ScaleNone && scalingProposal.Mem.Direction == s.ScaleNone {
		return nil
	}

	targetCpu := cluster.ResourceState.Cpu.CurrentCores + scalingProposal.Cpu.Amount
	targetMem := cluster.ResourceState.Memory.CurrentBytes + scalingProposal.Mem.Amount

	// Validate and scale cluster
	targetClusterProperties := *icDbaas.NewPatchClusterProperties()
	targetClusterProperties.Cores = &targetCpu
	targetClusterProperties.Ram = &targetMem
	targetCluster := *icDbaas.NewPatchClusterRequest()
	targetCluster.Properties = &targetClusterProperties

	validCluster := validateCluster(targetCluster, *i.Contract)
	if !validCluster {
		errorsTotalCounter.Inc()
		return fmt.Errorf("cluster is not valid")
	}

	slog.Info(fmt.Sprintf("Target for cluster %s: %d cores, %d bytes\n", cluster.ClusterName, *targetCluster.Properties.Cores, *targetCluster.Properties.Ram))
	_, _, err := i.DbaasApi.ClustersApi.ClustersPatch(context.TODO(), cluster.ClusterId).PatchClusterRequest(targetCluster).Execute()
	if err != nil {
		errorsTotalCounter.Inc()
		return fmt.Errorf("error while setting cluster resources: %s", err)
	}
	return nil
}

func validateCluster(cluster icDbaas.PatchClusterRequest, contract ic.Contract) bool {
	// TODO: Validate patch request against contract limits
	//if *cluster.Properties.Cores > *contract.Properties.ResourceLimits.CoresPerCluster || *cluster.Properties.Ram > *contract.Properties.ResourceLimits.RamPerCluster {
	//	return false
	//}

	return true
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
	if (i.Config.Token == "") && ((i.Config.Username == "") || i.Config.Password == "") {

		errorMessage := ""
		if i.Config.Token == "" {
			errorMessage += "Token "
		}
		if i.Config.Username == "" {
			errorMessage += "Username "
		}
		if i.Config.Password == "" {
			errorMessage += "Password "
		}

		errorMessage += "is empty"
		return fmt.Errorf(errorMessage)
	}

	if i.Config.ServerSource == nil && i.Config.ClusterSource == nil {
		return fmt.Errorf("no scaled object source provided")
	} else if i.Config.ServerSource != nil && i.Config.ClusterSource != nil {
		return fmt.Errorf("both server and cluster source provided, only one must be set")
	} else {
		if i.Config.ServerSource != nil {
			err := i.Config.ServerSource.Validate()
			if err != nil {
				return err
			}
		}
		if i.Config.ClusterSource != nil {
			err := i.Config.ClusterSource.Validate()
			if err != nil {
				return err
			}
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

func (i Ionos) GetScaledObjects() ([]s.ScaledObject, error) {
	var objects []s.ScaledObject
	if i.Config.ServerSource != nil {
		servers, err := i.getServers(1)
		if err != nil {
			return nil, fmt.Errorf("error while getting servers: %s", err)
		}
		for _, server := range servers {
			objects = append(objects, server)
		}
	}
	if i.Config.ClusterSource != nil {
		clusters, err := i.getClusters()
		if err != nil {
			return nil, fmt.Errorf("error while getting clusters: %s", err)
		}
		for _, cluster := range clusters {
			objects = append(objects, cluster)
		}
	}
	return objects, nil
}

func (i Ionos) UpdateScaledObject(object s.ScaledObject, scalingProposal s.ResourceScalingProposal) error {
	switch objectType := object.(type) {
	case *s.Server:
		server := objectType
		err := i.updateServer(*server, scalingProposal)
		if err != nil {
			return fmt.Errorf("error while updating server %s: %s", server.ServerName, err)
		}
	case *s.Cluster:
		cluster := objectType
		err := i.updateCluster(*cluster, scalingProposal)
		if err != nil {
			return fmt.Errorf("error while updating cluster %s: %s", cluster.ClusterName, err)
		}
	default:
		return fmt.Errorf("unsupported scaled object type: %s", object.GetType())
	}
	return nil
}
