package main

import (
	"encoding/json"
	"errors"
	"io"
	"os"
)

// Config represents the structure of the configuration file.
type Config struct {
	Verbose       bool             `json:"verbose"`
	Listener      string           `json:"listener"`
	AuthRequired  bool             `json:"auth-required"`
	Users         []User           `json:"users"`
	Blacklist     []BlacklistEntry `json:"blacklist"`
	Forwarders    []ForwardEntry   `json:"forwarders"`
	ResourcesPath string           `json:"resources-location"`
	Version       string           `json:"version"`

	TLS struct {
		CertPemPath string `json:"cert-pem-path"`
		KeyPemPath  string `json:"key-pem-path"`
		Enabled     bool   `json:"enabled"`
	}

	Path string `json:"-"`
}

// User represents the structure of a user in the configuration file.
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// BlacklistEntry represents the structure of a blacklisted host/regexp
type BlacklistEntry struct {
	Type       string `json:"type"`
	Value      string `json:"value"`
	ParseError error  `json:"-"`
}

// ForwardEntry represents the information of another proxy, and the rules which will
// forward traffic to it.
type ForwardEntry struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Destination string `json:"destination"`
	Rules       []struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"rules"`
	Checker ForwarderChecker `json:"checker"`
}

// ForwarderChecker represents rules used to check the health of a forwarder.
type ForwarderChecker struct {
	Type        string `json:"type"`
	Destination string `json:"destination"`
	Auth        string `json:"auth"`
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
