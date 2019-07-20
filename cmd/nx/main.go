package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hokiegeek/godemo"
	"github.com/hokiegeek/gonexus-private/iq"
	"github.com/sonatype-nexus-community/gonexus/iq"
	// "github.com/sonatype-nexus-community/gonexus/rm"

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

func listRepos(idx int) {
	// if repos, err := nexusrm.GetRepositories(demo.RM(0)); err == nil {
	if repos, err := demo.Repos(idx); err == nil {
		for _, r := range repos {
			fmt.Printf("%-15s (%6s : %s)\n", r.Name, r.Format, r.Type)
		}
	}
}

func createAndDeleteOrg() {
	orgID, err := nexusiq.CreateOrganization(demo.IQ(0), "arstarst")
	if err != nil {
		panic(err)
	}

	time.Sleep(15 * time.Second)

	if err := privateiq.DeleteOrganization(demo.IQ(0), orgID); err != nil {
		panic(err)
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
			/*
				Flags: []cli.Flag{
					cli.IntFlag{
						Name:  "i",
						Value: 0,
						Usage: "repo `idx`",
					},
				},
			*/
			Subcommands: []cli.Command{
				{
					Name:    "repos",
					Aliases: []string{"l", "ls"},
					Usage:   "lists all repos",
					Action: func(c *cli.Context) error {
						// fmt.Println(c.GlobalInt("i"))
						// fmt.Println(c.Int("i"))
						listRepos(0)
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

	fmt.Println("Detecting Nexus servers...")
	demo.Detect()

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
