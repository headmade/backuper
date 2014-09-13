package tasks

import (
	"fmt"
	"log"
)

type backupPostgresTask struct {
	*backupTask
}

func newBackupPostgresTask(config *Config) TaskInterface {
	return &backupPostgresTask{newBackupTask(config)}
}

func (self *backupPostgresTask) Run() error {
	log.Println("run backupPostgresTask")
	self.PrepareTmpDirectory()

	tables := self.config.Params["db_tables"]

	if len(tables) == 0 {
		tables = "\\*"
	}

	cmd := fmt.Sprintf("PGPASSWORD=%s pg_dump -h %s -p %s -U %s %s -t %s | bzip2 | %s >%s",
		self.config.Params["db_pass"],
		self.config.Params["db_host"],
		self.config.Params["db_port"],
		self.config.Params["db_user"],
		self.config.Params["db_name"],
		tables,
		self.EncryptCmd(self.config.Params["pass"]),
		self.tmpFilePath(),
	)
	log.Println(cmd)

	out, err := self.System(cmd)
	log.Println(string(out))
	self.CleanupTmpDirectory()

	return err
}

