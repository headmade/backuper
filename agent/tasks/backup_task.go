package tasks

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type BackupTaskInterface interface {
	TaskInterface
	PrepareTmpDirectory() error
//	GenerateBackupFile() error
	CleanupTmpDirectory() error
}

type backupTask struct {
	*task
	tmpFilePathCached string
}

func newBackupTask(config *Config) *backupTask {
	return &backupTask{newTask(config),""}
}

func (self *backupTask) tmpDirPath() string {
	return self.config.Params["tmp_path"]
}

func (self *backupTask) tmpFileName() string {

	return strings.Join([]string{
		self.config.Type,
		self.config.Id,
		time.Now().Format("20060102_1504"),
	},"_")
}

func (self *backupTask) tmpFilePath() string {
	if len(self.tmpFilePathCached) == 0 {
    self.tmpFilePathCached = filepath.Join(self.tmpDirPath(), self.tmpFileName())
	}
	return self.tmpFilePathCached
}

func (self *backupTask) PrepareTmpDirectory() error {
	log.Println("PrepareTmpDirectory():", self.tmpDirPath())
	return os.MkdirAll(self.tmpDirPath(), 0700)
}

func (self *backupTask) CleanupTmpDirectory() error {
	log.Println("CleanupTmpDirectory():", self.tmpFilePath())
	err := os.Remove(self.tmpFilePath())
	if err != nil {
		log.Println("ERR:", err.Error())
	}
	return err
}

func (self *backupTask) EncryptCmd(pass string) string {
	return fmt.Sprintf(
		"openssl aes-128-cbc -pass pass:%s",
	  pass,
	)
}

