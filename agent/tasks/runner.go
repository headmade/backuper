package tasks

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/headmade/backuper/backuper"
)

type Runner struct {
	agentConfig *backuper.AgentConfig
}

func NewRunner(config *backuper.AgentConfig) *Runner {
	return &Runner{agentConfig: config}
}

func (runner *Runner) tmpDirPath() string {
	return runner.agentConfig.TmpDir + "/backuper/"
}

func (runner *Runner) tmpFilePath(tmpFileName string) string {
	return filepath.Join(runner.tmpDirPath(), tmpFileName)
}

func (runner *Runner) prepareTmpDirectory() error {
	log.Println("prepareTmpDirectory():", runner.tmpDirPath())
	return os.MkdirAll(runner.tmpDirPath(), 0700)
}

func (runner *Runner) CleanupTmpDirectory() error {
	log.Println("cleanupTmpDirectory():", runner.tmpDirPath())
	return nil
	return os.RemoveAll(runner.tmpDirPath())
}

func (self *backupTask) EncryptCmd(pass string) string {
	return fmt.Sprintf(
		"openssl aes-128-cbc -pass pass:%s",
		pass,
	)
}

func (runner *Runner) Run() (res *backuper.BackupResult) {
	configs := &runner.agentConfig.Tasks

	res = &backuper.BackupResult{}

	err := runner.prepareTmpDirectory()

	res.Prepare = backuper.TmpDirResult{err, ""}

	res.Backup = make([]backuper.BackupTaskResult, 0, len(*configs))

	for _, config := range *configs {
		task, err := Get(&config)
		if err == nil {
			tmpFilePath := runner.tmpFilePath(task.tmpFileName())
			log.Printf("task type: %s, task object: %#v", config.Type, task)
			out, err := task.GenerateBackupFile(tmpFilePath)
			res.Backup = append(res.Backup, backuper.BackupTaskResult{err, tmpFilePath, string(out)})
		} else {
			log.Printf("task type: %s, no registered handler found", config.Type)
		}
	}
	//backupFileName, err := encryptTmpFiles()
	//res.encrypt
	//uploadBackup(backupFileName)
	runner.CleanupTmpDirectory()
	return
}

// TODO: move to some utils
func System(cmd string) ([]byte, error) {
	return exec.Command("sh", "-c", cmd).CombinedOutput()
}
