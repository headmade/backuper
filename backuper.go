package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"encoding/json"
	"strconv"
	"strings"
	"path/filepath"

	"github.com/rlmcpherson/s3gof3r"
	"github.com/codegangsta/cli"
	"github.com/headmade/backuper/agent"
	"github.com/headmade/backuper/backuper"
	"github.com/headmade/backuper/client"
	"github.com/headmade/backuper/config"
	"github.com/headmade/backuper/schedule"
	"github.com/headmade/backuper/hmutil"
)

const (
	version = "0.0.1"
)

// BackendAddr is server host
func BackendAddr() string {
	backend := os.Getenv("BACKEND")
	if backend == "" {
		backend = "https://gobackuper.com/v1"
		// backend = "http://localhost:3000/v1"
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
		client, err := client.NewClient(BackendAddr())
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
		client, err := client.NewClient(BackendAddr())
		if err != nil {
			log.Fatal(err)
		}
		if err := client.Backup(result); err != nil {
			log.Printf("Backup notification Error: %s", err)
		}
	}
}

func listAction(c *cli.Context) {
	conf, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	if !conf.Local {
		cl, err := client.NewClient(BackendAddr())
		if err != nil {
			log.Fatal(err)
		}

		var blist []backuper.BackupResult
		if err = cl.GetBackupsList(&blist); err != nil {
			log.Fatal(err)
		}

		tail, err := strconv.Atoi(c.Args()[0])
		if err != nil {
			log.Fatal(err)
		}

		r, err := json.MarshalIndent(blist[:tail], "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(r))
	}
}

func retrieveAction(c *cli.Context) {
	conf, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	if !conf.Local {
		cl, err := client.NewClient(BackendAddr())
		if err != nil {
			log.Fatal(err)
		}

		var backup backuper.BackupResult
		id, err := strconv.Atoi(c.Args()[0])
		if err != nil {
			log.Fatal(err)
		}

		if err = cl.GetBackup(&backup, id); err != nil {
			log.Fatal(err)
		}

		fmt.Println(backup.Upload)

		_, name := filepath.Split(backup.Upload.Path)
		name = strings.Split(name, ".encrypted")[0]
		switch backup.Upload.Type {
		case "ssh":
			addr := strings.Split(backup.Upload.Destination, "@")
			login, host := addr[0], addr[1]
			ssh := hmutil.SSHDownloader(
				22,
				conf.Secret["SSH"]["id_rsa"],
				host,
				login,
				backup.Upload.Path,
				conf.Secret["SSH"]["download_path"],
			)
			result, err := hmutil.SSHExec(ssh);

			if err != nil {
				log.Fatal(err)
			}

			decoded_result, err := hmutil.Decode(&result, []byte(conf.Secret["encryption"]["pass"]))

			if err != nil {
				log.Fatal(err)
			}

			if err = hmutil.WriteToFile(conf.Secret["SSH"]["download_path"] + name, decoded_result); err != nil {
				log.Fatal(err)
			}
		case "ftp":
			addr := strings.Split(backup.Upload.Destination, "@")
			login, host := addr[0], addr[1]
			result, err := hmutil.FTPDownload(21, host, login, conf.Secret["FTP"]["passw"], backup.Upload.Path)

			if err != nil {
				log.Fatal(err)
			}

			decoded_result, err := hmutil.Decode(&result, []byte(conf.Secret["encryption"]["pass"]))

			if err != nil {
				log.Fatal(err)
			}

			if err = hmutil.WriteToFile(conf.Secret["FTP"]["download_path"] + name, decoded_result); err != nil {
				log.Fatal(err)
			}
		case "S3":
			result, err := hmutil.DownloadFromS3(
				s3gof3r.Keys{
					AccessKey: conf.Secret["AWS"]["AWS_ACCESS_KEY_ID"], 
					SecretKey: conf.Secret["AWS"]["AWS_SECRET_ACCESS_KEY"],
				},
				backup.Upload.Destination,
				backup.Upload.Path,
			)

			if err != nil {
				log.Fatal(err)
			}

			decoded_result, err := hmutil.Decode(&result, []byte(conf.Secret["encryption"]["pass"]))

			if err != nil {
				log.Fatal(err)
			}

			if err = hmutil.WriteToFile(conf.Secret["AWS"]["download_path"] + name, decoded_result); err != nil {
				log.Fatal(err)
			}
		}
	}
}

func main() {
	/*
		// test crash debug info
		log.Println("GOTRACEBACK", os.Getenv("GOTRACEBACK"))
		if os.Getenv("GOTRACEBACK") == "" {
			os.Setenv("GOTRACEBACK", "0")
		}
		log.Println("GOTRACEBACK", os.Getenv("GOTRACEBACK"))
		var s *string
		*s = "" // crash
	*/
	app := cli.NewApp()
	app.Version = version
	app.Name = "gobackuper"
	app.Author = "Headmade LLC"
	app.Email = "support@gobackuper.com"
	app.Usage = "A client utility to perform backups and manage local config"
	app.Commands = []cli.Command{
		{
			Name:   "init",
			Usage:  "Configure agent with signed token",
			Action: initAction,
		},
		{
			Name:   "backup",
			Usage:  "Perform a backup [METHOD]",
			Action: backupAction,
		},
		{
			Name:   "check",
			Usage:  "Check if backup config changed",
			Action: checkAction,
		},
		{
			Name: "list",
			Usage: "Print list of backups [TAIL]",
			Action: listAction,
		},
		{
			Name: "get",
			Usage: "Download backup [ID PATH [ID_RSA]]",
			Action: retrieveAction,
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
					// Action:    providerAction,
				},
				{
					Name:      "encrypt",
					ShortName: "enc",
					Usage:     "encrypt [PASSWORD]",
					// Action:    providerAction,
				},
			},
		},
	}
	app.Run(os.Args)
}
