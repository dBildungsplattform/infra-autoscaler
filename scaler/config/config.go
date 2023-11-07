package config

import (
	"os"
	"testing"

	"gopkg.in/yaml.v3"
)

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

// Can't be in this package because of circular dependencies
// type Config struct {
// 	BBB *BBBConfig
// }

func LoadConfig[V Validater](path string) (*V, error) {
	r, ok := os.ReadFile(path)
	if ok != nil {
		return nil, ok
	}
	config, ok := ParseConfig[V](r)
	if ok != nil {
		return nil, ok
	}
	return config, nil
}

func ParseConfig[V Validater](data []byte) (*V, error) {
	var config V
	err := yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	err = config.Validate()
	if err != nil {
		return nil, err
	}
	return &config, nil
}
