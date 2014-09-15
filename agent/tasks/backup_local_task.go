package tasks

import (
	"os"
	"github.com/headmade/backuper/backuper"
)

type backupLocalTask struct {
	*backupTask
}

func newBackupLocalTask(config *backuper.TaskConfig) BackupTaskInterface {
	return &backupLocalTask{newBackupTask(config)}
}

func (self *backupLocalTask) GenerateBackupFile(tmpFilePath string) (string, []byte, error) {
	path := self.config.Params["path"]
	file, err := os.Open(path)
	if err == nil {
		err = file.Close()
	}
	return path, []byte{}, err
}

