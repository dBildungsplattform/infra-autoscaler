package shared

import (
	"os"

	"gopkg.in/yaml.v3"
)

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
