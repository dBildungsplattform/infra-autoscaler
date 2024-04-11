//Telekom Cloud OTC DNS management for the BBB autoscaler implementation
//Creation of a new OTC DNS record
//Deletion of a existing OTC DNS record
//Updating of the IP in a OTC DNS record
package otcdns

import (
	"fmt"

	otc "github.com/opentelekomcloud/gophertelekomcloud"
	otcos "github.com/opentelekomcloud/gophertelekomcloud/openstack"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dns/v2/recordsets"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dns/v2/zones"
)

const (
	dnsRecordTypeA     string = "A"
	dnsRecordDescription string = "BBB Autoscaler"
)

//
// The DNS client we use to trigger our DNS actions.
//
type OtcDnsClient struct {
	Sc *otc.ServiceClient
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

// ===========================================================================
// Zones
// ===========================================================================

//
// Retrieves a Zone data structure by its name.
// https://pkg.go.dev/github.com/opentelekomcloud/gophertelekomcloud@v0.3.2/openstack/dns/v2/zones
// github.com/opentelekomcloud/gophertelekomcloud/openstack/dns/v2/zones
//
func (dnsClient *OtcDnsClient) GetHostedZone(zoneName string) (*zones.Zone, error) {

	listOpts := zones.ListOpts{
		Name: zoneName,
	}

	allPages, err := zones.List(dnsClient.Sc, listOpts).AllPages()
	if err != nil {
		return nil, fmt.Errorf("zone %s not found: %s", zoneName, err)
	}

	allZones, err := zones.ExtractZones(allPages)
	if err != nil {
		return nil, fmt.Errorf("zone %s extraction failed: %s", zoneName, err)
	}

	// Debug
	//for _, zone := range allZones {
	//	fmt.Printf("%+v\n", zone)
	//}

	// We need exactly 1 zone to operate on
	if len(allZones) != 1 {
		return nil, fmt.Errorf("zone query with %s returned %d zones. Expected: 1", zoneName, len(allZones))
	}

	return &allZones[0], nil
}

// ===========================================================================
// RecordSets
// ===========================================================================

//Get Recordset by DNS Name
func (dnsClient *OtcDnsClient) GetARecordSet(zone *zones.Zone, dnsName string) (*recordsets.RecordSet, error) {
	listOpts := recordsets.ListOpts{
		Type: dnsRecordTypeA,
		Name: dnsName,
	}

	allPages, err := recordsets.ListByZone(dnsClient.Sc, zone.ID, listOpts).AllPages()
	if err != nil {
		return nil, fmt.Errorf("list records failed for dns entry %s: %s", dnsName, err)
	}

	allRRs, err := recordsets.ExtractRecordSets(allPages)
	if err != nil {
		return nil, fmt.Errorf("extract recordset failed for dns entry %s: %s", dnsName, err)
	}

	// Debug
	//for _, rr := range allRRs {
	//	fmt.Printf("%+v\n", rr)
	//}

	if len(allRRs) == 1 {
		// We need exactly 1 recordset to operate on
		return &allRRs[0], nil
	} else if len(allRRs) == 0 {
		// Query was successful, but no results
		return nil, nil
	} else {
		// More than 1 result.
		return nil, fmt.Errorf("query with %s returned %d recordsets. Expected: 1", dnsName, len(allRRs))
	}
}

//
//Create new A DNS RecordSet with DNS Name and IP
//
func (dnsClient *OtcDnsClient) CreateRecordSet (zone *zones.Zone, dnsName string, ipValue string) (*recordsets.RecordSet, error) {

	createOpts := recordsets.CreateOpts{
		Name:    dnsName,
		Description: dnsRecordDescription,
		Type:    dnsRecordTypeA,
		Records: []string{ipValue},
		TTL:     300,
	}

	var pCreatedRecordset *recordsets.RecordSet
	pCreatedRecordset, err := recordsets.Create(dnsClient.Sc, zone.ID, createOpts).Extract()
    if err != nil {
        fmt.Printf("Error creating DNS record: %v\n", err)
        return
    }
	return pCreatedRecordset, nil
}

//
//Delete a RecordSet by DNS Name
//
func (dnsClient *OtcDnsClient) DeleteRecordSet (zone *zones.Zone, dnsName string) {
	//get DNS details
	recordSet, err := dnsClient.GetARecordSet(zone.ID, dnsName)

	//Call delete function
	err := recordsets.Delete(dnsClient.Sc, zone.ID, recordSet.ID).ExtractErr()
	if err != nil {
		return fmt.Errorf("deletion of record with zoneId %s and recordsetId %s failed: %s", zone.ID, recordset.ID, err)
	}

	return nil
}

//
//Updating a RecordSet IP by DNS Name
//
func (dnsClient *OtcDnsClient) UpdateRecordSet(zone *zones.Zone, dnsName string, newIPValue string) (*recordsets.RecordSet, error) {
    //get DNS details
	recordSet, err := dnsClient.GetARecordSet(zone.ID, dnsName)

	updateOpts := recordsets.UpdateOpts{
        Records: []string{newIPValue},
    }

    updatedRecordSet, err := recordsets.Update(dnsClient, zone.ID, recordSet.ID, updateOpts).Extract()
    if err != nil {
        fmt.Printf("Error updating DNS record: %v\n", err)
        return nil, err
    }
    return updatedRecordSet, nil
}