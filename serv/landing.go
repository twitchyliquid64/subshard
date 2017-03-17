package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"html/template"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/elazarl/goproxy"
)

var gConfiguration *Config
var gTLSConfig *tls.Config
var gConfigReloads int
var gValidBlacklistEntries []BlacklistEntry

type landingPageInfo struct {
	Configuration *Config
	TLS           *tls.Config
	TLSInfo       map[string]string
	ReloadCount   int
	Blacklist     []BlacklistEntry
}

func makeLandingPageInfo() *landingPageInfo {
	servers := ""
	if gTLSConfig != nil {
		count := 0
		for name := range gTLSConfig.NameToCertificate {
			servers += name
			count++
			if count != len(gTLSConfig.NameToCertificate) {
				servers += ", "
			}
		}
	}
	return &landingPageInfo{
		Configuration: gConfiguration,
		TLS:           gTLSConfig,
		TLSInfo:       map[string]string{"Servers": servers},
		ReloadCount:   gConfigReloads,
		Blacklist:     gValidBlacklistEntries,
	}
}

func serveLandingPage(r *http.Request) (*http.Request, *http.Response) {
	landingPagePath := "web/landing.html"
	if gConfiguration.ResourcesPath != "" {
		landingPagePath = path.Join(gConfiguration.ResourcesPath, "web/landing.html")
	}

	buff := bytes.NewBufferString("")
	buff.Grow(4096) //pre-allocate
	t, err := template.ParseFiles(landingPagePath)
	if err != nil {
		log.Println(err)
		return r, goproxy.NewResponse(r, "text/html", 500, "Internal server error")
	}
	err = t.Execute(buff, makeLandingPageInfo())
	if err != nil {
		log.Println(err)
		return r, goproxy.NewResponse(r, "text/html", 500, "Internal server error")
	}

	return r, goproxy.NewResponse(r, "text/html", 500, buff.String())
}

func serveTestPage(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	guardPagePath := "web/guard_test.html"
	if gConfiguration.ResourcesPath != "" {
		guardPagePath = path.Join(gConfiguration.ResourcesPath, "web/guard_test.html")
	}

	buff := bytes.NewBufferString("")
	buff.Grow(4096) //pre-allocate
	t, err := template.ParseFiles(guardPagePath)
	if err != nil {
		log.Println(err)
		return r, goproxy.NewResponse(r, "text/html", 500, "Internal server error")
	}
	err = t.Execute(buff, map[string]interface{}{"REQ": r, "DUMP": spew.Sdump(r), "USR": spew.Sdump(ctx.UserData)})
	if err != nil {
		log.Println(err)
		return r, goproxy.NewResponse(r, "text/html", 500, "Internal server error")
	}

	return r, goproxy.NewResponse(r, "text/html", 500, buff.String())
}

func serveStatic(r *http.Request, fname, contentType string) (*http.Request, *http.Response) {
	if gConfiguration.ResourcesPath != "" {
		fname = path.Join(gConfiguration.ResourcesPath, fname)
	}
	d, _ := ioutil.ReadFile(fname)

	return r, goproxy.NewResponse(r, contentType, 200, string(d))
}

func serveForwarderStatus(r *http.Request) (*http.Request, *http.Response) {
	spl := strings.Split(r.URL.Path, "/")
	forwarderName := spl[len(spl)-1]
	for _, forwarder := range gConfiguration.Forwarders {
		if forwarder.Name == forwarderName {
			if forwarder.Checker.Destination != "" {
				out := forwarder.Checker.GetOutput()
				if out.Err != nil {
					out.Err2 = out.Err.Error()
				}
				b, err := json.Marshal(out)
				if err != nil {
					log.Println(err)
					return r, goproxy.NewResponse(r, "text/html", 500, "internal server error")
				}
				return r, goproxy.NewResponse(r, "application/json", 200, string(b))
			}
		}
	}

	return r, goproxy.NewResponse(r, "text/html", 404, "not found")
}

type forwarderCheckOutput struct {
	Ok   bool
	Err  error
	Err2 string
	Info map[string]interface{}
}

// GetOutput returns a structure describing the operational state of the forwarder.
func (f *ForwarderChecker) GetOutput() forwarderCheckOutput {
	tr := &http.Transport{}

	switch f.Type {
	case "HTTPS":
		tr = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		fallthrough
	case "HTTP":
		netClient := &http.Client{
			Transport: tr,
			Timeout:   time.Millisecond * time.Duration(f.ConnTimeout),
		}
		start := time.Now().Round(time.Millisecond)
		response, err := netClient.Head(f.Destination)
		end := time.Now().Round(time.Millisecond)
		if err != nil {
			if err, ok := err.(net.Error); ok && err.Timeout() {
				return forwarderCheckOutput{false, errors.New("Connection timeout"), "", nil}
			}
			return forwarderCheckOutput{false, err, "", nil}
		}
		return forwarderCheckOutput{
			true,
			nil,
			"",
			map[string]interface{}{
				"status":  response.Status,
				"latency": end.Sub(start).Nanoseconds() / 1000 / 1000,
			},
		}

	case "TOR":
		t := TorControl{}
		err := t.DialTimeout("tcp", f.Destination, f.ConnTimeout)
		if err != nil {
			return forwarderCheckOutput{false, err, "", nil}
		}
		defer t.Close()
		if f.Auth != "" {
			err = t.PasswordAuthenticate(f.Auth)
			if err != nil {
				return forwarderCheckOutput{false, err, "", nil}
			}
		}
		version, err := t.TorVersion()
		if err != nil {
			return forwarderCheckOutput{false, err, "", nil}
		}
		isDormant, _ := t.IsDormant()
		circuitsEstablished, _ := t.CircuitsEstablished()
		enoughDirInfo, _ := t.EnoughDirInfo()
		socksLocal, _ := t.SocksListenersAreLocal()
		controlLocal, _ := t.ControlListenersAreLocal()
		return forwarderCheckOutput{
			socksLocal && controlLocal,
			nil,
			"",
			map[string]interface{}{
				"version":             version,
				"isDormant":           isDormant,
				"circuitsEstablished": circuitsEstablished,
				"enoughDirInfo":       enoughDirInfo,
				"socksLocal":          socksLocal,
				"controlLocal":        controlLocal,
			},
		}
	}
	return forwarderCheckOutput{false, errors.New("Cannot handle ForwarderChecker: " + f.Type), "", nil}
}
