package main

import "fmt"

type Client struct {
	names []string
}

func main() {
	var client *Client
	namesSlice := []string{"denis"}

	client = &Client{
		namesSlice,
	}
	
	client.names = append(client.names, "vasya")

	for _, name := range client.names {
		fmt.Println(name)
	}
}
