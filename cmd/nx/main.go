package main

import (
	"fmt"

	// "github.com/hokiegeek/gonexus-private/iq"
	"github.com/hokiegeek/godemo"
)

func main() {
	// Identify all Repository Manager servers on the local machine
	fmt.Println("[RM Servers]")
	for _, s := range demo.DetectRMServers() {
		fmt.Println(s.Host)
	}

	/*
		// Identify all IQ Servers on the local machine
		fmt.Println("[IQ Servers]")
		for _, s := range demo.DetectIQServers() {
			fmt.Println(s.Host)
		}

		// Print all repositories in RM on localhost:8081
		fmt.Println("[RM Repos]")
		repos, _ := demo.Repos()
		for _, repo := range repos {
			fmt.Printf("%s (%s : %s)\n", repo.Name, repo.Format, repo.Type)
		}

		// Print all applications in IQ on localhost:8070
		fmt.Println("[IQ Apps]")
		apps, _ := demo.Apps()
		for _, app := range apps {
			fmt.Println(app.Name)
		}
	*/
	/*
		iq, err := nexusiq.New("http://localhost:8070", "admin", "admin123")
		if err != nil {
			panic(err)
		}

		orgID, err := nexusiq.CreateOrganization(iq, "arstarst")
		if err != nil {
			panic(err)
		}

		time.Sleep(15 * time.Second)

		if err := privateiq.DeleteOrganization(iq, orgID); err != nil {
			panic(err)
		}
	*/

}
