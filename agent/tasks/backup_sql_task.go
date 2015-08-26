package tasks

import (
  "strings"

  "github.com/headmade/backuper/backuper"
)

const (
  tmpFileSQLSuffix = ".sql"
)

type backupSQLTask struct {
  *backupTask
}

func newBackupSQLTask(config *backuper.TaskConfig) *backupSQLTask {
  return &backupSQLTask{&backupTask{config: config}}
}

func (sqlTask *backupSQLTask) compressionFilter() (cf string) {
  if sqlTask.needCompression() {
    cf = "| bzip2 -c"
  }
  return
}

func (sqlTask *backupSQLTask) TmpFileName() string {
  return strings.Join([]string{
    sqlTask.Type(),
    sqlTask.tmpFileBase,
  }, "_") + sqlTask.compressionSuffix()
}

func (sqlTask *backupSQLTask) compressionSuffix() (cs string) {
  if sqlTask.needCompression() {
    cs = ".bz2"
  }
  return
}
