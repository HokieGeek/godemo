package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"

	demo "github.com/hokiegeek/godemo"
	nexusrm "github.com/sonatype-nexus-community/gonexus/rm"
	nexuscli "github.com/sonatype-nexus-community/nexus-cli/cmd"
)

var (
	rmCommand        *cobra.Command
	rmIdx            int
	rmServer, rmAuth string
)

func createRmCommand() *cobra.Command {
	c := nexuscli.RmCommand
	c.Aliases = []string{"r"}
	c.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		splitAuth := strings.Split(rmAuth, ":")
		if rmServer != "" && len(splitAuth) == 2 {
			log.Printf("Connecting to %s\n", rmServer)
			demo.RMs = []demo.IdentifiedRM{demo.NewIdentifiedRM(rmServer, splitAuth[0], splitAuth[1])}
		}
	}

	return c
}

func init() {
	rmCommand = createRmCommand()

	rmCommand.AddCommand(&cobra.Command{
		Use:   "get",
		Short: "(beta) perform an http GET",
		Run: func(cmd *cobra.Command, args []string) {
			rmGet(rmIdx, args[0])
		},
	})

	rmCommand.AddCommand(&cobra.Command{
		Use:   "del",
		Short: "(beta) perform an http DELETE",
		Run: func(cmd *cobra.Command, args []string) {
			rmDel(rmIdx, args[0])
		},
	})

	rmCommand.AddCommand(&cobra.Command{
		Use:     "repos",
		Aliases: []string{"r"},
		Short:   "(beta) lists all repos",
		Run: func(cmd *cobra.Command, args []string) {
			rmListRepos(rmIdx)
		},
	})

	rmCommand.AddCommand(&cobra.Command{
		Use:     "ls",
		Aliases: []string{"l"},
		Short:   "(beta) lists all components in a repo",
		Run: func(cmd *cobra.Command, args []string) {
			rmListRepoComponents(rmIdx, args)
		},
	})

	rmCommand.AddCommand(
		func() *cobra.Command {
			var repo, coord, file string

			c := &cobra.Command{
				Use:     "up",
				Aliases: []string{"u"},
				Short:   "(beta) upload component",
				Run: func(cmd *cobra.Command, args []string) {
					rmUploadComponent(rmIdx, repo, coord, file)
				},
			}

			c.Flags().StringVarP(&repo, "repo", "r", "", "")
			c.Flags().StringVarP(&coord, "coord", "c", "", "")
			c.Flags().StringVarP(&file, "file", "f", "", "")

			return c
		}(),
	)

	rmCommand.AddCommand(
		func() *cobra.Command {
			c := &cobra.Command{
				Use:   "ro",
				Short: "(beta) read-only mode functions",
				Run: func(cmd *cobra.Command, args []string) {
					rmReadOnlyToggle(rmIdx)
					rmStatus(rmIdx)
				},
			}

			c.AddCommand(&cobra.Command{
				Use:     "enable",
				Aliases: []string{"e"},
				Short:   "enables read-only mode",
				Run: func(cmd *cobra.Command, args []string) {
					rmReadOnlyToggle(rmIdx)
					rmStatus(rmIdx)
				},
			})

			c.AddCommand(
				func() *cobra.Command {
					var force bool

					c := &cobra.Command{
						Use:     "release",
						Aliases: []string{"r"},
						Short:   "releases from read-only mode",
						Run: func(cmd *cobra.Command, args []string) {
							rmReadOnly(rmIdx, false, force)
							rmStatus(rmIdx)
						},
					}

					c.Flags().BoolVarP(&force, "force", "f", false, "")

					return c
				}(),
			)

			return c
		}(),
	)

	rmCommand.AddCommand(&cobra.Command{
		Use:   "zip",
		Short: "(beta) get support zip",
		Run: func(cmd *cobra.Command, args []string) {
			rmZip(rmIdx)
		},
	})
}

func rmGet(idx int, endpoint string) {
	body, resp, err := demo.RM(idx).Get(endpoint)
	if err != nil {
		panic(err)
	}

	log.Println(resp)
	fmt.Println(body)
}

func rmDel(idx int, endpoint string) {
	if _, err := demo.RM(idx).Del(endpoint); err != nil {
		panic(err)
	}
}

func rmListRepos(idx int) {
	w := csv.NewWriter(os.Stdout)

	w.Write([]string{"Name", "Format", "Type", "URL"})
	if repos, err := nexusrm.GetRepositories(demo.RM(idx)); err == nil {
		for _, r := range repos {
			w.Write([]string{r.Name, r.Format, r.Type, r.URL})
		}
	}

	w.Flush()

	if err := w.Error(); err != nil {
		panic(err)
	}
}

func rmListRepoComponents(idx int, repos []string) {
	w := csv.NewWriter(os.Stdout)

	w.Write([]string{"Repository", "Group", "Name", "Version", "Tags"})
	if len(repos) == 0 {
		all, _ := nexusrm.GetRepositories(demo.RM(idx))
		for _, r := range all {
			repos = append(repos, r.Name)
		}
	}

	for _, repo := range repos {
		if components, err := nexusrm.GetComponents(demo.RM(idx), repo); err == nil {
			for _, c := range components {
				w.Write([]string{c.Repository, c.Group, c.Name, c.Version, strings.Join(c.Tags, ";")})
			}
		}
	}

	w.Flush()

	if err := w.Error(); err != nil {
		panic(err)
	}
}

func rmUploadComponent(idx int, repo, coord, filePath string) {
	fmt.Println("Uploading component...")
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	// FIXME: something other than maven, please
	// if len(coord) > 0 {
	upload, err := nexusrm.NewUploadComponentMaven(coord, file)
	if err != nil {
		panic(err)
	}
	// }

	if err = nexusrm.UploadComponent(demo.RM(idx), repo, upload); err != nil {
		panic(err)
	}

	fmt.Println("Success!")
}

func rmStatus(idx int) {
	fmt.Println(demo.RM(idx).Info().Host)
	fmt.Printf("Readable: %v\n", nexusrm.StatusReadable(demo.RM(idx)))
	fmt.Printf("Writable: %v\n", nexusrm.StatusWritable(demo.RM(idx)))

	state, _ := nexusrm.GetReadOnlyState(demo.RM(idx))
	fmt.Println(state)
}

func rmZip(idx int) {
	zip, name, err := nexusrm.GetSupportZip(demo.RM(idx), nexusrm.NewSupportZipOptions())
	if err != nil {
		panic(err)
	}

	if err = ioutil.WriteFile(name, zip, 0644); err != nil {
		panic(err)
	}

	log.Printf("Created %s\n", name)
}

func rmReadOnly(idx int, enable, forceRelease bool) {
	if enable {
		nexusrm.ReadOnlyEnable(demo.RM(idx))
	} else {
		nexusrm.ReadOnlyRelease(demo.RM(idx), forceRelease)
	}
}

func rmReadOnlyToggle(idx int) {
	state, err := nexusrm.GetReadOnlyState(demo.RM(idx))
	if err != nil {
		return
	}
	if state.Frozen {
		rmReadOnly(idx, false, false)
	} else {
		rmReadOnly(idx, true, false)
	}
}
