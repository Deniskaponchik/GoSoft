package main

//https://gist.github.com/buroz/d4cc4b8016e5c21817b6704e93d54023
import (
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type SoapHeader struct {
	XMLName xml.Name `xml:"x:Header"`
}

type SoapBody struct {
	XMLName xml.Name `xml:"x:Body"`
	Request interface{}
}

type SoapRoot struct {
	XMLName xml.Name `xml:"x:Envelope"`
	X       string   `xml:"xmlns:x,attr"`
	Sch     string   `xml:"xmlns:sch,attr"`
	Header  SoapHeader
	Body    SoapBody
}

type GetCitiesRequest struct {
	XMLName xml.Name `xml:"sch:GetCitiesRequest"`
}

type GetCitiesResponse struct {
	XMLName xml.Name `xml:"ns3:GetCitiesResponse"`
	result  struct{} `xml:result`
	cities  struct{} `xml:cities`
}

func SoapCall(service string, request interface{}) string {
	var root = SoapRoot{}
	root.X = "http://schemas.xmlsoap.org/soap/envelope/"
	root.Sch = "http://www.n11.com/ws/schemas"
	root.Header = SoapHeader{}
	root.Body = SoapBody{}
	root.Body.Request = request

	out, _ := xml.MarshalIndent(&root, " ", "  ")
	body := string(out)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	response, err := client.Post(service, "text/xml", bytes.NewBufferString(body))

	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()

	content, _ := ioutil.ReadAll(response.Body)
	s := strings.TrimSpace(string(content))
	return s
}

/*
func main() {
	SoapCall("https://api.n11.com/ws/CityService.wsdl", GetCitiesRequest{})
}*/
