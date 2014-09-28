package tasks

import (
	"fmt"
	"log"
	"path/filepath"
	"os"
	"strings"
	"github.com/headmade/backuper/backuper"
)

const (
	tmpFileTarSuffix = ".tar"
)

type backupLocalTask struct {
	*backupTask
	pathParentDir string
	pathBaseName  string
}

func newBackupLocalTask(config *backuper.TaskConfig) BackupTaskInterface {
	localTask := backupLocalTask{backupTask: newBackupTask(config)}

	path := localTask.sourcePath()

	localTask.pathParentDir = filepath.Dir(path)
	localTask.pathBaseName = filepath.Base(path)

	//tmpFileBaseDir = filepath.Base(localTask.pathParentDir)
	tmpFileBase := strings.Replace(localTask.pathBaseName, ".", "_", -1)

	//localTask.tmpFileBase = strings.Join([]string{
	//	tmpFileBaseDir,
	//	tmpFileBase,
	//}, "_") + tmpFileTarSuffix

	localTask.tmpFileBase = tmpFileBase + tmpFileTarSuffix

	return &localTask
}

func (localTask *backupLocalTask) sourcePath() string {
	return localTask.config.Params["path"]
}

func (localTask *backupLocalTask) GenerateTmpFile(tmpFilePath string) (output []byte, err error) {

	file, err := os.Open(localTask.sourcePath())
	if err == nil {
		err = file.Close()
	}

	if err != nil {
		return
	}

	cmd := fmt.Sprintf(
		"tar -cf - -C %s %s %s >%s",
		localTask.pathParentDir,
		localTask.pathBaseName,
		localTask.compressionFlag(),
		tmpFilePath,
	)
	log.Println(cmd)

	return System(cmd)
}

func (localTask *backupLocalTask) compressionFlag() (cf string) {
	if localTask.needCompression() {
		cf = "--bzip"
	}
  return
}

