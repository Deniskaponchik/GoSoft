package main

import (
	"fmt"
	"net"
	"time"
)

func netDialTmt(ipString string) (status string) {
	timeout := 1 * time.Second

	url := ipString + ":" + "80"

	//	Dial("tcp", "golang.org:http")
	//	Dial("tcp", "192.0.2.1:http")
	//	Dial("tcp", "198.51.100.1:80")
	//conn, err := net.DialTimeout("tcp","mysyte:myport", timeout)
	//conn, err := net.DialTimeout("tcp", url, timeout)
	conn, err := net.DialTimeout("tcp", url, timeout)
	if err != nil {
		//log.Println("Site unreachable, error: ", err)
		fmt.Println("Visual не доступен по http")
		status = ""
	} else {
		//log.Println("Site reachable")
		status = "ok"
	}
	conn.Close()
	return
}
