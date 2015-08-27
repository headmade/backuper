package tasks

import (
	"errors"

	"github.com/headmade/backuper/backuper"
)

const (
	LOCAL_TASK_TYPE    = "local"
	POSTGRES_TASK_TYPE = "postgres"
	MYSQL_TASK_TYPE    = "mysql"
)

type BackupTaskBuilderFunc (func(*backuper.TaskConfig) BackupTaskInterface)

var newBackupTaskBuilders = map[string]BackupTaskBuilderFunc{
	LOCAL_TASK_TYPE:    newBackupLocalTask,
	POSTGRES_TASK_TYPE: newBackupPostgresTask,
	MYSQL_TASK_TYPE:    newBackupMySQLTask,
}

func RegisterBuilder(taskType string, taskBuilder BackupTaskBuilderFunc) {
	newBackupTaskBuilders[taskType] = taskBuilder
}

func GetTask(config *backuper.TaskConfig) (task BackupTaskInterface, err error) {
	taskBuilder := newBackupTaskBuilders[config.Type]
	if taskBuilder != nil {
		task = taskBuilder(config)
	} else {
		err = errors.New("Unsupported task type (" + config.Type + ")")
	}

	return task, err
}
