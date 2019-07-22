package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hokiegeek/godemo"
	"github.com/urfave/cli"
)

func listServers() {
	for i, s := range demo.RMs {
		fmt.Printf("RM[%d]: %s\n", i, s.Host)
	}
	for i, s := range demo.IQs {
		fmt.Printf("IQ[%d]: %s\n", i, s.Host)
	}
}

func main() {
	app := cli.NewApp()

	app.Commands = []cli.Command{
		{
			Name:    "list",
			Aliases: []string{"ls"},
			Usage:   "lists all detected Nexus servers",
			Action: func(c *cli.Context) error {
				listServers()
				return nil
			},
		},
		rmCommand,
		iqCommand,
	}

	app.Action = func(c *cli.Context) error {
		listServers()
		return nil
	}

	log.Println("Discovering Nexus servers...")
	demo.Detect()

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
