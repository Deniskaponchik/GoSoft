package ping

import (
	"fmt"
	"github.com/go-ping/ping"
)

func ipPing(ip string) {
	pinger, err := ping.NewPinger(ip)
	if err != nil {
		panic(err)
	}
	pinger.Count = 3
	err = pinger.Run() // Blocks until finished.
	if err != nil {
		panic(err)
	}
	stats := pinger.Statistics() // get send/receive/duplicate/rtt stats

	fmt.Println(stats)
}