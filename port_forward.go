package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {
	// Discover the tun0 IP
	var ip string
	intfs, _ := net.Interfaces()
	for _, intf := range intfs {
		//fmt.Printf("%s % #v\n", intf.HardwareAddr.String(), intf)
		if intf.Name == "tun0" {
			addrs, err := intf.Addrs()
			if err != nil {
				log.Fatal(err)
			}
			for _, addr := range addrs {
				ip = strings.Split(addr.String(), "/")[0]
			}
		}
	}

	// Check for passed username and password
	if len(os.Args) < 3 {
		log.Fatal("call with ./port-forward user pass")
	}
	user := os.Args[1]
	pass := os.Args[2]
	host, _ := os.Hostname()
	seed := user + pass + host
	hash := fmt.Sprintf("%x", md5.Sum([]byte(seed)))
	vals := url.Values{}
	vals.Add("user", user)
	vals.Add("pass", pass)
	vals.Add("client_id", hash)
	vals.Add("local_ip", ip)
	fmt.Fprintf(os.Stderr, "client_id: %s\n", hash)
	fmt.Fprintf(os.Stderr, "IP found: %s\n", ip)
	fmt.Println(vals.Encode())
	res, err := http.PostForm("https://www.privateinternetaccess.com/vpninfo/port_forward_assignment", vals)

	if err != nil {
		log.Fatal("Error on response", err)
	}

	if res.StatusCode != 200 {
		resp, _ := ioutil.ReadAll(res.Body)
		log.Fatal(res.StatusCode, resp)
	}

	//robots, err := ioutil.ReadAll(res.Body)
	//fmt.Println(string(robots))
	defer res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	dec := json.NewDecoder(res.Body)
	var r Response

	dec.Decode(&r)
	if r.Error != "" {
		log.Fatal("Error retrieving port: ", r.Error)
	}
	fmt.Println(r.Port)
}

type Response struct {
	Port  float64 `json:"port"`
	Error string  `json:"error"`
}
