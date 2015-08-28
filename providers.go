package main

import (
	"log"

	"github.com/codegangsta/cli"
	"github.com/headmade/backuper/config"
)

func providerAction(c *cli.Context) {
	var providerConfig config.Provider
	switch c.Command.Name {
	case "AWS":
		validateArgs(c, 2)
		providerConfig = config.Provider{"AWS_ACCESS_KEY_ID": c.Args()[0], "AWS_SECRET_ACCESS_KEY": c.Args()[1]}
	case "encryption":
		validateArgs(c, 1)
		providerConfig = config.Provider{"pass": c.Args()[0]}
	}
	providerCommand(c.Command.Name, providerConfig)
}

func validateArgs(c *cli.Context, length int) {
	if len(c.Args()) != length {
		log.Fatal("Bad arguments")
	}
}

func providerCommand(name string, providerConfig config.Provider) {
	conf, err := config.New()
	if err != nil {
		log.Fatal(err)
	}
	if conf.Secret == nil {
		conf.Secret = config.Providers{}
	}
	// conf.Secret[name] = providerConfig
	if conf.Secret[name] == nil {
		conf.Secret[name] = config.Provider{}
	}
	for k, v := range providerConfig {
		conf.Secret[name][k] = v
	}
	conf.Write(conf.Secret)
}
