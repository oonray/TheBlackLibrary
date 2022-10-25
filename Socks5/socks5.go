package main

import (
	"errors"
	"flag"
	"fmt"
	socks5 "github.com/armon/go-socks5"
	logR "github.com/sirupsen/logrus"
)

var (
	lp *string
)

func argparse() error {
	lp = flag.String("lp", "1080", "Port to listen to")
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

	server, err := socks5.New(&socks5.Config{})
	if err != nil {
		logR.Error(err)
		return
	}

	if err := server.ListenAndServe("tcp", fmt.Sprintf(":%s", *lp)); err != nil {
		logR.Fatalf("%s", err)
	}
}
