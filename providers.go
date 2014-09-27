package main

import (
	"log"

	"github.com/codegangsta/cli"
	"github.com/headmade/backuper/config"
)

func ProviderAction(c *cli.Context) {
	var providerConfig config.Provider
	switch c.Command.Name {
	case "AWS":
		validateArgs(c, 2)
		providerConfig = config.Provider{"AWS_ACCESS_KEY_ID": c.Args()[0], "AWS_SECRET_ACCESS_KEY": c.Args()[1]}
	case "encrypt":
		validateArgs(c, 1)
		providerConfig = config.Provider{"PASSWORD": c.Args()[0]}
	}
	providerAction(c.Command.Name, providerConfig)
}

func validateArgs(c *cli.Context, length int) {
	if len(c.Args()) != length {
		log.Fatal("Bad arguments")
	}
}

func providerAction(name string, providerConfig config.Provider) {
	conf, err := config.New()
	if err != nil {
		log.Fatal(err)
	}
	if conf.Secret == nil {
		conf.Secret = config.Providers{}
	}
	conf.Secret[name] = providerConfig
	conf.Write(conf.Secret)
}
