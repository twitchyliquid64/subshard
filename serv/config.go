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
	ResourcesPath string           `json:"resources-location"`

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
	Type  string `json:"type"`
	Value string `json:"value"`
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
