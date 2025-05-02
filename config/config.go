package config

import (
	"log"
	"os"
	"path"

	"github.com/ilyakaznacheev/cleanenv"
)

var cfg Config

type Config struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Database string `yaml:"database"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

func GetConfig() *Config {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal("Error while getting current working directory")
	}

	configPath := path.Join(dir[:len(dir)-4], "config.yaml") //TODO: переписать подъем из /cmd

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file doesn't exist: %v", configPath)
	}

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatal("Unable to read config file")
	}

	return &cfg
}
