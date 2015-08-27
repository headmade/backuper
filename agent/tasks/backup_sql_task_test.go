package tasks

import (
  "testing"
  "fmt"
  // "regexp"
  "strings"

  "github.com/headmade/backuper/backuper"
  // "github.com/headmade/backuper/hmutil"
)

var mysql_configuration *backuper.TaskConfig = &backuper.TaskConfig{
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

var postgres_configuration *backuper.TaskConfig = &backuper.TaskConfig{
  Type: "postgres",
  Name: "",
  Compression: "bzip2",
  Params: map[string]string{
    "db_base": "twizty",
    "db_host": "localhost",
    "db_pass": "",
    "db_port": "5432",
    "db_sock": "/tmp/.postgres.sock",
    "db_user": "twizty",
    "recreate": "",
    "tables": "ex_table",
  },
}

func TestMySQLFileGenetation(t *testing.T) {
  mysqlTask := newBackupMySQLTask(mysql_configuration)
  dumpname := strings.Join([]string{
    "mysql",
    (*mysql_configuration).Params["db_host"], 
    (*mysql_configuration).Params["db_port"],
    (*mysql_configuration).Params["db_base"],
  }, "_")
  dumpname = strings.Join([]string{dumpname, "sql", "bz2"}, ".")

  if dumpname != mysqlTask.TmpFileName() {
    t.Error("Sth is wrong with generating TmpFileName")
  }

  out, err := mysqlTask.GenerateTmpFile("/tmp/backuper/" + dumpname)
  
  if len(out) > 0 {
    t.Error(err)
  }
  
  // TODO: дописать тест с использованием раннера.
}

func TestPostgresFileGenetation(t *testing.T) {
  postgresTask := newBackupPostgresTask(postgres_configuration)
  dumpname := strings.Join([]string{
    "postgres",
    (*postgres_configuration).Params["db_host"], 
    (*postgres_configuration).Params["db_port"],
    (*postgres_configuration).Params["db_base"],
  }, "_")
  dumpname = strings.Join([]string{dumpname, "sql", "bz2"}, ".")

  if dumpname != postgresTask.TmpFileName() {
    t.Error("Sth is wrong with generating TmpFileName")
  }

  out, err := postgresTask.GenerateTmpFile("/tmp/backuper/" + dumpname)
  
  if len(out) > 0 {
    fmt.Println(fmt.Sprintf("%s", out))
    t.Error(err)
  }
}
