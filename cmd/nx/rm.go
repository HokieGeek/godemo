package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	demo "github.com/hokiegeek/godemo"
	nexusrm "github.com/sonatype-nexus-community/gonexus/rm"
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
		cli.StringFlag{
			Name:  "server, s",
			Value: "http://localhost:8081",
		},
		cli.StringFlag{
			Name:  "auth, a",
			Value: "admin:admin123",
		},
	},
	Before: func(c *cli.Context) error {
		host := c.String("server")
		auth := strings.Split(c.String("auth"), ":")
		if host != "" && len(auth) == 2 {
			log.Printf("Connecting to %s\n", host)
			demo.RMs = []demo.IdentifiedRM{demo.NewIdentifiedRM(host, auth[0], auth[1])}
		}
		return nil
	},
	Action: func(c *cli.Context) error {
		rmStatus(c.Int("idx"))
		return nil
	},
	Subcommands: []cli.Command{
		{
			Name:  "get",
			Usage: "perform an http GET",
			Action: func(c *cli.Context) error {
				rmGet(c.Parent().Int("idx"), c.Args().First())
				return nil
			},
		},
		{
			Name:  "del",
			Usage: "perform an http DELETE",
			Action: func(c *cli.Context) error {
				rmDel(c.Parent().Int("idx"), c.Args().First())
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
			Name:    "up",
			Aliases: []string{"u"},
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
				rmReadOnlyToggle(c.Parent().Int("idx"))
				rmStatus(c.Parent().Int("idx"))
				return nil
			},
			Subcommands: []cli.Command{
				{
					Name:    "enable",
					Aliases: []string{"e"},
					Usage:   "enables read-only mode",
					Action: func(c *cli.Context) error {
						rmReadOnly(c.Parent().Parent().Int("idx"), true, false)
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
						rmReadOnly(c.Parent().Parent().Int("idx"), false, c.Bool("force"))
						rmStatus(c.Parent().Int("idx"))
						return nil
					},
				},
			},
		},
		{
			Name:  "zip",
			Usage: "get support zip",
			Action: func(c *cli.Context) error {
				rmZip(c.Parent().Int("idx"))
				return nil
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
