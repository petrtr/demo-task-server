package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

type Config struct {
	Env      string `yaml:"env" env-default:"dev"`
	Host     string `yaml:"host" env-default:"localhost"`
	Port     int    `yaml:"port" env-default:"8080"`
	LogLevel string `yaml:"log_level" env-default:"info"`
}

func GetConfig() *Config {
	configPath := getConfigPath()
	if configPath == "" {
		panic("config path is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file not exist: " + configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}

	return &cfg
}

func getConfigPath() string {
	var configPath string

	// --config="./config.yaml"
	flag.StringVar(&configPath, "config", "", "path to config file")
	flag.Parse()

	if configPath == "" {
		return os.Getenv("CONFIG_PATH")
	}

	return configPath
}
