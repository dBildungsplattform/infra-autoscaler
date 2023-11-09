package config

import (
	"os"
	v "scaler/validater"

	"gopkg.in/yaml.v3"
)

// Can't be in this package because of circular dependencies
// type Config struct {
// 	BBB *BBBConfig
// }

func LoadConfig[V v.Validater](path string) (*V, error) {
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

func ParseConfig[V v.Validater](data []byte) (*V, error) {
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
