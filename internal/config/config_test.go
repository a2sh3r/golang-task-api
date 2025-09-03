package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig_Defaults(t *testing.T) {
	os.Clearenv()
	cfg, err := LoadConfig()
	assert.NoError(t, err)
	assert.Equal(t, "localhost:8080", cfg.RunAddress)
	assert.Equal(t, "/tmp/task-db.json", cfg.FileStoragePath)
	assert.Equal(t, 3, cfg.StoreInterval)
}

func TestLoadConfig_Env(t *testing.T) {
	os.Setenv("RUN_ADDRESS", "127.0.0.1:9999")
	os.Setenv("FILE_STORAGE_PATH", "/tmp/test.json")
	os.Setenv("STORE_INTERVAL", "42")
	defer os.Clearenv()

	cfg, err := LoadConfig()
	assert.NoError(t, err)
	assert.Equal(t, "127.0.0.1:9999", cfg.RunAddress)
	assert.Equal(t, "/tmp/test.json", cfg.FileStoragePath)
	assert.Equal(t, 42, cfg.StoreInterval)
}

func TestConfig_ParseFlags(t *testing.T) {

	os.Args = []string{"cmd", "-a", "0.0.0.0:1234", "-f", "/tmp/flag.json", "-i", "77"}
	cfg := &Config{}
	cfg.ParseFlags()
	assert.Equal(t, "0.0.0.0:1234", cfg.RunAddress)
	assert.Equal(t, "/tmp/flag.json", cfg.FileStoragePath)
	assert.Equal(t, 77, cfg.StoreInterval)
}
