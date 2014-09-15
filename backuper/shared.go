package backuper

type AgentConfig struct {
	Token  string
	TmpDir string
	Tasks  []TaskConfig
}

type TaskConfig struct {
	Id     string
	Type   string
	Name   string
	Params map[string]string
}

type BackupResult struct {
	Prepare PathResult
	Lock    PathResult
	Backup  []BackupTaskResult
	Encrypt PathResult
	Upload  PathResult
	Unlock  PathResult
	Cleanup PathResult
}

type PathResult struct {
	Err  string
	Path string
	Output string
}

type BackupTaskResult struct {
	PathResult
	TaskId string
}

func NewPathResult(err error, path, output string) (result PathResult) {
	result.Path = path
	result.Output = output
	if err != nil {
		result.Err = err.Error()
	}
	return
}
