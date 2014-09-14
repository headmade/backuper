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
	Prepare TmpDirResult
	Backup  []BackupTaskResult
	Encrypt BackupFileResult
	Upload  BackupFileResult
	Cleanup TmpDirResult
}

type BackupTaskResult struct {
	Err    error
	File   string
	Output string
}

type TmpDirResult struct {
	Err    error
	TmpDir string
}

type BackupFileResult struct {
	Err        error
	BackupFile string
}
