package main

import (
	"github.com/miekg/dns"
	"log"
	"net"
)

var mp []*Map

type Map struct {
	Handler  func(w dns.ResponseWriter, req *dns.Msg)
	Ip       string
	Hostname string
}

func (m *Map) default_handler(w dns.ResponseWriter, req *dns.Msg) {
	var rsp dns.Msg
	rsp.SetReply(req)

	for _, q := range req.Question {
		a := dns.A{
			Hdr: dns.RR_Header{
				Name:   q.Name,
				Rrtype: dns.TypeA,
				Class:  dns.ClassINET,
				Ttl:    0,
			},
			A: net.ParseIP(m.Ip).To4(),
		}
		resp.Answer = append(resp.Answer, &a)
	}
	w.WriteMsg(&rsp)
}

func main() {
	dns.HandleFunc(".", handle)
	log.Fatal(dns.ListenAndServe(":53", "udp", nil))
}
