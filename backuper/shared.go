package backuper

import (
	"time"
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
	Type   string            `json:"type"`
	Name   string            `json:"name"`
	Params map[string]string `json:"params"`
}

type BackupResult struct {
	Prepare   PathResult
	Lock      PathResult
	Backup    []BackupTaskResult
	Encrypt   PathResult
	Upload    PathResult
	Unlock    PathResult
	Cleanup   PathResult
	BeginTime time.Time
	EndTime   time.Time
	Size      int64
	Status    string
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
