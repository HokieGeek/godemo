package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/urfave/cli"

	demo "github.com/hokiegeek/godemo"
	privateiq "github.com/hokiegeek/gonexus-private/iq"
)

var iqWaiversCommand = cli.Command{
	Name:    "waivers",
	Aliases: []string{"pol", "p"},
	Usage:   "Do stuff with policies",
	Subcommands: []cli.Command{
		{
			Name:    "list",
			Aliases: []string{"ls, l"},
			Flags: []cli.Flag{
				cli.StringFlag{Name: "format, f"},
			},
			Usage: "Lists all policies configured on the instance",
			Action: func(c *cli.Context) error {
				listWaivers(c.Parent().Int("idx"), c.String("format"), c.Args().First())
				return nil
			},
		},
	},
}

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
		tmpl := template.Must(template.New("waivers").Funcs(template.FuncMap{"json": tmplJSONPretty}).Parse(format))
		tmpl.Execute(os.Stdout, waivers)
	} else {
		json, err := json.MarshalIndent(waivers, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(json))
	}
}
