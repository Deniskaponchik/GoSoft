package main

/*
import (
	"bytes"
	"encoding/xml"
	"log"
	"net/http"
)

type Envelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Header  *Header  `xml:"Header,omitempty"`
	Body    Body     `xml:"Body"`
}

type Header struct{}

type Body struct {
	XMLName    xml.Name `xml:"Body"`
	AddRequest AddRequest
}

type AddRequest struct {
	XMLName xml.Name `xml:"AddRequest"`
	//Arg1    int      `xml:"arg1"`
	Arg1 string `xml:"WebTutor"`
	//Arg2    int      `xml:"arg2"`
}

func main() {
	params := &AddRequest{
		//Arg1: 4,
		Arg1: "WebTutor",
		//Arg2: 5,
	}

	soapRequest := &Envelope{
		Body: Body{
			AddRequest: *params,
		},
	}

	xmlBytes, err := xml.Marshal(soapRequest)
	if err != nil {
		log.Fatal(err)
	}

	url := "http://10.246.37.15:8060/specs/aoi/tele2/bpm/bpmPortType"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(xmlBytes))
	if err != nil {
		log.Fatal(err)
	}

	soapAction := "readSystemsRequest"
	req.Header.Set("Content-Type", "text/xml")
	req.Header.Set("SOAPAction", soapAction)
	//req.Header.Set("SOAPAction", "readSystemsRequest")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	// handle response here
}*/
