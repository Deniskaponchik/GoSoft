package main

import (
	"fmt"
	"github.com/go-ping/ping"
	"net"
	"os"
	"os/signal"
)

func NOT432423main() {
	//ip4Ping("10.21.178.157")
	//ip4Ping("10.57.178.44")
	//unixPing("10.57.178.44")
	unixPing("10.21.178.157")
}

func ip4Ping(ipString string) {
	//pinger, err := ping.NewPinger("www.google.eytrtyrtyr")
	pinger, err := ping.NewPinger(ipString)
	if err != nil {
		panic(err)
	}

	addr, err := net.ResolveIPAddr("ip4:icmp", ipString)
	if err != nil {
		panic(err)
	}
	fmt.Println("Addr", addr.String())

	pinger.SetPrivileged(true) //for Windows
	pinger.SetIPAddr(addr)
	pinger.Count = 3
	//pinger.Interval = 1

	err = pinger.Run() // Blocks until finished.
	if err != nil {
		panic(err)
	}
	fmt.Sprintln("run пройден")

	stats := pinger.Statistics() // get send/receive/duplicate/rtt stats
	fmt.Println(stats)
	fmt.Sprintln(stats.PacketLoss)

	//if stats.
	//	fmt.Println("reachable")
	//} else {
	//	fmt.Println("unreachable")
	//}

}

func unixPing(ipString string) {
	//pinger, err := ping.NewPinger("www.google.com")
	pinger, err := ping.NewPinger(ipString)
	if err != nil {
		panic(err)
	}
	pinger.SetPrivileged(true) //for Windows
	pinger.Count = 4

	// Listen for Ctrl-C.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			pinger.Stop()
		}
	}()

	pinger.OnRecv = func(pkt *ping.Packet) {
		fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v\n",
			pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)
	}

	pinger.OnDuplicateRecv = func(pkt *ping.Packet) {
		fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v ttl=%v (DUP!)\n",
			pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt, pkt.Ttl)
	}

	pinger.OnFinish = func(stats *ping.Statistics) {
		fmt.Printf("\n--- %s ping statistics ---\n", stats.Addr)
		fmt.Printf("%d packets transmitted, %d packets received, %v%% packet loss\n",
			stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
		fmt.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
			stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
	}

	fmt.Printf("PING %s (%s):\n", pinger.Addr(), pinger.IPAddr())
	err = pinger.Run()
	if err != nil {
		panic(err)
	}
}
