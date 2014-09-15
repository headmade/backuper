package tasks

import (
	"os"
	"github.com/headmade/backuper/backuper"
)

type backupDirectoryTask struct {
	*backupTask
}

func newBackupDirectoryTask(config *backuper.TaskConfig) BackupTaskInterface {
	return &backupDirectoryTask{newBackupTask(config)}
}

func (self *backupDirectoryTask) GenerateBackupFile(tmpFilePath string) (string, []byte, error) {
	path := self.config.Params["path"]
	file, err := os.Open(path)
	if err == nil {
		err = file.Close()
	}
	return path, []byte{}, err
}

