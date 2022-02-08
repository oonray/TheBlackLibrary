package main

import (
	"fmt"
	"flag"
	"net"
	"bufio"
	"io"
	log "github.com/sirupsen/logrus"
)

var (
	in net.Listener
	out net.Conn

	lport *string

	host *string
	port *string

	cons string
	lcons string

	done bool
)

func argparse(){
	lport = flag.String("lp","","The port to listen on")
	host = flag.String("H","","The host to conenct to")
	port = flag.String("P","","The port to connect to")
	flag.Parse()
}

func main() {
	var err error
	defer func(){done=true}()

	argparse()

	cons = fmt.Sprintf("%s:%s",*host,*port)
	lcons = fmt.Sprintf(":%s",*lport)

	log.Infof("Listening on %s",lcons)
	in, err = net.Listen("tcp",lcons)
	if err != nil {
		log.Infof("Could not listen: %s",err)
		return
	}
	defer in.Close()

	for {
			con, err := in.Accept()
			if err != nil {
				log.Errorf("Could not accept: %s",err)
				continue}
			defer con.Close()

			out, err = net.Dial("tcp",cons)
			if err != nil {
				log.Errorf("Could not connect: %s",err)
				continue}

			cerr := make(chan error)

			reader := []io.Reader{}
			reader = append(reader,bufio.NewReader(con))
			reader = append(reader,bufio.NewReader(out))

			writer := []io.Writer{}
			writer = append(writer,bufio.NewWriter(out))
			writer = append(writer,bufio.NewWriter(con))

			go func(){
				for {
					_, err := io.Copy(writer[0],reader[0])
					if err != nil {
						cerr <- err
						break}
					if done {break}
				}
			}()

			go func(){
				for {
					_, err := io.Copy(writer[1],reader[1])
					if err != nil {
						cerr <- err
						break}
					if done {break}
				}
			}()

			go func(){
				for {
					err := <-cerr 
					log.Infof("IO Error: %s",err)
					if done {break}
				}
			}()
	}
}

