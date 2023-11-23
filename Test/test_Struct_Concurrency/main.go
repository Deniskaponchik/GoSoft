package main

import (
	"fmt"
	"strconv"
	"time"
)

type Counter struct {
	number int
}

func main() {
	counter := &Counter{
		number: 1,
	}
	go inc(counter)
	go read(counter)

	time.Sleep(58 * time.Second)
}

func inc(с *Counter) {
	for {
		с.number++
	}
}
func read(c *Counter) {
	for {
		fmt.Println(strconv.Itoa(c.number))
	}
}
