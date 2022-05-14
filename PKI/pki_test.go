package pki

import (
	"testing"
)

const folder string = "./certs"

var pki *PKI

func Test_PKI_New(t *testing.T) {
	var err error
	pki, err = New(folder, "NO", "NinjaTown", "localhost")
	if err != nil {
		t.Errorf("Could not create pki: %s", err)
	}
}

func Test_PKI_WRITE(t *testing.T) {
	err := pki.Write()
	if err != nil {
		t.Errorf("Could not write %s", err)
	}
}
