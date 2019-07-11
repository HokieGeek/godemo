package demo

import (
	"github.com/hokiegeek/gonexus/iq"
	"github.com/hokiegeek/gonexus/rm"
)

// RM returns an instance of Repository Manager with demo defaults
func RM() (nexusrm.RM, error) {
	return nexusrm.New("http://localhost:8081", "admin", "admin123")
}

// IQ returns an instance of IQ Server with demo defaults
func IQ() (nexusiq.IQ, error) {
	return nexusiq.New("http://localhost:8070", "admin", "admin123")
}

// Repos returns a list of all of the repositories in the demo RM
func Repos() ([]nexusrm.Repository, error) {
	rm, _ := RM()
	return nexusrm.GetRepositories(rm)
}

// Apps returns a list of all of the applications in the demo IQ
func Apps() ([]nexusiq.Application, error) {
	iq, _ := IQ()
	return nexusiq.GetAllApplications(iq)
}
