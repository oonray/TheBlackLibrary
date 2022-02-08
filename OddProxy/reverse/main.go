
package main

import (
	"os"
	"fmt"
	"flag"
	"errors"
	"bufio"
	"syscall"
	"strings"
	"github.com/armon/go-socks5"
	ssh "golang.org/x/crypto/ssh"
	sshterm "golang.org/x/crypto/ssh/terminal"
	knownhosts "golang.org/x/crypto/ssh/knownhosts"
	log "github.com/sirupsen/logrus"
	//pki "github.com/oonray/BlackHatProgramming/BlackHatGo/PKI"
)

var (
	host *string
	port *string
	user *string

	public *string

	khostfile *string

	s5conf *socks5.Config = &socks5.Config{}

	khosts ssh.HostKeyCallback
)

func getCreds()(string, string){
	var username string

	if *user == "" { 
		reader := bufio.NewReader(os.Stdin)

		fmt.Printf("Enter ssh Username:")
		username, _ = reader.ReadString('\n')
	} else {username = *user}

	fmt.Printf("Enter ssh Password:")
	bytePass, err := sshterm.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Errorf("Could not get password: %s",err)
		return username,""
	}

	return strings.TrimSpace(username), strings.TrimSpace(string(bytePass))
}

func Ssh_Config(user string, pass string) *ssh.ClientConfig{
	khosts, err := knownhosts.New(*khostfile)
	if err != nil {
		log.Fatalf("Could not read hosts file")
		return nil
	}

	return &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pass),
		},
		HostKeyCallback: khosts,
	}
}

func argparse() error {
	host = flag.String("H","","The host to conenct to")
	port = flag.String("P","","The port to connect to")
	user = flag.String("u","","The user to connect with")
	public = flag.String("pk","","Public key")
	khostfile = flag.String("hf","~/.ssh/known_hosts","Hosts file to use")
	flag.Parse()

	if *host == "" {
		return errors.New("Host Required")
	}

	if *port == "" {
		return errors.New("Port Required")
	}

	if *public != "" {
		return errors.New("PKI auth not implemented yet")
	}

	return nil
}

func main() {
	err := argparse()
	if err != nil {
		log.Errorf("%s",err)
		flag.PrintDefaults()
		return
	}

	user, pass := getCreds()
	fmt.Print("\n")

	conf := Ssh_Config(user,pass)
	con, err := ssh.Dial("tcp",fmt.Sprintf("%s:%s",*host,*port),conf)
	if err != nil {
		log.Fatalf("Could not connect %s:%s | %s",*host,*port,err)
		return
	}
	defer con.Close()

	log.Infof("Connected to %s@%s:%s",user,*host,*port)

	server, err := socks5.New(s5conf)
	if err != nil {
		log.Fatalf("Could not start socks server | %s",err)
		return 
	}

	if err := server.Serve(con); err != nil {
		log.Fatalf("Socks server error | %s",err)
		return 
	}
}
