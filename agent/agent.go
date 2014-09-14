package agent

import (
	"log"

	"github.com/headmade/backuper/backuper"
)

type Agent struct {
	Config *Config
}

func Get(config *Config) (*Agent, error) {
	config.TmpDir = "/tmp"
	return &Agent{config}, nil
}

func (agent *Agent) Backup() *backuper.BackupResult {
	log.Println("agent.Backup")
	return &backuper.BackupResult{}
}
