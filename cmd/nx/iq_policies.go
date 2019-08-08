package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/urfave/cli"

	demo "github.com/hokiegeek/godemo"
	privateiq "github.com/hokiegeek/gonexus-private/iq"
	nexusiq "github.com/sonatype-nexus-community/gonexus/iq"
)

var iqPoliciesCommand = cli.Command{
	Name:    "policies",
	Aliases: []string{"pol", "p"},
	Usage:   "Do stuff with policies",
	Subcommands: []cli.Command{
		{
			Name:    "import",
			Aliases: []string{"i"},
			Usage:   "Import the indicated policies",
			Action: func(c *cli.Context) error {
				importPolicies(c.Parent().Int("idx"), c.Args().First())
				return nil
			},
		},
		{
			Name:    "export",
			Aliases: []string{"a"},
			Usage:   "exports the policies of the indicated IQ",
			Action: func(c *cli.Context) error {
				exportPolicies(c.Parent().Int("idx"))
				return nil
			},
		},
		{
			Name:    "list",
			Aliases: []string{"ls, l"},
			Usage:   "Lists all policies configured on the instance",
			Action: func(c *cli.Context) error {
				listPolicies(c.Parent().Int("idx"))
				return nil
			},
		},
	},
}

func importPolicies(idx int, filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	err = privateiq.ImportPolicies(demo.IQ(idx), file)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Policies imported")
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

func listPolicies(idx int) {
	w := csv.NewWriter(os.Stdout)

	w.Write([]string{"Name", "PolicyType", "ThreatLevel", "OwnerID", "OwnerType", "ID"})
	if policies, err := nexusiq.GetPolicies(demo.IQ(idx)); err == nil {
		for _, p := range policies {
			w.Write([]string{p.Name, p.PolicyType, strconv.Itoa(p.ThreatLevel), p.OwnerID, p.OwnerType, p.ID})
		}
	}

	w.Flush()

	if err := w.Error(); err != nil {
		panic(err)
	}
}
