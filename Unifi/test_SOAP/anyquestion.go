package main

import "encoding/xml"

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
	Arg1    int      `xml:"arg1"`
	Arg2    int      `xml:"arg2"`
}

params := &AddRequest{
Arg1: 4,
Arg2: 5,
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

req, err := http.NewRequest("POST", url, bytes.NewBuffer(xmlBytes))
if err != nil {
log.Fatal(err)
}

req.Header.Set("Content-Type", "text/xml")
req.Header.Set("SOAPAction", soapAction)

client := &http.Client{}
resp, err := client.Do(req)
if err != nil {
log.Fatal(err)
}

defer resp.Body.Close()
// handle response here