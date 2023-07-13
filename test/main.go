package main

import (
	"fmt"
)

func main() {

	clientMacName := map[string]string{}

	clientMacName["00:01"] = "ClientOne"
	clientMacName["00:01"] = "ClientTwo"

	fmt.Println(clientMacName["00:01"])
}
