package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func execPing() {
	//out, _ := exec.Command("ping", "10.21.178.157", "-c 5", "-i 3", "-w 10").Output()
	out, _ := exec.Command("ping", "10.21.178.157").Output()
	if strings.Contains(string(out), "Destination Host Unreachable") {
		fmt.Println("TANGO DOWN")
	} else {
		fmt.Println("IT'S ALIVEEE")
	}
}
