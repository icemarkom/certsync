// Copyright 2021 CertSync Contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"bytes"
	"crypto/tls"
	"encoding/pem"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strings"
	"time"

	cs "github.com/icemarkom/certsync"
	"github.com/icemarkom/certsync/common"
)

var (
	cfg *cs.Config

	binaryName, version, gitCommit string
)

func init() {
	var v bool

	cfg = cs.NewConfig(binaryName, version, gitCommit)

	flag.Usage = func() { common.ProgramUsage(cfg) }

	flag.StringVar(&cfg.HostName, "host", "", "Server hostname")
	flag.IntVar(&cfg.Port, "port", cs.DefaultPort, "Server port")
	flag.StringVar(&cfg.CertFile, "clientcert", cs.DefaultCertFile, "Client certificate file")
	flag.StringVar(&cfg.CertKeyFile, "clientkey", cs.DefaultKeyFile, "Client private key file")
	flag.BoolVar(&cfg.DryRun, "dry_run", true, "Dry run - don't connect to the server")
	flag.StringVar(&cfg.NewCertFile, "newcert", cs.DefaultNewCertFile, "New certificate file")
	flag.StringVar(&cfg.NewCertKeyFile, "newkey", cs.DefaultNewKeyFile, "New key file")
	flag.DurationVar(&cfg.Timeout, "timeout", cs.DefaultTimeout*time.Second, "Server timeout in seconds")
	flag.BoolVar(&v, "version", false, "Print version and exit.")

	flag.Parse()

	if len(os.Args) == 1 {
		common.ProgramUsage(cfg)
		os.Exit(0)
	}

	if v {
		common.ProgramVersion(cfg)
		os.Exit(0)
	}

	if cfg.HostName == "" {
		log.Fatalf("Server hostname not specified.\n\n")
	}
	if cfg.Port < 1 && cfg.Port > math.MaxInt16 {
		log.Printf("Invalid port number: %d.", cfg.Port)
	}
}

func setupClient() (*http.Client, error) {
	cert, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.CertKeyFile)
	if err != nil {
		return nil, fmt.Errorf("cannot load x509 certificate or key: %v", err)
	}

	tc := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		MinVersion:   tls.VersionTLS13,
	}

	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tc,
		},
		Timeout: cfg.Timeout,
	}, nil
}

func parsePEM(data []byte) ([]*pem.Block, error) {
	var (
		pemList []*pem.Block
		err     error
	)

	for {
		block, rest := pem.Decode(data)
		if block == nil {
			break
		}
		pemList = append(pemList, block)
		data = rest
		if len(rest) == 0 {
			break
		}
	}
	if len(pemList) == 0 {
		err = fmt.Errorf("no PEM blocks found")
	}
	return pemList, err
}

func savePEM(blocks []*pem.Block, fileName string, fileMode os.FileMode) (err error) {

	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, fileMode)
	defer func() {
		err = f.Close()
	}()
	if err != nil {
		return fmt.Errorf("cannot write to %q: %v", fileName, err)
	}
	for _, block := range blocks {
		err = pem.Encode(f, block)
		if err != nil {
			return fmt.Errorf("cannnot write PEM block to %q: %v", fileName, err)
		}
	}
	return nil
}

func saveData(data []byte) error {
	var certs, keys, others []*pem.Block

	pemList, err := parsePEM(data)
	if err != nil {
		return fmt.Errorf("cannot parse PEM data: %v", err)
	}
	for _, block := range pemList {
		if block.Type == cs.PEMTypeCertificate {
			certs = append(certs, block)
			continue
		}
		if block.Type == cs.PEMTypePrivateKey {
			keys = append(keys, block)
			continue
		}
		others = append(others, block)
	}
	if len(others) != 0 {
		log.Printf("Ignoring %d PEM blocks that are not %q or %q.", len(others), cs.PEMTypeCertificate, cs.PEMTypePrivateKey)
	}
	err = savePEM(certs, cfg.NewCertFile, 0600)
	if err != nil {
		return fmt.Errorf("cannot save certs: %v", err)
	}
	err = savePEM(keys, cfg.NewCertKeyFile, 0600)
	if err != nil {
		return fmt.Errorf("cannot save keys: %v", err)
	}
	return nil
}

func main() {
	client, err := setupClient()
	if err != nil {
		log.Fatalf("Could not setup HTTPS client: %v", err)
	}
	if cfg.DryRun {
		log.Println("Dry run - not connecting to the server.")
		os.Exit(0)
	}
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://%s:%d", cfg.HostName, cfg.Port), bytes.NewBuffer([]byte{}))
	if err != nil {
		log.Fatalf("Cannot create HTTPS request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Cannot complete HTTPS request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp.StatusCode)
		log.Fatalf("Received %d (%q) from the server. Full response: %q", resp.StatusCode, http.StatusText(resp.StatusCode), strings.TrimSpace(string(body)))
	}

	err = saveData(body)
	if err != nil {
		log.Fatalf("Error saving data: %v", err)
	}
}
