package config

import (
	"encoding/json"
	"io/ioutil"
	"reflect"

	"github.com/headmade/backuper/backuper"
)

type Config struct {
	Client *backuper.ClientConfig
	Agent  *backuper.AgentConfig
}

func New() (*Config, error) {
	config, err := loadConfig()
	if err != nil {
	}
	return config, err
}

func configPath() string {
	//return "/etc/backuper/agent.json"
	return "/tmp/backuper.json"
}

func (c *Config) Write(value interface{}) error {

	switch reflect.ValueOf(value).Type() {
	case reflect.ValueOf(&backuper.AgentConfig{}).Type():
		c.Agent = value.(*backuper.AgentConfig)
	case reflect.TypeOf(&backuper.ClientConfig{}):
		c.Client = value.(*backuper.ClientConfig)
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
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(configPath(), data, 0644)
}
