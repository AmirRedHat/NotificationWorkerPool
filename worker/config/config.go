package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type WorkerDetail struct {
	Number int
}

type WorkersConfig struct {
	Sms   WorkerDetail `yaml:"sms"`
	Email WorkerDetail `yaml:"email"`
	Push  WorkerDetail `yaml:"push"`
}

type Config struct {
	Workers WorkersConfig `yaml:"workers"`
}

func LoadConfig(config *Config) {
	fileData, err := os.ReadFile("worker/config/config.yaml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(fileData, &config)
	if err != nil {
		panic(err)
	}
}
