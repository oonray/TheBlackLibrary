package tcp

import (
    "fmt"
    "net"
    "bytes"
    "sync"
    "context"
)

const (
    DEFAULT_TCP_ADDRESS = "0.0.0.0"
    DEFAULT_TCP_PORT = "42969"
    DEFAULT_TYPE = "tcp"
)

type (
    Coms struct {
        In  bytes.Buffer
        Out bytes.Buffer
        err chan error
    }

    Addr struct {
        Ip      *net.IPAddr
        Port    string
    }

    Client struct {
        Socket  *net.Conn
        Address *Addr
        Server  *Server
        Buffer  *Coms
    }

    Server struct {
        Listener net.Listener
        Address  *Addr
        Clients  *[]*Client
        wg       sync.WaitGroup
        ctx      context.Context
        Cancel   context.CancelFunc
    }
)

func NewComs() *Coms {
    return &Coms{}
}

func NewAddr(address string, port string) (*Addr,error) {
    ip, err := net.ResolveIPAddr(DEFAULT_TYPE,address)
    if(err!=nil){return nil,err}

    return &Addr{
        Ip: ip,
        Port: port,
    },nil
}

func (*Addr)   Network() string { return DEFAULT_TYPE }
func (s *Addr) String()  string { return fmt.Sprintf("%s:%s",s.Ip.String(),s.Port) }

func NewServer(ip string, port string) (*Server, error) {
    var err error

    addr, err := NewAddr(ip,port)
    if(err!=nil){return nil,err}

    srv := &Server{
        Listener:  nil,
        Address: addr,
        Clients: new([]*Client),
    }

    srv.ctx, srv.Cancel = context.WithCancel(context.Background())

    return srv,err
}

func NewServerFromAddr(addr *Addr) *Server{
    return &Server{
        Listener:  nil,
        Address: addr,
        Clients: new([]*Client),
    }
}

func (s *Server) Accept() error {
    var err error
    s.Listener, err = net.Listen(s.Address.Network(),s.Address.String())
    return err
}

func (s *Server) ListenAndServe() error {
    var err error
    s.Listener, err = net.Listen(s.Address.Network(),s.Address.String())
    if(err!=nil){return err}
    defer s.Listener.Close()

    s.wg.Go(func(){
        for {
            select {
            case <-s.ctx.Done():
                return
            default:
                con, err :=  s.Listener.Accept()
                if(err!=nil){return}

                s.Clients := append(NewClient())
            }
        }
    })

    return err
}

