package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"

	demo "github.com/hokiegeek/godemo"
	privateiq "github.com/hokiegeek/gonexus-private/iq"
	nexusiq "github.com/sonatype-nexus-community/gonexus/iq"
	nexuscli "github.com/sonatype-nexus-community/nexus-cli/cmd"
)

var (
	iqCommand        *cobra.Command
	iqIdx            int
	iqServer, iqAuth string
)

func createIqCommand() *cobra.Command {
	c := nexuscli.IqCommand
	c.Aliases = []string{"q"}
	c.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		splitAuth := strings.Split(iqAuth, ":")
		if iqServer != "" && len(splitAuth) == 2 {
			log.Printf("Connecting to %s\n", iqServer)
			demo.IQs = []demo.IdentifiedIQ{demo.NewIdentifiedIQ(iqServer, splitAuth[0], splitAuth[1])}
		}
	}

	return c
}

func init() {
	iqCommand = createIqCommand()

	iqCommand.AddCommand(iqScCommand)
	iqCommand.AddCommand(iqPoliciesCommand)
	iqCommand.AddCommand(iqWaiversCommand)
	iqCommand.AddCommand(iqReportCommand)

	iqCommand.AddCommand(&cobra.Command{
		Use:   "violations",
		Short: "(beta) List violations by policy name",
		Run: func(cmd *cobra.Command, args []string) {
			listViolatingApps(iqIdx, args...)
		},
	})

	iqCommand.AddCommand(&cobra.Command{
		Use:   "vulns",
		Short: "(beta) List vulnerability info",
		Run: func(cmd *cobra.Command, args []string) {
			iqVulnInfo(iqIdx, args...)
		},
	})

	iqCommand.AddCommand(&cobra.Command{
		Use:     "search",
		Short:   "(beta) Search for component",
		Aliases: []string{"q"},
		Run: func(cmd *cobra.Command, args []string) {
			iqSearch(iqIdx, args)
		},
	})

	iqCommand.AddCommand(&cobra.Command{
		Use:   "zip",
		Short: "(beta) get support zip",
		Run: func(cmd *cobra.Command, args []string) {
			iqZip(iqIdx)
		},
	})

	iqCommand.AddCommand(
		func() *cobra.Command {
			var secret, url, events string

			c := &cobra.Command{
				Use:     "webhook",
				Short:   "(beta) create a new webhook",
				Aliases: []string{"wh"},
				Run: func(cmd *cobra.Command, args []string) {
					webhook(iqIdx, url, secret, events)
				},
			}

			c.Flags().StringVarP(&secret, "secret", "s", "", "")
			c.Flags().StringVarP(&url, "url", "u", "", "")
			c.Flags().StringVarP(&events, "events", "e", "Application Evaluation", "")

			c.MarkFlagRequired("secret")
			c.MarkFlagRequired("url")

			return c
		}(),
	)

	iqCommand.AddCommand(
		func() *cobra.Command {
			var disable bool

			c := &cobra.Command{
				Use:     "auto-apps",
				Short:   "(beta) Manage auto applications",
				Aliases: []string{"auto"},
				Run: func(cmd *cobra.Command, args []string) {
					autoApps(iqIdx, disable, args[0])
				},
			}

			c.Flags().BoolVarP(&disable, "disable", "d", false, "")

			return c
		}(),
	)

	iqCommand.AddCommand(
		func() *cobra.Command {
			var disable bool

			c := &cobra.Command{
				Use:     "notice",
				Short:   "(beta) Set a message in IQ",
				Aliases: []string{"msg"},
				Run: func(cmd *cobra.Command, args []string) {
					systemNotice(iqIdx, disable, strings.Join(args, " "))
				},
			}

			c.Flags().BoolVarP(&disable, "disable", "d", false, "")

			return c
		}(),
	)

	iqCommand.AddCommand(
		func() *cobra.Command {
			c := &cobra.Command{
				Use:   "retention",
				Short: "(beta) Manage retention policies",
				Run: func(cmd *cobra.Command, args []string) {
					retentionList(iqIdx, args[0])
				},
			}

			c.AddCommand(&cobra.Command{
				Use:     "list",
				Short:   "(beta) List retention policies",
				Aliases: []string{"ls"},
				Run: func(cmd *cobra.Command, args []string) {
					retentionList(iqIdx, args[0])
				},
			})

			return c
		}(),
	)

	iqCommand.AddCommand(
		func() *cobra.Command {
			var app, org string

			c := &cobra.Command{
				Use:   "role",
				Short: "(beta) List roles",
				Run: func(cmd *cobra.Command, args []string) {
					rolesList(iqIdx, app, org)
				},
			}

			c.Flags().StringVarP(&app, "app", "a", "", "")
			c.Flags().StringVarP(&org, "org", "o", "", "")

			/*
				c.AddCommand(&cobra.Command{
						Use:    "user",
						Aliases: []string{"user"},
					Run: func(cmd *cobra.Command, args []string) {
						// TODO
					},
				})
			*/

			return c
		}(),
	)
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
