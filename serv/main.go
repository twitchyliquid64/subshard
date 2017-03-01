package main

import (
	"log"
	"net/http"

	"github.com/elazarl/goproxy"
)

const serverHost = "subshard"

func registerStatic(proxy *goproxy.ProxyHttpServer) {
	proxy.OnRequest(goproxy.UrlIs(serverHost + "/static/bootstrap.min.css")).DoFunc(func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		return serveStatic(r, "web/bootstrap.min.css", "text/css")
	})
}

func handleSubshardPage(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	if r.URL.Path == "/" {
		return serveLandingPage(r)
	}
	return r, goproxy.NewResponse(r, "text/html", 404, "Not found")
}

func main() {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true

	registerStatic(proxy)
	proxy.OnRequest(goproxy.UrlHasPrefix(serverHost + "/")).DoFunc(handleSubshardPage)

	log.Fatal(http.ListenAndServe(":8080", proxy))
}
