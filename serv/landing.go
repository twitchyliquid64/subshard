package main

import (
	"bytes"
	"crypto/tls"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path"

	"github.com/elazarl/goproxy"
)

var gConfiguration *Config
var gTLSConfig *tls.Config
var gConfigReloads int

type landingPageInfo struct {
	Configuration *Config
	TLS           *tls.Config
	TLSInfo       map[string]string
	ReloadCount   int
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

func serveStatic(r *http.Request, fname, contentType string) (*http.Request, *http.Response) {
	if gConfiguration.ResourcesPath != "" {
		fname = path.Join(gConfiguration.ResourcesPath, fname)
	}
	d, _ := ioutil.ReadFile(fname)

	return r, goproxy.NewResponse(r, contentType, 200, string(d))
}
