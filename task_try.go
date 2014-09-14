package main

import (
	"strconv"

	"github.com/headmade/backuper/agent"
	"github.com/headmade/backuper/agent/tasks"
)

func main() {
	types := []string{
		tasks.DIR_TASK_TYPE,
		tasks.POSTGRES_TASK_TYPE,
		"foo",
	}

	configs := make([]*tasks.Config, len(types))

	for i, t := range types {
		configs[i] = &tasks.Config{
			Type: t,
			Id: strconv.Itoa(i),
			Params: map[string]string{
				"tmp_path": "/tmp",
				"pass": "123",
				"dir": "/etc/openssl",
				"db_host": "localhost",
				"db_port": "5432",
				"db_user": "dev",
				"db_pass": "123",
				"db_name": "makerton_development",
				"db_tables": "",
			},
		}
	}

	r := tasks.NewRunner(&agent.Config{TmpDir: "/tmp"})
  r.Run(&configs)
}

