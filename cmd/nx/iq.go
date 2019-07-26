package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"text/template"

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
				cli.StringFlag{Name: "format, f"},
			},
			Action: func(c *cli.Context) error {
				iqEvaluate(c.Parent().Int("idx"), c.String("app"), c.String("format"), c.Args())
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
		{
			Name:  "sc",
			Usage: "source control actions",
			Action: func(c *cli.Context) error {
				scList(c.Parent().Int("idx"), "")
				return nil
			},
			Subcommands: []cli.Command{
				{
					Name:    "create",
					Aliases: []string{"c"},
					Usage:   "creates a source control entry",
					Action: func(c *cli.Context) error {
						/*
							appIDPtr := createCmd.String("app", "", "The identifier of the application in IQ")
							repoPtr := createCmd.String("repo", "", "The repo")
							tokenPtr := createCmd.String("token", "", "SC Token")
						*/
						scCreate(c.Parent().Parent().Int("idx"))
						return nil
					},
				},
				{
					Name:    "delete",
					Aliases: []string{"d"},
					Usage:   "deletes a source control entry",
					Action: func(c *cli.Context) error {
						/*
							appIDPtr := deleteCmd.String("app", "", "The identifier of the application in IQ")
							entryPtr := deleteCmd.String("entry", "", "The ID of the SC entry")

							deleteCmd.Parse(os.Args[2:])

							var scEntryID string
							if *entryPtr != "" {
								scEntryID = *entryPtr
							} else {
								scEntry, _ := get(iq, *appIDPtr)
								scEntryID = scEntry.ID
							}

							del(iq, *appIDPtr, scEntryID)
						*/
						scDelete(c.Parent().Parent().Int("idx"))
						return nil
					},
				},
				{
					Name:    "list",
					Aliases: []string{"l"},
					Usage:   "deletes a source control entry",
					Action: func(c *cli.Context) error {
						scList(c.Parent().Parent().Int("idx"), c.Args().First())
						return nil
					},
				},
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

func iqEvaluate(idx int, app, format string, components []string) {
	report, err := demo.Eval(idx, app, components...)
	if err != nil {
		log.Fatal(err)
	}

	if format != "" {
		tmpl := template.Must(template.New("report").Parse(format))
		tmpl.Execute(os.Stdout, report)
	} else {
		json, err := json.MarshalIndent(report, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(json))
	}
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

func scCreate(idx int) {
	app, repo, token := "", "", ""
	iq := demo.IQ(idx)
	err := nexusiq.CreateSourceControlEntry(iq, app, repo, token)
	if err != nil {
		panic(err)
	}

	entry, err := nexusiq.GetSourceControlEntry(iq, app)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%q\n", entry)
}

func scDelete(idx int) {
	app, id := "", ""
	nexusiq.DeleteSourceControlEntry(demo.IQ(idx), app, id)
}

func scList(idx int, appID string) {
	/*
		appIDPtr := listCmd.String("app", "", "The identifier of the application in IQ")

		listCmd.Parse(os.Args[2:])

		if *appIDPtr != "" {
			entry, _ := get(iq, *appIDPtr)
			fmt.Printf("%v\n", entry)
		} else {
			log.Println("listing all entries...")
			apps, err := nexusiq.GetAllApplications(iq)
			if err != nil {
				panic(err)
			}
			for _, app := range apps {
				if entry, err := get(iq, app.PublicID); err == nil {
					fmt.Printf("%s: %v\n", app.PublicID, entry)
				}
			}
		}
	*/
}
