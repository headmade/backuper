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
}

type BackupTaskResult struct {
	PathResult
	TaskId string
	Output string
}

func NewPathResult(err error, path string) (res PathResult) {
	res.Path = path
	if err != nil {
		res.Err = err.Error()
	}
	return
}
