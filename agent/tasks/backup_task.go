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
	config      *backuper.TaskConfig
	tmpFileBase string // expected to be initialized by enclosing struct
}

func newBackupTask(config *backuper.TaskConfig) *backupTask {
	return &backupTask{config: config}
}

func (btask *backupTask) Type() string {
	return btask.config.Type
}

func (btask *backupTask) TmpFileName() string {
	return strings.Join([]string{
		btask.Type(),
		btask.tmpFileBase,
	}, "_") + btask.compressionSuffix()
}

func (btask *backupTask) needCompression() bool {
	return len(btask.config.Compression) > 0
}

func (btask *backupTask) compressionSuffix() (cs string) {
	if btask.needCompression() {
		cs = ".bz2"
	}
	return
}
