package pki

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"time"
)

type PKI struct {
	private      *rsa.PrivateKey
	private_path string
	public       *rsa.PublicKey
	public_path  string
	certificate  *x509.Certificate
	cert         []byte
	folder       string
}

func New(folder string, country string, org string, name string) (*PKI, error) {
	var err error
	out := &PKI{
		folder: folder,
	}
	out.private, err = rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return out, err
	}

	out.public = &out.private.PublicKey
	out.certificate = &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{org},
			Country:      []string{country},
			CommonName:   name,
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(time.Hour * 24 * 180),

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	out.cert, err = x509.CreateCertificate(rand.Reader, out.certificate, out.certificate, out.public, out.private)
	if err != nil {
		return out, err
	}

	return out, nil
}

func (p *PKI) Write() error {
	var err error

	_, err = os.Stat(p.folder)
	if err != nil {
		os.Mkdir(p.folder, 0770)
	}

	file := fmt.Sprintf("%s/cert.crt", p.folder)
	crt, err := os.Create(file)
	if err != nil {
		return err
	}
	pem.Encode(crt, &pem.Block{Type: "CERTIFICATE", Bytes: p.cert})

	file = fmt.Sprintf("%s/key.pem", p.folder)
	key, err := os.Create(file)
	if err != nil {
		return err
	}
	pem.Encode(key, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(p.private)})

	file = fmt.Sprintf("%s/pub.pem", p.folder)
	_, err = os.Stat(file)
	pub, err := os.Create(file)
	if err != nil {
		return err
	}

	pem.Encode(pub, &pem.Block{Type: "RSA PUBLIC KEY", Bytes: x509.MarshalPKCS1PublicKey(p.public)})
	return nil
}
