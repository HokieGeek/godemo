package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/urfave/cli"

	demo "github.com/hokiegeek/godemo"
	// privateiq "github.com/hokiegeek/gonexus-private/iq"
	nexusiq "github.com/sonatype-nexus-community/gonexus/iq"
)

var iqComponentCommand = cli.Command{
	Name: "component",
	Flags: []cli.Flag{
		cli.StringFlag{Name: "format, f"},
	},
	Action: func(c *cli.Context) error {
		iqComponentDetails(c.Parent().Int("idx"), c.String("format"), c.Args()...)
		return nil
	},
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
		var components []nexusiq.ComponentDetail
		if err == nil {
			components, err = nexusiq.GetComponents(iq, []nexusiq.Component{*c})
		}
		if err != nil {
			errs = append(errs, catcher{id, err})
			continue
		}

		if format != "" {
			tmpl := template.Must(template.New("deets").Funcs(template.FuncMap{"json": tmplJSONPretty}).Parse(format))
			tmpl.Execute(os.Stdout, components)
		} else {
			buf, err := json.MarshalIndent(components, "", "  ")
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

func iqAllComponents(idx int, format string) {
	iq := demo.IQ(idx)

	components, err := nexusiq.GetAllComponents(iq)
	if err != nil {
		log.Printf("error listing components: %v\n", err)
		return
	}

	if format != "" {
		tmpl := template.Must(template.New("deets").Funcs(template.FuncMap{"json": tmplJSONPretty}).Parse(format))
		tmpl.Execute(os.Stdout, components)
	} else {
		buf, err := json.MarshalIndent(components, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(buf))
	}
}
