package backuper

import (
	"time"
)

const (
	BackupErrorCrytical = "error_crytical"
	BackupErrorCleanup  = "error_cleanup"
	BackupErrorTask     = "error_task"
	BackupErrorTaskAll  = "error_task_all"
	BackupErrorNo       = "success"
)

type AgentConfig struct {
	StartNow    bool         `json:"start_now"`
	Destination Destination  `json:"destination"`
	TmpDir      string       `json:"tmp_dir"`
	Tasks       []TaskConfig `json:"tasks"`
	Period      Period       `json:"period"`
}

type ClientConfig struct {
	Token string
}

type Period struct {
	Type       string   `json:"type"`
	Time       string   `json:"time"`
	DaysOfWeek []string `json:"days_of_week"`
}

type Destination struct {
	Type   string            `json:"type"`
	Params map[string]string `json:"params"`
}

type TaskConfig struct {
	Type        string            `json:"type"`
	Name        string            `json:"name"`
	Compression string            `json:"compression"`
	Params      map[string]string `json:"params"`
}

type BackupResult struct {
	Prepare   PathResult			`json:"prepare"`
	Lock      PathResult            `json:"lock"`
	Backup    []BackupTaskResult	`json:"backup"`
	Encrypt   PathResult            `json:"encrypt"`
	Upload    PathResult			`json:"upload"`
	Unlock    PathResult            `json:"unlock"`
	Cleanup   PathResult            `json:"cleanup"`
	BeginTime time.Time             `json:"begin_time"`
	EndTime   time.Time             `json:"end_time"`
	Size      int64                 `json:"size"`
	Status    string                `json:"status"`
}

type PathResult struct {
	Err       string		`json:"error"`
	Path      string    	`json:"path"`
	Output    string    	`json:"output"`
	BeginTime time.Time 	`json:"begin_time"`
	EndTime   time.Time 	`json:"end_time"`
}

type BackupTaskResult struct {
	PathResult
	TaskId string
}

func NewPathResult(err error, path, output string, beginTime, endTime time.Time) (result PathResult) {
	result.Path = path
	result.Output = output
	result.BeginTime = beginTime
	result.EndTime = endTime
	if err != nil {
		result.Err = err.Error()
	}
	return
}

func (pathResult *PathResult) IsSuccess() bool {
	return len(pathResult.Err) == 0
}
