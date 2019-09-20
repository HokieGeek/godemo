package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/urfave/cli"

	demo "github.com/hokiegeek/godemo"
	privateiq "github.com/hokiegeek/gonexus-private/iq"
	nexusiq "github.com/sonatype-nexus-community/gonexus/iq"
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
		iqScCommand,
		iqPoliciesCommand,
		iqWaiversCommand,
		iqReportCommand,
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
			Name:    "remediation",
			Aliases: []string{"rem"},
			Flags: []cli.Flag{
				cli.StringFlag{Name: "app, a"},
				cli.StringFlag{Name: "org, o"},
				cli.StringFlag{Name: "stage, s", Value: "build"},
				cli.StringFlag{Name: "format, f"},
				// cli.StringFlag{Name: "report, r"},
			},
			Action: func(c *cli.Context) error {
				remediation(c.Parent().Int("idx"), c.String("format"), c.String("stage"), c.String("app"), c.String("org"), c.Args().First())
				return nil
			},
		},
		{
			Name:    "license",
			Aliases: []string{"lic"},
			Action: func(c *cli.Context) error {
				installLicense(c.Parent().Int("idx"), c.Args().First())
				return nil
			},
		},
		{
			Name:  "zip",
			Usage: "get support zip",
			Action: func(c *cli.Context) error {
				iqZip(c.Parent().Int("idx"))
				return nil
			},
		},
		{
			Name:    "webhook",
			Aliases: []string{"wh"},
			Flags: []cli.Flag{
				cli.StringFlag{Name: "secret, s"},
				cli.StringFlag{Name: "url, u"},
				cli.StringFlag{Name: "events, e", Value: "Application Evaluation"},
			},
			Action: func(c *cli.Context) error {
				webhook(c.Parent().Int("idx"), c.String("url"), c.String("secret"), c.String("events"))
				return nil
			},
		},
		{
			Name:    "auto-apps",
			Aliases: []string{"auto"},
			Flags: []cli.Flag{
				cli.BoolFlag{Name: "disable, d"},
			},
			Action: func(c *cli.Context) error {
				autoApps(c.Parent().Int("idx"), c.Bool("disable"), c.Args().First())
				return nil
			},
		},
		{
			Name:  "violations",
			Usage: "List violations by policy name",
			Action: func(c *cli.Context) error {
				listViolatingApps(c.Parent().Int("idx"), c.Args()...)
				return nil
			},
		},
		{
			Name:    "notice",
			Aliases: []string{"msg"},
			Flags: []cli.Flag{
				cli.BoolFlag{Name: "disable, d"},
			},
			Usage: "Set a message in IQ",
			Action: func(c *cli.Context) error {
				systemNotice(c.Parent().Int("idx"), c.Bool("disable"), strings.Join(c.Args(), " "))
				return nil
			},
		},
		{
			Name:    "search",
			Aliases: []string{"q"},
			Action: func(c *cli.Context) error {
				iqSearch(c.Parent().Int("idx"), c.Args())
				return nil
			},
		},
		{
			Name: "retention",
			Action: func(c *cli.Context) error {
				retentionList(c.Parent().Int("idx"), c.Args().First())
				return nil
			},
			Subcommands: []cli.Command{
				{
					Name:    "list",
					Aliases: []string{"ls"},
					Action: func(c *cli.Context) error {
						retentionList(c.Parent().Parent().Int("idx"), c.Args().First())
						return nil
					},
				},
			},
		},
		{
			Name: "role",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "app, a"},
				cli.StringFlag{Name: "org, o"},
			},
			Action: func(c *cli.Context) error {
				rolesList(c.Parent().Int("idx"), c.String("app"), c.String("org"))
				return nil
			},
			Subcommands: []cli.Command{
				{
					Name:    "user",
					Aliases: []string{"user"},
					Action: func(c *cli.Context) error {
						retentionList(c.Parent().Int("idx"), c.Args().First())
						return nil
					},
				},
			},
		},
		{
			Name: "vulns",
			Action: func(c *cli.Context) error {
				iqVulnInfo(c.Parent().Int("idx"), c.Args()...)
				return nil
			},
		},
		{
			Name: "component",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "format, f"},
			},
			Action: func(c *cli.Context) error {
				iqComponentDetails(c.Parent().Int("idx"), c.String("format"), c.Args()...)
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
	if orgs, err := nexusiq.GetAllOrganizations(demo.IQ(idx)); err == nil {
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
		buf, err := json.MarshalIndent(remediation, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(buf))
	}
}

func iqZip(idx int) {
	zip, name, err := privateiq.GetSupportZip(demo.IQ(idx))
	if err != nil {
		panic(err)
	}

	if err = ioutil.WriteFile(name, zip, 0644); err != nil {
		panic(err)
	}

	log.Printf("Created %s\n", name)
}

