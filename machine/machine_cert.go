package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"net"
	"time"
)

func main() {
	certPrivKey, _ := rsa.GenerateKey(rand.Reader, 4096)
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			Organization:  []string{"Rezidencija"},
			Country:       []string{"US"},
			Province:      []string{"CA"},
			Locality:      []string{"Los Gatos"},
			StreetAddress: []string{"120 Carlton Ave Unit 36"},
			PostalCode:    []string{"95032"},
		},
		IsCA:        false,
		IPAddresses: []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().AddDate(1, 0, 0),
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature,
		PublicKey:   certPrivKey.PublicKey,
	}
	cert.SubjectKeyId, _ = ComputeSKI(cert)

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, nil, &certPrivKey.PublicKey, nil)

	if err != nil {
		log.Fatal(err)
	}

	certPEM := new(bytes.Buffer)
	pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	certPrivKeyPEM := new(bytes.Buffer)
	pem.Encode(certPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(certPrivKey),
	})

	fmt.Println(certPEM)
}

func ComputeSKI(template *x509.Certificate) ([]byte, error) {
	pub := template.PublicKey
	encodedPub, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return nil, err
	}

	var subPKI subjectPublicKeyInfo
	_, err = asn1.Unmarshal(encodedPub, &subPKI)
	if err != nil {
		return nil, err
	}

	pubHash := sha1.Sum(subPKI.SubjectPublicKey.Bytes)
	return pubHash[:], nil
}

type subjectPublicKeyInfo struct {
	Algorithm        pkix.AlgorithmIdentifier
	SubjectPublicKey asn1.BitString
}
