package main

import (
	"encoding/xml"
)

type FooResponse struct {
	Bar string
}
type FooRequest struct {
	XMLName xml.Name `xml:"fooRequest"`
	//XMLName xml.Name `xml:"SmacWiFiКeadServiceFields"`
	Foo string
}

/*
func main() {
	//client := soap.NewClient("http://127.0.0.1:8080/", nil)
	client := soap.NewClient("http://10.246.37.15:8060/specs/aoi/tele2/bpm/bpmPortType", nil)

	response := &FooResponse{}
	//httpResponse, err := client.Call( "operationFoo", &FooRequest{Foo: "hello i am foo"}, response)
	httpResponse, err := client.Call(context.TODO(), "operationFoo", &FooRequest{Foo: "WebTutor"}, response)
	if err != nil {
		panic(err)
	}
	log.Println(response.Bar, httpResponse.Status)
}*/
