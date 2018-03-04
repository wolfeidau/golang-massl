package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Server echo server
type Server struct {
	addr string
}

// NewServer return a new echo server
func NewServer(addr string) *Server {
	return &Server{addr: addr}
}

// Listen listen for http requests
func (svr *Server) Listen() error {

	handler := http.HandlerFunc(svr.EchoHandler)

	http.Handle("/echo", logMiddleware(handler))

	caCert, err := ioutil.ReadFile("./certs/ca.pem")
	if err != nil {
		return err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		ClientCAs:  caCertPool,                     // used to verify the client cert is signed by the CA and is therefore valid
		ClientAuth: tls.RequireAndVerifyClientCert, // this requires a valid client certificate to be supplied during handshake
	}

	server := &http.Server{
		Addr:      svr.addr,
		TLSConfig: tlsConfig,
	}

	// listen using the server certificate which is validated by the client
	return server.ListenAndServeTLS("./certs/server.pem", "./certs/server-key.pem")
}

// EchoHandler handle the echo request
func (svr *Server) EchoHandler(w http.ResponseWriter, req *http.Request) {

	if req.Method != http.MethodPost {
		w.WriteHeader(500)
		w.Write([]byte("bad request post only"))
		return
	}

	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("bad request post only"))
		return
	}

	defer req.Body.Close()

	tag := fmt.Sprintf("[%s -> %s]", req.URL, req.RemoteAddr)

	log.Printf("%s line: %s", tag, data)

	w.WriteHeader(200)
	w.Write(data)
}

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		tag := fmt.Sprintf("[%s -> %s]", req.URL, req.RemoteAddr)
		log.Printf("%s accept", tag)

		if len(req.TLS.PeerCertificates) > 0 {
			log.Printf("%s client common name: %+v", tag, req.TLS.PeerCertificates[0].Subject.CommonName)
		}
		next.ServeHTTP(w, req)
	})
}
