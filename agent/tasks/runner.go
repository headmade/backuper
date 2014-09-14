package tasks

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/headmade/backuper/backuper"
)

type Runner struct {
	agentConfig *backuper.AgentConfig
	timestamp   string
}

func NewRunner(config *backuper.AgentConfig) *Runner {
	return &Runner{
		agentConfig: config,
		timestamp:   time.Now().Format("20060102_1504"),
	}
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

func (runner *Runner) backupFileName() string {
	return runner.appendTimestamp("backup")
}

func (runner *Runner) appendTimestamp(str string) string {
	return strings.Join([]string{
		str,
		runner.timestamp,
	}, "_")
}

func (runner *Runner) encryptTmpFiles(fileNames []string) (backupFilePath string, err error) {

	if len(fileNames) == 0 {
		return "", errors.New("No files to encrypt")
	}

	backupFilePath = runner.tmpFilePath(runner.backupFileName())
	cmd := fmt.Sprintf(
		"tar -cf - -C %s %s | %s >%s",
		runner.tmpDirPath(),
		strings.Join(fileNames, " "),
		EncryptCmd("PASS"),
		backupFilePath,
	)
	log.Println(cmd)

	_, err = System(cmd)
	return backupFilePath, err
}

func (runner *Runner) uploadBackupFile(nackupFileName string) error {
	return nil
}

func (runner *Runner) Run() (backupResult *backuper.BackupResult) {
	configs := &runner.agentConfig.Tasks

	backupResult = &backuper.BackupResult{}

	err := runner.prepareTmpDirectory()

	tmpDirPath := runner.tmpDirPath()
	backupResult.Prepare = backuper.TmpDirResult{err, tmpDirPath}

	len_configs := len(*configs)
	backupResult.Backup = make([]backuper.BackupTaskResult, 0, len_configs)
	tmpFiles := make([]string, 0, len_configs)

	for _, config := range *configs {
		task, err := GetTask(&config)
		if err == nil {
			tmpFileName := runner.appendTimestamp(task.tmpFileName())
			tmpFilePath := runner.tmpFilePath(tmpFileName)
			log.Printf("task type: %s, task object: %#v", config.Type, task)
			out, err := task.GenerateBackupFile(tmpFilePath)
			if err == nil {
				tmpFiles = append(tmpFiles, tmpFileName)
			}
			backupResult.Backup = append(backupResult.Backup, backuper.BackupTaskResult{err, tmpFilePath, string(out)})
		} else {
			log.Printf("task type: %s, no registered handler found", config.Type)
		}
	}

	backupFileName, err := runner.encryptTmpFiles(tmpFiles)
	backupResult.Encrypt = backuper.BackupFileResult{err, backupFileName}

	err = runner.uploadBackupFile(backupFileName)
	backupResult.Upload = backuper.BackupFileResult{err, backupFileName}

	err = runner.CleanupTmpDirectory()
	backupResult.Cleanup = backuper.TmpDirResult{err, tmpDirPath}
	return
}

// TODO: move to some utils
func System(cmd string) ([]byte, error) {
	return exec.Command("sh", "-c", cmd).CombinedOutput()
}

func EncryptCmd(pass string) string {
	return fmt.Sprintf(
		"openssl aes-128-cbc -pass pass:%s",
		pass,
	)
}
