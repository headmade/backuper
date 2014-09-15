package tasks

import (
	"github.com/headmade/backuper/backuper"
)

type backupDirectoryTask struct {
	*backupTask
}

func newBackupDirectoryTask(config *backuper.TaskConfig) BackupTaskInterface {
	return &backupDirectoryTask{newBackupTask(config)}
}

func (self *backupDirectoryTask) GenerateBackupFile(tmpFilePath string) (string, []byte, error) {
	// TODO: validate that dir exists
	return self.config.Params["dir"], []byte{}, nil
}

