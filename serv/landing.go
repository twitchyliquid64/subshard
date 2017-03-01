package main

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"net/http"

	"github.com/elazarl/goproxy"
)

func serveLandingPage(r *http.Request) (*http.Request, *http.Response) {
	buff := bytes.NewBufferString("")
	t, _ := template.ParseFiles("web/landing.html")
	t.Execute(buff, nil)

	return r, goproxy.NewResponse(r, "text/html", 200, buff.String())
}

func serveStatic(r *http.Request, fname, contentType string) (*http.Request, *http.Response) {
	d, _ := ioutil.ReadFile(fname)

	return r, goproxy.NewResponse(r, contentType, 200, string(d))
}
