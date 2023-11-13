package shared

import "testing"

type Validater interface {
	Validate() error
}

func ValidatePass(t *testing.T, v Validater) {
	err := v.Validate()
	if err != nil {
		t.Error(err)
	}
}

func ValidateFail(t *testing.T, v Validater) {
	err := v.Validate()
	if err == nil {
		t.Error("expected error")
	}
}
