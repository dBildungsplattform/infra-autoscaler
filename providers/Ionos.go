package providers

import (
	"context"
	"fmt"
	"regexp"
	s "scaler/shared"
	"time"

	icDbaas "github.com/ionos-cloud/sdk-go-dbaas-postgres"
	ic "github.com/ionos-cloud/sdk-go/v6"
)

type ProviderConfig struct {
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
	Servers  []s.Server        `yaml:"-"`
	Clusters []s.Cluster       `yaml:"-"`
}

func (i *Ionos) Init() error {
	// go sdk api client
	i.Api = *ic.NewAPIClient(ic.NewConfiguration(
		string(i.Config.Username),
		string(i.Config.Password),
		"",
		""))

	// dbaas api client
	i.DbaasApi = *icDbaas.NewAPIClient(icDbaas.NewConfiguration(
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

func (i Ionos) getServers(depth int) ([]s.Server, error) {
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
		Ready:           *dcServer.Properties.VmState == "RUNNING" && *dcServer.Metadata.State == "AVAILABLE",
	})
}

func (i Ionos) SetScaledObject(obj s.ScaledObject, proposal s.ScaleResource) error {
	switch obj.GetType() {
	case s.ServerType:
		server := obj.(s.Server)
		err := i.setServerResources(server, proposal)
		if err != nil {
			return fmt.Errorf("error while setting resources for server %s: %s", server.ServerName, err)
		}
	case s.ClusterType:
		cluster := obj.(s.Cluster)
		err := i.setClusterResources(cluster, proposal)
		if err != nil {
			return fmt.Errorf("error while setting resources for cluster %s: %s", cluster.ClusterName, err)
		}
	}
	return nil
}

func (i Ionos) setServerResources(server s.Server, scalingProposal s.ScaleResource) error {
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

	targetCpu := server.ServerCpu + scalingProposal.Cpu.Amount
	targetMem := server.ServerRam + scalingProposal.Mem.Amount

	// Validate and scale server
	targetServer := *ic.NewServer(ic.ServerProperties{
		Cores: &targetCpu,
		Ram:   &targetMem,
	})
	validServer := validateServer(targetServer, *i.Contract)
	if !validServer {
		errorsTotalCounter.Inc()
		return fmt.Errorf("server is not valid")
	}

	fmt.Printf("Target for server %s: %d cores, %d bytes\n", server.ServerName, *targetServer.Properties.Cores, *targetServer.Properties.Ram)
	//_, _, err := i.Api.ServersApi.DatacentersServersPut(context.TODO(), server.DatacenterId, server.ServerId).Server(targetServer).XContractNumber(int32(i.Config.ContractId)).Execute()
	//if err != nil {
	//	errorsTotalCounter.Inc()
	//	return fmt.Errorf("error while setting server resources: %s", err)
	//}
	return nil
}

func (i Ionos) getFilteredClusters(clusters *icDbaas.ClusterList) error {
	var err error

	filter := i.Config.ClusterSource.ClusterFilterName
	if filter == "" {
		*clusters, _, err = i.DbaasApi.ClustersApi.ClustersGet(context.TODO()).Execute()
		if err != nil {
			return fmt.Errorf("error while getting clusters: %s", err)
		}
	} else {
		*clusters, _, err = i.DbaasApi.ClustersApi.ClustersGet(context.TODO()).FilterName(filter).Execute()
		if err != nil {
			return fmt.Errorf("error while getting filtered clusters: %s", err)
		}
	}
	return nil
}

func (i Ionos) applyClusterRegexFilter(clusters *icDbaas.ClusterList, filteredClusters *[]s.Cluster) error {
	if i.Config.ClusterSource.ClusterNameRegex != "" {
		for _, cluster := range *clusters.Items {
			if match, _ := regexp.MatchString(i.Config.ClusterSource.ClusterNameRegex, *cluster.Properties.DisplayName); match {
				fmt.Println("Matched cluster: ", *cluster.Properties.DisplayName)
				fmt.Println("Cluster regex: ", i.Config.ClusterSource.ClusterNameRegex)
				addCluster(filteredClusters, cluster)
			}
		}
	}
	return nil
}

func addCluster(scaledClusters *[]s.Cluster, cluster icDbaas.ClusterResponse) {
	*scaledClusters = append(*scaledClusters, s.Cluster{
		ClusterId:          *cluster.Id,
		ClusterName:        *cluster.Properties.DisplayName,
		ClusterCpu:         *cluster.Properties.Cores,
		ClusterRam:         *cluster.Properties.Ram,
		ClusterStorageSize: *cluster.Properties.StorageSize,
		ClusterStorageType: string(*cluster.Properties.StorageType),
		LastUpdated:        time.Now(),
		Ready:              *cluster.Metadata.State == "AVAILABLE",
	})
}

func (i Ionos) getClusters() ([]s.Cluster, error) {
	var filteredClusters icDbaas.ClusterList
	err := i.getFilteredClusters(&filteredClusters)
	if err != nil {
		return nil, err
	}

	var scaledClusters []s.Cluster
	err = i.applyClusterRegexFilter(&filteredClusters, &scaledClusters)
	if err != nil {
		return nil, err
	}

	return scaledClusters, err
}

func (i Ionos) setClusterResources(cluster s.Cluster, scalingProposal s.ScaleResource) error {
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

	targetCpu := cluster.ClusterCpu + scalingProposal.Cpu.Amount
	targetMem := cluster.ClusterRam + scalingProposal.Mem.Amount

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

	fmt.Printf("Target for cluster %s: %d cores, %d bytes\n", cluster.ClusterName, *targetCluster.Properties.Cores, *targetCluster.Properties.Ram)
	//_, _, err := i.DbaasApi.ClustersApi.ClustersPatch(context.TODO(), cluster.ClusterId).PatchClusterRequest(targetCluster).Execute()
	//if err != nil {
	//	errorsTotalCounter.Inc()
	//	return fmt.Errorf("error while setting cluster resources: %s", err)
	//}
	return nil
}

func validateCluster(cluster icDbaas.PatchClusterRequest, contract ic.Contract) bool {
	// TODO: Validate patch request against contract limits
	//if *cluster.Properties.Cores > *contract.Properties.ResourceLimits.CoresPerCluster || *cluster.Properties.Ram > *contract.Properties.ResourceLimits.RamPerCluster {
	//	return false
	//}

	return true
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
	var scaledObjects []s.ScaledObject
	if i.Config.ServerSource != nil {
		servers, err := i.getServers(1)
		if err != nil {
			return nil, fmt.Errorf("error while getting servers: %s", err)
		}
		for _, server := range servers {
			scaledObjects = append(scaledObjects, server)
		}
	}
	if i.Config.ClusterSource != nil {
		clusters, err := i.getClusters()
		fmt.Printf("Clusters: %+v\n", clusters)
		if err != nil {
			return nil, fmt.Errorf("error while getting clusters: %s", err)
		}
		for _, cluster := range clusters {
			scaledObjects = append(scaledObjects, cluster)
		}
	}
	return scaledObjects, nil
}
