package common

func Load_provider[Name InfrastructureProviderName, Provider InfrastructureProvider](name Name) Provider {
	// TODO load config from file, env, CLI args, etc.
	switch name {
	case Ionos:
		Ionos := IonosProvider{
			ProviderName:  Ionos,
			Username:      "username",
			Password:      "password",
			DatacenterIds: []string{"datacenter-id-1", "datacenter-id-2"},
		}
		return Provider(Ionos)
	}
	// TODO how to handle invalid provider name? panic?
	panic("not implemented")
}

//type ServiceConfig[p InfrastructureProvider] struct {
//	Name     string
//	Type     InfrastructureType
//	Provider p
//}
