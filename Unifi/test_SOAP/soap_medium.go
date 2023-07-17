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
	/*
		url := fmt.Sprintf("%s%s%s",
		"https://12.34.56.78:9443",
		"/services/RemoteUserStoreManagerService",
		".RemoteUserStoreManagerServiceHttpsSoap11Endpoint")
	*/
	url := "http://10.246.37.15:8060/specs/aoi/tele2/bpm/bpmPortType"
	userlogin := "denis.tirskikh"
	pcName := "wsir-tirskikh"
	anomalies := []string{
		"anomal1",
		"anomaly2",
		"anomaly3",
	}
	apName := "IRK-CO-1FL"
	desAnomalies := strings.Join(anomalies, "\n")

	//description := "<![CDATA[Tootsie roll tiramisu macaroon wafer carrot cake. <br /> Danish topping sugar plum tart bonbon caramels cake.]]>"
	//description := "Tootsie roll tiramisu macaroon wafer carrot cake. \n Danish topping sugar plum tart bonbon caramels cake."
	//description := "Tootsie roll tiramisu maca" + "\n" + "Danish topping sugar plum tart bonbon "
	//description := "У клиента зафиксированы следующие Аномалии:" + "\n" + desAnomalies + "\n" + ""
	description := "На ноутбуке:" + "\n" + pcName + "\n" + "" + "\n" + "зафиксированы следующие Аномалии:" + "\n" + desAnomalies + "\n" + "" + "\n" + "Предполагаемое, но не на 100% точное имя точки:" + "\n" + apName + "\n" + "" + "\n" + "Рекомендации по выполнению таких инцидентов собраны на страничке корпоративной wiki" + "\n" + "https://wiki.tele2.ru/display/ITKB/%5BHelpdesk+IT%5D+System+Monitoring" + "\n" + ""
	region := "Москва ЦФ"

	strBefore := "<soapenv:Envelope xmlns:soapenv=\"http://schemas.xmlsoap.org/soap/envelope/\" xmlns:bpm=\"http://www.bercut.com/specs/aoi/tele2/bpm\"><soapenv:Header/><soapenv:Body><bpm:createRequestRequest><SystemId>5594b877-3bb7-46db-99f5-3c75b3e46556</SystemId><ServiceId>ed84a37f-4b31-4dab-85fe-ba4fe87325b1</ServiceId><Subject>Description</Subject><UserName>UserLogin</UserName><RequestType>Request</RequestType><Priority>Normal</Priority><Filds><ID>5c8dee23-e48a-45bc-a084-573e1a6cc5ca</ID><Value>Region</Value></Filds></bpm:createRequestRequest></soapenv:Body></soapenv:Envelope>"
	//replacer := strings.NewReplacer("Description", "My des", "UserLogin", "denis.tirskikh", "Region", "Москва ЦФ")
	replacer := strings.NewReplacer("Description", description, "UserLogin", userlogin, "Region", region)
	strAfter := replacer.Replace(strBefore)
	//fmt.Println(strAfter)
	payload := []byte(strAfter)
	//os.Exit(0)

	/*ОРИГИНАЛ
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
	//Запрос на Получение id Системы из bpm
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
	/*Запрос на создание заявки в SMAC.Wi-Fi
	//payload := []byte(strings.TrimSpace(`
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

	//<bpm:Top></bpm:Top>
	//<bpm:Skip></bpm:Skip>

	//soapAction := "urn:listUsers"          // The format is `urn:<soap_action>`
	//soapAction := "urn:readSystemsRequest" // The format is `urn:<soap_action>`

	//username := "admin"	//password := "admin"

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
	/*ORIGINAL
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
