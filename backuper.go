package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/codegangsta/cli"
	"github.com/headmade/backuper/agent"
	"github.com/headmade/backuper/client"
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
	if err := initServer(c); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Success! This server is ready to backup.")
}

func getAgentConfig(client *client.Client) (*agent.Config, error) {
	var config *agent.Config
	response, err := client.GetConfig()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(response, &config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func BackupAction(c *cli.Context) {
	client, err := client.Get(BackendAddr())
	if err != nil {
		log.Fatal(err)
	}
	config, err := getAgentConfig(client)
	if err != nil {
		log.Fatal(err)
	}
	agent, err := agent.Get(config)
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Backup(agent.Backup()); err != nil {
		log.Printf("Backup notification Error: %s", err)
	}
}

func main() {
	app := cli.NewApp()
	//app.Version = Version
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
	}
	app.Run(os.Args)
}
