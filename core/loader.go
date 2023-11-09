package core

import (
	c "scaler/common"
	"scaler/providers/Ionos"
)

func load_provider[P c.Provider](def c.ProviderDefinition) P {
	// TODO load config from file, env, CLI args, etc.
	switch t := def.Type; t {
	case c.Ionos:
		Ionos := Ionos.Provider{
			ProviderName:  def.Name,
			Username:      "username",
			Password:      "password",
			DatacenterIds: []string{"datacenter-id-1", "datacenter-id-2"},
		}
		return c.Provider(Ionos).(P) // Panics if Ionos.Provider does not implement P
	}
	// TODO how to handle invalid provider name? panic?
	panic("not implemented")
}
