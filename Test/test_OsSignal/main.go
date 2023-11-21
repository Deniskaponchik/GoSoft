package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {

	fmt.Printf("Process PID : %v\n", os.Getpid())

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)
	go counter(c)

	// Block until a signal is received.
	s := <-c
	fmt.Println("Got signal:", s)
}

func counter(c chan os.Signal) {
	count := 1
	for true {
		count++
		fmt.Println(strconv.Itoa(int(count)))
		time.Sleep(10 * time.Second)
	}
}
