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
	"fmt"
	"net"
	"net/http"
	"strings"
)

func validReverse(cfg *Config, ip net.IP, host string) bool {
	addrList, err := cfg.Resolver.LookupAddr(context.Background(), ip.String())
	if err != nil {
		return false
	}
	for _, a := range addrList {
		if strings.TrimSuffix(a, ".") == host {
			return true
		}
	}
	return false
}

func ValidateAddresses(cfg *Config, host string, hostAddr net.IP) error {
	addrList, err := LookupAddresses(cfg, host)
	if err != nil {
		return err
	}
	found := false
	for _, addr := range addrList {
		if hostAddr.Equal(addr.IP) {
			if validReverse(cfg, addr.IP, host) {
				found = true
				break
			}
		}
	}
	if !found {
		return fmt.Errorf("address %q is not valid for host %q", hostAddr, host)
	}
	return nil
}

// LookupAddresses ...
func LookupAddresses(cfg *Config, hostName string) ([]net.IPAddr, error) {
	a, err := cfg.Resolver.LookupIPAddr(context.Background(), hostName)
	if err != nil {
		return nil, err
	}
	return a, nil
}

// IPFromRequest
func IPFromRequest(r *http.Request) (net.IP, error) {
	if r == nil {
		return nil, fmt.Errorf("requestor host:port is empty")
	}
	if r.Header.Get("X-Forwarded-For") != "" {
		ip := net.ParseIP(strings.Split(r.Header.Get("X-Forwarded-For"), " ")[0])
		return ip, nil
	} else {
		h, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			return nil, fmt.Errorf("invalid requestor host:port combination: %v", err)
		}
		ip := net.ParseIP(h)
		return ip, nil
	}
}
