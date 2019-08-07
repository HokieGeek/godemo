package main

import (
	"fmt"

	"github.com/urfave/cli"

	demo "github.com/hokiegeek/godemo"
	nexusiq "github.com/sonatype-nexus-community/gonexus/iq"
)

var iqScCommand = cli.Command{
	Name:  "sc",
	Usage: "source control actions",
	Subcommands: []cli.Command{
		scCreate(),
		scDelete(),
		scList(),
	},
}

func scCreate() cli.Command {
	return cli.Command{
		Name:    "create",
		Aliases: []string{"c"},
		Usage:   "create <appId> <repositoryUrl> <accessToken>",
		Action: func(c *cli.Context) error {
			idx := c.Parent().Parent().Int("idx")
			app, repo, token := c.Args()[0], c.Args()[1], c.Args()[2]

			// app, repo, token := args[0], args[1], args[2]
			iq := demo.IQ(idx)
			err := nexusiq.CreateSourceControlEntry(iq, app, repo, token)
			if err != nil {
				panic(err)
			}

			entry, err := nexusiq.GetSourceControlEntry(iq, app)
			if err != nil {
				panic(err)
			}

			fmt.Printf("%q\n", entry)
			return nil
		},
	}
}

func scDelete() cli.Command {
	return cli.Command{
		Name:    "delete",
		Aliases: []string{"d"},
		Flags: []cli.Flag{
			cli.StringFlag{Name: "app, a"},
			cli.StringFlag{Name: "id, i"},
		},
		Usage: "deletes a source control entry",
		Action: func(c *cli.Context) error {
			idx := c.Parent().Parent().Int("idx")
			appID, entryID := c.String("app"), c.String("id")

			iq := demo.IQ(idx)
			var scEntryID string
			if entryID != "" {
				scEntryID = entryID
			} else {
				scEntry, err := nexusiq.GetSourceControlEntry(iq, appID)
				if err != nil {
					panic(err)
				}
				scEntryID = scEntry.ID
			}

			nexusiq.DeleteSourceControlEntry(iq, appID, scEntryID)

			fmt.Println("Deleted")
			return nil
		},
	}
}

func scList() cli.Command {
	return cli.Command{
		Name:    "list",
		Aliases: []string{"l"},
		Usage:   "deletes a source control entry",
		Action: func(c *cli.Context) error {
			idx, appID := c.Parent().Parent().Int("idx"), c.Args().First()
			iq := demo.IQ(idx)
			if appID != "" {
				entry, _ := nexusiq.GetSourceControlEntry(iq, appID)
				fmt.Printf("%v\n", entry)
			} else {
				apps, err := nexusiq.GetAllApplications(iq)
				if err != nil {
					panic(err)
				}
				for _, app := range apps {
					if entry, err := nexusiq.GetSourceControlEntry(iq, app.PublicID); err == nil {
						fmt.Printf("%s: %v\n", app.PublicID, entry)
					}
				}
			}
			return nil
		},
	}
}
