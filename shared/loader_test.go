package shared

import (
	"os"
	"testing"
)

type TestStruct struct {
	Test1 StringFromEnv `yaml:"test1"`
	Test2 StringFromEnv `yaml:"test2"`
	Test3 StringFromEnv `yaml:"test3"`
}

func (t TestStruct) Validate() error {
	return nil
}

func TestLoadEnv(t *testing.T) {
	os.Setenv("TEST_ENV_1", "value1")

	testConfig, err := OpenConfig("test_files/env.yaml")
	if err != nil {
		t.Fatal(err)
	}
	testStruct, err := LoadConfig[TestStruct](testConfig)
	if err != nil {
		t.Errorf("unexpected error unmarshalling the config yaml: %s", err)
	}

	if testStruct.Test1 != "value1" {
		t.Errorf("expected Test1 to be value1, got %s", testStruct.Test1)
	}
	if testStruct.Test2 != "$Test_env_2" {
		t.Errorf("expected Test2 to be $Test_env_2, got %s", testStruct.Test2)
	}
	if testStruct.Test3 != "not_an_env_var" {
		t.Errorf("expected Test3 to be not_an_env_var, got %s", testStruct.Test3)
	}

	defer os.Unsetenv("TEST_ENV_1")
}

func TestLoadUnsetEnv(t *testing.T) {
	testConfig, err := OpenConfig("test_files/env.yaml")
	if err != nil {
		t.Fatal(err)
	}
	_, err = LoadConfig[TestStruct](testConfig)
	if err == nil {
		t.Error("expected error: environment variable TEST_ENV_1 not set")
	}
}
