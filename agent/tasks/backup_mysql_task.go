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
	tmpFileSqlSuffix = ".sql"
)

type backupMysqlTask struct {
	*backupTask
}

func newBackupMysqlTask(config *backuper.TaskConfig) BackupTaskInterface {
	mysqlTask := backupMysqlTask{newBackupTask(config)}

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

	mysqlTask.tmpFileBase = tmpFileBase + tmpFileSqlSuffix
	return &mysqlTask
}

func (mysqlTask *backupMysqlTask) GenerateTmpFile(tmpFilePath string) ([]byte, error) {
	params := &mysqlTask.config.Params
	
	cmd := fmt.Sprintf("mysqldump --skip-add-drop-table -p%s %s -h %s -P %s -u %s %s %s %s >%s",
		(*params)["db_pass"],
		mysqlTask.recreateFlag(),
		(*params)["db_host"],
		(*params)["db_port"],
		(*params)["db_user"],
		(*params)["db_base"],
		mysqlTask.tablesFlag(),
		mysqlTask.compressionFilter(),
		tmpFilePath,
	)

	out, err := hmutil.System(cmd)
	if len(out) > 0 {
		err = errors.New("mysqldump failed")
	}

	return out, err
}

func (mysqlTask *backupMysqlTask) compressionFilter() (cf string) {
	if mysqlTask.needCompression() {
		cf = " | bzip2 -c "
	}
	return
}

func (mysqlTask *backupMysqlTask) recreateFlag() (rf string) {
	if len(mysqlTask.config.Params["recreate"]) > 0 {
		rf = "--add-drop-table"
	}
	return
}

func (mysqlTask *backupMysqlTask) tablesFlag() (tf string) {
	tables := mysqlTask.config.Params["tables"]
	if len(tables) > 0 {
		tf = tables
	}
	return
}