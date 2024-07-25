package util

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

func LoadYAML(path string, ptr interface{}) error {
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read config file error: %w", err)
	}

	err = yaml.Unmarshal(yamlFile, ptr)
	if err != nil {
		return fmt.Errorf("unmarshal config file error: %w", err)
	}

	return nil
}
