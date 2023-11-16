package shared

import (
	"os"

	"gopkg.in/yaml.v3"
)

func OpenConfig(path string) ([]byte, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func LoadConfig[V Validater](data []byte) (*V, error) {
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
