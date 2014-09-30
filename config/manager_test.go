package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfigPath(t *testing.T) {
	assert.Equal(t, configPath(), "/tmp/backuper.json")
}

func TestWriteAndLoadConfig(t *testing.T) {
	config := Config{Local: false, Secret: Providers{"AWS": Provider{}}}
	WriteConfig(&config)
	loadConfig, _ := loadConfig()
	assert.Equal(t, config, *loadConfig)
}
