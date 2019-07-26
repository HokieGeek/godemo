package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/hokiegeek/godemo"
	"github.com/urfave/cli"
)

var tmplJSONPretty = func(v interface{}) string {
	a, _ := json.MarshalIndent(v, "", "  ")
	return string(a)
}

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
	app.Usage = "CLI to interact with Repository Manager and IQ"
	app.HideVersion = true

	defaultAction := func(c *cli.Context) error {
		listServers()
		return nil
	}

	app.Commands = []cli.Command{
		{
			Name:   "ls",
			Usage:  "lists all detected Nexus servers",
			Action: defaultAction,
		},
		rmCommand,
		iqCommand,
	}

	app.Action = defaultAction

	app.Before = func(c *cli.Context) error {
		log.Println("Discovering Nexus servers...")
		demo.Detect()
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
