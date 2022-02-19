package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
)

var (
	command   *string
	host      *string
	container *string
	pod       *string
	port      *int
	namespace *string
)

func init() {
	host = flag.String("H", "localhost", "The host to query")
	command = flag.String("c", "ls", "the command to run")
	pod = flag.String("pod", "nginx", "the pod to use")
	container = flag.String("C", "nginx", "the container to use")
	port = flag.Int("p", 10250, "the port to connect to")
	namespace = flag.String("n", "default", "the namespace to use")
	flag.Parse()

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
}

func main() {
	log.Infof("Running %s on %s", *host, *command)

	t_url := url.URL{
		Scheme:   "https",
		Host:     fmt.Sprintf("%s:%d", *host, *port),
		Path:     fmt.Sprintf("/run/%s/%s/%s", *namespace, *pod, *container),
		RawQuery: fmt.Sprintf("cmd=%s", *command),
	}

	resp, err := http.Get(t_url.String())
	if err != nil {
		log.Errorf("Could not run %s : %s", *command, err)
		return
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Could get output from %s : %s", *command, err)
		return
	}

	log.Infof("==== OUT ====\n\n %s \n\n ==== OUT ====\n", data)
}
