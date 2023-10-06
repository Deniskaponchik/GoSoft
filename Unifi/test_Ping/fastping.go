package main

import (
	"fmt"
	"github.com/tatsushid/go-fastping"
	"net"
	"os"
	"time"
)

func main546456546456454() {
	fastPing("10.21.178.157")
	//fastPing("10.57.178.41")
}

func fastPing(ipString string) {
	var success bool

	p := fastping.NewPinger()

	//ra, err := net.ResolveIPAddr("ip4:icmp", os.Args[1])
	ra, err := net.ResolveIPAddr("ip4:icmp", ipString)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	p.AddIPAddr(ra)
	//p.Size = 4

	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		fmt.Printf("IP Addr: %s receive, RTT: %v\n", addr.String(), rtt)
		success = true
	}

	p.OnIdle = func() {
		fmt.Println("finish")
	}

	err = p.Run()
	if err != nil {
		fmt.Println(err)
	}

	if success == true {
		fmt.Println("success")
	} else {
		fmt.Println("провал")
	}

}
