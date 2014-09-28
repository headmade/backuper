package tasks

import (
	"fmt"
	"log"

	"github.com/headmade/backuper/backuper"
)

type backupPostgresTask struct {
	*backupTask
}

func newBackupPostgresTask(config *backuper.TaskConfig) BackupTaskInterface {
	return &backupPostgresTask{newBackupTask(config)}
}

func (postgresTask *backupPostgresTask) GenerateBackupFile(tmpFilePath string) (string, []byte, error) {

	tables := postgresTask.config.Params["db_tables"]

	if len(tables) == 0 {
		tables = "\\*"
	}

	cmd := fmt.Sprintf("PGPASSWORD=%s pg_dump -h %s -p %s -U %s %s -t %s | bzip2 -c >%s",
		postgresTask.config.Params["db_pass"],
		postgresTask.config.Params["db_host"],
		postgresTask.config.Params["db_port"],
		postgresTask.config.Params["db_user"],
		postgresTask.config.Params["db_base"],
		tables,
		tmpFilePath,
	)
	log.Println(cmd)

	out, lastErr := System(cmd)
	return tmpFilePath, out, lastErr
}

