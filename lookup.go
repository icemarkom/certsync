package certsync

import (
	"fmt"
	"net"
	"net/http"
	"strings"
)

func validReverse(ip net.IP, host string) bool {
	addrList, err := net.LookupAddr(ip.String())
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

func ValidateAddresses(host string, hostAddr net.IP) (bool, error) {
	addrList, err := LookupAddresses(host)
	if err != nil {
		return false, err
	}
	found := false
	for _, addr := range addrList {
		if hostAddr.Equal(addr) {
			if validReverse(addr, host) {
				found = true
				break
			}
		}
	}
	if !found {
		return false, fmt.Errorf("address %q is not valid for host %q", hostAddr, host)
	}
	return true, nil
}

// LookupAddresses ...
func LookupAddresses(hostName string) ([]net.IP, error) {
	a, err := net.LookupIP(hostName)
	if err != nil {
		return nil, err
	}
	return a, nil
}

// IPFromRequest
func IPFromRequest(r *http.Request) net.IP {
	if r.Header.Get("X-Forwarded-For") != "" {
		ip := net.ParseIP(strings.Split(r.Header.Get("X-Forwarded-For"), " ")[0])
		return ip
	} else {
		h, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			return nil
		}
		ip := net.ParseIP(h)
		return ip
	}
}
