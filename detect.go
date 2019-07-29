package demo

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/sonatype-nexus-community/gonexus"
	"github.com/sonatype-nexus-community/gonexus/iq"
	"github.com/sonatype-nexus-community/gonexus/rm"
	// "github.com/hokiegeek/gonexus-private/iq"
)

// These variables contain the detected servers
var (
	RMs []IdentifiedRM
	IQs []IdentifiedIQ
)

const (
	defaultRMUser = "admin"
	defaultRMPass = "admin123"
	defaultIQUser = "admin"
	defaultIQPass = "admin123"

	defaultRMPort = 8081
	defaultIQPort = 8070

	detectionHost = "http://localhost"
)

type identifiedServer struct {
	nexus.ServerInfo
}

func (d *identifiedServer) auth(username, password string) *identifiedServer {
	d.Username = username
	d.Password = password
	return d
}

// IdentifiedRM provides a helper for a found RM instance
type IdentifiedRM struct {
	identifiedServer
}

// Auth adds username and password to the identified thingy
func (d *IdentifiedRM) Auth(username, password string) *IdentifiedRM {
	d.auth(username, password)
	return d
}

// Client returns a client of this RM instance
func (d IdentifiedRM) Client() nexusrm.RM {
	rm, _ := nexusrm.New(d.Host, d.Username, d.Password)
	return rm
}

// IdentifiedIQ provides a helper for a found ÍQ instance
type IdentifiedIQ struct {
	identifiedServer
}

// Client returns a client of this IQ instance
func (d IdentifiedIQ) Client() nexusiq.IQ {
	iq, _ := nexusiq.New(d.Host, d.Username, d.Password)
	return iq
}

// Auth adds username and password to the identified thingy
func (d *IdentifiedIQ) Auth(username, password string) *IdentifiedIQ {
	d.auth(username, password)
	return d
}

// NewIdentifiedRM creates a new instance of a IdentifiedRM
func NewIdentifiedRM(host, username, password string) (rm IdentifiedRM) {
	rm.Host = host
	rm.auth(username, password)
	return
}

// NewIdentifiedIQ creates a new instance of a IdentifiedIQ
func NewIdentifiedIQ(host, username, password string) (iq IdentifiedIQ) {
	iq.Host = host
	iq.auth(username, password)
	return
}

func detectServers(host string, sniff func(string, http.Header)) {
	portInUse := func(p int) (_ bool) {
		l, err := net.Listen("tcp4", ":"+strconv.Itoa(p))
		if err == nil {
			l.Close()
			return
		}
		return true
	}

	var httpc = &http.Client{
		Timeout: 1 * time.Second,
	}

	var wg sync.WaitGroup
	ports := make(chan int, 200)
	for w := 1; w <= 120; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for p := range ports {
				if portInUse(p) {
					url := fmt.Sprintf("%s:%d", host, p)
					if resp, err := httpc.Head(url); err == nil {
						sniff(url, resp.Header)
					}
				}
			}
		}()
	}

	for p := 1; p < 65535; p++ {
		ports <- p
	}
	close(ports)

	wg.Wait()
}

// DetectRMServers returns all instances of Repository Manager detected on the local machine
func DetectRMServers() (servers []IdentifiedRM) {
	found := make(chan string, 10)
	detectServers(detectionHost, func(url string, headers http.Header) {
		if v, ok := headers["Server"]; ok && strings.HasPrefix(v[0], "Nexus") {
			found <- url
		}
	})
	close(found)

	portSuffix := fmt.Sprintf(":%d", defaultRMPort)
	for url := range found {
		if strings.HasSuffix(url, portSuffix) {
			rms := []IdentifiedRM{NewIdentifiedRM(url, defaultRMUser, defaultRMPass)}
			rms = append(rms, servers...)
			servers = rms
		} else {
			servers = append(servers, NewIdentifiedRM(url, defaultRMUser, defaultRMPass))
		}
	}

	return
}

// DetectIQServers returns all instances of IQ detected on the local machine
func DetectIQServers() (servers []IdentifiedIQ) {
	found := make(chan string, 10)
	detectServers(detectionHost, func(url string, headers http.Header) {
		if v, ok := headers["Set-Cookie"]; ok && strings.HasPrefix(v[0], "CLM-CSRF-TOKEN") {
			found <- url
		}
	})
	close(found)

	portSuffix := fmt.Sprintf(":%d", defaultIQPort)
	for url := range found {
		if strings.HasSuffix(url, portSuffix) {
			iqs := []IdentifiedIQ{NewIdentifiedIQ(url, defaultIQUser, defaultIQPass)}
			iqs = append(iqs, servers...)
			servers = iqs
		} else {
			servers = append(servers, NewIdentifiedIQ(url, defaultIQUser, defaultIQPass))
		}
	}

	return
}

// Detect populates globals and returns any IQ and RM servers found on the machine
func Detect() ([]IdentifiedRM, []IdentifiedIQ) {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		RMs = DetectRMServers()
		wg.Done()
	}()

	go func() {
		IQs = DetectIQServers()
		wg.Done()
	}()

	wg.Wait()
	return RMs, IQs
}
