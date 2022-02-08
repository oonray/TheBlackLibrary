package main

import (
	"os"
	"log"
	"fmt"
	"flag"
	"errors"
    "github.com/things-go/go-socks5"
	logR "github.com/sirupsen/logrus"
)

var (
	lp *string
)

func argparse() error {
	lp = flag.String("lp","","Port to listen to")
	flag.Parse()

	if *lp == "" {
		return errors.New("Listen port must be declared")
	}	
	return nil
}

func main() {
	err := argparse()
	if err != nil {
		logR.Error(err)
		return
	}

	server := socks5.NewServer(
		socks5.WithLogger(socks5.NewLogger(log.New(os.Stdout, "socks5: ", log.LstdFlags))),
	)

	if err := server.ListenAndServe("tcp",fmt.Sprintf(":%s",*lp));err != nil {
		logR.Fatalf("%s",err)
	}
}
