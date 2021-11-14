package certsync

import "time"

const (
	DefaultPort        = 15000
	DefaultCertFile    = "cert.pem"
	DefaultCertKeyFile = "key.pem"
	DefaultCACertFile  = "ca.pem"
)

type Config struct {
	HostName, CertFile, CertKeyFile, CACertFile string
	Port                                        int
	Timeout                                     time.Duration
}
