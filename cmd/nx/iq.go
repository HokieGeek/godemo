package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hokiegeek/godemo"
	"github.com/hokiegeek/gonexus-private/iq"
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
		{
			Name:  "policies",
			Usage: "Do stuff with policies",
			Subcommands: []cli.Command{
				{
					Name:    "export",
					Aliases: []string{"a"},
					Usage:   "exports the policies of the indicated IQ",
					Action: func(c *cli.Context) error {
						exportPolicies(c.Parent().Int("idx"))
						return nil
					},
				},
				/*
					{
						Name:    "import",
						Aliases: []string{"i"},
						Usage:   "Import the indicated policies",
						Action: func(c *cli.Context) error {
							// iqListOrgs(c.Parent().Int("idx"))
							return nil
						},
					},
				*/
			},
		},
	},
}

func iqListApps(idx int) {
	fmt.Printf("%s, %s, %s, %s\n", "Name", "Public ID", "ID", "Organization ID")
	orgsID2Name, _, _ := demo.OrgsIDMap(idx)
	if apps, err := nexusiq.GetAllApplications(demo.IQ(idx)); err == nil {
		for _, a := range apps {
			fmt.Printf("%s, %s, %s, %s\n", a.Name, a.PublicID, a.ID, orgsID2Name[a.OrganizationID])
		}
	}
}

func iqListOrgs(idx int) {
	fmt.Printf("%s, %s\n", "Name", "ID")
	if orgs, err := demo.Orgs(idx); err == nil {
		for _, o := range orgs {
			fmt.Printf("%s, %s\n", o.Name, o.ID)
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

func exportPolicies(idx int) {
	policies, err := privateiq.ExportPolicies(demo.IQ(idx))
	if err != nil {
		log.Fatal(err)
	}

	json, err := json.MarshalIndent(policies, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(json))
}
