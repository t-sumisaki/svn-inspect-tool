package main

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

type ConfigFile struct {
	NotifyLockInfo NotifyLockInfoConfigSet `yaml:"notifylockinfo"`
}

func loadConfig(config *ConfigFile) error {

	data, err := os.ReadFile("config.yml")
	if err != nil {
		return fmt.Errorf("readfile error: %v", err)
	}

	if err := yaml.Unmarshal(data, config); err != nil {
		return fmt.Errorf("data parse error: %v", err)
	}

	return nil
}
