package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/hokiegeek/godemo"
	"github.com/sonatype-nexus-community/gonexus/rm"
	"github.com/urfave/cli"
)

var rmCommand = cli.Command{
	Name:    "rm",
	Aliases: []string{"r"},
	Usage:   "repository-specific commands",
	Flags: []cli.Flag{
		cli.IntFlag{
			Name:  "idx, i",
			Value: 0,
			Usage: "rm `idx`",
		},
	},
	Action: func(c *cli.Context) error {
		rmStatus(c.Int("idx"))
		return nil
	},
	Subcommands: []cli.Command{
		{
			Name:  "get",
			Usage: "perform an http get",
			Action: func(c *cli.Context) error {
				rmGet(c.Parent().Int("idx"), c.Args().First())
				return nil
			},
		},
		{
			Name:    "repos",
			Aliases: []string{"r"},
			Usage:   "lists all repos",
			Action: func(c *cli.Context) error {
				rmListRepos(c.Parent().Int("idx"))
				return nil
			},
		},
		{
			Name:    "ls",
			Aliases: []string{"l"},
			Usage:   "lists all components in a repo",
			Action: func(c *cli.Context) error {
				rmListRepoComponents(c.Parent().Int("idx"), c.Args())
				return nil
			},
		},
		{
			Name:    "upload",
			Aliases: []string{"u", "up"},
			Usage:   "upload component",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "repo, r"},
				cli.StringFlag{Name: "coord, c"},
				cli.StringFlag{Name: "file, f"},
			},
			Action: func(c *cli.Context) error {
				rmUploadComponent(c.Parent().Int("idx"), c.String("repo"), c.String("coord"), c.String("file"))
				return nil
			},
		},
		{
			Name:  "ro",
			Usage: "read-only mode functions",
			Action: func(c *cli.Context) error {
				demo.RmReadOnlyToggle(c.Parent().Int("idx"))
				rmStatus(c.Parent().Int("idx"))
				return nil
			},
			Subcommands: []cli.Command{
				{
					Name:    "enable",
					Aliases: []string{"e"},
					Usage:   "enables read-only mode",
					Action: func(c *cli.Context) error {
						demo.RmReadOnly(c.Parent().Parent().Int("idx"), true, false)
						rmStatus(c.Parent().Int("idx"))
						return nil
					},
				},
				{
					Name:    "release",
					Aliases: []string{"r"},
					Usage:   "releases from read-only mode",
					Flags: []cli.Flag{
						cli.BoolFlag{Name: "force, f"},
					},
					Action: func(c *cli.Context) error {
						demo.RmReadOnly(c.Parent().Parent().Int("idx"), false, c.Bool("force"))
						rmStatus(c.Parent().Int("idx"))
						return nil
					},
				},
			},
		},
	},
}

func rmGet(idx int, endpoint string) {
	body, resp, err := demo.RM(idx).Get(endpoint)
	if err != nil {
		panic(err)
	}

	log.Println(resp)
	fmt.Println(body)
}

func rmListRepos(idx int) {
	format := "%s, %s, %s, %s\n"
	fmt.Printf(format, "Name", "Format", "Type", "URL")
	if repos, err := demo.Repos(idx); err == nil {
		for _, r := range repos {
			fmt.Printf(format, r.Name, r.Format, r.Type, r.URL)
		}
	}
}

func rmListRepoComponents(idx int, repos []string) {
	format := "%s, %s, %s, %s, %s\n"
	fmt.Printf(format, "Repository", "Group", "Name", "Version", "Tags")

	if len(repos) == 0 {
		all, _ := demo.Repos(idx)
		for _, r := range all {
			repos = append(repos, r.Name)
		}
	}

	for _, repo := range repos {
		if components, err := demo.Components(idx, repo); err == nil {
			for _, c := range components {
				fmt.Printf(format, c.Repository, c.Group, c.Name, c.Version, strings.Join(c.Tags, ";"))
			}
		}
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
