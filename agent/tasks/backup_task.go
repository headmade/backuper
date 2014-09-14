package tasks

import (
	"strings"
	"time"

	"github.com/headmade/backuper/backuper"
)

type BackupTaskInterface interface {
	GenerateBackupFile(string) ([]byte, error)
	tmpFileName() string
}

type backupTask struct {
	config            *backuper.TaskConfig
	tmpFilePathCached string
}

func newBackupTask(config *backuper.TaskConfig) *backupTask {
	return &backupTask{config, ""}
}

func (self *backupTask) tmpFileName() string {
	return strings.Join([]string{
		self.config.Type,
		self.config.Id,
		time.Now().Format("20060102_1504"),
	}, "_")
}

