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
		{
			Name:    "rm",
			Aliases: []string{"r"},
			Usage:   "repository-specific commands",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "idx, i",
					Value: 0,
					Usage: "repo `idx`",
				},
			},
			/*
				Action: func(c *cli.Context) error {
					TODO: list server info
					fmt.Println(c.Int("idx"))
					listRepos(0)
					return nil
				},
			*/
			Subcommands: []cli.Command{
				{
					Name:    "repos",
					Aliases: []string{"r"},
					Usage:   "lists all repos",
					Action: func(c *cli.Context) error {
						rmListRepos(c.Parent().Int("idx"))
						return nil
					},
				},
				{
					Name:    "ls",
					Aliases: []string{"l"},
					Usage:   "lists all components in a repo",
					Action: func(c *cli.Context) error {
						rmListRepoComponents(c.Parent().Int("idx"), c.Args())
						return nil
					},
				},
				{
					Name:    "upload",
					Aliases: []string{"u", "up"},
					Usage:   "upload component",
					Flags: []cli.Flag{
						cli.StringFlag{Name: "repo, r"},
						cli.StringFlag{Name: "coord, c"},
						cli.StringFlag{Name: "file, f"},
					},
					Action: func(c *cli.Context) error {
						rmUploadComponent(c.Parent().Int("idx"), c.String("repo"), c.String("coord"), c.String("file"))
						return nil
					},
				},
			},
		},
	}

	app.Action = func(c *cli.Context) error {
		listServers()
		return nil
	}

	log.Println("Detecting Nexus servers...")
	demo.Detect()

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
