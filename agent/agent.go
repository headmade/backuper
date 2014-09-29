package agent

import (
	"log"

	"github.com/headmade/backuper/agent/tasks"
	"github.com/headmade/backuper/backuper"
	"github.com/headmade/backuper/config"
)

type Agent struct {
	Config *backuper.AgentConfig
	Secret *config.Providers
}

func Get(config *backuper.AgentConfig, secretConfig *config.Providers) (*Agent, error) {
	return &Agent{config, secretConfig}, nil
}

func (agent *Agent) Backup() (error, *backuper.BackupResult) {
	log.Println("agent.Backup")
	runner := tasks.NewRunner(agent.Config, agent.Secret)
	return runner.Run()
}
