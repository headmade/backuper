package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/codegangsta/cli"
	"github.com/headmade/backuper/agent"
	"github.com/headmade/backuper/backuper"
	"github.com/headmade/backuper/client"
	"github.com/headmade/backuper/config"
)

const (
	Version = "0.0.1"
)

func BackendAddr() string {
	backend := os.Getenv("BACKEND")
	if backend == "" {
		backend = "localhost:3000"
	}
	return backend
}

func initServer(c *cli.Context) error {
	token := c.Args().First()
	if len(token) < 1 {
		return errors.New("Invalid auth token")
	}
	err := client.InitServer(BackendAddr(), token)
	return err
}

func CheckUid(commandName string) {
	if false { //os.Getuid() != 0 {
		fmt.Printf("FAILED! Are you root? Please, run `sudo rollbackup %s [ARGS]`\n", commandName)
		os.Exit(0)
	}
}

func InitAction(c *cli.Context) {
	CheckUid(c.Command.Name)
	if c.Args().First() == "local" {
		// conf := config.Config{Local: true}
		conf, _ := config.New()
		conf.Local = true
		conf.Write(true)
	} else {
		if err := initServer(c); err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("Success! This server is ready to backup.")
}

func CheckAction(c *cli.Context) {
	conf, err := config.New()
	if err != nil {
		log.Fatal(err)
	}
	if !conf.Local {
		client, err := client.Get(BackendAddr())
		if err != nil {
			log.Fatal(err)
		}
		var agentConfig *backuper.AgentConfig
		err = client.GetConfig(&agentConfig)
		if err != nil {
			log.Fatal(err)
		}
		conf.Write(agentConfig)

		if conf.Agent.StartNow {
			log.Println("StartNow")
			BackupAction(c)
			conf.Agent.StartNow = !conf.Agent.StartNow
			conf.Write(agentConfig)
		}
	}
}

func BackupAction(c *cli.Context) {
	conf, err := config.New()
	if err != nil {
		log.Fatal("This server is not ready to backup. Please exec 'backuper init'")
	}
	agent, err := agent.Get(conf.Agent)
	if err != nil {
		log.Fatal(err)
	}
	lastErr, result := agent.Backup()
	if lastErr != nil {
		// do smth clever or useful, or both
	}

	if !conf.Local {
		client, err := client.Get(BackendAddr())
		if err != nil {
			log.Fatal(err)
		}
		if err := client.Backup(result); err != nil {
			log.Printf("Backup notification Error: %s", err)
		}
	}
}

func main() {
	app := cli.NewApp()
	app.Version = Version
	app.Name = "backuper"
	app.Author = "Headmade LLC"
	app.Email = "backuper@headmade.pro"
	app.Usage = "A client utility for manage backuper"
	app.Commands = []cli.Command{
		{
			Name:   "init",
			Usage:  "Configure agent with signed token",
			Action: InitAction,
		},
		{
			Name:   "backup",
			Usage:  "Make a backup",
			Action: BackupAction,
		},
		{
			Name:   "check",
			Usage:  "Check server change",
			Action: CheckAction,
		},
	}
	app.Run(os.Args)
}
