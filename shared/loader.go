package shared

import (
	"fmt"
	"os"
	"regexp"
	"strconv"

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

// StringFromEnv is a string that can be loaded from an environment variable
// It implements the yaml.Unmarshaler interface
type StringFromEnv string

// Regex to match environment variables
// Variable must start with "$" and contain only uppercase letters, numbers, and underscores
var envRegex = regexp.MustCompile(`^\$([A-Z_1-9]+)$`)

func (s *StringFromEnv) UnmarshalYAML(node *yaml.Node) error {
	var nodeVal string
	err := node.Decode(&nodeVal)
	if err != nil {
		return err
	}
	env := envRegex.FindString(nodeVal)
	if env == "" {
		*s = StringFromEnv(nodeVal) // not an environment variable
	} else {
		envVal, bool := os.LookupEnv(env[1:]) // remove the "$" from the environment variable
		if !bool {
			return fmt.Errorf("environment variable %s not set", env)
		}
		*s = StringFromEnv(envVal)
	}
	return nil
}

type IntFromEnv int

// Could be made generic as TypeFromEnv[T any]
func (i *IntFromEnv) UnmarshalYAML(node *yaml.Node) error {
	var nodeVal string
	err := node.Decode(&nodeVal)
	if err != nil {
		return err
	}
	env := envRegex.FindString(nodeVal)
	if env == "" {
		*i = IntFromEnv(0) // not an environment variable
	} else {
		envVal, bool := os.LookupEnv(env[1:]) // remove the "$" from the environment variable
		if !bool {
			return fmt.Errorf("environment variable %s not set", env)
		}
		intRepr, err := strconv.Atoi(envVal)
		if err != nil {
			return fmt.Errorf("environment variable %s is not an integer", env)
		}
		*i = IntFromEnv(intRepr)
	}
	return nil
}
