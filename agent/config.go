package agent

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Token string
	TmpDir string
}

func ConfigPath() string {
	//return "/etc/backuper/agent.json"
	return "/tmp/agent.json"
}

func LoadConfig(configPath string) (*Config, error) {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func WriteConfig(c *Config, configPath string) error {
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(configPath, data, 0644)
}
