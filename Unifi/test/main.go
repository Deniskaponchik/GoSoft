package main

import (
	"fmt"
	"time"
)

func main() {

	currentTime := time.Now()

	fmt.Println(time.Now().Hour(), ":", time.Now().Minute())
	//fmt.Println(time.Now().Format("dd-MONTH hh:mm"))
	fmt.Println(currentTime.Format("02 January, 15:04:05"))
}
