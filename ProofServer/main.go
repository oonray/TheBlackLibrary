package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/oonray/TheBlackLibrary/PKI@1.1.14"
	"log"
)

const (
	folder string = "./certs"
)

var (
	pki PKI.PKI
	err error
)

func init() {
	pki, err = PKI.New(folder, "NO", "NinjaTown", "loclahost")
	if err != nil {
		log.Paincf("Could not Init PKI")
	}
	err = pki.Write()
	if err != nil {
		log.Paincf("Could not Write PKI")
	}
}

func main() {

}
