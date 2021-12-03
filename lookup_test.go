package certsync

import (
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

			"1.0.0.10.in-addr.arpa.": {
				PTR: []string{"valid.example.com"},
			},
			"2.0.0.10.in-addr.arpa.": {
				PTR: []string{"badly-mismatched.example.com"},
			},
			"1.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.8.b.d.0.1.0.0.2.ip6.arpa.": {
				PTR: []string{"valid.example.com"},
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
	}
	for _, tc := range tests {
		got, err := IPFromRequest(tc.r)
		if !got.Equal(tc.want) || (err == nil && tc.wantErr) {
			t.Errorf("IPFromRequest(%q): want: %v (got: %v), wantErr: %v (gotErr: %v)", tc.r.RemoteAddr, tc.want, got, tc.wantErr, err != nil)
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
