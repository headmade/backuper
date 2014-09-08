package main

import (
	"fmt"
	"log"
	"os"

	"github.com/codegangsta/cli"
)

var Version = "development"

func initHost(c *cli.Context) error {
	authId := c.Args().First()
	fmt.Println(authId)
	return nil
}
func CheckUid(commandName string) {
	if false { //os.Getuid() != 0 {
		fmt.Printf("FAILED! Are you root? Please, run `sudo rollbackup %s [ARGS]`\n", commandName)
		os.Exit(0)
	}
}

func InitAction(c *cli.Context) {
	CheckUid(c.Command.Name)
	if err := initHost(c); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Success! Host ready to backup.")
}

func BackupAction(c *cli.Context) {
	fmt.Println(c.Command.Name)

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
	}
	app.Run(os.Args)
}
