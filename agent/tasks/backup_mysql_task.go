package tasks

import (
  "errors"
  "fmt"
  "regexp"
  "strings"

  "github.com/headmade/backuper/backuper"
  "github.com/headmade/backuper/hmutil"
)

const (
  tmpFileMySQLSuffix = ".sql"
)

type backupMySQLTask struct {
  *backupSQLTask
}

func newBackupMySQLTask(config *backuper.TaskConfig) BackupTaskInterface {
  mysqlTask := backupMySQLTask{newBackupSQLTask(config)}

  params := &config.Params
  tmpFileBase := strings.Join([]string{
    (*params)["db_host"],
    (*params)["db_port"],
    (*params)["db_base"],
  }, "_")
  
  tmpFileBase = strings.Join(
    regexp.MustCompile(`[^\d\w]+`).Split(tmpFileBase, -1),
    "_",
  )

  mysqlTask.tmpFileBase = tmpFileBase + tmpFileSQLSuffix

  return &mysqlTask
}

func (mysqlTask *backupMySQLTask) GenerateTmpFile(tmpFilePath string) ([]byte, error) {
  database := mysqlTask.config.Params["db_base"]
  tables := mysqlTask.config.Params["db_tables"]
  password := mysqlTask.config.Params["db_pass"]

  if len(database) == 0 && len(tables) == 0 {
    database = "--all-databases" 
  }

  if len(password) != 0 {
    password = fmt.Sprintf("MYSQL_PWD=%s", password)
  }

  params := &mysqlTask.config.Params

  cmd := fmt.Sprintf("%s mysqldump -h %s -P %s -u %s %s %s %s > %s",
    password,
    (*params)["db_host"],
    (*params)["db_port"],
    (*params)["db_user"],
    database,
    tables,
    mysqlTask.compressionFilter(),
    tmpFilePath,
  )

  out, err := hmutil.System(cmd)
  if len(out) > 0 {
    err = errors.New("mysqldump failed")
  }

  return out, err
}

// func (mysqlTask *backupMySQLTask) compressionFilter() (cf string) {
//   if mysqlTask.needCompression() {
//     cf = "| bzip2 -с"
//   }
//   return
// }

// func (mysqlTask *backupMySQLTask) TmpFileName() string {
//   return strings.Join([]string{
//     mysqlTask.Type(),
//     mysqlTask.tmpFileBase,
//   }, "_") + mysqlTask.compressionSuffix()
// }

// func (mysqlTask *backupMySQLTask) compressionSuffix() (cs string) {
//   if mysqlTask.needCompression() {
//     cs = ".bz2"
//   }
//   return
// }