func installLicense(idx int, licensePath string) {
	license, err := os.Open(licensePath)
	if err != nil {
		panic(err)
	}

	err = privateiq.InstallLicense(demo.IQ(idx), license)
	if err != nil {
		panic(err)
	}

	log.Println("Installed license")
}

func webhook(idx int, url, secret, events string) {
	err := privateiq.CreateWebhook(demo.IQ(idx), url, secret, strings.Split(events, ","))
	if err != nil {
		panic(err)
	}
	log.Println("Created webhook")
}

func autoApps(idx int, disable bool, orgName string) {
	if disable {
		err := privateiq.DisableAutomaticApplications(demo.IQ(idx))
		if err != nil {
			panic(err)
		}
		log.Println("Disabled automatic applications")
	} else {
		err := privateiq.EnableAutomaticApplications(demo.IQ(idx), orgName)
		if err != nil {
			panic(err)
		}
		log.Println("Enabled automatic applications")
	}
}

func listViolatingApps(idx int, policyNames ...string) {
	var violations []nexusiq.ApplicationViolation
	var err error
	if len(policyNames) > 0 {
		violations, err = nexusiq.GetPolicyViolationsByName(demo.IQ(idx), policyNames...)
	} else {
		violations, err = nexusiq.GetAllPolicyViolations(demo.IQ(idx))
	}
	if err != nil {
		panic(err)
	}

	fmt.Println(violations)
}

func systemNotice(idx int, disable bool, message string) {
	var err error
	if disable {
		err = privateiq.DisableNotice(demo.IQ(idx))
	} else {
		err = privateiq.EnableNotice(demo.IQ(idx), message)
	}
	if err != nil {
		panic(err)
	}
}

func iqSearch(idx int, criteria []string) {
	query := nexusiq.NewSearchQueryBuilder()
	for _, c := range criteria {
		key := strings.Split(c, "=")[0]
		val := strings.Split(c, "=")[1]
		switch key {
		case "stage":
			query = query.Stage(val)
		case "hash":
			query = query.Hash(val)
		case "format":
			query = query.Format(val)
		case "purl":
			query = query.PackageURL(val)
			/*
						case "coord":
							var c nexusiq.Coordinates
							if err := json.Unmarshal([]byte(val), &c); err != nil {
								panic(err)
							}
							query = query.Coordinates(c)
				case "id":
					var c nexusiq.ComponentIdentifier
					if err := json.Unmarshal([]byte(val), &c); err != nil {
						panic(err)
					}
					query = query.ComponentIdentifier(c)
			*/
		}
	}

	components, err := nexusiq.SearchComponents(demo.IQ(idx), query)
	if err != nil {
		log.Fatalf("Did not complete search: %v", err)
	}

	fmt.Printf("%q\n", components)
}

func retentionList(idx int, org string) {
	policies, err := nexusiq.GetRetentionPolicies(demo.IQ(idx), org)
	if err != nil {
		panic(err)
	}

	buf, err := json.MarshalIndent(policies, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(buf))
}

func rolesList(idx int, app, org string) {
	var mappings []nexusiq.MemberMapping
	if app != "" && org == "" {
		panic("TODO")
	}

	if org != "" {
		var err error
		mappings, err = nexusiq.OrganizationAuthorizations(demo.IQ(idx), org)
		if err != nil {
			panic(err)
		}
	}

	buf, err := json.MarshalIndent(mappings, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(buf))
}

func iqVulnInfo(idx int, ids ...string) {
	iq := demo.IQ(idx)

	type catcher struct {
		id  string
		err error
	}

	errs := make([]catcher, 0)
	for _, id := range ids {
		info, err := privateiq.VulnerabilityInfoHTML(iq, id)
		if err != nil {
			errs = append(errs, catcher{id, err})
			continue
		}
		fmt.Println(info)
	}

	for _, e := range errs {
		log.Printf("error with %s: %v\n", e.id, e.err)
	}
}

func iqComponentDetails(idx int, format string, ids ...string) {
	iq := demo.IQ(idx)

	type catcher struct {
		id  string
		err error
	}

	errs := make([]catcher, 0)
	for _, id := range ids {
		c, err := nexusiq.NewComponentFromString(id)
		var infos []nexusiq.ComponentDetail
		if err == nil {
			infos, err = nexusiq.GetComponent(iq, []nexusiq.Component{*c})
		}
		if err != nil {
			errs = append(errs, catcher{id, err})
			continue
		}

		if format != "" {
			tmpl := template.Must(template.New("deets").Funcs(template.FuncMap{"json": tmplJSONPretty}).Parse(format))
			tmpl.Execute(os.Stdout, remediation)
		} else {
			buf, err := json.MarshalIndent(infos[0], "", "  ")
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(string(buf))
		}
	}

	for _, e := range errs {
		log.Printf("error with %s: %v\n", e.id, e.err)
	}
}
