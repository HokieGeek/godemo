package demo

import (
	privateiq "github.com/hokiegeek/gonexus-private/iq"
	"github.com/hokiegeek/gonexus/iq"
	"github.com/hokiegeek/gonexus/rm"
)

func RM() (*nexusrm.RM, error) {
	return nexusrm.New("http://localhost:8081", "admin", "admin123")
}

func IQ() (*nexusiq.IQ, error) {
	return nexusiq.New("http://localhost:8070", "admin", "admin123")
}

func IQp() (*privateiq.IQ, error) {
	return privateiq.New("http://localhost:8070", "admin", "admin123")
}

func IQ2p(iq *nexusiq.IQ) (*privateiq.IQ, error) {
	return privateiq.FromPublic(iq)
}

func Repos() ([]nexusrm.Repository, error) {
	rm, _ := RM()
	return nexusrm.GetRepositories(rm)
}

func Apps() ([]nexusiq.ApplicationDetails, error) {
	iq, _ := IQ()
	return nexusiq.GetAllApplications(iq)
}
