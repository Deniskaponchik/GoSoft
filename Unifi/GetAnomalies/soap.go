package main

/*
import (
	"encoding/xml"
	"github.com/globusdigital/soap"
	"log"
)

func notMainmain() {
	client := soap.NewClient("http://127.0.0.1:8080/", nil)
	response := &FooResponse{}
	httpResponse, err := client.Call("operationFoo", &FooRequest{Foo: "hello i am foo"}, response)
	if err != nil {
		panic(err)
	}
	log.Println(response.Bar, httpResponse.Status)
}

// FooRequest a simple request
type FooRequest struct {
	XMLName xml.Name `xml:"fooRequest"`
	Foo     string
}

// FooResponse a simple response
type FooResponse struct {
	Bar string
}

type Soap struct {
	ServerName     string
	login          string
	password       string
	BpmSystemName  string
	BpmServiceName string
}

*/
