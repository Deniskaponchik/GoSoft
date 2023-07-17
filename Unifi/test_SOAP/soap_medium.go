package main

import (
	"encoding/xml"
	"strings"
)
import (
	"bytes"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
)

func main() {
	/*			url := fmt.Sprintf("%s%s%s",
				"https://12.34.56.78:9443",
				"/services/RemoteUserStoreManagerService",
				".RemoteUserStoreManagerServiceHttpsSoap11Endpoint")*/
	url := "http://10.246.37.15:8060/specs/aoi/tele2/bpm/bpmPortType"

	//region := "БиДВ"
	//regionField := "<Value>" + region + "</Value>"

	/*Тема с массивом строк не зашла
	strArray := [17]string{
		//"'",
		"'<soapenv:Envelope xmlns:soapenv=\"http://schemas.xmlsoap.org/soap/envelope/\" xmlns:bpm=\"http://www.bercut.com/specs/aoi/tele2/bpm\">",
		"<soapenv:Header/>",
		"soapenv:Body>",
		"<bpm:createRequestRequest>",
		"<SystemId>5594b877-3bb7-46db-99f5-3c75b3e46556</SystemId>",
		"<ServiceId>ed84a37f-4b31-4dab-85fe-ba4fe87325b1</ServiceId>",
		"<Subject>Тестовая заявка medium</Subject>",
		"<UserName>denis.tirskikh</UserName>",
		"<RequestType>Request</RequestType>",
		"<Priority>Normal</Priority>",
		"<Filds>",
		"<ID>5c8dee23-e48a-45bc-a084-573e1a6cc5ca</ID>",
		regionField,
		"</Filds>",
		"</bpm:createRequestRequest>",
		"</soapenv:Body>",
		"</soapenv:Envelope>'",
		//"'",
	}
	fmt.Println(strArray)
	buf := &bytes.Buffer{}
	gob.NewEncoder(buf).Encode(strArray)
	mypayload := buf.Bytes()
	fmt.Println(mypayload)
	*/
	strBefore := "<soapenv:Envelope xmlns:soapenv=\"http://schemas.xmlsoap.org/soap/envelope/\" xmlns:bpm=\"http://www.bercut.com/specs/aoi/tele2/bpm\"><soapenv:Header/><soapenv:Body><bpm:createRequestRequest><SystemId>5594b877-3bb7-46db-99f5-3c75b3e46556</SystemId><ServiceId>ed84a37f-4b31-4dab-85fe-ba4fe87325b1</ServiceId><Subject>Description</Subject><UserName>UserLogin</UserName><RequestType>Request</RequestType><Priority>Normal</Priority><Filds><ID>5c8dee23-e48a-45bc-a084-573e1a6cc5ca</ID><Value>Region</Value></Filds></bpm:createRequestRequest></soapenv:Body></soapenv:Envelope>"
	replacer := strings.NewReplacer("Description", "My des", "UserLogin", "denis.tirskikh", "Region", "Москва ЦФ")
	//re := regexp.MustCompile("fox|dog")
	//newStr := re.ReplaceAllString(str, "cat")
	strAfter := replacer.Replace(strBefore)
	fmt.Println(strAfter)
	//os.Exit(0)
	/*
			payload := []byte(strings.TrimSpace(`
				<soapenv:Envelope
					xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/"
					xmlns:ser="http://service.ws.um.carbon.wso2.org"
				>
				<soapenv:Body>
						<ser:listUsers>
							<ser:filter></ser:filter>
							<ser:maxItemLimit>100</ser:maxItemLimit>
						</ser:listUsers>
					</soapenv:Body>
				</soapenv:Envelope>`,
			))
		payload := []byte(strings.TrimSpace(`
			<soapenv:Envelope
				xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/"
				xmlns:bpm="http://www.bercut.com/specs/aoi/tele2/bpm"
			>
			<soapenv:Header/>
			<soapenv:Body>
					<bpm:readSystemsRequest>
						<bpm:Filter>WebTutor</bpm:Filter>
					</bpm:readSystemsRequest>
				</soapenv:Body>
			</soapenv:Envelope>`,
		))*/
	//payload := []byte(strings.TrimSpace(`
	payload := []byte(strAfter)
	/*
		payload := []byte(`
				<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:bpm="http://www.bercut.com/specs/aoi/tele2/bpm">
			   <soapenv:Header/>
			   <soapenv:Body>
			      <bpm:createRequestRequest>
			          <SystemId>5594b877-3bb7-46db-99f5-3c75b3e46556</SystemId>
			         <ServiceId>ed84a37f-4b31-4dab-85fe-ba4fe87325b1</ServiceId>
			         <Subject>Тестовая заявка Denis Tirskikh</Subject>
			         <UserName>denis.tirskikh</UserName>
			        <RequestType>Request</RequestType>
			        <Priority>Normal</Priority>
			        <Filds>
			            <ID>5c8dee23-e48a-45bc-a084-573e1a6cc5ca</ID>
			            <Value>Воронеж ОЦО</Value>
			         </Filds>
			      </bpm:createRequestRequest>
			   </soapenv:Body>
					</soapenv:Envelope>`,
		)*/
	//fmt.Println(payload)
	//os.Exit(0)

	//<bpm:Top></bpm:Top>
	//<bpm:Skip></bpm:Skip>
	//fmt.Println(payload)

	//soapAction := "urn:listUsers"        // The format is `urn:<soap_action>`
	//soapAction := "urn:readSystemsRequest" // The format is `urn:<soap_action>`

	//username := "admin"
	//password := "admin"

	//httpMethod := "POST"
	httpMethod := "POST"

	req, err :=
		http.NewRequest(httpMethod, url, bytes.NewReader(payload))
	if err != nil {
		log.Fatal("Error on creating request object. ", err.Error())
		return
	}

	//req.Header.Set("Content-type", "application/xml")
	//req.Header.Set("SOAPAction", soapAction)
	//req.SetBasicAuth(username, password)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	res, err := client.Do(req)
	if err != nil {
		log.Fatal("Error on dispatching request. ", err.Error())
		return
	}
	/*
		type UserList struct {
			XMLName xml.Name
			Body    struct {
				XMLName           xml.Name
				ListUsersResponse struct {
					XMLName xml.Name
					Return  []string `xml:"return"`
				} `xml:"listUsersResponse"`
			}
		}*/
	type SystemList struct {
		XMLName xml.Name
		Body    struct {
			XMLName             xml.Name
			readSystemsResponse struct {
				XMLName xml.Name
				Table   struct {
					//XMLName xml.Name
					head struct {
						XMLName xml.Name
						Head    []string `xml:",innerxml"`
					} `xml:"head"`
					row struct {
						XMLName xml.Name
						Cell    []string `xml:",innerxml"`
					} `xml:"row"`
				} `xml:"Table"`
			} `xml:"readSystemsResponse"`
		}
	}

	//result := new(UserList)
	result := new(SystemList)
	err = xml.NewDecoder(res.Body).Decode(result)
	if err != nil {
		log.Fatal("Error on unmarshaling xml. ", err.Error())
		return
	}

	//users := result.Body.ListUsersResponse.Return
	//systems := result.Body.readSystemsResponse.Table.row.Cell[0]

	//fmt.Println(strings.Join(users, ", "))
	//fmt.Println(strings.Join(systems, ", "))
	fmt.Println(res)

}
