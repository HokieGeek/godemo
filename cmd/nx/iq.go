package main

import (
	"time"

	"github.com/hokiegeek/godemo"
	"github.com/hokiegeek/gonexus-private/iq"
	"github.com/sonatype-nexus-community/gonexus/iq"
)

func iqCreateAndDeleteOrg(idx int) {
	orgID, err := nexusiq.CreateOrganization(demo.IQ(idx), "arstarst")
	if err != nil {
		panic(err)
	}

	time.Sleep(15 * time.Second)

	if err := privateiq.DeleteOrganization(demo.IQ(idx), orgID); err != nil {
		panic(err)
	}
}
