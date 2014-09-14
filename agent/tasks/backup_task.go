package tasks

import (
	"path/filepath"
	"strings"
	"time"
)

type BackupTaskInterface interface {
	GenerateBackupFile() ([]byte, error)
	tmpFileName() string
}

type backupTask struct {
	config            *Config
	tmpFilePathCached string
}

func newBackupTask(config *Config) *backupTask {
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
