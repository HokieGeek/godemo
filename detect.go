package demo

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hokiegeek/gonexus"
	// "github.com/hokiegeek/gonexus-private/iq"
	"github.com/hokiegeek/gonexus/iq"
	"github.com/hokiegeek/gonexus/rm"
)

// These variables contain the detected servers
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
	for w := 1; w <= 50; w++ {
		go func() {
			wg.Add(1)
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
func DetectRMServers() (servers []DetectedRM) {
	host := "http://localhost"

	found := make(chan string, 10)
	detectServers(host, func(url string, headers http.Header) {
		if v, ok := headers["Server"]; ok && strings.HasPrefix(v[0], "Nexus") {
			found <- url
		}
	})
	close(found)

	for url := range found {
		servers = append(servers, newDetectedRM(url))
	}

	return
}

// DetectIQServers returns all instances of IQ detected on the local machine
func DetectIQServers() (servers []DetectedIQ) {
	host := "http://localhost"

	found := make(chan string, 10)
	detectServers(host, func(url string, headers http.Header) {
		if v, ok := headers["Set-Cookie"]; ok && strings.HasPrefix(v[0], "CLM-CSRF-TOKEN") {
			found <- url
		}
	})
	close(found)

	for url := range found {
		servers = append(servers, newDetectedIQ(url))
	}

	return
}

// Detect populates globals and returns any IQ and RM servers found on the machine
func Detect() ([]DetectedRM, []DetectedIQ) {
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
