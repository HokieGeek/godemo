package main

import (
	"encoding/json"
	"fmt"
	"os"

	nexuscli "github.com/sonatype-nexus-community/nexus-cli/cmd"

	"github.com/spf13/cobra"

	demo "github.com/hokiegeek/godemo"
)

var templateJSONPretty = func(v interface{}) string {
	a, _ := json.MarshalIndent(v, "", "  ")
	return string(a)
}

func listServers() {
	demo.Detect()

	for i, s := range demo.RMs {
		fmt.Printf("RM[%d]: %s\n", i, s.Host)
	}

	for i, s := range demo.IQs {
		fmt.Printf("IQ[%d]: %s\n", i, s.Host)
	}
}

// RootCmd TODO
// var RootCmd = &cobra.Command{}
var RootCmd = nexuscli.RootCmd

func main() {
	RootCmd.Use = "nx"
	RootCmd.Run = func(cmd *cobra.Command, args []string) {
		listServers()
	}

	RootCmd.AddCommand(&cobra.Command{
		Use:   "ls",
		Short: "List all nexus servers",
		Run: func(cmd *cobra.Command, args []string) {
			listServers()
		},
	})

	// RootCmd.AddCommand(&cobra.Command{

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
