package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	cs "github.com/icemarkom/certsync"
	"github.com/icemarkom/certsync/common"
)

var (
	cfg cs.Config

	binaryName, version, gitHash string
)

func init() {
	var v bool

	cfg.Version = version
	cfg.GitHash = gitHash
	cfg.BinaryName = binaryName
	if cfg.BinaryName == "" {
		cfg.BinaryName = os.Args[0]
	}

	flag.Usage = func() { common.ProgramUsage(cfg) }

	flag.StringVar(&cfg.HostName, "host", "", "Server hostname")
	flag.IntVar(&cfg.Port, "port", cs.DefaultPort, "Server port")
	flag.StringVar(&cfg.CertFile, "cert", cs.DefaultCertFile, "Certificate file")
	flag.StringVar(&cfg.CertKeyFile, "key", cs.DefaultCertKeyFile, "Private key file")
	flag.StringVar(&cfg.CACertFile, "ca", cs.DefaultCACertFile, "Client CA certificate file")
	flag.DurationVar(&cfg.Timeout, "timeout", 30*time.Second, "Server timeout.")
	flag.BoolVar(&v, "version", false, "Print version and exit.")

	flag.Parse()

	if v {
		common.ProgramVersion(cfg)
		os.Exit(0)
	}

	if cfg.HostName == "" {
		h, err := os.Hostname()
		if err != nil {
			log.Fatalf("cannot get local hostname: %v", err)
		}
		cfg.HostName = h
		log.Printf("Hostname not specified, using default local name: %q.", cfg.HostName)
	}
	if cfg.Port < 1 && cfg.Port > math.MaxInt16 {
		log.Printf("Invalid port number: %d.", cfg.Port)
	}

	log.Printf("Configuration: host: %q, port: %d, cert file: %q, key file: %q, CA cert: %q.", cfg.HostName, cfg.Port, cfg.CertFile, cfg.CertKeyFile, cfg.CACertFile)
}

func setupServer() (*http.Server, error) {
	ca, err := ioutil.ReadFile(cfg.CACertFile)
	if err != nil {
		return nil, fmt.Errorf("error opening CA certificate file %q: %v", cfg.CACertFile, err)
	}

	caPool := x509.NewCertPool()
	caPool.AppendCertsFromPEM(ca)

	tc := &tls.Config{
		ServerName: cfg.HostName,
		ClientAuth: tls.RequireAndVerifyClientCert,
		ClientCAs:  caPool,
		MinVersion: tls.VersionTLS13,
	}

	return &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		TLSConfig:    tc,
	}, nil
}

func validRequest(r *http.Request) bool {
	cn := r.TLS.VerifiedChains[0][0].Subject.CommonName

	_, err := cs.ValidateAddresses(cn, cs.IPFromRequest(r))
	if err != nil {
		log.Printf("Client %q IPs did not validate: %v", cn, err)
		return false
	}
	return true
}

func logRequest(r *http.Request) {
	log.Printf("%s %q request for host %q from client address %q (X-Forwarded-For: %q)",
		r.Method, r.URL.Path, r.Host, r.RemoteAddr, r.Header.Get("X-Forwarded-For"))
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	if !validRequest(r) {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
	data, err := ioutil.ReadFile(cfg.CertFile)
	if err != nil {
		log.Printf("Error serving cert file: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Write(data)
	log.Println("Certificate sent.")
	data, err = ioutil.ReadFile(cfg.CertKeyFile)
	if err != nil {
		log.Printf("Error serving cert key file: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Write(data)
	log.Println("Certificate key sent.")
	data = nil
}

func main() {
	server, err := setupServer()
	if err != nil {
		log.Fatalf("Unable to configure HTTPS server: %v", err)
	}
	http.HandleFunc("/", handleRoot)

	log.Printf("Starting HTTPS server on host %s:%d", cfg.HostName, cfg.Port)
	if err := server.ListenAndServeTLS(cfg.CertFile, cfg.CertKeyFile); err != nil {
		log.Fatal(err)
	}
}
