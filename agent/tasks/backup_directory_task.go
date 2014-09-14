package tasks

import (
	"path/filepath"
	"fmt"
	"log"
	"os"

	"github.com/headmade/backuper/backuper"
)

type backupDirectoryTask struct {
	*backupTask
}

func newBackupDirectoryTask(config *backuper.TaskConfig) BackupTaskInterface {
	return &backupDirectoryTask{newBackupTask(config)}
}

func (self *backupDirectoryTask) GenerateBackupFile(tmpFilePath string) ([]byte, error) {

	err := os.MkdirAll(tmpFilePath, 0700)
	if err != nil {
		return nil, err
	}

	dir := self.config.Params["dir"]
	parentDir := filepath.Dir(dir)
	baseName := filepath.Base(dir)

	// copy entire tree/file, preserve soft/hard links and other special files
	cmd := fmt.Sprintf(
		"tar -cf - -C %s %s | tar -xf - -C %s",
		parentDir,
		baseName,
		tmpFilePath,
	)
	log.Println(cmd)

	return System(cmd)
}
