package agent

import (
	"encoding/json"
	"io/ioutil"

	"github.com/headmade/backuper/backuper"
)

func ConfigPath() string {
	//return "/etc/backuper/agent.json"
	return "/tmp/agent.json"
}

func LoadConfig(configPath string) (*backuper.AgentConfig, error) {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config backuper.AgentConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func WriteConfig(c *backuper.AgentConfig, configPath string) error {
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(configPath, data, 0644)
}
