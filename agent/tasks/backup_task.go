package tasks

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/headmade/backuper/backuper"
)

type BackupTaskInterface interface {
	GenerateBackupFile() ([]byte, error)
	tmpFileName() string
}

type backupTask struct {
	config            *backuper.TaskConfig
	tmpFilePathCached string
}

func newBackupTask(config *backuper.TaskConfig) *backupTask {
	return &backupTask{config, ""}
}

func (self *backupTask) tmpDirPath() string {
	return self.config.Params["tmp_path"]
}

func (self *backupTask) tmpFileName() string {
	return strings.Join([]string{
		self.config.Type,
		self.config.Id,
		time.Now().Format("20060102_1504"),
	}, "_")
}

func (self *backupTask) tmpFilePath() string {
	if len(self.tmpFilePathCached) == 0 {
		self.tmpFilePathCached = filepath.Join(self.tmpDirPath(), self.tmpFileName())
	}
	return self.tmpFilePathCached
}
