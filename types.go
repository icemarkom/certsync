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

package certsync

import (
	"context"
	"net"
	"time"
)

// Const ...
const (
	DefaultPort        = 82 // XFER; mostly unused.
	DefaultCertFile    = "cert.pem"
	DefaultKeyFile     = "key.pem"
	DefaultCACertFile  = "ca.pem"
	DefaultDryRun      = true
	DefaultNewCertFile = "newcert.pem"
	DefaultNewKeyFile  = "newkey.pem"
	DefaultTimeout     = 30

	PEMTypeCertificate = "CERTIFICATE"
	PEMTypePrivateKey  = "PRIVATE KEY"

	headerXFF = "X-Forwarded-For"
)

// Resolver ...
type Resolver interface {
	LookupAddr(context.Context, string) ([]string, error)
	LookupIPAddr(context.Context, string) ([]net.IPAddr, error)
	LookupHost(context.Context, string) ([]string, error)
}

// NewConfig ...
func NewConfig(b, v, g string) *Config {
	return &Config{
		CertFile:       DefaultCertFile,
		CertKeyFile:    DefaultKeyFile,
		NewCertFile:    DefaultNewCertFile,
		NewCertKeyFile: DefaultNewKeyFile,
		CACertFile:     DefaultCACertFile,
		DryRun:         DefaultDryRun,
		Port:           DefaultPort,
		Timeout:        DefaultTimeout * time.Second,
		BinaryName:     b,
		Version:        v,
		GitCommit:      g,
		Resolver:       &net.Resolver{},
	}
}

// Config ...
type Config struct {
	HostName                       string
	CertFile, CertKeyFile          string
	NewCertFile, NewCertKeyFile    string
	CACertFile                     string
	DryRun                         bool
	Port                           int
	Timeout                        time.Duration
	BinaryName, Version, GitCommit string
	Resolver                       Resolver
}
