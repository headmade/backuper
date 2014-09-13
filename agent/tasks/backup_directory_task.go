package tasks

import (
	"log"
)

type backupDirectoryTask struct {
	*backupTask
}

func newBackupDirectoryTask(config *Config) TaskInterface {
	return &backupDirectoryTask{newBackupTask(config)}
}

func (self *backupDirectoryTask) Run() error {
	log.Println("run backupDirectoryTask")
	return nil
}

