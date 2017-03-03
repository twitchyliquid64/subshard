package main

import (
	"errors"
	"log"
	"net"
	"net/http"
	"regexp"
	"strings"

	"github.com/elazarl/goproxy"
	"github.com/elazarl/goproxy/ext/auth"
)

// register URL handlers to handle static files.
func registerStatic(proxy *goproxy.ProxyHttpServer) {
	proxy.OnRequest(goproxy.UrlIs(serverHost + "/static/bootstrap.min.css")).DoFunc(func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		return serveStatic(r, "web/bootstrap.min.css", "text/css")
	})
}

func registerBlacklistHandlers(configuration *Config, proxy *goproxy.ProxyHttpServer) {
	var hostEntries []string
	var hostRegexps []*regexp.Regexp
	for _, entry := range configuration.Blacklist {
		switch entry.Type {
		case "host":
			hostEntries = append(hostEntries, entry.Value)
		case "host-regexp":
			re, err := regexp.Compile(entry.Value)
			if err != nil {
				log.Printf("Omitting invalid blacklist regex %s - %s.\n", entry.Value, err.Error())
				continue
			}
			hostRegexps = append(hostRegexps, re)

		case "prefix":
			proxy.OnRequest(goproxy.UrlHasPrefix(entry.Value)).DoFunc(handleBlacklistedHost)

		case "regexp":
			re, err := regexp.Compile(entry.Value)
			if err != nil {
				log.Printf("Omitting invalid blacklist regex %s - %s.\n", entry.Value, err.Error())
				continue
			}
			proxy.OnRequest(goproxy.UrlMatches(re)).DoFunc(handleBlacklistedHost)
		}
	}

	// To support HTTPS, block hosts via intercept of dial.
	proxy.Tr.Dial = func(network, addr string) (c net.Conn, err error) {
		testHost := strings.Split(addr, ":")[0]
		for _, host := range hostEntries {
			if host == testHost {
				return nil, errors.New("Entry blacklisted: " + host)
			}
		}
		for _, host := range hostRegexps {
			if host.MatchString(testHost) {
				return nil, errors.New("Entry blacklisted: " + host.String())
			}
		}
		c, err = net.Dial(network, addr)
		return
	}
}

// handle a request to subshard/
func handleSubshardPage(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	if r.URL.Path == "/" {
		return serveLandingPage(r)
	}
	return r, goproxy.NewResponse(r, "text/html", 404, "Not found")
}

// handle a request to a host which is in configuration.BlasklistedHosts
func handleBlacklistedHost(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	return r, goproxy.NewResponse(r, "text/html", 403, "Forbidden")
}

func makeProxyServer(configuration *Config) (*goproxy.ProxyHttpServer, error) {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = configuration.Verbose

	// setup auth
	if configuration.AuthRequired {
		auth.ProxyBasic(proxy, "subshard", func(user, passwd string) bool {
			for _, usr := range configuration.Users {
				if usr.Username == user && usr.Password == passwd {
					return true
				}
			}
			return false
		})
	}

	// setup blacklists
	registerBlacklistHandlers(configuration, proxy)

	registerStatic(proxy)
	gConfiguration = configuration
	proxy.OnRequest(goproxy.UrlHasPrefix(serverHost + "/")).DoFunc(handleSubshardPage)
	return proxy, nil
}
