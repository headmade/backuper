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
			"$hostname":  hostname,
			"$timestamp": runner.timestamp,
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

	awsProvider := (*runner.secretConfig)["AWS"]
	envPath := os.Getenv("PATH")
	envGopath := os.Getenv("GOPATH")

	cmd := fmt.Sprintf(
		"AWS_ACCESS_KEY_ID=%s AWS_SECRET_ACCESS_KEY=%s PATH=%s:%s/bin gof3r put -p %s -b %s -k %s",
		awsProvider["AWS_ACCESS_KEY_ID"],
		awsProvider["AWS_SECRET_ACCESS_KEY"],
		envPath,
		envGopath,
		backupFilePath,
		bucket,
		runner.formatDstPath(dstPath),
	)
	log.Println(cmd)

	return System(cmd)
}

func (runner *Runner) runTasks(configs *[]backuper.TaskConfig) (results []backuper.PathResult) {

	results = make([]backuper.PathResult, 0, len(*configs))

	for _, config := range *configs {
		beginTime := time.Now()
		task, err := GetTask(&config)
		if err == nil {
			taskTmpFileName := task.TmpFileName()
			tmpFilePath := runner.tmpFilePath(taskTmpFileName)

			out, err := task.GenerateTmpFile(tmpFilePath)

			results = append(results,
				backuper.NewPathResult(
					err,
					taskTmpFileName,
					string(out),
					beginTime,
					time.Now(),
			))
		} else {
			log.Printf("task type: %s, no registered handler found", config.Type)
		}
	}
	return
}

func (runner *Runner) Run() (err error, backupResult *backuper.BackupResult) {
	configs := &runner.agentConfig.Tasks

	backupResult = &backuper.BackupResult{
		BeginTime: time.Now(),
		Size:      -1,
		Status:    backuper.BackupErrorCrytical,
	}

	beginTime := time.Now()
	err = runner.prepareTmpDir()

	tmpDirPath := runner.tmpDirPath()
	backupResult.Prepare = backuper.NewPathResult(
		err,
		tmpDirPath,
		"",
		beginTime,
		time.Now(),
	)

	if err != nil {
		log.Println("ERR: prepare:", err.Error())
		return
	}

	defer func() {
		beginTime = time.Now()
		err = runner.CleanupTmpDir()
		backupResult.Cleanup = backuper.NewPathResult(
			err,
			tmpDirPath,
			"",
			beginTime,
			time.Now(),
		)

		cryticalErrors := []*backuper.PathResult{
			&backupResult.Prepare,
			&backupResult.Lock,
			&backupResult.Encrypt,
			&backupResult.Upload,
		}

		backupResult.Status = backuper.BackupErrorNo

		hasCryticalError := false
		for _, br := range cryticalErrors {
			if !ResultSuccess(br) {
				hasCryticalError = true
				backupResult.Status = backuper.BackupErrorCrytical
				break
			}
		}

		if !hasCryticalError {
			numTaskErrors := 0
			for _, br := range backupResult.Backup {
				if !ResultSuccess(&br) {
					numTaskErrors++
					backupResult.Status = backuper.BackupErrorTask
				}
			}
			if numTaskErrors == len(backupResult.Backup) {
				backupResult.Status = backuper.BackupErrorTaskAll
			} else {
				if !ResultSuccess(&backupResult.Unlock) || !ResultSuccess(&backupResult.Cleanup) {
					backupResult.Status = backuper.BackupErrorCleanup
				}
			}
		}

		backupResult.EndTime = time.Now()

		log.Printf("%#v", backupResult)
	}()

	beginTime = time.Now()
	pidFilePath := runner.pidFilePath()
	err = runner.lockPidFile(pidFilePath)
	backupResult.Lock = backuper.NewPathResult(
		err,
		pidFilePath,
		"",
		beginTime,
		time.Now(),
	)

	if err != nil {
		log.Println("ERR: lock:", err.Error())
		return
	}

	backupResult.Backup = runner.runTasks(configs)

	tmpFiles := make([]string, 0, len(*configs))
	for _, result := range backupResult.Backup {
		if ResultSuccess(&result) {
			tmpFiles = append(tmpFiles, result.Path)
		}
	}

	beginTime = time.Now()

	backupFilePath := runner.backupFilePath()

	output, err := runner.encryptTmpFiles(backupFilePath, tmpFiles)
	backupResult.Encrypt = backuper.NewPathResult(
		err,
		backupFilePath,
		string(output),
		beginTime,
		time.Now(),
	)

	if err != nil {
		log.Println("ERR: encrypt:", err.Error())
		return
	}

	fi, err := os.Stat(backupFilePath)
	if err == nil {
		backupResult.Size = fi.Size()
	}

	beginTime = time.Now()
	output, err = runner.uploadBackupFile(backupFilePath, "headmade", "backup/$hostname/$timestamp")
	backupResult.Upload = backuper.NewPathResult(
		err,
		backupFilePath,
		string(output),
		beginTime,
		time.Now(),
	)

	if err != nil {
		log.Println("ERR: upload:", err.Error())
		return
	}

	beginTime = time.Now()
	// TODO: remember to uncomment the unlock()!
	//err = runner.unlockPidFile()
	backupResult.Unlock = backuper.NewPathResult(
		err,
		pidFilePath,
		"",
		beginTime,
		time.Now(),
	)

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

func ResultSuccess(pathResult *backuper.PathResult) bool {
	return pathResult.Err == nil
}

func errString(err error) (s *string) {
  if err != nil {
		tmp := err.Error()
		s = &tmp
	}
	return
}
