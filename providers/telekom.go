//Telekom Cloud OTC DNS management
package otcdns

import (
	"fmt"

	otc "github.com/opentelekomcloud/gophertelekomcloud"
	otcos "github.com/opentelekomcloud/gophertelekomcloud/openstack"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dns/v2/recordsets"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dns/v2/zones"
)

const (
	dnsRecordTypeTxt     string = "TXT"
	dnsRecordDescription string = "ACME Challenge"
	acmeChallengePrefix  string = "_acme-challenge."
)

//
// The DNS client we use to trigger our DNS actions.
//
type OtcDnsClient struct {
	Sc *otc.ServiceClient

	//
	// Optional subdomain, which will be inserted between "_acme-challenge." and the zone name.
	//
	Subdomain string
}

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

type RecordSetConfig struct {
	ZoneID   string
	Name     string
	Description string
	Type     string
	Records  []string
	TTL      int
}

//IPs sollen durch IONOS DHCP auto erzeugt werden.
// func (i Ionos) CreateIPAddress () string{
// 	Ipblocks := *openapiclient.NewIpBlock(*openapiclient.NewIpBlockProperties("eu/txl", int32(5)))

//     createdIP, _, err := i.Api.IPBlocksApi.IpblocksPost(ctx, datacenterId).Ipblocks(ipblock).Pretty(pretty).Depth(depth).XContractNumber(xContractNumber).Execute()
//     if err != nil {
//         return nil, fmt.Errorf("error when calling `IPBlocksApi.IpblocksPost`: %w", err)
//     }
//     return &createdIP, nil
// }

// func (i Ionos) DeleteIPAddress (IPValue string) {
// 	deleteIP, _, err := i.Api.IPBlocksApi.IpblocksDelete(ctx, IPValue).Pretty(pretty).Depth(depth).XContractNumber(xContractNumber).Execute()
//     if err != nil {
//         return nil, fmt.Errorf("error when calling `IPBlocksApi.IpblocksDelete`: %w", err)
//     }
// }

//Nutzung von SDK
func (dnsClient *OtcDnsClient) CreateDNSRecord (dnsClient *golangsdk.ServiceClient, config RecordSetConfig) (*recordsets.RecordSet, error) {

	createOpts := recordsets.CreateOpts{
		Name:    config.Name,
		Description: config.Description
		Type:    config.Type,
		Records: config.Records,
		TTL:     &config.TTL,
	}

	recordset, err := recordsets.Create(dnsClient, config.ZoneID, createOpts).Extract()
    if err != nil {
        fmt.Printf("Error creating DNS record: %v\n", err)
        return
    }
	return recordset, nil //Includes the recordSetID
}

func (dnsClient *OtcDnsClient) DeleteDNSRecord(dnsClient *golangsdk.ServiceClient, zoneID string, recordSetID string) error {
    err := recordsets.Delete(dnsClient, zoneID, recordSetID).ExtractErr()
    if err != nil {
        fmt.Printf("Error deleting DNS record: %v\n", err)
        return err
    }
    return nil
}

func (dnsClient *OtcDnsClient) UpdateDNSRecord(dnsClient *golangsdk.ServiceClient, zoneID string, recordSetID string, challengeValues []string) (*recordsets.RecordSet, error) {
    if len(challengeValues) == 0 {
		return nil, fmt.Errorf("update TXT records failed. The challengeValue records must have at least one entry")
	}
	updateOpts := recordsets.UpdateOpts{
        Records: challengeValues,
    }

    updatedRecordSet, err := recordsets.Update(dnsClient, zoneID, recordSetID, updateOpts).Extract()
    if err != nil {
        fmt.Printf("Error updating DNS record: %v\n", err)
        return nil, err
    }
    return updatedRecordSet, nil
}
