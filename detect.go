package demo

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/hokiegeek/gonexus"
	// "github.com/hokiegeek/gonexus-private/iq"
	"github.com/hokiegeek/gonexus/iq"
	"github.com/hokiegeek/gonexus/rm"
)

// DetectedRM provides a helper for a found RM instance
type DetectedRM struct {
	nexus.ServerInfo
}

// Client returns a client of this RM instance
func (d DetectedRM) Client() (rm nexusrm.RM) {
	// TODO: ask for user and password
	rm, _ = nexusrm.New(d.Host, d.Username, d.Password)
	return
}

// DetectedIQ provides a helper for a found ÍQ instance
type DetectedIQ struct {
	nexus.ServerInfo
}

// Client returns a client of this IQ instance
func (d DetectedIQ) Client() (iq nexusiq.IQ) {
	// TODO: ask for user and password
	iq, _ = nexusiq.New(d.Host, d.Username, d.Password)
	return
}

func newDetectedRM(host string) (rm DetectedRM) {
	rm.Host = host
	return
}

func newDetectedIQ(host string) (iq DetectedIQ) {
	iq.Host = host
	return
}

func portInUse(p int) bool {
	l, err := net.Listen("tcp", ":"+strconv.Itoa(p))
	if err == nil {
		l.Close()
	}
	return err != nil
}

// DetectRMServers returns all instances of Repository Manager detected on the local machine
func DetectRMServers() (servers []DetectedRM) {
	host := "http://localhost"

	isRM := func(resp *http.Response) bool {
		if v, ok := resp.Header["Server"]; ok {
			return strings.HasPrefix(v[0], "Nexus")
		}
		return false
	}

	for p := 1; p < 65535; p++ {
		if portInUse(p) {
			url := fmt.Sprintf("%s:%d", host, p)
			fmt.Printf("in use: %d\n", p)
			if resp, err := http.Head(url); err == nil && isRM(resp) {
				servers = append(servers, newDetectedRM(url))
			}
		}
	}

	return
}

// DetectIQServers returns all instances of IQ detected on the local machine
func DetectIQServers() (servers []DetectedIQ) {
	host := "http://localhost"

	isIQ := func(resp *http.Response) bool {
		if v, ok := resp.Header["Set-Cookie"]; ok {
			return strings.HasPrefix(v[0], "CLM-CSRF-TOKEN")
		}
		return false
	}

	for p := 1; p < 65535; p++ {
		if portInUse(p) {
			url := fmt.Sprintf("%s:%d", host, p)
			if resp, err := http.Head(url); err == nil && isIQ(resp) {
				servers = append(servers, newDetectedIQ(url))
			}
		}
	}

	return
}
