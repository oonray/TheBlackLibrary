package pki

import (
	"fmt"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"log"
	"math/big"
	"os"
	"time"
	log "github.com/sirupsen/logrus"
)

type PKI struct {
	cert *x509.Certificate
	priv *rsa.PrivateKey
	pub  *rsa.PublicKey
	name *pkix.Name
	ca   *string
	org  string
	addr string
	code string
}

func NewPKI(org string, addr string, code string)(*PKI,error){
	var err error

	out := &PKI{org:org,addr:addr,code:code}
	out.cert = &x509.Certificate{
		SerialNumber: big.NewInt(1653),
		Subject: Pkix(),
		NotBefore: time.Now(),
		NotAfter: time.Now().AddDate(10,0,0),
		IsCA: true,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth,x509.ExtKeyUsageServerAuth},
		KeyUsage: x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

    out.priv, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		logrus.Errorf("Could not Create PKI")
		return out, err
	}

	out.pub = &out.priv.PublicKey
	out.ca,err = x509.CreateCertificate(rand.Reader,out.cert,out.cert,out.pub,out.priv)
	if err != nil {
		logrus.Errorf("Could not Create Cert")
		return out, err
	}

	return out, nil
	
}

func (pki *PKI)Pkix(){
	pki.name = pkix.Name{
		Organization: []string{pki.org},
		Country: []string{"NO"},
		Province: []string{"Oslo"},
		Locality: []string{"Oslo"},
		StreetAddress: []string{pku.addr},
		PostalCode: []string{pki.code},
	}
}

func (pki *PKI)Save() error {
	out, err := os.Create("ca.crt")
	if(err!=nil){return err}
	pem.Encode(out,&pem.Block{Type:"CERTIFICATE",Bytes:x509.MarshallPKCS1PrivateKey(pki.ca)})
	out.Close()

	out, err := os.Create("ca.key", os.O_WERONLY|os.O_CREATE|os.O_TRUNK,0600)
	if(err!=nil){return err}
	pem.Encode(out,&pem.Block{Type:"RSA PRIVATE KEY",Bytes:x509.MarshallPKCS1PrivateKey(pki.priv)})
	out.Close()

	out, err := os.Create("ca.pub", os.O_WERONLY|os.O_CREATE|os.O_TRUNK,0600)
	if(err!=nil){return err}
	out.Write(pki.pub)
	out.Close()
}
