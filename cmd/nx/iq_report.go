package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/urfave/cli"

	demo "github.com/hokiegeek/godemo"
	privateiq "github.com/hokiegeek/gonexus-private/iq"
	nexusiq "github.com/sonatype-nexus-community/gonexus/iq"
)

var iqReportCommand = cli.Command{
	Name:    "report",
	Aliases: []string{"r"},
	Flags: []cli.Flag{
		cli.StringFlag{Name: "format, f"},
	},
	Action: func(c *cli.Context) error {
		appReport(c.Parent().Int("idx"), c.String("format"), c.Args()...)
		return nil
	},
	Subcommands: []cli.Command{
		{
			Name:    "reevaluate",
			Aliases: []string{"rv"},
			Usage:   "reevaluate [appID:stage] [appID:stage] [appID]",
			Action: func(c *cli.Context) error {
				reportReevaluate(c.Parent().Parent().Int("idx"), c.Args()...)
				return nil
			},
		},
		{
			Name:    "diff",
			Aliases: []string{"d"},
			Usage:   "diff [appID] [report1ID] [report2ID]",
			Action: func(c *cli.Context) error {
				reportDiff(c.Parent().Parent().Int("idx"), c.Args().Get(0), c.Args().Get(1), c.Args().Get(2))
				return nil
			},
		},
	},
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
