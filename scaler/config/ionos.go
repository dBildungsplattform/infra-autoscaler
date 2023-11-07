package config

import (
	"fmt"
	"regexp"
)

type IonosServerInstancesSource struct {
	DatacenterIds   []string `yaml:"datacenter_ids"`
	ServerNameRegex string   `yaml:"server_name_regex"`
}

type InlineIonosServerInstancesSource []struct {
	DatacenterId string `yaml:"datacenter_id"`
	ServerId     string `yaml:"server_id"`
}

func (ionos IonosServerInstancesSource) Validate() error {
	if len(ionos.DatacenterIds) == 0 {
		return fmt.Errorf("ionos.datacenter_ids is empty")
	}
	if ionos.ServerNameRegex == "" {
		return fmt.Errorf("ionos.server_name_regex is empty")
	}
	if _, err := regexp.Compile(ionos.ServerNameRegex); err != nil {
		return fmt.Errorf("ionos.server_name_regex is invalid: %s", err)
	}
	return nil
}

func (inlineIonos InlineIonosServerInstancesSource) Validate() error {
	if len(inlineIonos) == 0 {
		return fmt.Errorf("inline_ionos is empty")
	}
	for index, server := range inlineIonos {
		if server.DatacenterId == "" {
			return fmt.Errorf("inline_ionos[%d].datacenter_id is empty", index)
		}
		if server.ServerId == "" {
			return fmt.Errorf("inline_ionos[%d].server_id is empty", index)
		}
	}
	return nil
}
