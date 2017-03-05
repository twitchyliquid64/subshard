package main

import (
	"net"
	"regexp"

	socks "github.com/fangdingjun/socks-go"
)

type hostMatcher interface {
	shouldHandleHost(host string) bool
}

type forwardinghostMatcher interface {
	hostMatcher
	Dial(network, addr string) (net.Conn, error)
}

type socksForwarder struct {
	Destination string
	MatchRules  []hostMatcher
}

func (f *socksForwarder) shouldHandleHost(host string) bool {
	for _, rule := range f.MatchRules {
		if rule.shouldHandleHost(host) {
			return true
		}
	}
	return false
}

func (f *socksForwarder) Dial(network, addr string) (net.Conn, error) {
	socksConn, errSocksDial := net.Dial("tcp", f.Destination)
	if errSocksDial != nil {
		return nil, errSocksDial
	}

	sc := &socks.Client{Conn: socksConn}
	return sc.Dial(network, addr)
}

type blacklistedhostMatcher struct {
	Host string
}

func (b *blacklistedhostMatcher) shouldHandleHost(host string) bool {
	if b.Host == host {
		return true
	}
	return false
}

type blacklistedRegexhostMatcher struct {
	HostRegex string
	Regex     *regexp.Regexp
}

func (b *blacklistedRegexhostMatcher) shouldHandleHost(host string) bool {
	if b.Regex.MatchString(host) {
		return true
	}
	return false
}

func makeBlacklistedHostRegexpHandler(regex string) (*blacklistedRegexhostMatcher, error) {
	re, err := regexp.Compile(regex)
	if err != nil {
		return nil, err
	}
	return &blacklistedRegexhostMatcher{HostRegex: regex, Regex: re}, nil
}
