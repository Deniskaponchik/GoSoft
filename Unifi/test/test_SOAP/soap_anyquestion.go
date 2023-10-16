package main

//Вроде что-то отправляет, но ОТВЕТ не обработан
import (
	"bytes"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
)

type Envelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Header  *Header  `xml:"Header,omitempty"`
	Body    Body     `xml:"Body"`
}

type Header struct{}

/*
	type Body struct {
		XMLName    xml.Name `xml:"Body"`
		AddRequest AddRequest
	}
*/
type Body struct {
	XMLName            xml.Name `xml:"Body"`
	readSystemsRequest readSystemsRequest
}

/*
	type AddRequest struct {
		XMLName xml.Name `xml:"AddRequest"`
		Arg1    int      `xml:"arg1"`
		Arg2    int      `xml:"arg2"`
	}
*/
type readSystemsRequest struct {
	XMLName xml.Name `xml:"readSystemsRequest"`
	Top     int      `xml:"Top"`
	Skip    int      `xml:"Skip"`
	Filter  string   `xml:"Filter"`
}

func NOT32423442main() {
	/*
		params := &AddRequest{
			Arg1: 4,
			Arg2: 5,
		}*/
	params := &readSystemsRequest{
		//Top: 0,
		//Skip: 5,
		Filter: "WebTutor",
	}

	soapRequest := &Envelope{
		Body: Body{
			//AddRequest: *params,
			readSystemsRequest: *params,
		},
	}

	xmlBytes, err := xml.Marshal(soapRequest)
	if err != nil {
		log.Fatal(err)
	}

	url := "http://10.246.37.15:8060/specs/aoi/tele2/bpm/bpmPortType"
	//req, err := http.NewRequest("POST", url, bytes.NewBuffer(xmlBytes))
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(xmlBytes))
	if err != nil {
		log.Fatal(err)
	}

	soapAction := "urn:readSystemsRequest"
	req.Header.Set("Content-Type", "text/xml")
	req.Header.Set("SOAPAction", soapAction)
	//req.Header.Set("SOAPAction", "readSystemsRequest")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp)

	defer resp.Body.Close()
	// handle response here

}
