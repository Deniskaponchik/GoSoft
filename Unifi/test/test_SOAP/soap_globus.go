package main

//"github.com/globusdigital/http"
import (
	"encoding/xml"
)

type FooResponse struct {
	Bar string
}
type FooRequest struct {
	//XMLName xml.Name `xml:"fooRequest"`
	XMLName xml.Name `xml:"readSystemsRequest"`
	//Foo     string
	Filter string
}

/*
func main() {
	//client := http.NewClient("http://127.0.0.1:8080/", nil)
	client := http.NewClient("http://10.246.37.15:8060/specs/aoi/tele2/bpm/bpmPortType", nil)

	response := &FooResponse{}
	//httpResponse, err := client.Call( "operationFoo", &FooRequest{Foo: "hello i am foo"}, response)
	httpResponse, err := client.Call(context.TODO(),
		"readSystemsRequest",
		&FooRequest{Filter: "WebTutor"},
		response)
	if err != nil {
		panic(err)
	}
	log.Println(response.Bar, httpResponse.Status)
}*/
