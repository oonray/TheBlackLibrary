package main

import (
	"fmt"
	"flag"
	"net/http"
	"os"
	"io"
	"regexp"
	"strings"
)

const (
	url string = "http://standards-oui.ieee.org/oui/oui.txt"
	chache string = "~/.cache/macid_cache"
	separator string = ":"
)

var (
	mac string
	vendors map[string]string
)


func main(){
	flag.StringVar(&mac,"mac","","The Mac addr to be checked")
	flag.Parse()

	if(mac == "") {
		flag.Usage()
		os.Exit(2)
	}
	
	rsp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Could not get Mac IDs %s",err)
		os.Exit(1)
	}

	data, err := io.ReadAll(rsp.Body)
	if err != nil {
		fmt.Printf("Could not get Mac IDs %s",err)
		os.Exit(1)
	}
	
	mac_s := strings.Split(mac,separator)
	if len(mac_s) < 3 {
		fmt.Printf("Invalid mac addr: %s", mac)
		os.Exit(1)
	}
	mcd := fmt.Sprintf("%s-%s-%s\\W+\\(hex\\)\\W+(.*)",mac_s[0],mac_s[1],mac_s[2])
	re, err := regexp.Compile(mcd)

	if err != nil {
		fmt.Printf("Could not compile regex %s",err)
		os.Exit(1)
	}

	found := re.FindAll(data,-1)
	if(found==nil){
		fmt.Printf("Nothing found.")
		os.Exit(1)
	}

	for _,item := range found {
		fmt.Printf("Found:\n%s\n", string(item))
	}
}
