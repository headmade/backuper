package tasks

import (
	"errors"
)

const (
	DIR_TASK_TYPE      = "directory"
	POSTGRES_TASK_TYPE = "postgres"
	MYSQL_TASK_TYPE    = "mysql"
)

type BackupTaskBuilderFunc (func(*Config) BackupTaskInterface)

var newBackupTaskBuilders = map[string]BackupTaskBuilderFunc{
	DIR_TASK_TYPE:      newBackupDirectoryTask,
	POSTGRES_TASK_TYPE: newBackupPostgresTask,
}

func RegisterBuilder(taskType string, taskBuilder BackupTaskBuilderFunc) {
	newBackupTaskBuilders[taskType] = taskBuilder
}

func Get(config *Config) (task BackupTaskInterface, err error) {
	taskBuilder := newBackupTaskBuilders[config.Type]
	if taskBuilder != nil {
		task = taskBuilder(config)
	} else {
		err = errors.New("Unsupported task type (" + config.Type + ")")
	}

	return task, err
}
