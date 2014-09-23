package backuper

type AgentConfig struct {
	Token  string       `json:"token"`
	TmpDir string       `json:"tmp_dir"`
	Tasks  []TaskConfig `json:"tasks"`
}

type TaskConfig struct {
	Id     string            `json:"id"`
	Type   string            `json:"type"`
	Name   string            `json:"name"`
	Params map[string]string `json:"params"`
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
	Err    string
	Path   string
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
