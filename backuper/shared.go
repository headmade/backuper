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
	Err    error
	Path string
}

type BackupTaskResult struct {
	PathResult
	Output string
}

