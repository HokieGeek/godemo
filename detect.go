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

var (
	RMs []DetectedRM
	IQs []DetectedIQ
)

type detectedServer struct {
	nexus.ServerInfo
}

func (d *detectedServer) login(username, password string) *detectedServer {
	d.Username = username
	d.Password = password
	return d
}

// DetectedRM provides a helper for a found RM instance
type DetectedRM struct {
	detectedServer
}

// Login adds username and password to the detected thingy
func (d *DetectedRM) Login(username, password string) *DetectedRM {
	d.login(username, password)
	return d
}

// Client returns a client of this RM instance
func (d DetectedRM) Client() nexusrm.RM {
	rm, _ := nexusrm.New(d.Host, d.Username, d.Password)
	return rm
}

// DetectedIQ provides a helper for a found ÍQ instance
type DetectedIQ struct {
	detectedServer
}

// Client returns a client of this IQ instance
func (d DetectedIQ) Client() nexusiq.IQ {
	iq, _ := nexusiq.New(d.Host, d.Username, d.Password)
	return iq
}

// Login adds username and password to the detected thingy
func (d *DetectedIQ) Login(username, password string) *DetectedIQ {
	d.login(username, password)
	return d
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

// Detect populates globals and returns any IQ and RM servers found on the machine
func Detect() ([]DetectedRM, []DetectedIQ) {
	RMs = DetectRMServers()
	IQs = DetectIQServers()
	return RMs, IQs
}
