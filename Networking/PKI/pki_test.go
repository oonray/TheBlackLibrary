package pki

import (
	"testing"
)

const folder string = "./certs"

func Test_PKI_New(t *testing.T) {
	_, err := New(folder, "NO", "NinjaTown", "localhost")
	if err != nil {
		t.Errorf("Could not create pki: %s", err)
	}
}

func Test_PKI_WRITE(t *testing.T) {
	pki, err := New(folder, "NO", "NinjaTown", "localhost")
	if err != nil {
		t.Errorf("Could not create pki: %s", err)
	}
	err = pki.Write()
	if err != nil {
		t.Errorf("Could not write %s", err)
	}
}
