package backuper

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

