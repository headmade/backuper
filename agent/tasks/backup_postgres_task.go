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

func (self *backupPostgresTask) GenerateBackupFile(tmpFilePath string) ([]byte, error) {

	tables := self.config.Params["db_tables"]

	if len(tables) == 0 {
		tables = "\\*"
	}

	cmd := fmt.Sprintf("PGPASSWORD=%s pg_dump -h %s -p %s -U %s %s -t %s | bzip2 -c >%s",
		self.config.Params["db_pass"],
		self.config.Params["db_host"],
		self.config.Params["db_port"],
		self.config.Params["db_user"],
		self.config.Params["db_base"],
		tables,
		tmpFilePath,
	)
	log.Println(cmd)

	return System(cmd)
}
