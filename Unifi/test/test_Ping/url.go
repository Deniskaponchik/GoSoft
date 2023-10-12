package main

import (
	"log"
	"net"
	"time"
)

func main() {
	netDialTmt("10.21.178.157")
	//netDialTmt("10.57.178.147")
}

func netDialTmt(ipString string) {
	timeout := 1 * time.Second

	url := ipString + ":" + "80"

	//	Dial("tcp", "golang.org:netdial")
	//	Dial("tcp", "192.0.2.1:netdial")
	//	Dial("tcp", "198.51.100.1:80")
	//conn, err := netdial.DialTimeout("tcp","mysyte:myport", timeout)
	//conn, err := netdial.DialTimeout("tcp", url, timeout)
	_, err := net.DialTimeout("tcp", url, timeout)
	if err != nil {
		log.Println("Site unreachable, error: ", err)
	} else {
		log.Println("Site reachable")
	}
}
