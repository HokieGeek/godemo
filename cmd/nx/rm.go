package main

import (
	"fmt"
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
	/*
		Action: func(c *cli.Context) error {
			TODO: list server info
			fmt.Println(c.Int("idx"))
			listRepos(0)
			return nil
		},
	*/
	Subcommands: []cli.Command{
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
	},
}

func rmListRepos(idx int) {
	format := "%s, %s, %s, %s\n"
	fmt.Printf(format, "Name", "Format", "Type", "URL")
	// if repos, err := demo.Repos(idx); err == nil {
	if repos, err := nexusrm.GetRepositories(demo.RM(idx)); err == nil {
		for _, r := range repos {
			fmt.Printf(format, r.Name, r.Format, r.Type, r.URL)
		}
	}
}

func rmListRepoComponents(idx int, repos []string) {
	format := "%s, %s, %s, %s, %s\n"
	fmt.Printf(format, "Repository", "Group", "Name", "Version", "Tags")

	if len(repos) == 0 {
		all, _ := nexusrm.GetRepositories(demo.RM(idx))
		for _, r := range all {
			repos = append(repos, r.Name)
		}
	}

	for _, repo := range repos {
		// if components, err := demo.Components(idx, repo); err == nil {
		if components, err := nexusrm.GetComponents(demo.RM(idx), repo); err == nil {
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

	upload, err := nexusrm.NewUploadComponentMaven(coord, file)
	if err != nil {
		panic(err)
	}

	if err = nexusrm.UploadComponent(demo.RM(idx), repo, upload); err != nil {
		panic(err)
	}

	fmt.Println("Success!")
}
