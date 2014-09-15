package agent

import (
	"log"

	"github.com/headmade/backuper/agent/tasks"
	"github.com/headmade/backuper/backuper"
)

type Agent struct {
	Config *backuper.AgentConfig
}

func Get(config *backuper.AgentConfig) (*Agent, error) {
	return &Agent{config}, nil
}

func (agent *Agent) Backup() (error, *backuper.BackupResult) {
	log.Println("agent.Backup")
	runner := tasks.NewRunner(agent.Config)
	return runner.Run()
}
