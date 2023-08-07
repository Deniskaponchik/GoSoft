package main

import "fmt"

func main() {
	myMap := map[string]bool{}
	myMap["a"] = true
	myMap["b"] = true
	if myMap["c"] {
		fmt.Println("c is exist")
	} else {
		fmt.Println("c is NOT exist")
	}

}
