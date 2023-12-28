package shared

import "fmt"

/*** Provider definition ***/
type Provider interface {
	Validate() error
	GetScaledObjects() ([]ScaledObject, error)
	UpdateScaledObject(scaledObject ScaledObject, targetRes ScaleResource) error
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
