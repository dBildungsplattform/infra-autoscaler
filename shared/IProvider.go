package shared

import "fmt"

// Interface to get the scaled objects and update them
type Provider interface {
	Validate() error
	GetScaledObjects() ([]ScaledObject, error)
	UpdateScaledObject(scaledObject ScaledObject, targetRes ResourceScalingProposal) error
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
