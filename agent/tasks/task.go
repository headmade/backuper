package tasks

import (
	"os/exec"
)

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

func (self *task) System(cmd string) ([]byte, error) {
	return exec.Command("sh", "-c", cmd).CombinedOutput()
}

