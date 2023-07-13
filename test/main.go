package main

import (
	"fmt"
	"time"
)

func main() {
	/*
		clientMacName := map[string]string{}
		clientMacName["00:01"] = "ClientOne"
		clientMacName["00:01"] = "ClientTwo"
		fmt.Println(clientMacName["00:01"])
	*/

	for true {
		if time.Now().Minute() == 19 {
			fmt.Println(19)
		} else {
			fmt.Println("wait 1 minute")
		}
		time.Sleep(60 * time.Second)
	}
}
