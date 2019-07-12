package demo

import (
	"github.com/hokiegeek/gonexus/iq"
	"github.com/hokiegeek/gonexus/rm"
)

// RM returns an instance of Repository Manager with demo defaults
func RM(idx int) nexusrm.RM {
	return RMs[idx].Login("admin", "admin123").Client()
}

// IQ returns an instance of IQ Server with demo defaults
func IQ(idx int) nexusiq.IQ {
	return IQs[idx].Login("admin", "admin123").Client()
}

// Repos returns a list of all of the repositories in the demo RM
func Repos(idx int) ([]nexusrm.Repository, error) {
	return nexusrm.GetRepositories(RM(idx))
}

// Apps returns a list of all of the applications in the demo IQ
func Apps(idx int) ([]nexusiq.Application, error) {
	return nexusiq.GetAllApplications(IQ(idx))
}
