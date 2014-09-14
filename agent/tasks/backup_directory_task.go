package tasks

import (
	"fmt"
	"log"
)

type backupDirectoryTask struct {
	*backupTask
}

func newBackupDirectoryTask(config *Config) BackupTaskInterface {
	return &backupDirectoryTask{newBackupTask(config)}
}

func (self *backupDirectoryTask) GenerateBackupFile() ([]byte, error) {
	cmd := fmt.Sprintf("tar --bzip -cf - -C %s . | %s >%s",
		self.config.Params["dir"],
		self.EncryptCmd(self.config.Params["pass"]),
		self.tmpFilePath(),
	)
	log.Println(cmd)

	return System(cmd)
}

