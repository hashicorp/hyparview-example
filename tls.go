// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"

	"google.golang.org/grpc/credentials"
)

func serverCreds(c *clientConfig) (credentials.TransportCredentials, error) {
	certificate, err := tls.LoadX509KeyPair(c.serverCert, c.serverKey)
	if err != nil {
		return nil, fmt.Errorf("could not load server key pair: %s", err)
	}

	certPool, err := caPool(c.caCert)
	if err != nil {
		return nil, err
	}

	creds := credentials.NewTLS(&tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{certificate},
		ClientCAs:    certPool,
	})

	return creds, nil
}

func caPool(caCert string) (*x509.CertPool, error) {
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(caCert)
	if err != nil {
		return nil, fmt.Errorf("could not read ca certificate: %s", err)
	}

	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		return nil, errors.New("failed to append client certs")
	}

	return certPool, nil
}

func clientCreds(c *clientConfig, addr string) (credentials.TransportCredentials, error) {
	certificate, err := tls.LoadX509KeyPair(c.clientCert, c.clientKey)
	if err != nil {
		return nil, fmt.Errorf("could not load client key pair: %s", err)
	}

	certPool, err := caPool(c.caCert)
	if err != nil {
		return nil, err
	}

	// Create the TLS credentials
	creds := credentials.NewTLS(&tls.Config{
		ServerName:   addr,
		Certificates: []tls.Certificate{certificate},
		RootCAs:      certPool,
	})

	return creds, nil
}
