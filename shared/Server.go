package shared

import (
	"fmt"
	"regexp"
	"time"
)

type Server struct {
	DatacenterId    string
	ServerId        string
	ServerName      string
	CpuArchitecture string
	ServerCpu       int32
	ServerRam       int32
	ServerCpuUsage  float32
	ServerRamUsage  float32
	LastUpdated     time.Time
}

type ServerSource struct {
	Static  *ServerStaticSource  `yaml:"static"`
	Dynamic *ServerDynamicSource `yaml:"dynamic"`
}

type ServerDynamicSource struct {
	DatacenterIds   []string `yaml:"datacenter_ids"`
	ServerNameRegex string   `yaml:"server_name_regex"`
}

type ServerStaticSource []struct {
	DatacenterId string `yaml:"datacenter_id"`
	ServerId     string `yaml:"server_id"`
}

func (ionos ServerSource) Validate() error {
	static, dynamic := ionos.Static, ionos.Dynamic
	if static == nil && dynamic == nil {
		return fmt.Errorf("ionos.static and ionos.dynamic are nil, one must be set")
	}
	if static != nil && dynamic != nil {
		return fmt.Errorf("ionos.static and ionos.dynamic are both set, only one must be set")
	}
	if static != nil {
		err := static.Validate()
		if err != nil {
			return err
		}
	}
	if dynamic != nil {
		err := dynamic.Validate()
		if err != nil {
			return err
		}
	}
	return nil
}

func (ionos ServerDynamicSource) Validate() error {
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

func (inlineIonos ServerStaticSource) Validate() error {
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
