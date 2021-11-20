package certsync

import "time"

const (
	DefaultPort        = 15000
	DefaultCertFile    = "cert.pem"
	DefaultCertKeyFile = "key.pem"
	DefaultCACertFile  = "ca.pem"
	DefaultNewCertFile = "newcert.pem"
	DefaultNewKeyFile  = "newkey.pem"

	PEMTypeCertificate = "CERTIFICATE"
	PEMTypePrivateKey  = "PRIVATE KEY"
)

type Config struct {
	HostName                     string
	CertFile, CertKeyFile        string
	NewCertFile, NewCertKeyFile  string
	CACertFile                   string
	Port                         int
	Timeout                      time.Duration
	BinaryName, Version, GitHash string
}
