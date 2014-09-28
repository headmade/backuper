package tasks

import (
	"strings"

	"github.com/headmade/backuper/backuper"
)

type BackupTaskInterface interface {
  Type() string
	GenerateTmpFile(tmpFilePath string) (output []byte, err error)
	TmpFileName() string
}

type backupTask struct {
	config *backuper.TaskConfig
	tmpFileBase string  // expected to be initialized by enclosing struct
}


func newBackupTask(config *backuper.TaskConfig) *backupTask {
	return &backupTask{config: config}
}

func (self *backupTask) Type() string {
	return self.config.Type
}

func (self *backupTask) TmpFileName() string {
	return strings.Join([]string{
		self.Type(),
		self.tmpFileBase,
	}, "_") + self.compressionSuffix()
}

func (self *backupTask) needCompression() bool {
	return len(self.config.Params["compression"]) > 0
}

func (self *backupTask) compressionSuffix() (cs string) {
	if self.needCompression() {
		cs = ".bz2"
	}
	return
}

