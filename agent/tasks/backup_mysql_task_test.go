package tasks

import (
  "testing"
  // "fmt"
  // "regexp"
  "strings"

  "github.com/headmade/backuper/backuper"
  // "github.com/headmade/backuper/hmutil"
)

var configuration *backuper.TaskConfig = &backuper.TaskConfig{
  Type: "mysql",
  Name: "",
  Compression: "bzip2",
  Params: map[string]string{
    "db_base": "example",
    "db_host": "localhost",
    "db_pass": "1995",
    "db_port": "3306",
    "db_sock": "/tmp/mysql.sock",
    "db_user": "twizty",
    "recreate": "",
    "db_tables": "ex_table",
  },
}

func TestFileGenetation(t *testing.T) {
  mysqlTask := newBackupMySQLTask(configuration)
  dumpname := strings.Join([]string{
    "mysql",
    (*configuration).Params["db_host"], 
    (*configuration).Params["db_port"],
    (*configuration).Params["db_user"],
    (*configuration).Params["db_base"],
  }, "_")
  dumpname = strings.Join([]string{dumpname, "sql", "bz2"}, ".")

  if dumpname != mysqlTask.TmpFileName() {
    t.Error("Sth is wrong with generating TmpFileName")
  }
  // TODO: дописать тест с использованием раннера.
}
