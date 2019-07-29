package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
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
		cli.StringFlag{
			Name:  "server, s",
			Value: "http://localhost:8070",
		},
		cli.StringFlag{
			Name:  "auth, a",
			Value: "admin:admin123",
		},
	},
	Before: func(c *cli.Context) error {
		host := c.String("server")
		auth := strings.Split(c.String("auth"), ":")
		if host != "" && len(auth) == 2 {
			log.Printf("Connecting to %s\n", host)
			demo.IQs = []demo.IdentifiedIQ{demo.NewIdentifiedIQ(host, auth[0], auth[1])}
		}
		return nil
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
		{
			Name:    "report",
			Aliases: []string{"r"},
			Usage:   "r [appID:stage] [appID:stage] [appID]",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "format, f"},
			},
			Action: func(c *cli.Context) error {
				appReport(c.Parent().Int("idx"), c.String("format"), c.Args()...)
				return nil
			},
		},
		{
			Name:    "remediation",
			Aliases: []string{"rem"},
			Flags: []cli.Flag{
				cli.StringFlag{Name: "app, a"},
				cli.StringFlag{Name: "org, o"},
				cli.StringFlag{Name: "stage, s", Value: "build"},
				cli.StringFlag{Name: "format, f"},
			},
			Action: func(c *cli.Context) error {
				remediation(c.Parent().Int("idx"), c.String("format"), c.String("stage"), c.String("app"), c.String("org"), c.Args().First())
				return nil
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
	if orgs, err := nexusiq.GetAllOrganizations(demo.RM(idx)); err == nil {
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
		tmpl := template.Must(template.New("report").Funcs(template.FuncMap{"json": tmplJSONPretty}).Parse(format))
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
	iq := demo.IQ(idx)
	if appID != "" {
		entry, _ := nexusiq.GetSourceControlEntry(iq, appID)
		fmt.Printf("%v\n", entry)
	} else {
		apps, err := nexusiq.GetAllApplications(iq)
		if err != nil {
			panic(err)
		}
		for _, app := range apps {
			if entry, err := nexusiq.GetSourceControlEntry(iq, app.PublicID); err == nil {
				fmt.Printf("%s: %v\n", app.PublicID, entry)
			}
		}
	}
}

func appReport(idx int, format string, apps ...string) {
	for _, app := range apps {
		splitPos := strings.LastIndex(app, ":")
		appID := app[:splitPos]
		stage := app[splitPos+1:]

		report, err := nexusiq.GetReportByAppID(demo.IQ(idx), appID, stage)
		if err != nil {
			log.Printf("did not find report for '%s' at '%s' build stage: %v", appID, stage, err)
		}

		if format != "" {
			tmpl := template.Must(template.New("report").Funcs(template.FuncMap{"json": tmplJSONPretty}).Parse(format))
			tmpl.Execute(os.Stdout, report)
		} else {
			json, err := json.MarshalIndent(report, "", "  ")
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(string(json))
		}
	}
}

func remediation(idx int, format, stage, app, org, comp string) {
	c, _ := nexusiq.NewComponentFromString(comp)
	var err error
	var remediation nexusiq.Remediation
	switch {
	case app != "":
		remediation, err = nexusiq.GetRemediationByApp(demo.IQ(idx), *c, stage, app)
	case org != "":
		remediation, err = nexusiq.GetRemediationByOrg(demo.IQ(idx), *c, stage, org)
	default:
		panic("Need either an app or an org")
	}
	if err != nil {
		panic(err)
	}

	if format != "" {
		tmpl := template.Must(template.New("remediation").Funcs(template.FuncMap{"json": tmplJSONPretty}).Parse(format))
		tmpl.Execute(os.Stdout, remediation)
	} else {
		json, err := json.MarshalIndent(remediation, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(json))
	}
}
