package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hokiegeek/godemo"
	"github.com/sonatype-nexus-community/gonexus/iq"
	"github.com/urfave/cli"
)

var iqCommand = cli.Command{
	Name:    "iq",
	Aliases: []string{"q"},
	Usage:   "iq-specific commands",
	Flags: []cli.Flag{
		cli.IntFlag{
			Name:  "idx, i",
			Value: 0,
			Usage: "iq `idx`",
		},
	},
	Subcommands: []cli.Command{
		{
			Name:    "app",
			Aliases: []string{"a"},
			Usage:   "lists all applications",
			Action: func(c *cli.Context) error {
				iqListApps(c.Parent().Int("idx"))
				return nil
			},
		},
		{
			Name:    "org",
			Aliases: []string{"o"},
			Usage:   "lists all organizations",
			Action: func(c *cli.Context) error {
				iqListOrgs(c.Parent().Int("idx"))
				return nil
			},
		},
		{
			Name:    "eval",
			Aliases: []string{"e"},
			Usage:   "evaluate the components",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "app, a"},
			},
			Action: func(c *cli.Context) error {
				iqEvaluate(c.Parent().Int("idx"), c.String("app"), c.Args())
				return nil
			},
		},
	},
}

func iqListApps(idx int) {
	format := "%s, %s, %s, %s\n"
	fmt.Printf(format, "Name", "Public ID", "ID", "Organization ID")
	orgsId2Name, _, _ := demo.OrgsIdMap(idx)
	if apps, err := nexusiq.GetAllApplications(demo.IQ(idx)); err == nil {
		for _, a := range apps {
			fmt.Printf(format, a.Name, a.PublicID, a.ID, orgsId2Name[a.OrganizationID])
		}
	}
}

func iqListOrgs(idx int) {
	format := "%s, %s\n"
	fmt.Printf(format, "Name", "ID")
	if orgs, err := nexusiq.GetAllOrganizations(demo.IQ(idx)); err == nil {
		for _, o := range orgs {
			fmt.Printf(format, o.Name, o.ID)
		}
	}
}

func iqEvaluate(idx int, app string, components []string) {
	report, err := demo.Eval(idx, app, components...)
	if err != nil {
		log.Fatal(err)
	}

	json, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(json))
}
