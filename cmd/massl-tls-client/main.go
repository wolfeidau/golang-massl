package main

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
)

func main() {

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

	conn, err := tls.Dial("tcp", "localhost:2222", tlsConfig)
	if err != nil {
		log.Fatalf("failed to open conn: %s", err)
	}

	// this is required to complete the handshake and populate the connection state
	// we are doing this so we can print the peer certificates prior to reading / writing to the connection
	err = conn.Handshake()
	if err != nil {
		log.Fatalf("failed to complete handshake: %s", err)
	}

	tag := fmt.Sprintf("[%s -> %s]", conn.LocalAddr(), conn.RemoteAddr())
	log.Printf("%s connect", tag)

	if len(conn.ConnectionState().PeerCertificates) > 0 {
		log.Printf("%s client common name: %+v", tag, conn.ConnectionState().PeerCertificates[0].Subject.CommonName)
	}

	_, err = conn.Write([]byte("abc123\n"))
	if err != nil {
		log.Fatalf("failed to write line: %s", err)
	}

	log.Printf("%s write", tag)

	b := bufio.NewReader(conn)

	line, err := b.ReadBytes('\n')
	if err != nil {
		log.Fatalf("failed to read echoed line: %s", err)
	}

	log.Printf("%s line: %s", tag, string(line))
}
