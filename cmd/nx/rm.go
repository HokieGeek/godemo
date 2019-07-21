package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/hokiegeek/godemo"
	"github.com/sonatype-nexus-community/gonexus/rm"
)

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
