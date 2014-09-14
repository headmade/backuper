package tasks

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/headmade/backuper/agent"
)

type Runner struct {
	agentConfig *agent.Config
}

type backupTaskResult struct {
	err error
	file string
	output string
}

type tmpDirResult struct {
	err error
	tmpDir string
}

type backupFileResult struct {
	err error
	backupFile string
}

type BackupResult struct {
	Prepare tmpDirResult
	Backup []backupTaskResult
	Encrypt backupFileResult
	Upload backupFileResult
	Cleanup tmpDirResult
}

func NewRunner(config *agent.Config) *Runner {
	return &Runner{agentConfig: config}
}

func (runner *Runner) tmpDirPath() string {
	return runner.agentConfig.TmpDir + "/backuper/"
}

func (runner *Runner) tmpFilePath(task BackupTaskInterface) string {
	return runner.tmpDirPath() + "/" + task.tmpFileName()
}

func (runner *Runner) prepareTmpDirectory() error {
	log.Println("prepareTmpDirectory():", runner.tmpDirPath())
	return os.MkdirAll(runner.tmpDirPath(), 0700)
}

func (runner *Runner) CleanupTmpDirectory() error {
	log.Println("cleanupTmpDirectory():", runner.tmpDirPath())
	return os.RemoveAll(runner.tmpDirPath())
}


func (self *backupTask) EncryptCmd(pass string) string {
	return fmt.Sprintf(
		"openssl aes-128-cbc -pass pass:%s",
	  pass,
	)
}

func (runner *Runner) Run(configs *[]*Config) (res *BackupResult){

	res = &BackupResult{}

	err := runner.prepareTmpDirectory()

	res.Prepare = tmpDirResult{err, ""}

  res.Backup = make([]backupTaskResult, 0, len(*configs))

	for _, config := range *configs {
		task, err := Get(config)
		if err == nil {
			log.Printf("task type: %s, task object: %#v", config.Type, task)
			out, err := task.GenerateBackupFile()
			res.Backup = append(res.Backup, backupTaskResult{err,runner.tmpFilePath(task),string(out)})
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

