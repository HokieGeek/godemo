package demo

import (
	"github.com/hokiegeek/gonexus-private/iq"
	"github.com/sonatype-nexus-community/gonexus/iq"
	"github.com/sonatype-nexus-community/gonexus/rm"
)

// RM returns an instance of Repository Manager with demo defaults
func RM(idx int) nexusrm.RM {
	if len(RMs) == 0 {
		RMs = DetectRMServers()
	}
	return RMs[idx].Client()
}

// IQ returns an instance of IQ Server with demo defaults
func IQ(idx int) nexusiq.IQ {
	return IQs[idx].Client()
}

// Repos returns a list of all of the repositories in the demo RM
func Repos(idx int) ([]nexusrm.Repository, error) {
	return nexusrm.GetRepositories(RM(idx))
}

// Components returns a list of all of the repositories in the demo RM
func Components(idx int, repo string) ([]nexusrm.RepositoryItem, error) {
	return nexusrm.GetComponents(RM(idx), repo)
}

// Apps returns a list of all of the applications in the demo IQ
func Apps(idx int) ([]nexusiq.Application, error) {
	return nexusiq.GetAllApplications(IQ(idx))
}

// Orgs returns a list of all of the organizations in the demo IQ
func Orgs(idx int) ([]nexusiq.Organization, error) {
	return nexusiq.GetAllOrganizations(IQ(idx))
}

// OrgsIDMap returns a map of organization ids by name and the reverse
func OrgsIDMap(idx int) (id2name map[string]string, name2id map[string]string, err error) {
	if orgs, err := Orgs(idx); err == nil {
		id2name = make(map[string]string)
		name2id = make(map[string]string)
		for _, o := range orgs {
			id2name[o.ID] = o.Name
			id2name[o.Name] = o.ID
		}
	}
	return
}

// IQComponentSliceFromStringSlice converts slice of : delimeted component names into a slice of IQ Component
func IQComponentSliceFromStringSlice(components []string) []nexusiq.Component {
	comps := make([]nexusiq.Component, len(components))
	for i, c := range components {
		comp, _ := nexusiq.NewComponentFromString(c)
		comps[i] = *comp
	}
	return comps
}

// Eval performs a Lifecycle evaluation of the indicated component
func Eval(idx int, appID string, components ...string) (report *nexusiq.Evaluation, err error) {
	c := IQComponentSliceFromStringSlice(components)

	if appID != "" {
		report, err = nexusiq.EvaluateComponents(IQ(idx), c, appID)
	} else {
		report, err = privateiq.EvaluateComponentsWithRootOrg(IQ(idx), c)
	}
	if err != nil {
		return
	}

	return
}

// RmReadOnly allows to control the read-only state of an RM instance
func RmReadOnly(idx int, enable, forceRelease bool) {
	if enable {
		nexusrm.ReadOnlyEnable(RM(idx))
	} else {
		nexusrm.ReadOnlyRelease(RM(idx), forceRelease)
	}
}

// RmReadOnlyToggle Toggles read-only mode
func RmReadOnlyToggle(idx int) {
	state, err := nexusrm.GetReadOnlyState(RM(idx))
	if err != nil {
		return
	}
	if state.Frozen {
		RmReadOnly(idx, false, false)
	} else {
		RmReadOnly(idx, true, false)
	}
}
