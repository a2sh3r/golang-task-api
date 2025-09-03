package config

import (
	"flag"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	RunAddress      string `env:"RUN_ADDRESS" envDefault:"localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:"/tmp/task-db.json"`
	StoreInterval   int    `env:"STORE_INTERVAL" envDefault:"3"`
	DatabaseURI     string `env:"DATABASE_URI" envDefault:"postgres://postgres:postgres@localhost:5432/golangtaskapi?sslmode=disable"`
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (cfg *Config) ParseFlags() {
	var (
		runAddress    string
		filePath      string
		storeInterval int
		databaseURI   string
	)

	flag.StringVar(&runAddress, "a", "", "address host:port")
	flag.StringVar(&filePath, "f", "", "file storage path for saving/loading inmemory data")
	flag.IntVar(&storeInterval, "i", 0, "file storage time interval")
	flag.StringVar(&databaseURI, "d", "", "file storage time interval")

	flag.Parse()

	if runAddress != "" {
		cfg.RunAddress = runAddress
	}

	if filePath != "" {
		cfg.FileStoragePath = filePath
	}

	if storeInterval != 0 {
		cfg.StoreInterval = storeInterval
	}
	if databaseURI != "" {
		cfg.DatabaseURI = databaseURI
	}
}
