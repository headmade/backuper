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
	"github.com/headmade/backuper/schedule"
)

const (
	version = "0.0.1"
)

// BackendAddr is server host
func BackendAddr() string {
	backend := os.Getenv("BACKEND")
	if backend == "" {
		backend = "api.backuper.headmade.pro"
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

func checkUID(commandName string) {
	if false { //os.Getuid() != 0 {
		fmt.Printf("FAILED! Are you root? Please, run `sudo rollbackup %s [ARGS]`\n", commandName)
		os.Exit(0)
	}
}

func initAction(c *cli.Context) {
	checkUID(c.Command.Name)
	scheduler, err := schedule.New()
	if err != nil {
		log.Fatal(err)
	}
	if err := scheduler.UpdateCheck(); err != nil {
		log.Fatal(err)
	}
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

func checkAction(c *cli.Context) {
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
		scheduler, err := schedule.New()
		if err != nil {
			log.Fatal(err)
		}
		if err := scheduler.UpdateBackup(agentConfig.Period); err != nil {
			log.Println(err)
		}

		if conf.Agent.StartNow {
			log.Println("StartNow")
			backupAction(c)
			conf.Agent.StartNow = !conf.Agent.StartNow
			conf.Write(agentConfig)
		}
	}
}

func backupAction(c *cli.Context) {
	conf, err := config.New()
	if err != nil {
		log.Fatal("This server is not ready to backup. Please exec 'backuper init'")
	}
	agent, err := agent.Get(conf.Agent, &conf.Secret)
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
	app.Version = version
	app.Name = "backuper"
	app.Author = "Headmade LLC"
	app.Email = "backuper@headmade.pro"
	app.Usage = "A client utility for manage backuper"
	app.Commands = []cli.Command{
		{
			Name:   "init",
			Usage:  "Configure agent with signed token",
			Action: initAction,
		},
		{
			Name:   "backup",
			Usage:  "Make a backup",
			Action: backupAction,
		},
		{
			Name:   "check",
			Usage:  "Check server change",
			Action: checkAction,
		},
		{
			Name:      "provider",
			ShortName: "p",
			Usage:     "Add provider params",
			//Action:    ProviderAction,
			Subcommands: []cli.Command{
				{
					Name:      "AWS",
					ShortName: "aws",
					Usage:     "AWS [AWS_ACCESS_KEY_ID] [AWS_SECRET_ACCESS_KEY]",
					Action:    providerAction,
				},
				{
					Name:      "encrypt",
					ShortName: "enc",
					Usage:     "encrypt [PASSWORD]",
					Action:    providerAction,
				},
			},
		},
	}
	app.Run(os.Args)
}
