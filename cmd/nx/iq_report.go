package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	demo "github.com/hokiegeek/godemo"
	privateiq "github.com/hokiegeek/gonexus-private/iq"
	nexusiq "github.com/sonatype-nexus-community/gonexus/iq"
	"github.com/spf13/cobra"
)

var iqReportCommand = func() *cobra.Command {
	var format string

	c := &cobra.Command{
		Use:     "report",
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
