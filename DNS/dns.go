package main

import (
	//	"flag"
	"flag"
	"fmt"
	"github.com/miekg/dns"
	"log"
	"net"
	"strings"
	"sync"
)

var (
	hsts      Hosts = Hosts{Hosts: map[string]string{}}
	hsts_lock sync.RWMutex
	dnsserver string
	host      string
	port      string
	udp       bool
	proto     string = "tcp"
)

type Hosts struct {
	Hosts map[string]string
}

func (h *Hosts) Set(value string) error {
	data := strings.Split(value, ":")
	h.Hosts[data[0]] = data[1]
	log.Printf("%v", h.Hosts)
	return nil
}

func default_handler(w dns.ResponseWriter, req *dns.Msg) {
	var a dns.Msg = dns.Msg{}
	var A dns.A = dns.A{
		Hdr: dns.RR_Header{
			Name:   req.Question[0].Name,
			Rrtype: dns.TypeA,
			Class:  dns.ClassINET,
			Ttl:    0,
		},
	}

	a.Answer = append(a.Answer, &A)
	a.SetReply(req)

	if len(req.Question) < 1 {
		dns.HandleFailed(w, req)
	}

	name := strings.Trim(req.Question[0].Name, ".")

	hsts_lock.RLock()
	out, in := hsts.Hosts[name]
	hsts_lock.RUnlock()

	if !in {
		resp, err := dns.Exchange(req, dnsserver)
		if err != nil {
			dns.HandleFailed(w, req)
			return
		}
		a = *resp
	} else {
		A.A = net.ParseIP(out).To4()
	}

	if err := w.WriteMsg(&a); err != nil {
		dns.HandleFailed(w, req)
		return
	}
}

func init() {
	flag.StringVar(&host, "H", "127.0.0.1", "The host to bind the server to")
	flag.StringVar(&port, "P", "53", "The port to bind the server to")
	flag.StringVar(&dnsserver, "DNS", "8.8.8.8:53", "The upstream DNS server")
	flag.Parse()
}

func main() {
	for _, data := range flag.Args() {
		hsts.Set(data)
	}
	dns.HandleFunc(".", default_handler)
	log.Fatal(dns.ListenAndServe(fmt.Sprintf("%s:%s", host, port), "udp", nil))
}
