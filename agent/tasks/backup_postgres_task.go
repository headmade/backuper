package tasks

import (
	"log"
)

type backupPostgresTask struct {
	*backupTask
}

func newBackupPostgresTask(config *Config) TaskInterface {
	return &backupPostgresTask{newBackupTask(config)}
}

func (self *backupPostgresTask) Run() error {
	log.Println("run backupPostgresTask")
	return nil
}

