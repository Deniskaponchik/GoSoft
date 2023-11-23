package main

import "fmt"

type Client struct {
	name string
}

func main() {
	macClient := make(map[string]*Client)
	hostClient := make(map[string]*Client)

	client := &Client{
		"denis",
	}

	macClient["abcd"] = client
	hostClient["wsir"] = client

	client1 := macClient["abcd"]
	client2 := hostClient["wsir"]

	fmt.Println(client1.name)
	fmt.Println(client2.name)

	client1.name = "vasya"

	fmt.Println(client1.name)
	fmt.Println(client2.name)
}
