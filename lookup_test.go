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
	"errors"
	"net"
	"net/http"
	"testing"

	"github.com/foxcpp/go-mockdns"
	// "github.com/icemarkom/certsync"
)

func setUp() *Config {
	cfg := NewConfig("test_binary", "test_version", "test_commit")
	cfg.Resolver = &mockdns.Resolver{
		Zones: map[string]mockdns.Zone{
			"valid.example.com.": {
				A:    []string{"10.0.0.1"},
				AAAA: []string{"2001:db8::1"},
			},
			"mismatched.example.com.": {
				A: []string{"10.0.0.2"},
			},
			"cname-for-valid.example.com.": {
				CNAME: "valid.example.com.",
			},
			"1.0.0.10.in-addr.arpa.": {
				PTR: []string{"valid.example.com."},
			},
			"2.0.0.10.in-addr.arpa.": {
				PTR: []string{"badly-mismatched.example.com."},
			},
			"1.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.8.b.d.0.1.0.0.2.ip6.arpa.": {
				PTR: []string{"valid.example.com."},
			},
		},
	}
	return cfg
}

func setupBadResolver() *Config {
	cfg := NewConfig("test_binary", "test_version", "test_commit")
	cfg.Resolver = &mockdns.Resolver{
		Zones: map[string]mockdns.Zone{
			"valid.example.com": {
				Err: errors.New("error"),
			},
		},
	}
	return cfg
}

func TestIPFromRequest(t *testing.T) {
	var tests = []struct {
		r       *http.Request
		want    net.IP
		wantErr bool
	}{
		{
			r:       nil,
			want:    nil,
			wantErr: true,
		},
		{
			r: &http.Request{
				RemoteAddr: "10.0.0.1",
			},
			want:    nil,
			wantErr: true,
		},
		{
			r: &http.Request{
				RemoteAddr: "2001:db8::1",
			},
			want:    nil,
			wantErr: true,
		},
		{
			r: &http.Request{
				RemoteAddr: "2001:db8::5000",
			},
			want:    nil,
			wantErr: true,
		},
		{
			r: &http.Request{
				RemoteAddr: "[2001:db8::1]",
			},
			want:    nil,
			wantErr: true,
		},
		{
			r: &http.Request{
				RemoteAddr: "[2001:db8::1]:5000",
			},
			want:    net.ParseIP("2001:db8::1"),
			wantErr: false,
		},
		{
			r: &http.Request{
				RemoteAddr: "10.0.0.1:5000",
			},
			want:    net.ParseIP("10.0.0.1"),
			wantErr: false,
		},
		{
			r: &http.Request{
				RemoteAddr: "10.0.0.5:5000",
				Header: http.Header{
					headerXFF: {"10.0.0.1"},
				},
			},
			want:    net.ParseIP("10.0.0.1"),
			wantErr: false,
		},
		{
			r: &http.Request{
				RemoteAddr: "10.0.0.5:5000",
				Header: http.Header{
					headerXFF: {"10.0.0.1", "10.0.0.2", "10.0.0.3"},
				},
			},
			want:    net.ParseIP("10.0.0.1"),
			wantErr: false,
		},
		{
			r: &http.Request{
				RemoteAddr: "10.0.0.5:5000",
				Header: http.Header{
					headerXFF: {},
				},
			},
			want:    net.ParseIP("10.0.0.5"),
			wantErr: false,
		},
		{
			r: &http.Request{
				RemoteAddr: "10.0.0.5:5000",
				Header: http.Header{
					headerXFF: {"a"},
				},
			},
			want:    net.ParseIP("10.0.0.5"),
			wantErr: false,
		},
	}
	for _, tc := range tests {
		got, err := IPFromRequest(tc.r)
		if !got.Equal(tc.want) || (err == nil && tc.wantErr) {
			t.Errorf("IPFromRequest(%q [XFF: %q]): want: %v (got: %v), wantErr: %v (gotErr: %v)", tc.r.RemoteAddr, tc.r.Header.Get(headerXFF), tc.want, got, tc.wantErr, err != nil)
		}
	}

}

func TestValidReverse(t *testing.T) {
	var tests = []struct {
		host string
		ip   net.IP
		want bool
	}{
		{
			host: "",
			ip:   nil,
			want: false,
		},
		{
			host: "invalid.example.com",
			ip:   net.ParseIP("10.0.0.2"),
			want: false,
		},
		{
			host: "valid.example.com",
			ip:   net.ParseIP("10.0.0.3"),
			want: false,
		},
		{
			host: "mismatched.example.com",
			ip:   net.ParseIP("10.0.0.2"),
			want: false,
		},
		{
			host: "valid.example.com",
			ip:   net.ParseIP("10.0.0.1"),
			want: true,
		},
		{
			host: "valid.example.com",
			ip:   net.ParseIP("2001:db8::1"),
			want: true,
		},
	}

	cfg := setUp()

	for _, tc := range tests {
		got := validReverse(cfg, tc.ip, tc.host)
		if got != tc.want {
			t.Errorf("validReverse(%q, %q): want: %v, got: %v", tc.ip, tc.host, tc.want, got)
		}
	}
}

func TestValidateAddresses(t *testing.T) {
	var tests = []struct {
		host    string
		ip      net.IP
		wantErr bool
	}{
		{
			host:    "",
			ip:      nil,
			wantErr: true,
		},
		{
			host:    "invalid.example.com",
			ip:      net.ParseIP("10.0.0.1"),
			wantErr: true,
		},
		{
			host:    "cname-for-valid.example.com",
			ip:      net.ParseIP("10.0.0.1"),
			wantErr: true,
		},
		{
			host:    "valid.example.com",
			ip:      net.ParseIP("10.0.0.1"),
			wantErr: false,
		},
		{
			host:    "valid.example.com",
			ip:      net.ParseIP("2001:db8::1"),
			wantErr: false,
		},
	}

	cfg := setUp()

	for _, tc := range tests {
		if gotErr := ValidateAddresses(cfg, tc.host, tc.ip) != nil; gotErr != tc.wantErr {
			t.Errorf("ValidateAddresses(%q, %q): want: %v, got: %v", tc.host, tc.ip, tc.wantErr, gotErr)
		}
	}

	// Special case, when the resolver errors out.
	cfg = setupBadResolver()
	if gotErr := ValidateAddresses(cfg, "valid.example.com", net.ParseIP("10.0.0.1")) != nil; !gotErr {
		t.Errorf("[resolver error case] ValidateAddresses(%q, %q): want: %v, got: %v", "valid.example.com", "10.0.0.1", true, gotErr)
	}

}
