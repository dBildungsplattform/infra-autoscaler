package core

import (
	p "scaler/providers"
	s "scaler/shared"
)

func load_provider[P s.Provider](def s.ProviderDefinition) P {
	// TODO load config from file, env, CLI args, etc.
	switch t := def.Type; t {
	case s.Ionos:
		Ionos := p.Ionos{
			ProviderName:  def.Name,
			Username:      "username",
			Password:      "password",
			DatacenterIds: []string{"datacenter-id-1", "datacenter-id-2"},
		}
		return s.Provider(Ionos).(P) // Panics if Ionos.Provider does not implement P
	}
	// TODO how to handle invalid provider name? panic?
	panic("not implemented")
}
