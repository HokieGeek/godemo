# godemo

Silly little time-saver that takes [github.com/hokiegeek/gonexus](//github.com/hokiegeek/gonexus) and provides some shortcuts to it with an emphasis on my demo instances on localhost with standard ports and logins.

An example program:
```go
package main

import (
	"fmt"

	"demo"
)

func main() {
    // Print all repositories in RM on localhost:8081
	repos, _ := demo.Repos()
	fmt.Printf("%v\n", repos)

    // Print all applications in IQ on localhost:8070
	apps, _ := demo.Apps()
	fmt.Printf("%v\n", apps)

	// Create a new organization in IQ on localhost:8070
	iq, _ := demo.IQ()
	iq.CreateOrganization("foobar")
}
```
