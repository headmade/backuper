package tasks

import (
	"fmt"
	"log"

	"github.com/headmade/backuper/backuper"
)

type backupDirectoryTask struct {
	*backupTask
}

func newBackupDirectoryTask(config *backuper.TaskConfig) BackupTaskInterface {
	return &backupDirectoryTask{newBackupTask(config)}
}

func (self *backupDirectoryTask) GenerateBackupFile(tmpFilePath string) ([]byte, error) {
	cmd := fmt.Sprintf(
		"tar --bzip -cf - -C %s . >%s",
		self.config.Params["dir"],
		tmpFilePath,
	)
	log.Println(cmd)

	return System(cmd)
}
