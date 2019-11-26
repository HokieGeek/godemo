package main

import (
	"fmt"

	"github.com/spf13/cobra"

	demo "github.com/hokiegeek/godemo"
	nexusiq "github.com/sonatype-nexus-community/gonexus/iq"
)

var iqScCommand = func() *cobra.Command {
	c := &cobra.Command{
		Use:   "sc",
		Short: "source control actions",
	}

	c.AddCommand(scCreate())
	c.AddCommand(scDelete())
	c.AddCommand(scList())

	return c
}()

func scCreate() *cobra.Command {
	return &cobra.Command{
		Use:     "create",
		Aliases: []string{"c"},
		Short:   "create <appId> <repositoryUrl> <accessToken>",
		Run: func(cmd *cobra.Command, args []string) {
			app, repo, token := args[0], args[1], args[2]

			// app, repo, token := args[0], args[1], args[2]
			iq := demo.IQ(iqIdx)
			err := nexusiq.CreateSourceControlEntry(iq, app, repo, token)
			if err != nil {
				panic(err)
			}

			entry, err := nexusiq.GetSourceControlEntry(iq, app)
			if err != nil {
				panic(err)
			}

			fmt.Printf("%q\n", entry)
		},
	}
}

func scDelete() *cobra.Command {
	var appID, entryID string

	c := &cobra.Command{
		Use:     "delete",
		Aliases: []string{"d"},
		Short:   "deletes a source control entry",
		Run: func(cmd *cobra.Command, args []string) {
			iq := demo.IQ(iqIdx)
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
		},
	}

	c.Flags().StringVarP(&appID, "app", "a", "", "")
	c.Flags().StringVarP(&entryID, "id", "i", "", "")

	return c
}

func scList() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Aliases: []string{"l"},
		Short:   "deletes a source control entry",
		Run: func(cmd *cobra.Command, args []string) {
			appID := args[0]
			iq := demo.IQ(iqIdx)
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
		},
	}
}
