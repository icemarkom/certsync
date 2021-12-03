package certsync

import (
	"net"
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
