package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/spf13/cobra"

	demo "github.com/hokiegeek/godemo"
	privateiq "github.com/hokiegeek/gonexus-private/iq"
	nexuscli "github.com/sonatype-nexus-community/nexus-cli/cmd"
)

var iqWaiversCommand = func() *cobra.Command {
	c := &cobra.Command{
		Use:     "waivers",
		Short:   "(beta) Do stuff with waivers",
		Aliases: []string{"waive", "w"},
	}

	c.AddCommand(func() *cobra.Command {
		var format string
		c := &cobra.Command{
			Use:     "list",
			Short:   "Lists all waivers configured on the instance",
			Aliases: []string{"ls, l"},
			Run: func(cmd *cobra.Command, args []string) {
				listWaivers(iqIdx, format, args[0])
			},
		}

		c.Flags().StringVarP(&format, "format", "f", "", "")
		return c
	}())

	return c
}()

func listWaivers(idx int, format, appID string) {
	var (
		waivers []privateiq.Waiver
		err     error
	)
	if appID != "" {
		waivers, err = privateiq.WaiversByAppID(demo.IQ(idx), appID)
	} else {
		waivers, err = privateiq.Waivers(demo.IQ(idx))
	}
	if err != nil {
		panic(err)
	}

	if format != "" {
		tmpl := template.Must(template.New("waivers").Funcs(template.FuncMap{"json": nexuscli.TemplateJSONPretty}).Parse(format))
		tmpl.Execute(os.Stdout, waivers)
	} else {
		json, err := json.MarshalIndent(waivers, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(json))
	}
}
