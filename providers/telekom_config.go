// Config of Telekom Cloud
package providers

import (
	"encoding/json"
	"fmt"

	otc "github.com/opentelekomcloud/gophertelekomcloud"
	otcos "github.com/opentelekomcloud/gophertelekomcloud/openstack"
)

const (
	envPrefix          string = "OS_"
	OtcProfileNameUser string = "otcuser"
	OtcProfileNameAkSk string = "otcaksk"
)

var EnvOS = otcos.NewEnv(envPrefix)

// ===========================================================================
// Configs
// ===========================================================================

// Creates a new DNSv2 ServiceClient.
// See also gophertelekomcloud/acceptance/clients/clients.go
func NewDNSV2ClientWithAuth(authOpts otc.AuthOptionsProvider, endpointOpts otc.EndpointOpts) (*OtcDnsClient, error) {

	providerClient, err := getProviderClientWithAccessKeyAuth(authOpts)
	if err != nil {
		return nil, fmt.Errorf("cannot create providerClient. %s", err)
	}

	serviceClient, err := otcos.NewDNSV2(providerClient, endpointOpts)
	if err != nil {
		return nil, fmt.Errorf("cannot create serviceClient. %s", err)
	}

	return &OtcDnsClient{Sc: serviceClient}, nil
}

func getProviderClientWithAccessKeyAuth(authOpts otc.AuthOptionsProvider) (*otc.ProviderClient, error) {
	provider, err := otcos.AuthenticatedClient(authOpts)
	if err != nil {
		return nil, fmt.Errorf("provider creation has failed: %s", err)
	}
	return provider, nil
}

func getCloud() (*otcos.Cloud, error) {
	return getCloudProfile(OtcProfileNameUser)
}

func copyCloud(src *otcos.Cloud) (*otcos.Cloud, error) {
	srcJson, err := json.Marshal(src)
	if err != nil {
		return nil, fmt.Errorf("error marshalling cloud: %s", err)
	}

	res := new(otcos.Cloud)
	if err := json.Unmarshal(srcJson, res); err != nil {
		return nil, fmt.Errorf("error unmarshalling cloud: %s", err)
	}

	return res, nil
}

func getCloudProfile(otcProfileName string) (*otcos.Cloud, error) {

	cloud, err := EnvOS.Cloud(otcProfileName)
	if err != nil {
		return nil, fmt.Errorf("error constructing cloud configuration: %s", err)
	}

	cloud, err = copyCloud(cloud)
	if err != nil {
		return nil, fmt.Errorf("error copying cloud: %s", err)
	}

	return cloud, nil
}

func getProviderClient() (*otc.ProviderClient, error) {
	return getProviderClientProfile(OtcProfileNameUser)
}

func getProviderClientProfile(otcProfileName string) (*otc.ProviderClient, error) {

	client, err := EnvOS.AuthenticatedClient(otcProfileName)
	if err != nil {
		return nil, fmt.Errorf("cloud and provider creation has failed: %s", err)
	}

	return client, nil
}

// Creates a new DNSv2 ServiceClient.
// See also gophertelekomcloud/acceptance/clients/clients.go
func NewDNSV2Client() (*OtcDnsClient, error) {
	cloudsConfig, err := getCloud()
	if err != nil {
		return nil, err
	}
	endpointOpts := otc.EndpointOpts{
		Region: cloudsConfig.RegionName,
	}

	providerClient, err := getProviderClient()
	if err != nil {
		return nil, err
	}

	serviceClient, err := otcos.NewDNSV2(providerClient, endpointOpts)
	if err != nil {
		return nil, err
	}

	return &OtcDnsClient{Sc: serviceClient}, nil
}
