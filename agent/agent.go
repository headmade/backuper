package agent

type Agent struct {
	Config Config
}

func Get(config *Config) (*Agent, error) {
	return &Agent{}, nil
}

func (agent *Agent) Backup() error {
	return nil
}
