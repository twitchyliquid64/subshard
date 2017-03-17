package main

import (
	"encoding/json"
	"errors"
	"io"
	"os"
)

// Config represents the structure of the configuration file.
type Config struct {
	Verbose       bool             `json:"verbose"`            //If true, logs warnings to STDOUT.
	Listener      string           `json:"listener"`           //Address to listen on, formatted <interface>:<port>. Interface can be omitted.
	AuthRequired  bool             `json:"auth-required"`      //When set to true, a valid user/pass or client certificate must be provided by each client.
	Users         []User           `json:"users"`              //List of user objects.
	Blacklist     []BlacklistEntry `json:"blacklist"`          //List of blacklist entries.
	Forwarders    []ForwardEntry   `json:"forwarders"`         //List of forwarders.
	ResourcesPath string           `json:"resources-location"` //Path to the directory with resources for the server. The web/ dir should be in this dir.
	Version       string           `json:"version"`            //Version to show on the landing page to users.

	TLS struct {
		CertPemPath string `json:"cert-pem-path"`
		KeyPemPath  string `json:"key-pem-path"`
		Enabled     bool   `json:"enabled"`
	}

	Path string `json:"-"`
}

// User represents the structure of a user in the configuration file.
type User struct {
	Username             string `json:"username"`
	Password             string `json:"password"` //SHA256 hash of password
	DisallowPasswordAuth bool   //NOT IMPLEMENTED: Stop password logins if set. Intended to be used in conjunction with a client certificate.
}

// BlacklistEntry represents the structure of a blacklisted host/regexp
type BlacklistEntry struct {
	Type       string `json:"type"`  //Supported values: host, host-regexp, prefix (HTTP only), regexp (HTTP only)
	Value      string `json:"value"` //Hostname for type host, URL (without scheme) for prefix, regexp for regexp (obviously).
	ParseError error  `json:"-"`
}

// ForwardEntry represents the information of another proxy, and the rules which will
// forward traffic to it.
type ForwardEntry struct {
	Name                    string `json:"name"`        //Name shown on the UI for this forwarder.
	Type                    string `json:"type"`        //Type of forwarder: HTTP, HTTPS, SOCKS. HTTPS will make a HTTPS connection to the destination on your behalf.
	Destination             string `json:"destination"` //Address of site/proxy to forward to. Formatted as host:port.
	HostFieldForHTTPProxies string `json:"host"`        //Only valid for HTTPS/HTTP forwarders and optional. If set, it will rewrite the Host header to the value specified.
	Rules                   []struct {
		Type  string `json:"type"`  //Only prefix is supported for HTTP/HTTPS proxies. host or host-regexp is supported for SOCKS proxies.
		Value string `json:"value"` //For prefix: The url prefix to be matched on (eg 'ok' will match the URL 'ok/a'). For host: exact match of the host. Likewise host-regexp, except a regex can be provided.
	} `json:"rules"` //If any rule in this section matches, traffic will be handled by this forwarder.
	Checker ForwarderChecker `json:"checker"` //Optional specification for rules for health checking the forwarder, to be shown in the UI.
}

// ForwarderChecker represents rules used to check the health of a forwarder.
type ForwarderChecker struct {
	Type        string `json:"type"`         //Type of health check: HTTP, HTTPS, TOR. HTTP* check HTTP(S) sites and proxies, HTTPS does not verify cert. TOR to connect to a tor control port.
	Destination string `json:"destination"`  //Address of the endpoint to check. For HTTP*: Formatted like 'http(s)://host[:port]'. For TOR: 'host:port'
	Auth        string `json:"auth"`         //If a password is required to connect to the tor control port, specify here. Otherwise, do not populate.
	ConnTimeout int    `json:"conn-timeout"` //Milliseconds till the connection attempt times out.
}

func readConfig(fpath string) (*Config, error) {
	var m = &Config{}

	confF, err := os.Open(fpath)

	if err != nil {
		return nil, err
	}
	defer confF.Close()

	dec := json.NewDecoder(confF)

	if err := dec.Decode(&m); err == io.EOF {
	} else if err != nil {
		return nil, err
	}
	m.Path = fpath
	return m, validateConfig(m)
}

func validateConfig(configuration *Config) error {
	if configuration.Listener == "" {
		return errors.New("listener not specified")
	}
	return nil
}
