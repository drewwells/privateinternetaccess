package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
)

const url = "http://ifconfig.me:80"

func main() {

	tr := &http.Transport{
		Dial: func(network, addr string) (net.Conn, error) {
			ief, err := net.InterfaceByName("tun0")
			if err != nil {
				log.Fatal(err)
			}
			addrs, err := ief.Addrs()
			fmt.Println(addrs)
			if err != nil {
				log.Fatal(err)
			}
			tcpAddr := &net.TCPAddr{
				IP: addrs[0].(*net.IPNet).IP,
			}
			d := net.Dialer{LocalAddr: tcpAddr}
			c, err := d.Dial("tcp", "ifconfig.me:80")
			return c, err

			// la, err := net.ResolveTCPAddr("tcp", "10.109.1.6:0")
			// if err != nil {
			// 	log.Fatal("error creating local address", err, la)
			// }
			// ra, err := net.ResolveTCPAddr("tcp", addr)
			// if err != nil {
			// 	log.Fatal("error creating receiving address", err, ra)
			// }
			// c, err := net.DialTCP("tcp", la, ra)
			// fmt.Printf("% #v\n", c)
			// return c, err
		},
	}
	client := &http.Client{Transport: tr}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent:", "curl/7.37.1")
	res, err := client.Do(req)
	if err != nil {
		log.Fatal("failed to request: ", err)
	}
	defer res.Body.Close()
	bs, _ := ioutil.ReadAll(res.Body)
	fmt.Println("Response", string(bs))
}
