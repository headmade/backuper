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
	"github.com/headmade/backuper/config"
	"github.com/nightlyone/lockfile"
)

type Runner struct {
	agentConfig  *backuper.AgentConfig
	secretConfig *config.Providers
	timestamp    string
	lockfile     lockfile.Lockfile
}

func NewRunner(config *backuper.AgentConfig, secretConfig *config.Providers) *Runner {
	return &Runner{
		agentConfig:  config,
		secretConfig: secretConfig,
		timestamp:    time.Now().Format("20060102_1504"),
	}
}

func (runner *Runner) tmpDirPath() string {
	return runner.agentConfig.TmpDir + "/backuper/"
}

func (runner *Runner) tmpFilePath(tmpFileName string) string {
	return filepath.Join(runner.tmpDirPath(), tmpFileName)
}

func (runner *Runner) prepareTmpDir() error {
	log.Println("prepareTmpDir():", runner.tmpDirPath())
	return os.MkdirAll(runner.tmpDirPath(), 0700)
}

func (runner *Runner) pidFilePath() string {
	return runner.tmpFilePath("backuper.pid")
}

func (runner *Runner) lockPidFile(pidFilePath string) (err error) {
	log.Println("lockPidFile():", pidFilePath)
	runner.lockfile, err = lockfile.New(pidFilePath)
	if err == nil {
		err = runner.lockfile.TryLock()
	}
	return
}

func (runner *Runner) unlockPidFile() error {
	return runner.lockfile.Unlock()
}

func (runner *Runner) CleanupTmpDir() error {
	log.Println("cleanupTmpLocal():", runner.tmpDirPath())
	return nil
	return os.RemoveAll(runner.tmpDirPath())
}

func (runner *Runner) backupFileName() string {
	return runner.appendTimestamp("backup")
}

func (runner *Runner) backupFilePath() string {
	return runner.tmpFilePath(runner.backupFileName())
}

func (runner *Runner) appendTimestamp(str string) string {
	return strings.Join([]string{
		str,
		runner.timestamp,
	}, "_")
}

func (runner *Runner) formatDstPath(path string) string {
	hostname, _ := os.Hostname()
	return ReplaceVars(
		path,
		map[string]string{
			"$server_name": hostname,
			"$timestamp":   runner.timestamp,
		},
	)
}

func (runner *Runner) encryptTmpFiles(backupFilePath string, tmpFiles []string) (output []byte, err error) {

	if len(tmpFiles) == 0 {
		return []byte{}, errors.New("No files to encrypt")
	}

	cmd := fmt.Sprintf(
		"tar -cf - -C %s %s | %s > %s",
		runner.tmpDirPath(),
		strings.Join(tmpFiles, " "),
		EncryptCmd("PASS"),
		backupFilePath,
	)
	log.Println(cmd)

	return System(cmd)
}

func (runner *Runner) uploadBackupFile(backupFilePath, bucket, dstPath string) (output []byte, err error) {

	cmd := fmt.Sprintf(
		"gof3r put -p %s -b %s -k %s",
		backupFilePath,
		bucket,
		runner.formatDstPath(dstPath),
	)
	log.Println(cmd)

	return System(cmd)
}

func (runner *Runner) runTasks(configs *[]backuper.TaskConfig) (results []backuper.BackupTaskResult) {

	results = make([]backuper.BackupTaskResult, 0, len(*configs))

	for _, config := range *configs {
		task, err := GetTask(&config)
		if err == nil {
			taskTmpFileName := task.TmpFileName()
			tmpFilePath := runner.tmpFilePath(taskTmpFileName)

			out, err := task.GenerateTmpFile(tmpFilePath)

			results = append(results, backuper.BackupTaskResult{
				PathResult: backuper.NewPathResult(err, taskTmpFileName, string(out)),
			})
		} else {
			log.Printf("task type: %s, no registered handler found", config.Type)
		}
	}
	return
}

func (runner *Runner) Run() (err error, backupResult *backuper.BackupResult) {
	configs := &runner.agentConfig.Tasks

	backupResult = &backuper.BackupResult{}

	err = runner.prepareTmpDir()

	tmpDirPath := runner.tmpDirPath()
	backupResult.Prepare = backuper.NewPathResult(err, tmpDirPath, "")

	if err != nil {
		log.Println("ERR: prepare:", err.Error())
		return
	}

	pidFilePath := runner.pidFilePath()
	err = runner.lockPidFile(pidFilePath)
	backupResult.Lock = backuper.NewPathResult(err, pidFilePath, "")

	if err != nil {
		log.Println("ERR: lock:", err.Error())
		return
	}

	backupResult.Backup = runner.runTasks(configs)

	tmpFiles := make([]string, 0, len(*configs))
	for _, result := range backupResult.Backup {
		if len(result.Err) == 0 {
			tmpFiles = append(tmpFiles, result.Path)
		}
	}

	backupFilePath := runner.backupFilePath()

	output, err := runner.encryptTmpFiles(backupFilePath, tmpFiles)
	backupResult.Encrypt = backuper.NewPathResult(err, backupFilePath, string(output))

	if err != nil {
		log.Println("ERR: encrypt:", err.Error())
		return
	}

	output, err = runner.uploadBackupFile(backupFilePath, "headmade", "backup/$hostname/$timestamp")
	backupResult.Upload = backuper.NewPathResult(err, backupFilePath, string(output))

	if err != nil {
		log.Println("ERR: upload:", err.Error())
		return
	}

	// TODO: remember to uncomment the unlock()!
	//err = runner.unlockPidFile()
	backupResult.Unlock = backuper.NewPathResult(err, pidFilePath, "")

	err = runner.CleanupTmpDir()
	backupResult.Cleanup = backuper.NewPathResult(err, tmpDirPath, "")

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
func ReplaceVars(str string, replacements map[string]string) string {
	for from, to := range replacements {
		str = strings.Replace(str, from, to, -1)
	}
	return str
}
