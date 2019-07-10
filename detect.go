package demo

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"

	// "github.com/hokiegeek/gonexus"
	// "github.com/hokiegeek/gonexus-private/iq"
	"github.com/hokiegeek/gonexus/iq"
	"github.com/hokiegeek/gonexus/rm"
)

func portInUse(p int) bool {
	l, err := net.Listen("tcp", ":"+strconv.Itoa(p))
	if err == nil {
		l.Close()
	}
	return err != nil
}

// DetectRMServers returns all instances of Repository Manager detected on the local machine
func DetectRMServers() (servers []*nexusrm.RM) {
	host := "http://localhost"

	isRM := func(resp *http.Response) bool {
		if v, ok := resp.Header["Server"]; ok {
			return strings.HasPrefix(v[0], "Nexus")
		}
		return false
	}

	for p := 1; p < 65535; p++ {
		url := fmt.Sprintf("%s:%d", host, p)
		if portInUse(p) {
			if resp, err := http.Head(url); err == nil && isRM(resp) {
				if s, err := nexusrm.New(url, "admin", "admin123"); err == nil {
					servers = append(servers, s)
				}
			}
		}
	}

	return
}

// DetectIQServers returns all instances of IQ detected on the local machine
func DetectIQServers() (servers []*nexusiq.IQ) {
	host := "http://localhost"

	isIQ := func(resp *http.Response) bool {
		if v, ok := resp.Header["Set-Cookie"]; ok {
			return strings.HasPrefix(v[0], "CLM-CSRF-TOKEN")
		}
		return false
	}

	for p := 1; p < 65535; p++ {
		url := fmt.Sprintf("%s:%d", host, p)
		if portInUse(p) {
			if resp, err := http.Head(url); err == nil && isIQ(resp) {
				if s, err := nexusiq.New(url, "admin", "admin123"); err == nil {
					servers = append(servers, s)
				}
				// if s, err := privateiq.New(url, "admin", "admin123"); err == nil {
				// 	servers = append(servers, s)
				// }
			}
		}
	}

	return
}
