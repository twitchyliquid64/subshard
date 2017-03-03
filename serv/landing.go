package main

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/elazarl/goproxy"
)

var gConfiguration *Config

func serveLandingPage(r *http.Request) (*http.Request, *http.Response) {
	buff := bytes.NewBufferString("")
	t, _ := template.ParseFiles("web/landing.html")
	err := t.Execute(buff, gConfiguration)
	if err != nil {
		log.Println(err)
	}

	return r, goproxy.NewResponse(r, "text/html", 200, buff.String())
}

func serveStatic(r *http.Request, fname, contentType string) (*http.Request, *http.Response) {
	d, _ := ioutil.ReadFile(fname)

	return r, goproxy.NewResponse(r, contentType, 200, string(d))
}
