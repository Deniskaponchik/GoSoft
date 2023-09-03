package main

import (
	"fmt"
	"os"
)

func main() {
	// Store the PATH environment variable in a variable
	path, exists := os.LookupEnv("UNIFI_USERNAME")

	if exists {
		// Print the value of the environment variable
		fmt.Print(path)
	}
}

func test12() {
	var a int
	var b int

	fmt.Scan(&a)
	fmt.Scan(&b)

	fmt.Println(a + b)
}

func test11() {
	var a int
	var b int

	_, err1 := fmt.Scanf("%d", &a)
	_, err2 := fmt.Scanf("%d", &b)

	if err1 == nil && err2 == nil {
		fmt.Println(a + b)
	}
}

func test10() {
	var a int
	var b int

	fmt.Scanf("%d", &a)
	fmt.Scanf("%d", &b)

	fmt.Println(a + b)

	//_, err := fmt.Scanf("%d", &a)
	//_, err := fmt.Scanf("%d", &b)
}
