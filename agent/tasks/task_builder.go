package tasks

import (
  "errors"
)

const (
	DIR_TASK_TYPE = "directory"
	POSTGRES_TASK_TYPE = "postgres"
	MYSQL_TASK_TYPE = "mysql"
)

type TaskBuilderFunc (func(*Config)TaskInterface)

var newTaskBuilders = map[string]TaskBuilderFunc{
  DIR_TASK_TYPE: newBackupDirectoryTask,
  POSTGRES_TASK_TYPE: newBackupPostgresTask,
}

func RegisterBuilder(taskType string, taskBuilder TaskBuilderFunc) {
	newTaskBuilders[taskType] = taskBuilder
}

func Get(config *Config) (task TaskInterface, err error) {
	taskBuilder := newTaskBuilders[config.Type]
	if taskBuilder != nil {
		task = taskBuilder(config)
	} else {
		err = errors.New("Unsupported task type (" + config.Type + ")")
	}

	return task, err
}

