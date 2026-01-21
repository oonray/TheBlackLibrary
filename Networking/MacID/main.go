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
	cache string = "~/.cache/macid_cache"
	separator string = ":"
)

var (
	mac string
	local bool
	file_path string
	vendors map[string]string
	data []byte
)

func get_http() error {
	var err error

	rsp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Could not get Mac IDs %s",err)
		os.Exit(1)
	}

	if rsp.StatusCode != 200 {
		return fmt.Errorf("Error status %s",rsp.Status)
	}

	data, err = io.ReadAll(rsp.Body)
	if err != nil {
		return err
	}

	return nil
}

func get_file() error {
	var err error
	_, perr := os.Stat(file_path)
	if perr != nil {
		return perr
	}
	data, err =  os.ReadFile(file_path)
	return err
}
func write_file() error {
	return os.WriteFile(file_path,data,0744)
}


func main(){
	flag.StringVar(&mac,"mac","","The Mac addr to be checked")
	flag.StringVar(&file_path,"file",cache,"The file to use")
	flag.BoolVar(&local,"local",false,"Use local file")
	flag.Parse()

	if(mac == "") {
		flag.Usage()
		os.Exit(2)
	}
	
	mac_s := strings.Split(mac,separator)
	if len(mac_s) < 3 {
		fmt.Printf("Invalid mac addr: %s", mac)
		os.Exit(1)
	}

	err := get_http()
	if err != nil {
		fmt.Printf("Could not fetch data: %s\n", err)
		err := get_file()
		if err != nil {
			fmt.Printf("Could not fetch file data: %s\n", err)
			os.Exit(1)
		}
	}

	err = write_file()
	if err != nil {
		fmt.Printf("Could not write to file: %s\n", err)
		os.Exit(1)
	}

	if len(data) < 500 {
		fmt.Printf("%s",data)
	}

	mcd := fmt.Sprintf("%s-%s-%s\\W+\\(hex\\)\\W+.*",mac_s[0],mac_s[1],mac_s[2])
	re, err := regexp.Compile(mcd)
	if err != nil {
		fmt.Printf("Could not compile regex %s",err)
		os.Exit(1)
	}

	found := re.Find(data)
	if(found == nil){
		fmt.Printf("Nothing found.")
		os.Exit(1)
	}

	fmt.Printf("%s\n",found)
}
