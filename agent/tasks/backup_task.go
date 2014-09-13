package tasks

type backupTask struct {
	*task
}

func newBackupTask(config *Config) *backupTask {
	return &backupTask{newTask(config)}
}

