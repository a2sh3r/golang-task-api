package config

import (
	"flag"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	RunAddress string `env:"RUN_ADDRESS" envDefault:"localhost:8080"`
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (cfg *Config) ParseFlags() {
	var runAddress string

	flag.StringVar(&runAddress, "a", "", "address host:port")

	flag.Parse()

	if runAddress != "" {
		cfg.RunAddress = runAddress
	}
}
