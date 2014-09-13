package main

import (
	"log"
	"agent/tasks"
)

func main() {
	types := []string{
		tasks.DIR_TASK_TYPE,
		tasks.POSTGRES_TASK_TYPE,
		"foo",
	}

	for _, t := range types {
		c := tasks.Config{Type: t}
		task, err := tasks.Get(&c)
		if err == nil {
			log.Printf("task type: %s, task object: %#v", t, task)
			task.Run()
		} else {
			log.Printf("task type: %s, no registered handler found", t)
		}
	}
}
