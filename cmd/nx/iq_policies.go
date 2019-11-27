package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	demo "github.com/hokiegeek/godemo"
	privateiq "github.com/hokiegeek/gonexus-private/iq"
	nexusiq "github.com/sonatype-nexus-community/gonexus/iq"
)

var iqPoliciesCommand = func() *cobra.Command {
	c := &cobra.Command{
		Use:     "policies",
		Short:   "(beta) Do stuff with policies",
		Aliases: []string{"pol", "p"},
	}

	c.AddCommand(func() *cobra.Command {
		c := &cobra.Command{
			Use:     "import",
			Short:   "Import the indicated policies",
			Aliases: []string{"i"},
			Run: func(cmd *cobra.Command, args []string) {
				importPolicies(iqIdx, args[0])
			},
		}

		return c
	}())

	c.AddCommand(func() *cobra.Command {
		c := &cobra.Command{
			Use:     "export",
			Aliases: []string{"a"},
			Short:   "exports the policies of the indicated IQ",
			Run: func(cmd *cobra.Command, args []string) {
				exportPolicies(iqIdx)
			},
		}

		return c
	}())

	c.AddCommand(func() *cobra.Command {
		c := &cobra.Command{
			Use:     "list",
			Aliases: []string{"ls, l"},
			Short:   "Lists all policies configured on the instance",
			Run: func(cmd *cobra.Command, args []string) {
				listPolicies(iqIdx)
			},
		}

		return c
	}())

	return c
}()

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
