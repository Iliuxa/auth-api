package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env      string        `yaml:"env" env-default:"local"`
	Database string        `env:"storage_path" env-default:"host=localhost user=postgres password=postgres dbname=postgres sslmode=disable"`
	GRPC     GRPCConfig    `yaml:"grpc"`
	TokenTTL time.Duration `yaml:"token_ttl" env-default:"1h"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config path does not exist")
	}

	var config Config

	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		panic(err.Error())
	}

	return &config
}

func fetchConfigPath() string {
	var res string
	flag.StringVar(&res, "config", "config/config.yaml", "config path")
	flag.Parse()
	return res
}
