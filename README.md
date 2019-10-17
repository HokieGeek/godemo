# godemo [![DepShield Badge](https://depshield.sonatype.org/badges/HokieGeek/godemo/depshield.svg)](https://depshield.github.io) [![CircleCI](https://circleci.com/gh/HokieGeek/godemo.svg?style=svg)](https://circleci.com/gh/HokieGeek/godemo)

Silly little time-saver that takes [github.com/sonatype-nexus-community/gonexus](//github.com/sonatype-nexus-community/gonexus) and provides some shortcuts to it with an emphasis on my demo instances on localhost with standard ports and logins.

## Using library

An example program:

```go
package main

import (
	"fmt"
	"github.com/hokiegeek/godemo"
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

	// Print all repositories in first detected RM
	fmt.Println("[RM Repos]")
	repos, _ := demo.Repos(0)
	for _, repo := range repos {
		fmt.Printf("%s (%s : %s)\n", repo.Name, repo.Format, repo.Type)
	}

	// Print all applications in first detected IQ
	fmt.Println("[IQ Apps]")
	apps, _ := demo.Apps(0)
	for _, app := range apps {
		fmt.Println(app.Name)
	}
}
```

## Using CLI

A few example commands:

| Description                                                 | Example                                                                    |
| ----------------------------------------------------------- | -------------------------------------------------------------------------- |
| List all Nexus servers on the local machine                 | `nx ls`                                                                    |
| List all component in the listed repositories               | `nx rm ls maven-releases npm-proxy golang-group`                           |
| List all applications in an IQ instance                     | `nx iq app`                                                                |
| Evaluate the indicated components against Root Organization | `nx iq eval "maven:jackson-databind:com.fasterxml.jackson.core:2.6.1:jar"` |
