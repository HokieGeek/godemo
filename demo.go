package demo

import (
	privateiq "github.com/hokiegeek/gonexus-private/iq"
	nexusiq "github.com/sonatype-nexus-community/gonexus/iq"
	nexusrm "github.com/sonatype-nexus-community/gonexus/rm"
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
	if len(IQs) == 0 {
		IQs = DetectIQServers()
	}
	return IQs[idx].Client()
}

// OrgsIDMap returns a map of organization ids by name and the reverse
func OrgsIDMap(idx int) (id2name map[string]string, name2id map[string]string, err error) {
	if orgs, err := nexusiq.GetAllOrganizations(IQ(idx)); err == nil {
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
