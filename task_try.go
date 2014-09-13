package main

import (
	"log"
	"strconv"
	"time"

	"github.com/headmade/backuper/agent/tasks"
)

func main() {
	log.Println(time.Now().Format("2006-01-02 15:04:05"))

	types := []string{
		tasks.DIR_TASK_TYPE,
		tasks.POSTGRES_TASK_TYPE,
		"foo",
	}

	for i, t := range types {
		c := tasks.Config{Type: t, Id: strconv.Itoa(i), Params: map[string]string{
			"tmp_path":  "/tmp",
			"pass":      "123",
			"dir":       "/Users/relevv/git/hm/backuper",
			"db_host":   "localhost",
			"db_port":   "5432",
			"db_user":   "dev",
			"db_pass":   "123",
			"db_name":   "makerton_development",
			"db_tables": "",
		}}
		task, err := tasks.Get(&c)
		if err == nil {
			log.Printf("task type: %s, task object: %#v", t, task)
			btask := task.(tasks.BackupTaskInterface)
			btask.Run()
		} else {
			log.Printf("task type: %s, no registered handler found", t)
		}
	}
}
