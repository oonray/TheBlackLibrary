package main

import (
	"errors"
	"flag"
	"fmt"
	logR "github.com/sirupsen/logrus"
	"github.com/things-go/go-socks5"
	"log"
	"os"
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

	server := socks5.NewServer(
		socks5.WithLogger(socks5.NewLogger(log.New(os.Stdout, "socks5: ", log.LstdFlags))),
	)

	if err := server.ListenAndServe("tcp", fmt.Sprintf(":%s", *lp)); err != nil {
		logR.Fatalf("%s", err)
	}
}
