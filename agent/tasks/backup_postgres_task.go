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

type backupPostgresTask struct {
	*backupSQLTask
}

func newBackupPostgresTask(config *backuper.TaskConfig) BackupTaskInterface {
	postgresTask := backupPostgresTask{newBackupSQLTask(config)}

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

	postgresTask.tmpFileBase = tmpFileBase + tmpFileSQLSuffix
	return &postgresTask
}

func (postgresTask *backupPostgresTask) GenerateTmpFile(tmpFilePath string) ([]byte, error) {

	tables := postgresTask.config.Params["db_tables"]

	if len(tables) == 0 {
		tables = "\\*"
	}

	params := &postgresTask.config.Params

	cmd := fmt.Sprintf("PGPASSWORD=%s pg_dump %s -h %s -p %s -U %s %s %s %s >%s",
		(*params)["db_pass"],
		postgresTask.recreateFlag(),
		(*params)["db_host"],
		(*params)["db_port"],
		(*params)["db_user"],
		(*params)["db_base"],
		postgresTask.tablesFlag(),
		postgresTask.compressionFilter(),
		tmpFilePath,
	)

	out, err := hmutil.System(cmd)
	if len(out) > 0 {
		err = errors.New("pg_dump failed")
	}

	return out, err
}

func (postgresTask *backupPostgresTask) recreateFlag() (rf string) {
	if len(postgresTask.config.Params["recreate"]) > 0 {
		rf = "-c"
	}
	return
}

func (postgresTask *backupPostgresTask) tablesFlag() (tf string) {
	tables := postgresTask.config.Params["tables"]
	if len(tables) > 0 {
		tf = "-t " + tables
	}
	return
}

// func (postgresTask *backupPostgresTask) compressionFilter() (cf string) {
// 	if postgresTask.needCompression() {
// 		cf = " | bzip2 -c "
// 	}
// 	return
// }

// func (postgresTask *backupPostgresTask) TmpFileName() string {
// 	return strings.Join([]string{
// 		postgresTask.Type(),
// 		postgresTask.tmpFileBase,
// 	}, "_") + postgresTask.compressionSuffix()
// }

// func (postgresTask *backupPostgresTask) compressionSuffix() (cs string) {
// 	if postgresTask.needCompression() {
// 		cs = ".bz2"
// 	}
// 	return
// }
