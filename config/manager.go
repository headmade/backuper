package config

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
	// "runtime"
	// "path/filepath"
	"github.com/headmade/backuper/backuper"
)

type Config struct {
	Client *backuper.ClientConfig
	Agent  *backuper.AgentConfig
	Local  bool
	Secret Providers
}
type Providers map[string]Provider
type Provider map[string]string

func New() (*Config, error) {
	config, err := loadConfig()
	return config, err
}

// func configPath() string {
// 	if runtime.GOOS == "windows" { 
// 		return filepath.Join("C:", "gobackuper", "agent.json")
// 	} else { 
// 		return filepath.Join("/", "etc", "gobackuper", "agent.json")
// 	}
// 	// return "/tmp/backuper.json"
// }

func (c *Config) Write(value interface{}) error {

	switch reflect.ValueOf(value).Type() {
	case reflect.TypeOf(&backuper.AgentConfig{}):
		c.Agent = value.(*backuper.AgentConfig)
	case reflect.TypeOf(&backuper.ClientConfig{}):
		c.Client = value.(*backuper.ClientConfig)
	case reflect.TypeOf(Providers{}):
		c.Secret = value.(Providers)
	}
	err := WriteConfig(c)
	return err
}

func loadConfig() (*Config, error) {
	data, err := ioutil.ReadFile(configPath())
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func WriteConfig(c *Config) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(configPath(), data, 0644)
}
