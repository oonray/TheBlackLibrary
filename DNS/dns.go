package main

import (
	//	"flag"
	"flag"
	"fmt"
	"github.com/miekg/dns"
	"log"
	"net"
	"strings"
)

var (
	hsts  Hosts = Hosts{Hosts: map[string]string{}}
	host  string
	port  string
	udp   bool
	proto string = "tcp"
)

type Hosts struct {
	Hosts map[string]string
}

func (h *Hosts) Set(value string) error {
	data := strings.Split(value, ":")
	log.Printf("%v", data)
	h.Hosts[data[0]] = data[1]
	log.Printf("%v", h.Hosts)
	return nil
}

func default_handler(w dns.ResponseWriter, req *dns.Msg) {
	var a dns.Msg = dns.Msg{}
	a.SetReply(req)

	for _, data := range req.Question {
		out, in := hsts.Hosts[strings.Trim(data.Name, ".")]
		if in {
			log.Printf("found %s", data.Name)
			answer := dns.A{
				Hdr: dns.RR_Header{
					Name:   data.Name,
					Rrtype: dns.TypeA,
					Class:  dns.ClassINET,
					Ttl:    0,
				},
				A: net.ParseIP(out).To4(),
			}
			a.Answer = append(a.Answer, &answer)
		}
		if !in {
			log.Printf("not found %s", data.Name)
			ip, err := net.ResolveIPAddr("ip", data.Name)
			if err == nil {
				answer := dns.A{
					Hdr: dns.RR_Header{
						Name:   data.Name,
						Rrtype: dns.TypeA,
						Class:  dns.ClassINET,
						Ttl:    0,
					},
					A: ip.IP,
				}
				a.Answer = append(a.Answer, &answer)
			}
		}
	}
	w.WriteMsg(&a)
}

func init() {
	flag.StringVar(&host, "H", "127.0.0.1", "The host to bind the server to")
	flag.StringVar(&port, "P", "53", "The port to bind the server to")
	flag.Parse()
}

func main() {
	for _, data := range flag.Args() {
		hsts.Set(data)
	}
	dns.HandleFunc(".", default_handler)
	log.Fatal(dns.ListenAndServe(fmt.Sprintf("%s:%s", host, port), "udp", nil))
}
