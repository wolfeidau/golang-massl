package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Client echo http client
type Client struct {
	url string
}

// NewClient return a new echo client
func NewClient(url string) *Client {
	return &Client{url: url}
}

// Do make the request to the server
func (svr *Client) Do() error {

	cert, err := tls.LoadX509KeyPair("./certs/client.pem", "./certs/client-key.pem")

	caCert, err := ioutil.ReadFile("./certs/ca.pem")
	if err != nil {
		log.Fatalf("failed to load cert: %s", err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert}, // this certificate is used to sign the handshake
		RootCAs:      caCertPool,              // this is used to validate the server certificate
	}
	tlsConfig.BuildNameToCertificate()

	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: transport}

	req, err := http.NewRequest(http.MethodPost, svr.url, bytes.NewBufferString("abc123"))
	if err != nil {
		log.Fatalf("failed to build req: %s", err)
	}
	req.Header.Set("Content-Type", "plain/text")

	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("failed to post req: %s", err)
	}

	tag := fmt.Sprintf("[-> %s]", req.URL.Host)
	log.Printf("%s accept", tag)

	if len(res.TLS.PeerCertificates) > 0 {
		log.Printf("%s client common name: %+v", tag, res.TLS.PeerCertificates[0].Subject.CommonName)
	}

	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("failed to read body content: %s", err)
	}

	log.Printf("%s line: %s", tag, string(contents))

	return nil
}
