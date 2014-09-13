package tasks

import (
	"fmt"
	"log"
)

type backupDirectoryTask struct {
	*backupTask
}

func newBackupDirectoryTask(config *Config) TaskInterface {
	return &backupDirectoryTask{newBackupTask(config)}
}

func (self *backupDirectoryTask) Run() error {
	self.PrepareTmpDirectory()
	log.Println("run backupDirectoryTask")

	cmd := fmt.Sprintf("tar --bzip -cf - -C %s . | %s >%s",
		self.config.Params["dir"],
		self.EncryptCmd(self.config.Params["pass"]),
		self.tmpFilePath(),
	)
	log.Println(cmd)

	out, err := self.System(cmd)
	log.Println(string(out))

	self.CleanupTmpDirectory()
	return err
}

