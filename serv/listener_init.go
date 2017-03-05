package main

import (
	"crypto/tls"
	"errors"
	"net"
)

func tlsConfig(configuration *Config) (*tls.Config, error) {
	tlsConfig := new(tls.Config)

	tlsConfig.MinVersion = tls.VersionTLS12
	tlsConfig.CurvePreferences = []tls.CurveID{tls.CurveP384, tls.CurveP256}
	tlsConfig.PreferServerCipherSuites = true
	tlsConfig.CipherSuites = []uint16{
		tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
		tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
		tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		tls.TLS_RSA_WITH_AES_128_CBC_SHA,
		tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA}

	if configuration.TLS.CertPemPath == "" || configuration.TLS.KeyPemPath == "" {
		return nil, errors.New("Must specify path to PEM encoded certificate and key file in TLS mode")
	}

	cert, err := tls.LoadX509KeyPair(configuration.TLS.CertPemPath, configuration.TLS.KeyPemPath)
	if err != nil {
		return nil, err
	}

	tlsConfig.Certificates = append(tlsConfig.Certificates, cert)
	tlsConfig.BuildNameToCertificate()
	return tlsConfig, nil
}

func setupListener(configuration *Config) (net.Listener, error) {
	if configuration.TLS.Enabled {
		tlsConf, err := tlsConfig(configuration)
		if err != nil {
			return nil, err
		}

		listener, err := tls.Listen("tcp", configuration.Listener, tlsConf)
		if err != nil {
			return nil, err
		}

		gTLSConfig = tlsConf
		return listener, nil
	}

	listener, err := net.Listen("tcp", configuration.Listener)
	if err != nil {
		return nil, err
	}
	return listener, nil
}

// ErrInterrupt is raised if SIGINT is recieved.
var ErrInterrupt = errors.New("Interrupt")
