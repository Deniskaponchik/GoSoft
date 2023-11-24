package main

import "fmt"

type Client struct {
	name string
}

func main() {
	macClient := make(map[string]*Client)
	//hostClient := make(map[string]*Client)

	var client *Client
	sliceName := []string{"denis", "vasya", "petya"}

	for i := 0; i < 3; i++ {
		client = &Client{
			sliceName[i],
		}
		macClient[sliceName[i]] = client
	}

	for _, mapClient := range macClient {
		fmt.Println(mapClient.name)
	}

}
