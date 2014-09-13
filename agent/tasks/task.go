package tasks

type TaskInterface interface {
//	Init(config *Config)
	Run() error
}

type task struct {
	config *Config
}

func newTask(config *Config) *task {
	return &task{config}
}

