package tasks

import (
	"strings"

	"github.com/headmade/backuper/backuper"
)

type BackupTaskInterface interface {
	GenerateBackupFile(string) (string, []byte, error)
	tmpFileName() string
}

type backupTask struct {
	config *backuper.TaskConfig
}

func newBackupTask(config *backuper.TaskConfig) *backupTask {
	return &backupTask{config}
}

func (self *backupTask) tmpFileName() string {
	return strings.Join([]string{
		self.config.Type,
	}, "_")
}
