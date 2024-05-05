package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type ApplicationConfigManager struct {
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`
	PostgreSql struct {
		Url string `yaml:"url"`
	} `yaml:"postgresql"`
	Redis struct {
		Host      string `yaml:"host"`
		Port      int    `yaml:"port"`
		EventName string `yaml:"eventName"`
	} `yaml:"redis"`
}

func LoadConfig(filename string) (*ApplicationConfigManager, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var config ApplicationConfigManager
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config data: %v", err)
	}

	return &config, nil
}

func NewApplicationConfigManager() ApplicationConfigManager {
	var applicationConfigManager, _ = LoadConfig("resource/application_config.yaml")
	return *applicationConfigManager
}
