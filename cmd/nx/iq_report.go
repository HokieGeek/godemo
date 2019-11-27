package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/spf13/cobra"

	demo "github.com/hokiegeek/godemo"
	privateiq "github.com/hokiegeek/gonexus-private/iq"
	nexusiq "github.com/sonatype-nexus-community/gonexus/iq"
	nexuscli "github.com/sonatype-nexus-community/nexus-cli/cmd"
)

var iqReportCommand = func() *cobra.Command {
	var format string

	c := &cobra.Command{
		Use:     "report",
		Short:   "(beta) Manage IQ app reports",
		Aliases: []string{"r"},
		Run: func(cmd *cobra.Command, args []string) {
			appReport(iqIdx, format, args...)
		},
	}

	c.Flags().StringVarP(&format, "format", "f", "", "")

	c.AddCommand(func() *cobra.Command {
		c := &cobra.Command{
			Use:     "reevaluate",
			Aliases: []string{"rv"},
			Short:   "reevaluate [appID:stage] [appID:stage] [appID]",
			Run: func(cmd *cobra.Command, args []string) {
				reportReevaluate(iqIdx, args...)
			},
		}

		return c
	}())

	c.AddCommand(func() *cobra.Command {
		c := &cobra.Command{
			Use:     "diff",
			Aliases: []string{"d"},
			Short:   "diff [appID] [report1ID] [report2ID]",
			Run: func(cmd *cobra.Command, args []string) {
				reportDiff(iqIdx, args[0], args[1], args[2])
			},
		}

		return c
	}())

	return c
}()

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
			tmpl := template.Must(template.New("report").Funcs(template.FuncMap{"json": nexuscli.TemplateJSONPretty}).Parse(format))
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

func reportReevaluate(idx int, apps ...string) {
	if len(apps) == 0 {
		if err := privateiq.ReevaluateAllReports(demo.IQ(idx)); err != nil {
			log.Printf("could not re-evaluate reports: %v", err)
		}
		return
	}

	for _, app := range apps {
		splitPos := strings.LastIndex(app, ":")
		appID := app[:splitPos]
		stage := app[splitPos+1:]

		if err := privateiq.ReevaluateReport(demo.IQ(idx), appID, stage); err != nil {
			log.Printf("could not re-evaluate report for '%s' at '%s' stage: %v", appID, stage, err)
		}
	}
}

func reportDiff(idx int, appID, report1ID, report2ID string) {
	diff, err := nexusiq.ReportsDiff(demo.IQ(idx), appID, report1ID, report2ID)
	if err != nil {
		panic(err)
	}

	buf, err := json.MarshalIndent(diff, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(buf))
}
