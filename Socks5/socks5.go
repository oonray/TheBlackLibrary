package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	socks5 "github.com/armon/go-socks5"
	logR "github.com/sirupsen/logrus"
	"net"
	"time"
)

var (
	lp    *string
	reslv *strin

type DirectResolver struct {
	Reslv net.Resolver
}

func (d DirectResolver) Resolve(ctx context.Context, name string) (context.Context, net.IP, error) {
	addr, err := d.Reslv.LookupIP(ctx, "ip", name)
	if err != nil {
		return ctx, nil, err
	}
	return ctx, addr[0], err
}

func argparse() error {
	lp = flag.String("lp", "1080", "Port to listen to")
	reslv = flag.String("rslv", "", "The addr/port of the resolver to use")
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

	if *reslv != "" {
		conf.Resolver = DirectResolver{
			Reslv: net.Resolver{
				PreferGo: true,
				Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
					d := net.Dialer{
						Timeout: time.Millisecond * time.Duration(10000),
					}
					return d.DialContext(ctx, network, *reslv)
				},
			},
		}
	}

	server, err := socks5.New(&conf)
	if err != nil {
		logR.Error(err)
		return
	}

	if err := server.ListenAndServe("tcp", fmt.Sprintf(":%s", *lp)); err != nil {
		logR.Fatalf("%s", err)
	}
}
