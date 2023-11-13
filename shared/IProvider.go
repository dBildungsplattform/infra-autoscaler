package shared

import "fmt"

/*** Provider definition ***/
type Provider interface {
	Get_login_id() string
	Get_login_secret() string
	Get_type() ProviderType
	Get_name() string
	Validate() error
}

type ProviderDefinition struct {
	Name string
	Type ProviderType
}

type ProviderType string

const (
	Ionos = "Ionos"
)

func (p ProviderType) Validate() error {
	switch p {
	case Ionos:
		return nil
	default:
		return fmt.Errorf("unknown provider type: %s", p)
	}
}
