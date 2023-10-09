package main

/*
import (
	"encoding/xml"
	"github.com/tiaguinho/gosoap"
	"log"
	"net/http"
	"time"
)

// GetIPLocationResponse will hold the Soap response
type GetIPLocationResponse struct {
	GetIPLocationResult string `xml:"GetIpLocationResult"`
}
type ReadSystemsResponse struct {
	ReadSystemsResult string `xml:"ReadSystemsResult"`
}

// GetIPLocationResult will
type GetIPLocationResult struct {
	XMLName xml.Name `xml:"GeoIP"`
	Country string   `xml:"Country"`
	State   string   `xml:"State"`
}
type ReadSystemsResult struct {
	XMLName xml.Name `xml:"ber-ns0:row"`
	Name    string   `xml:"ber-ns0:cell"`
	Id      string   `xml:"ber-ns0:cell"`
}

var (
	//r GetIPLocationResponse
	r ReadSystemsResponse
)

func main4353454353443() {
	httpClient := &http.Client{
		Timeout: 1500 * time.Millisecond,
	}
	gosoap.SetCustomEnvelope("soapenv", map[string]string{
		"xmlns:soapenv": "http://schemas.xmlsoap.org/soap/envelope/",
		//"xmlns:tem": "http://tempuri.org/",
		"xmlns:bpm": "http://www.bercut.com/specs/aoi/tele2/bpm",
	})

	//http, err := gosoap.SoapClient("http://wsgeoip.lavasoft.com/ipservice.asmx?WSDL", httpClient)
	http, err := gosoap.SoapClient("http://10.246.37.15:8060/specs/aoi/tele2/bpm/bpmPortType?wsdl", httpClient)
	if err != nil {
		log.Fatalf("SoapClient error: %s", err)
	}

	// Use gosoap.ArrayParams to support fixed position params
	params := gosoap.Params{
		//"sIp": "8.8.8.8",
		"Filter": "WebTutor",
	}

	//res, err := http.Call("GetIpLocation", params)
	res, err := http.Call("ReadSystems", params)
	if err != nil {
		log.Fatalf("Call error: %s", err)
	}

	res.Unmarshal(&r)

	// GetIpLocationResult will be a string. We need to parse it to XML
	//result := GetIPLocationResult{}
	result := ReadSystemsResult{}

	//err = xml.Unmarshal([]byte(r.GetIPLocationResult), &result)
	err = xml.Unmarshal([]byte(r.ReadSystemsResult), &result)
	if err != nil {
		log.Fatalf("xml.Unmarshal error: %s", err)
	}

	if result.Country != "US" {
		log.Fatalf("error: %+v", r)
	}

	log.Println("Country: ", result.Country)
	log.Println("State: ", result.State)
}
*/
