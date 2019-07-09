# godemo

Silly little time-saver that takes [github.com/hokiegeek/gonexus](//github.com/hokiegeek/gonexus) and provides some shortcuts to it with an emphasis on my demo instances on localhost with standard ports and logins.

An example program:
```go
package main

import (
	"demo"
	"fmt"
)

func main() {
	fmt.Println("[RM Servers]")
	for _, s := range demo.DetectRMServers() {
		fmt.Println(s.Host)
	}

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
