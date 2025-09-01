package tcp

import (
    "net"
    "testing"
)

func BenchmarkServer(b *testing.B){


}

func TestServer(t *testing.T){


}

func TestListener(t *testing.T){
    addr,err := NewAddr(DEFAULT_TCP_ADDRESS, DEFAULT_TCP_PORT)
    if(err!=nil){t.Fatalf("Could not Parse address %s %s",DEFAULT_TCP_ADDRESS,err)}

    con, err := net.Listen(addr.Network(), addr.String())
    if(err!=nil){t.Fatalf("Could not Listen on %s %s",addr.String(),err)}
    defer con.Close()

    t.Logf("Listening on %s:%s",DEFAULT_TCP_ADDRESS,DEFAULT_TCP_PORT)
}
