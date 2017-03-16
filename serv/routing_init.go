package main

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"regexp"
	"strings"

	"github.com/elazarl/goproxy"
	"github.com/elazarl/goproxy/ext/auth"
)

// register URL handlers to handle static files.
func registerStatic(configuration *Config, proxy *goproxy.ProxyHttpServer) {
	proxy.OnRequest(goproxy.UrlIs(serverHost + "/static/bootstrap.min.css")).DoFunc(func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		return serveStatic(r, "web/bootstrap.min.css", "text/css")
	})
}

func makeForwarderHandler(entry ForwardEntry) (forwardinghostMatcher, error) {
	switch entry.Type {
	case "HTTP":
		return &httpForwarder{Destination: entry.Destination, Scheme: "http", Host: entry.HostFieldForHTTPProxies}, nil
	case "HTTPS":
		return &httpForwarder{Destination: entry.Destination, Scheme: "https", Host: entry.HostFieldForHTTPProxies}, nil
	case "SOCKS":
		out := &socksForwarder{Destination: entry.Destination}
		for _, rule := range entry.Rules {
			matcher, err := makeHostBasedBlacklistHandler(rule.Type, rule.Value)
			if err != nil {
				log.Printf("Ommitting invalid forwarder rule (%s): %s\n", rule.Value, err)
			} else {
				out.MatchRules = append(out.MatchRules, matcher)
			}
		}
		return out, nil
	}
	return nil, errors.New("Could not recognise forwarder type")
}

func makeHostBasedBlacklistHandler(entryType, entryValue string) (hostMatcher, error) {
	switch entryType {
	case "host":
		return &blacklistedhostMatcher{Host: entryValue}, nil
	case "host-regexp":
		return makeBlacklistedHostRegexpHandler(entryValue)
	}
	return nil, errors.New("Could not handler blacklist type " + entryType)
}

// Blacklist match handlers + Forwarder handlers
func registerURLHandlers(configuration *Config, proxy *goproxy.ProxyHttpServer) {
	var blacklisthostMatchers []hostMatcher
	var forwardingHandlers []forwardinghostMatcher

	for _, entry := range configuration.Forwarders {
		handler, err := makeForwarderHandler(entry)
		if err != nil {
			log.Println("Forwarder err: ", err)
		} else {
			forwardingHandlers = append(forwardingHandlers, handler)
			if hForwarder, ok := handler.(*httpForwarder); ok {
				for _, matcher := range entry.Rules {
					switch matcher.Type {
					case "prefix":
						proxy.OnRequest(goproxy.UrlHasPrefix(matcher.Value)).DoFunc(hForwarder.Handle)
					}
				}
			}
		}
	}

	for _, entry := range configuration.Blacklist {
		switch entry.Type {
		case "host":
		case "host-regexp":
			handler, err := makeHostBasedBlacklistHandler(entry.Type, entry.Value)
			if err != nil {
				log.Printf("Omitting invalid blacklist %s - %s.\n", entry.Value, err.Error())
				entry.ParseError = err
			} else {
				blacklisthostMatchers = append(blacklisthostMatchers, handler)
			}

		case "prefix":
			proxy.OnRequest(goproxy.UrlHasPrefix(entry.Value)).DoFunc(handleBlacklistedHost)

		case "regexp":
			re, err := regexp.Compile(entry.Value)
			if err != nil {
				log.Printf("Omitting invalid blacklist regex %s - %s.\n", entry.Value, err.Error())
				entry.ParseError = err
			} else {
				proxy.OnRequest(goproxy.UrlMatches(re)).DoFunc(handleBlacklistedHost)
			}
		}
		gValidBlacklistEntries = append(gValidBlacklistEntries, entry)
	}

	// Intercept in the Dial method for host-based blacklist and forward rules
	proxy.Tr.Dial = func(network, addr string) (c net.Conn, err error) {
		host, _, err := net.SplitHostPort(addr)
		if err != nil {
			return nil, err
		}

		if blacklisted, msg := isHostBlacklisted(blacklisthostMatchers, host); blacklisted {
			return nil, errors.New(msg)
		}

		for _, forwardingHandler := range forwardingHandlers {
			if forwardingHandler.shouldHandleHost(host) {
				return forwardingHandler.Dial(network, addr)
			}
		}

		c, err = net.Dial(network, addr)
		return
	}
	proxy.ConnectDial = proxy.Tr.Dial
}

func isHostBlacklisted(blacklisthostMatchers []hostMatcher, host string) (bool, string) {
	for _, blacklistHandler := range blacklisthostMatchers {
		if blacklistHandler.shouldHandleHost(host) {
			return true, "Entry blacklisted: " + host
		}
	}
	return false, ""
}

// handle a request to subshard/
func handleSubshardPage(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	if r.URL.Path == "/" {
		return serveLandingPage(r)
	}
	if strings.HasPrefix(r.URL.Path, "/test") {
		return serveTestPage(r)
	}

	return r, goproxy.NewResponse(r, "text/html", 404, "Not found")
}

//handle a request to subshard.onion or subshard/test
func serveOnionPage(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	return serveTestPage(r)
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
			sha1Hash := fmt.Sprintf("%x", sha256.Sum256([]byte(passwd)))

			for _, usr := range configuration.Users {
				if usr.Username == user && usr.Password == sha1Hash {
					return true
				}
			}
			return false
		})
	}

	// setup blacklists + forwarders
	registerURLHandlers(configuration, proxy)

	registerStatic(configuration, proxy)
	gConfiguration = configuration
	proxy.OnRequest(goproxy.UrlHasPrefix(serverHost + "/")).DoFunc(handleSubshardPage)
	proxy.OnRequest(goproxy.UrlHasPrefix("subshard.onion/")).DoFunc(serveOnionPage)
	return proxy, nil
}
