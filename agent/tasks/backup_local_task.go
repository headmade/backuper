package tasks

import (
	// "fmt"
	"github.com/headmade/backuper/backuper"
	"github.com/headmade/backuper/hmutil"
	"os"
	"path/filepath"
	"strings"
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

func (localTask *backupLocalTask) excludePath() []string {
	filenames := strings.Split(localTask.config.Params["exclude_path"], ", ")

	for i, e := range filenames {
		filenames[i] = filepath.Join(localTask.sourcePath(), e)
	}

	return filenames
}

func (localTask *backupLocalTask) GenerateTmpFile(tmpFilePath string) (output []byte, err error) {

	file, err := os.Open(localTask.sourcePath())
	if err == nil {
		err = file.Close()
	}

	if err != nil {
		return
	}

	// cmd := fmt.Sprintf(
	// 	"tar -cf - %s -C %s %s >%s",
	// 	localTask.compressionFlag(),
	// 	localTask.pathParentDir,
	// 	localTask.pathBaseName,
	// 	tmpFilePath,
	// )

	// return hmutil.System(cmd)

	err = hmutil.PackAndCompress(
		localTask.pathParentDir,
		[]string{localTask.pathBaseName},
		localTask.excludePath(),
		tmpFilePath,
		[]byte{},
		false,
	)
	return []byte{}, err
}

func (localTask *backupLocalTask) compressionFlag() (cf string) {
	if localTask.needCompression() {
		cf = "--bzip"
	}
	return
}
