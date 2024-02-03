package main

import (
	"fmt"
	"time"
)

func main() {
	//timeZone := +5
	timeZone := 0

	//timeNowU := time.Now()
	timeNowU := time.Now().Add(time.Duration(-5) * time.Hour)

	//fmt.Println(timeNowU.Format("2006-01-02 15:04:05"))
	fmt.Println(timeNowU.Format("2006-01-02 15:00:00"))

	//before30days := timeNowU.Add(time.Duration(-720) * time.Hour).Format("2006-01-02 15:04:05")
	rostovTime := timeNowU.Add(time.Duration(-timeZone) * time.Hour).Format("2006-01-02 15:00:00")
	fmt.Println(rostovTime) //before30days

	novosibTime := timeNowU.Add(time.Duration(-timeZone+4) * time.Hour).Format("2006-01-02 15:04:05")
	fmt.Println(novosibTime) //before30days

}
