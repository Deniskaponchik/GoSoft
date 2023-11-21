package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func counter1(c chan os.Signal) {
	count := 1
	for true {
		count++
		fmt.Println(strconv.Itoa(int(count)))
		time.Sleep(10 * time.Second)
	}

	// Block until a signal is received.
	s := <-c
	fmt.Println("Got signal:", s)
}
