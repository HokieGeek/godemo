# godemo

Silly little time-saver that takes [github.com/sonatype-nexus-community/gonexus](//github.com/sonatype-nexus-community/gonexus) and provides some shortcuts to it with an emphasis on my demo instances on localhost with standard ports and logins.

An example program:
```go
package main

import (
	"demo"
	"fmt"
)

func main() {
	// Identify all Repository Manager servers on the local machine
	fmt.Println("[RM Servers]")
	for _, s := range demo.DetectRMServers() {
		fmt.Println(s.Host)
	}

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
}
```
