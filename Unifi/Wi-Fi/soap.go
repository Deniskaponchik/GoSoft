package main

//test_SOAP/soap_medium
import (
	"encoding/xml"
	"fmt"
	"io"
	"strings"
	"time"
)
import (
	"bytes"
	"crypto/tls"
	"log"
	"net/http"
)

func CreateWiFiTicketErr(
	soapServer string, bpmUrl string, userLogin string, description string, noutName string, region string, apName string, incidentType string) (
	srSlice []string) {

	if userLogin != "" {
		strBefore :=
			"<soapenv:Envelope " +
				"xmlns:soapenv=\"http://schemas.xmlsoap.org/soap/envelope/\" " +
				"xmlns:bpm=\"http://www.bercut.com/specs/aoi/tele2/bpm\">" +
				"<soapenv:Header/>" +
				"<soapenv:Body>" +
				"<bpm:createRequestRequest>" +
				"<SystemId>5594b877-3bb7-46db-99f5-3c75b3e46556</SystemId>" +
				"<ServiceId>ed84a37f-4b31-4dab-85fe-ba4fe87325b1</ServiceId>" +
				"<Subject>Description</Subject>" +
				"<UserName>UserLogin</UserName>" +
				"<RequestType>Request</RequestType>" +
				"<Priority>Normal</Priority>" +
				"<Filds>" +
				"<ID>28bbdcc4-ed50-4bcd-ac06-eeea667d62ac</ID>" +
				"<Value>Reason</Value>" +
				"</Filds>" +
				"<Filds>" +
				"<ID>5c8dee23-e48a-45bc-a084-573e1a6cc5ca</ID>" +
				"<Value>Region</Value>" +
				"</Filds>" +
				"<Filds>" +
				"<ID>f01f84be-b8f1-454f-a947-2c7f832bbb88</ID>" +
				"<Value>Monitoring</Value>" +
				"</Filds>" +
				"<Filds>" +
				"<ID>bde054e7-2b91-41c1-abba-2dcbe3a8f3f4</ID>" +
				"<Value>incidentType</Value>" +
				"</Filds>" +
				"</bpm:createRequestRequest>" +
				"</soapenv:Body>" +
				"</soapenv:Envelope>"
		//replacer := strings.NewReplacer("Description", "My des", "UserLogin", "denis.tirskikh", "Region", "Москва ЦФ")
		//replacer := strings.NewReplacer("Description", description, "UserLogin", userLogin, "incidentType", incidentType, "Region", region)
		replacer := strings.NewReplacer("Description", description, "UserLogin", userLogin, "Reason", noutName, "Region", region,
			"Monitoring", apName, "incidentType", incidentType)
		strAfter := replacer.Replace(strBefore)
		payload := []byte(strAfter)
		//os.Exit(0)
		httpMethod := "POST"

		//Вбиваем результат запроса из постмана сюда: https://tool.hiofd.com/en/xml-to-go/
		type Envelope struct {
			XMLName xml.Name `xml:"Envelope"`
			Text    string   `xml:",chardata"`
			SOAPENV string   `xml:"SOAP-ENV,attr"`
			Body    struct {
				Text                  string `xml:",chardata"`
				BerNs0                string `xml:"ber-ns0,attr"`
				CreateRequestResponse struct {
					Text        string `xml:",chardata"`
					Code        int    `xml:"Code"`
					ID          string `xml:"ID"`
					Number      string `xml:"Number"`
					SystemName  string `xml:"SystemName"`
					Description string `xml:"Description"`
				} `xml:"createRequestResponse"`
			} `xml:"Body"`
		}

		myError := 1
		for myError != 0 {
			//req, err :=	http.NewRequest(httpMethod, url, bytes.NewReader(payload))
			req, errHttpReq := http.NewRequest(httpMethod, soapServer, bytes.NewReader(payload))
			if errHttpReq == nil {
				client := &http.Client{
					Transport: &http.Transport{
						TLSClientConfig: &tls.Config{
							InsecureSkipVerify: true,
						},
					},
				}
				res, errClientDo := client.Do(req)
				if errClientDo == nil {
					/*Посмотреть response Body, если понадобится
					defer res.Body.Close()
					b, err := io.ReadAll(res.Body)
					if err != nil {
						log.Fatalln(err)
					}
					fmt.Println(string(b))
					//os.Exit(0)
					*/

					// Смог победить только через unmarshal. Кривенько косо, но работает и куча времени угрохано даже на это
					envelope := &Envelope{}
					bodyByte, errIOread := io.ReadAll(res.Body)
					if errIOread == nil {
						erXmlUnmarshal := xml.Unmarshal(bodyByte, envelope)
						if erXmlUnmarshal == nil {
							if envelope.Body.CreateRequestResponse.Code == 0 {
								srID := envelope.Body.CreateRequestResponse.ID
								srNumber := envelope.Body.CreateRequestResponse.Number
								bpmLink := bpmUrl + srID
								srSlice = append(srSlice, srID)
								srSlice = append(srSlice, srNumber)
								srSlice = append(srSlice, bpmLink)
								myError = 0
							} else {
								fmt.Println("Заявка НЕ создалась на ФИНАЛЬНОМ этапе")
								fmt.Println(envelope.Body.CreateRequestResponse.Description)
								fmt.Println("Проверь корректность:")
								fmt.Println("SOAP-сервер: " + soapServer)
								fmt.Println("User login: " + userLogin)
								fmt.Println("Регион: " + region)
								fmt.Println("Тип инцидента: " + incidentType)
								fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
								fmt.Println("")
								time.Sleep(60 * time.Second)
								myError++
							}
						} else {
							//log.Fatalln(erXmlUnmarshal)
							fmt.Println("Ошибка перекодировки ответа в xml")
							fmt.Println(erXmlUnmarshal.Error())
							fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
							time.Sleep(60 * time.Second)
							myError++
						}
					} else {
						fmt.Println("Ошибка чтения байтов из ответа")
						fmt.Println(errIOread.Error())
						fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
						time.Sleep(60 * time.Second)
						myError++
					}
				} else {
					//log.Fatal("Error on dispatching request. ", errClientDo.Error())
					//return
					fmt.Println("Ошибка отправки запроса")
					fmt.Println(errClientDo.Error())
					fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
					time.Sleep(60 * time.Second)
					myError++
				}
			} else {
				//log.Fatal("Error on creating request object. ", errHttpReq.Error())
				//return
				fmt.Println("Ошибка создания объекта запроса")
				fmt.Println(errHttpReq.Error())
				fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
				time.Sleep(60 * time.Second)
				myError++
			}
		}
	} else {
		//Для аномальных заявок
		srSlice = append(srSlice, "Заявка не была создана. User Login пустой")
		srSlice = append(srSlice, "Заявка не была создана. User Login пустой")
		srSlice = append(srSlice, "Заявка не была создана. User Login пустой")
	}
	return srSlice
}

func CreateSmacWiFiTicketErr(
	soapServer string, bpmUrl string, userLogin string, description string, region string, incidentType string) (
	srSlice []string) {

	if userLogin != "" {
		strBefore :=
			"<soapenv:Envelope " +
				"xmlns:soapenv=\"http://schemas.xmlsoap.org/soap/envelope/\" " +
				"xmlns:bpm=\"http://www.bercut.com/specs/aoi/tele2/bpm\">" +
				"<soapenv:Header/>" +
				"<soapenv:Body>" +
				"<bpm:createRequestRequest>" +
				"<SystemId>5594b877-3bb7-46db-99f5-3c75b3e46556</SystemId>" +
				"<ServiceId>ed84a37f-4b31-4dab-85fe-ba4fe87325b1</ServiceId>" +
				"<Subject>Description</Subject>" +
				"<UserName>UserLogin</UserName>" +
				"<RequestType>Request</RequestType>" +
				"<Priority>Normal</Priority>" +
				"<Filds>" +
				"<ID>5c8dee23-e48a-45bc-a084-573e1a6cc5ca</ID>" +
				"<Value>Region</Value>" +
				"</Filds>" +
				"<Filds>" +
				"<ID>bde054e7-2b91-41c1-abba-2dcbe3a8f3f4</ID>" +
				"<Value>incidentType</Value>" +
				"</Filds>" +
				"</bpm:createRequestRequest>" +
				"</soapenv:Body>" +
				"</soapenv:Envelope>"
		//replacer := strings.NewReplacer("Description", "My des", "UserLogin", "denis.tirskikh", "Region", "Москва ЦФ")
		replacer := strings.NewReplacer("Description", description, "UserLogin", userLogin, "incidentType", incidentType, "Region", region)
		strAfter := replacer.Replace(strBefore)
		payload := []byte(strAfter)
		//os.Exit(0)
		httpMethod := "POST"

		//Вбиваем результат запроса из постмана сюда: https://tool.hiofd.com/en/xml-to-go/
		type Envelope struct {
			XMLName xml.Name `xml:"Envelope"`
			Text    string   `xml:",chardata"`
			SOAPENV string   `xml:"SOAP-ENV,attr"`
			Body    struct {
				Text                  string `xml:",chardata"`
				BerNs0                string `xml:"ber-ns0,attr"`
				CreateRequestResponse struct {
					Text        string `xml:",chardata"`
					Code        int    `xml:"Code"`
					ID          string `xml:"ID"`
					Number      string `xml:"Number"`
					SystemName  string `xml:"SystemName"`
					Description string `xml:"Description"`
				} `xml:"createRequestResponse"`
			} `xml:"Body"`
		}

		myError := 1
		for myError != 0 {
			//req, err :=	http.NewRequest(httpMethod, url, bytes.NewReader(payload))
			req, errHttpReq := http.NewRequest(httpMethod, soapServer, bytes.NewReader(payload))
			if errHttpReq == nil {
				client := &http.Client{
					Transport: &http.Transport{
						TLSClientConfig: &tls.Config{
							InsecureSkipVerify: true,
						},
					},
				}
				res, errClientDo := client.Do(req)
				if errClientDo == nil {
					/*Посмотреть response Body, если понадобится
					defer res.Body.Close()
					b, err := io.ReadAll(res.Body)
					if err != nil {
						log.Fatalln(err)
					}
					fmt.Println(string(b))
					//os.Exit(0)
					*/

					// Смог победить только через unmarshal. Кривенько косо, но работает и куча времени угрохано даже на это
					envelope := &Envelope{}
					bodyByte, errIOread := io.ReadAll(res.Body)
					if errIOread == nil {
						erXmlUnmarshal := xml.Unmarshal(bodyByte, envelope)
						if erXmlUnmarshal == nil {
							if envelope.Body.CreateRequestResponse.Code == 0 {
								srID := envelope.Body.CreateRequestResponse.ID
								srNumber := envelope.Body.CreateRequestResponse.Number
								bpmLink := bpmUrl + srID
								srSlice = append(srSlice, srID)
								srSlice = append(srSlice, srNumber)
								srSlice = append(srSlice, bpmLink)
								myError = 0
							} else {
								fmt.Println("Заявка НЕ создалась на ФИНАЛЬНОМ этапе")
								fmt.Println(envelope.Body.CreateRequestResponse.Description)
								fmt.Println("Проверь корректность:")
								fmt.Println("SOAP-сервер: " + soapServer)
								fmt.Println("User login: " + userLogin)
								fmt.Println("Регион: " + region)
								fmt.Println("Тип инцидента: " + incidentType)
								fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
								fmt.Println("")
								time.Sleep(60 * time.Second)
								myError++
							}
						} else {
							//log.Fatalln(erXmlUnmarshal)
							fmt.Println("Ошибка перекодировки ответа в xml")
							fmt.Println(erXmlUnmarshal.Error())
							fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
							time.Sleep(60 * time.Second)
							myError++
						}
					} else {
						fmt.Println("Ошибка чтения байтов из ответа")
						fmt.Println(errIOread.Error())
						fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
						time.Sleep(60 * time.Second)
						myError++
					}
				} else {
					//log.Fatal("Error on dispatching request. ", errClientDo.Error())
					//return
					fmt.Println("Ошибка отправки запроса")
					fmt.Println(errClientDo.Error())
					fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
					time.Sleep(60 * time.Second)
					myError++
				}
			} else {
				//log.Fatal("Error on creating request object. ", errHttpReq.Error())
				//return
				fmt.Println("Ошибка создания объекта запроса")
				fmt.Println(errHttpReq.Error())
				fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
				time.Sleep(60 * time.Second)
				myError++
			}
		}
	} else {
		//Для аномальных заявок
		srSlice = append(srSlice, "Заявка не была создана. User Login пустой")
		srSlice = append(srSlice, "Заявка не была создана. User Login пустой")
		srSlice = append(srSlice, "Заявка не была создана. User Login пустой")
	}
	return srSlice
}

func CreateSmacWiFiTicket(
	bpmServer string, userLogin string, description string, region string, incidentType string) (
	srSlice []string) {

	if userLogin != "" {
		url := bpmServer
		/*desAps := strings.Join(aps, "\n")
		description := "Зафиксировано отключение точек:" + "\n" +
			desAps + "\n" +
			"" + "\n" +
			"Рекомендации по выполнению таких инцидентов собраны на страничке корпоративной wiki" + "\n" +
			"https://wiki.tele2.ru/display/ITKB/%5BHelpdesk+IT%5D+System+Monitoring" + "\n" +
			""
		//fmt.Println(description)
		//region := "Москва ЦФ"
		//incidentType := "Недоступна точка доступа"
		*/
		strBefore :=
			"<soapenv:Envelope " +
				"xmlns:soapenv=\"http://schemas.xmlsoap.org/soap/envelope/\" " +
				"xmlns:bpm=\"http://www.bercut.com/specs/aoi/tele2/bpm\">" +
				"<soapenv:Header/>" +
				"<soapenv:Body>" +
				"<bpm:createRequestRequest>" +
				"<SystemId>5594b877-3bb7-46db-99f5-3c75b3e46556</SystemId>" +
				"<ServiceId>ed84a37f-4b31-4dab-85fe-ba4fe87325b1</ServiceId>" +
				"<Subject>Description</Subject>" +
				"<UserName>UserLogin</UserName>" +
				"<RequestType>Request</RequestType>" +
				"<Priority>Normal</Priority>" +
				"<Filds>" +
				"<ID>5c8dee23-e48a-45bc-a084-573e1a6cc5ca</ID>" +
				"<Value>Region</Value>" +
				"</Filds>" +
				"<Filds>" +
				"<ID>bde054e7-2b91-41c1-abba-2dcbe3a8f3f4</ID>" +
				"<Value>incidentType</Value>" +
				"</Filds>" +
				"</bpm:createRequestRequest>" +
				"</soapenv:Body>" +
				"</soapenv:Envelope>"
		//replacer := strings.NewReplacer("Description", "My des", "UserLogin", "denis.tirskikh", "Region", "Москва ЦФ")
		replacer := strings.NewReplacer("Description", description, "UserLogin", userLogin, "incidentType", incidentType, "Region", region)
		strAfter := replacer.Replace(strBefore)
		//fmt.Println(strAfter)
		payload := []byte(strAfter)
		//os.Exit(0)
		httpMethod := "POST"
		req, err :=
			http.NewRequest(httpMethod, url, bytes.NewReader(payload))
		if err != nil {
			log.Fatal("Error on creating request object. ", err.Error())
			return
		}
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
		/*Посмотреть response Body, если понадобится
		defer res.Body.Close()
		b, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(string(b))
		//os.Exit(0)
		*/

		//Вбиваем результат запроса из постмана сюда: https://tool.hiofd.com/en/xml-to-go/
		type Envelope struct {
			XMLName xml.Name `xml:"Envelope"`
			Text    string   `xml:",chardata"`
			SOAPENV string   `xml:"SOAP-ENV,attr"`
			Body    struct {
				Text                  string `xml:",chardata"`
				BerNs0                string `xml:"ber-ns0,attr"`
				CreateRequestResponse struct {
					Text       string `xml:",chardata"`
					Code       string `xml:"Code"`
					ID         string `xml:"ID"`
					Number     string `xml:"Number"`
					SystemName string `xml:"SystemName"`
				} `xml:"createRequestResponse"`
			} `xml:"Body"`
		}
		// Смог победить только через unmarshal. Кривенько косо, но работает и куча времени угрохано даже на это
		//res := &MyRespEnvelope{}
		envelope := &Envelope{}
		bodyByte, err := io.ReadAll(res.Body)
		error := xml.Unmarshal(bodyByte, envelope)
		if error != nil {
			log.Fatalln(err)
		}

		srID := envelope.Body.CreateRequestResponse.ID
		srNumber := envelope.Body.CreateRequestResponse.Number
		var bpmLink string
		if bpmServer == "http://10.246.37.15:8060/specs/aoi/tele2/bpm/bpmPortType" {
			bpmLink = "https://t2ru-tr-tst-01.corp.tele2.ru/0/Nui/ViewModule.aspx#CardModuleV2/CasePage/edit/" + srID
		} else {
			bpmLink = "https://bpm.tele2.ru/0/Nui/ViewModule.aspx#CardModuleV2/CasePage/edit/" + srID
		}

		srSlice = append(srSlice, srID)
		srSlice = append(srSlice, srNumber)
		srSlice = append(srSlice, bpmLink)
		//return srNumber, srID, bpmLink
	} else {
		srSlice = append(srSlice, "Заявка не была создана. User Login пустой")
		srSlice = append(srSlice, "Заявка не была создана. User Login пустой")
		srSlice = append(srSlice, "Заявка не была создана. User Login пустой")
	}
	return srSlice
}

func CreateApTicketErr(
	soapServer string, bpmUrl string, userLogin string, description string, region string, incidentType string) (
	srSlice []string) {

	if userLogin != "" {
		strBefore :=
			"<soapenv:Envelope " +
				"xmlns:soapenv=\"http://schemas.xmlsoap.org/soap/envelope/\" " +
				"xmlns:bpm=\"http://www.bercut.com/specs/aoi/tele2/bpm\">" +
				"<soapenv:Header/>" +
				"<soapenv:Body>" +
				"<bpm:createRequestRequest>" +
				"<SystemId>5594b877-3bb7-46db-99f5-3c75b3e46556</SystemId>" +
				"<ServiceId>ed84a37f-4b31-4dab-85fe-ba4fe87325b1</ServiceId>" +
				"<Subject>Description</Subject>" +
				"<UserName>UserLogin</UserName>" +
				"<RequestType>Request</RequestType>" +
				"<Priority>Normal</Priority>" +
				"<Filds>" +
				"<ID>5c8dee23-e48a-45bc-a084-573e1a6cc5ca</ID>" +
				"<Value>Region</Value>" +
				"</Filds>" +
				"<Filds>" +
				"<ID>bde054e7-2b91-41c1-abba-2dcbe3a8f3f4</ID>" +
				"<Value>incidentType</Value>" +
				"</Filds>" +
				"</bpm:createRequestRequest>" +
				"</soapenv:Body>" +
				"</soapenv:Envelope>"
		//replacer := strings.NewReplacer("Description", "My des", "UserLogin", "denis.tirskikh", "Region", "Москва ЦФ")
		replacer := strings.NewReplacer("Description", description, "UserLogin", userLogin, "incidentType", incidentType, "Region", region)
		strAfter := replacer.Replace(strBefore)
		payload := []byte(strAfter)
		//os.Exit(0)
		httpMethod := "POST"

		//Вбиваем результат запроса из постмана сюда: https://tool.hiofd.com/en/xml-to-go/
		type Envelope struct {
			XMLName xml.Name `xml:"Envelope"`
			Text    string   `xml:",chardata"`
			SOAPENV string   `xml:"SOAP-ENV,attr"`
			Body    struct {
				Text                  string `xml:",chardata"`
				BerNs0                string `xml:"ber-ns0,attr"`
				CreateRequestResponse struct {
					Text        string `xml:",chardata"`
					Code        int    `xml:"Code"`
					ID          string `xml:"ID"`
					Number      string `xml:"Number"`
					SystemName  string `xml:"SystemName"`
					Description string `xml:"Description"`
				} `xml:"createRequestResponse"`
			} `xml:"Body"`
		}

		myError := 1
		for myError != 0 {
			//req, err :=	http.NewRequest(httpMethod, url, bytes.NewReader(payload))
			req, errHttpReq := http.NewRequest(httpMethod, soapServer, bytes.NewReader(payload))
			if errHttpReq == nil {
				client := &http.Client{
					Transport: &http.Transport{
						TLSClientConfig: &tls.Config{
							InsecureSkipVerify: true,
						},
					},
				}
				res, errClientDo := client.Do(req)
				if errClientDo == nil {
					/*Посмотреть response Body, если понадобится
					defer res.Body.Close()
					b, err := io.ReadAll(res.Body)
					if err != nil {
						log.Fatalln(err)
					}
					fmt.Println(string(b))
					//os.Exit(0)
					*/

					// Смог победить только через unmarshal. Кривенько косо, но работает и куча времени угрохано даже на это
					envelope := &Envelope{}
					bodyByte, errIOread := io.ReadAll(res.Body)
					if errIOread == nil {
						erXmlUnmarshal := xml.Unmarshal(bodyByte, envelope)
						if erXmlUnmarshal == nil {
							if envelope.Body.CreateRequestResponse.Code == 0 {
								srID := envelope.Body.CreateRequestResponse.ID
								srNumber := envelope.Body.CreateRequestResponse.Number
								bpmLink := bpmUrl + srID
								srSlice = append(srSlice, srID)
								srSlice = append(srSlice, srNumber)
								srSlice = append(srSlice, bpmLink)
								myError = 0
							} else {
								fmt.Println("Заявка НЕ создалась на ФИНАЛЬНОМ этапе")
								fmt.Println(envelope.Body.CreateRequestResponse.Description)
								fmt.Println("Проверь корректность:")
								fmt.Println("SOAP-сервер: " + soapServer)
								fmt.Println("User login: " + userLogin)
								fmt.Println("Регион: " + region)
								fmt.Println("Тип инцидента: " + incidentType)
								fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
								fmt.Println("")
								time.Sleep(60 * time.Second)
								myError++
							}
						} else {
							//log.Fatalln(erXmlUnmarshal)
							fmt.Println("Ошибка перекодировки ответа в xml")
							fmt.Println(erXmlUnmarshal.Error())
							fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
							time.Sleep(60 * time.Second)
							myError++
						}
					} else {
						fmt.Println("Ошибка чтения байтов из ответа")
						fmt.Println(errIOread.Error())
						fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
						time.Sleep(60 * time.Second)
						myError++
					}
				} else {
					//log.Fatal("Error on dispatching request. ", errClientDo.Error())
					//return
					fmt.Println("Ошибка отправки запроса")
					fmt.Println(errClientDo.Error())
					fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
					time.Sleep(60 * time.Second)
					myError++
				}
			} else {
				//log.Fatal("Error on creating request object. ", errHttpReq.Error())
				//return
				fmt.Println("Ошибка создания объекта запроса")
				fmt.Println(errHttpReq.Error())
				fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
				time.Sleep(60 * time.Second)
				myError++
			}
		}
	} else {
		//Для аномальных заявок
		srSlice = append(srSlice, "Заявка не была создана. User Login пустой")
		srSlice = append(srSlice, "Заявка не была создана. User Login пустой")
		srSlice = append(srSlice, "Заявка не была создана. User Login пустой")
	}
	return srSlice
}

func CreateApTicket(
	bpmServer string, userLogin string, aps []string, region string) (
	srSlice []string) {

	//url := "http://10.246.37.15:8060/specs/aoi/tele2/bpm/bpmPortType" //TEST
	//url := "http://10.12.15.148/specs/aoi/tele2/bpm/bpmPortType"   //PROD
	url := bpmServer
	desAps := strings.Join(aps, "\n")

	description := "Зафиксировано отключение точек:" + "\n" +
		desAps + "\n" +
		"" + "\n" +
		"Рекомендации по выполнению таких инцидентов собраны на страничке корпоративной wiki" + "\n" +
		"https://wiki.tele2.ru/display/ITKB/%5BHelpdesk+IT%5D+System+Monitoring" + "\n" +
		""
	//fmt.Println(description)
	//region := "Москва ЦФ"
	incidentType := "Недоступна точка доступа"
	strBefore :=
		"<soapenv:Envelope " +
			"xmlns:soapenv=\"http://schemas.xmlsoap.org/soap/envelope/\" " +
			"xmlns:bpm=\"http://www.bercut.com/specs/aoi/tele2/bpm\">" +
			"<soapenv:Header/>" +
			"<soapenv:Body>" +
			"<bpm:createRequestRequest>" +
			"<SystemId>5594b877-3bb7-46db-99f5-3c75b3e46556</SystemId>" +
			"<ServiceId>ed84a37f-4b31-4dab-85fe-ba4fe87325b1</ServiceId>" +
			"<Subject>Description</Subject>" +
			"<UserName>UserLogin</UserName>" +
			"<RequestType>Request</RequestType>" +
			"<Priority>Normal</Priority>" +
			"<Filds>" +
			"<ID>5c8dee23-e48a-45bc-a084-573e1a6cc5ca</ID>" +
			"<Value>Region</Value>" +
			"</Filds>" +
			"<Filds>" +
			"<ID>bde054e7-2b91-41c1-abba-2dcbe3a8f3f4</ID>" +
			"<Value>incidentType</Value>" +
			"</Filds>" +
			"</bpm:createRequestRequest>" +
			"</soapenv:Body>" +
			"</soapenv:Envelope>"
	//replacer := strings.NewReplacer("Description", "My des", "UserLogin", "denis.tirskikh", "Region", "Москва ЦФ")
	replacer := strings.NewReplacer("Description", description, "UserLogin", userLogin, "incidentType", incidentType, "Region", region)
	strAfter := replacer.Replace(strBefore)
	//fmt.Println(strAfter)
	payload := []byte(strAfter)
	//os.Exit(0)
	httpMethod := "POST"
	req, err :=
		http.NewRequest(httpMethod, url, bytes.NewReader(payload))
	if err != nil {
		log.Fatal("Error on creating request object. ", err.Error())
		return
	}
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
	//Вбиваем результат запроса из постмана сюда: https://tool.hiofd.com/en/xml-to-go/
	type Envelope struct {
		XMLName xml.Name `xml:"Envelope"`
		Text    string   `xml:",chardata"`
		SOAPENV string   `xml:"SOAP-ENV,attr"`
		Body    struct {
			Text                  string `xml:",chardata"`
			BerNs0                string `xml:"ber-ns0,attr"`
			CreateRequestResponse struct {
				Text       string `xml:",chardata"`
				Code       string `xml:"Code"`
				ID         string `xml:"ID"`
				Number     string `xml:"Number"`
				SystemName string `xml:"SystemName"`
			} `xml:"createRequestResponse"`
		} `xml:"Body"`
	}
	// Смог победить только через unmarshal. Кривенько косо, но работает и куча времени угрохано даже на это
	//res := &MyRespEnvelope{}
	envelope := &Envelope{}
	bodyByte, err := io.ReadAll(res.Body)
	error := xml.Unmarshal(bodyByte, envelope)
	if error != nil {
		log.Fatalln(err)
	}

	srID := envelope.Body.CreateRequestResponse.ID
	srNumber := envelope.Body.CreateRequestResponse.Number
	var bpmLink string
	if bpmServer == "http://10.246.37.15:8060/specs/aoi/tele2/bpm/bpmPortType" {
		bpmLink = "https://t2ru-tr-tst-01.corp.tele2.ru/0/Nui/ViewModule.aspx#CardModuleV2/CasePage/edit/" + srID
	} else {
		bpmLink = "https://bpm.tele2.ru/0/Nui/ViewModule.aspx#CardModuleV2/CasePage/edit/" + srID
	}

	srSlice = append(srSlice, srID)
	srSlice = append(srSlice, srNumber)
	srSlice = append(srSlice, bpmLink)
	//return srNumber, srID, bpmLink
	return srSlice
}

func CreateAnomalyTicketErr(
	soapServer string, bpmUrl string, userLogin string, description string, noutName string, region string, apName string, incidentType string) (
	srSlice []string) {

	if userLogin != "" {
		strBefore :=
			"<soapenv:Envelope " +
				"xmlns:soapenv=\"http://schemas.xmlsoap.org/soap/envelope/\" " +
				"xmlns:bpm=\"http://www.bercut.com/specs/aoi/tele2/bpm\">" +
				"<soapenv:Header/>" +
				"<soapenv:Body>" +
				"<bpm:createRequestRequest>" +
				"<SystemId>5594b877-3bb7-46db-99f5-3c75b3e46556</SystemId>" +
				"<ServiceId>ed84a37f-4b31-4dab-85fe-ba4fe87325b1</ServiceId>" +
				"<Subject>Description</Subject>" +
				"<UserName>UserLogin</UserName>" +
				"<RequestType>Request</RequestType>" +
				"<Priority>Normal</Priority>" +
				"<Filds>" +
				"<ID>28bbdcc4-ed50-4bcd-ac06-eeea667d62ac</ID>" +
				"<Value>Reason</Value>" +
				"</Filds>" +
				"<Filds>" +
				"<ID>5c8dee23-e48a-45bc-a084-573e1a6cc5ca</ID>" +
				"<Value>Region</Value>" +
				"</Filds>" +
				"<Filds>" +
				"<ID>f01f84be-b8f1-454f-a947-2c7f832bbb88</ID>" +
				"<Value>Monitoring</Value>" +
				"</Filds>" +
				"<Filds>" +
				"<ID>bde054e7-2b91-41c1-abba-2dcbe3a8f3f4</ID>" +
				"<Value>incidentType</Value>" +
				"</Filds>" +
				"</bpm:createRequestRequest>" +
				"</soapenv:Body>" +
				"</soapenv:Envelope>"
		//replacer := strings.NewReplacer("Description", "My des", "UserLogin", "denis.tirskikh", "Region", "Москва ЦФ")
		//replacer := strings.NewReplacer("Description", description, "UserLogin", userLogin, "incidentType", incidentType, "Region", region)
		replacer := strings.NewReplacer("Description", description, "UserLogin", userLogin, "Reason", noutName, "Region", region,
			"Monitoring", apName, "incidentType", incidentType)
		strAfter := replacer.Replace(strBefore)
		payload := []byte(strAfter)
		//os.Exit(0)
		httpMethod := "POST"

		//Вбиваем результат запроса из постмана сюда: https://tool.hiofd.com/en/xml-to-go/
		type Envelope struct {
			XMLName xml.Name `xml:"Envelope"`
			Text    string   `xml:",chardata"`
			SOAPENV string   `xml:"SOAP-ENV,attr"`
			Body    struct {
				Text                  string `xml:",chardata"`
				BerNs0                string `xml:"ber-ns0,attr"`
				CreateRequestResponse struct {
					Text        string `xml:",chardata"`
					Code        int    `xml:"Code"`
					ID          string `xml:"ID"`
					Number      string `xml:"Number"`
					SystemName  string `xml:"SystemName"`
					Description string `xml:"Description"`
				} `xml:"createRequestResponse"`
			} `xml:"Body"`
		}

		myError := 1
		for myError != 0 {
			//req, err :=	http.NewRequest(httpMethod, url, bytes.NewReader(payload))
			req, errHttpReq := http.NewRequest(httpMethod, soapServer, bytes.NewReader(payload))
			if errHttpReq == nil {
				client := &http.Client{
					Transport: &http.Transport{
						TLSClientConfig: &tls.Config{
							InsecureSkipVerify: true,
						},
					},
				}
				res, errClientDo := client.Do(req)
				if errClientDo == nil {
					/*Посмотреть response Body, если понадобится
					defer res.Body.Close()
					b, err := io.ReadAll(res.Body)
					if err != nil {
						log.Fatalln(err)
					}
					fmt.Println(string(b))
					//os.Exit(0)
					*/

					// Смог победить только через unmarshal. Кривенько косо, но работает и куча времени угрохано даже на это
					envelope := &Envelope{}
					bodyByte, errIOread := io.ReadAll(res.Body)
					if errIOread == nil {
						erXmlUnmarshal := xml.Unmarshal(bodyByte, envelope)
						if erXmlUnmarshal == nil {
							if envelope.Body.CreateRequestResponse.Code == 0 {
								srID := envelope.Body.CreateRequestResponse.ID
								srNumber := envelope.Body.CreateRequestResponse.Number
								bpmLink := bpmUrl + srID
								srSlice = append(srSlice, srID)
								srSlice = append(srSlice, srNumber)
								srSlice = append(srSlice, bpmLink)
								myError = 0
							} else {
								fmt.Println("Заявка НЕ создалась на ФИНАЛЬНОМ этапе")
								fmt.Println(envelope.Body.CreateRequestResponse.Description)
								fmt.Println("Проверь корректность:")
								fmt.Println("SOAP-сервер: " + soapServer)
								fmt.Println("User login: " + userLogin)
								fmt.Println("Регион: " + region)
								fmt.Println("Тип инцидента: " + incidentType)
								fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
								fmt.Println("")
								time.Sleep(60 * time.Second)
								myError++
							}
						} else {
							//log.Fatalln(erXmlUnmarshal)
							fmt.Println("Ошибка перекодировки ответа в xml")
							fmt.Println(erXmlUnmarshal.Error())
							fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
							time.Sleep(60 * time.Second)
							myError++
						}
					} else {
						fmt.Println("Ошибка чтения байтов из ответа")
						fmt.Println(errIOread.Error())
						fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
						time.Sleep(60 * time.Second)
						myError++
					}
				} else {
					//log.Fatal("Error on dispatching request. ", errClientDo.Error())
					//return
					fmt.Println("Ошибка отправки запроса")
					fmt.Println(errClientDo.Error())
					fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
					time.Sleep(60 * time.Second)
					myError++
				}
			} else {
				//log.Fatal("Error on creating request object. ", errHttpReq.Error())
				//return
				fmt.Println("Ошибка создания объекта запроса")
				fmt.Println(errHttpReq.Error())
				fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
				time.Sleep(60 * time.Second)
				myError++
			}
		}
	} else {
		//Для аномальных заявок
		srSlice = append(srSlice, "Заявка не была создана. User Login пустой")
		srSlice = append(srSlice, "Заявка не была создана. User Login пустой")
		srSlice = append(srSlice, "Заявка не была создана. User Login пустой")
	}
	return srSlice
}

func CreateAnomalyTicket(
	bpmServer string, userLogin string, pcName string, anomalies []string, apName string, region string) (
	//srNumber string, srID string, bpmLink string){
	srSlice []string) {

	//url := "http://10.246.37.15:8060/specs/aoi/tele2/bpm/bpmPortType" //TEST
	//url := "http://10.12.15.148/specs/aoi/tele2/bpm/bpmPortType"   //PROD
	url := bpmServer
	desAnomalies := strings.Join(anomalies, "\n")
	description := "На ноутбуке:" + "\n" +
		pcName + "\n" + "" + "\n" +
		"зафиксированы следующие Аномалии:" + "\n" +
		desAnomalies + "\n" +
		"" + "\n" +
		"Предполагаемое, но не на 100% точное имя точки:" + "\n" +
		apName + "\n" +
		"" + "\n" +
		"Рекомендации по выполнению таких инцидентов собраны на страничке корпоративной wiki" + "\n" +
		"https://wiki.tele2.ru/display/ITKB/%5BHelpdesk+IT%5D+System+Monitoring" + "\n" +
		""
	//fmt.Println(description)
	incidentType := "Плохое качество соединения клиента"
	//region := "Москва ЦФ"
	strBefore :=
		"<soapenv:Envelope " +
			"xmlns:soapenv=\"http://schemas.xmlsoap.org/soap/envelope/\" " +
			"xmlns:bpm=\"http://www.bercut.com/specs/aoi/tele2/bpm\">" +
			"<soapenv:Header/>" +
			"<soapenv:Body>" +
			"<bpm:createRequestRequest>" +
			"<SystemId>5594b877-3bb7-46db-99f5-3c75b3e46556</SystemId>" +
			"<ServiceId>ed84a37f-4b31-4dab-85fe-ba4fe87325b1</ServiceId>" +
			"<Subject>Description</Subject>" +
			"<UserName>UserLogin</UserName>" +
			"<RequestType>Request</RequestType>" +
			"<Priority>Normal</Priority>" +
			"<Filds>" +
			"<ID>5c8dee23-e48a-45bc-a084-573e1a6cc5ca</ID>" +
			"<Value>Region</Value>" +
			"</Filds>" +
			"<Filds>" +
			"<ID>bde054e7-2b91-41c1-abba-2dcbe3a8f3f4</ID>" +
			"<Value>incidentType</Value>" +
			"</Filds>" +
			"</bpm:createRequestRequest>" +
			"</soapenv:Body>" +
			"</soapenv:Envelope>"
	//replacer := strings.NewReplacer("Description", "My des", "UserLogin", "denis.tirskikh", "Region", "Москва ЦФ")
	replacer := strings.NewReplacer("Description", description, "UserLogin", userLogin, "incidentType", incidentType, "Region", region)
	strAfter := replacer.Replace(strBefore)
	//fmt.Println(strAfter)
	//time.Sleep(60 * time.Second)

	payload := []byte(strAfter)
	//os.Exit(0)

	httpMethod := "POST"

	req, err :=
		http.NewRequest(httpMethod, url, bytes.NewReader(payload))
	if err != nil {
		log.Fatal("Error on creating request object. ", err.Error())
		return
	}

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

	//Вбиваем результат запроса из постмана сюда: https://tool.hiofd.com/en/xml-to-go/
	type Envelope struct {
		XMLName xml.Name `xml:"Envelope"`
		Text    string   `xml:",chardata"`
		SOAPENV string   `xml:"SOAP-ENV,attr"`
		Body    struct {
			Text                  string `xml:",chardata"`
			BerNs0                string `xml:"ber-ns0,attr"`
			CreateRequestResponse struct {
				Text       string `xml:",chardata"`
				Code       string `xml:"Code"`
				ID         string `xml:"ID"`
				Number     string `xml:"Number"`
				SystemName string `xml:"SystemName"`
			} `xml:"createRequestResponse"`
		} `xml:"Body"`
	}
	// Смог победить только через unmarshal. Кривенько косо, но работает и куча времени угрохано даже на это
	//res := &MyRespEnvelope{}
	envelope := &Envelope{}
	bodyByte, err := io.ReadAll(res.Body)
	error := xml.Unmarshal(bodyByte, envelope)
	if error != nil {
		log.Fatalln(err)
	}

	srID := envelope.Body.CreateRequestResponse.ID
	srNumber := envelope.Body.CreateRequestResponse.Number
	var bpmLink string
	if bpmServer == "http://10.246.37.15:8060/specs/aoi/tele2/bpm/bpmPortType" {
		bpmLink = "https://t2ru-tr-tst-01.corp.tele2.ru/0/Nui/ViewModule.aspx#CardModuleV2/CasePage/edit/" + srID
	} else {
		bpmLink = "https://bpm.tele2.ru/0/Nui/ViewModule.aspx#CardModuleV2/CasePage/edit/" + srID
	}

	srSlice = append(srSlice, srID)
	srSlice = append(srSlice, srNumber)
	srSlice = append(srSlice, bpmLink)
	//return srNumber, srID, bpmLink
	return srSlice
}

func CheckTicketStatusErr(soapServer string, srID string) (statusSlice []string) {
	if len(srID) == 36 {
		//Убрать из строки \n
		strBefore := "<soapenv:Envelope xmlns:soapenv=\"http://schemas.xmlsoap.org/soap/envelope/\" xmlns:bpm=\"http://www.bercut.com/specs/aoi/tele2/bpm\"><soapenv:Header/><soapenv:Body><bpm:getStatusRequest><CaseID>SRid</CaseID></bpm:getStatusRequest></soapenv:Body></soapenv:Envelope>"
		replacer := strings.NewReplacer("SRid", srID)
		strAfter := replacer.Replace(strBefore)
		payload := []byte(strAfter)
		httpMethod := "POST" // GET запрос не срабатывает

		type Envelope struct {
			XMLName xml.Name `xml:"Envelope"`
			Text    string   `xml:",chardata"`
			SOAPENV string   `xml:"SOAP-ENV,attr"`
			Body    struct {
				Text              string `xml:",chardata"`
				BerNs0            string `xml:"ber-ns0,attr"`
				GetStatusResponse struct {
					Text        string `xml:",chardata"`
					Code        int    `xml:"Code"`
					Status      string `xml:"Status"`
					StatisId    string `xml:"StatisId"`
					Description string `xml:"Description"`
				} `xml:"getStatusResponse"`
			} `xml:"Body"`
		}

		myError := 1
		for myError != 0 {
			//req, errHttpReq := http.NewRequest(httpMethod, url, bytes.NewReader(payload))
			req, errHttpReq := http.NewRequest(httpMethod, soapServer, bytes.NewReader(payload))
			if errHttpReq == nil {
				client := &http.Client{
					Transport: &http.Transport{
						TLSClientConfig: &tls.Config{
							InsecureSkipVerify: true,
						},
					},
				}
				res, errClientDo := client.Do(req)
				if errClientDo == nil {
					/*Посмотреть response Body, если понадобится
					defer res.Body.Close() //ОСТОРОЖНЕЕ с этой штукой. Дальше могут данные не пойти
					b, err := io.ReadAll(res.Body)
					if err != nil {
						log.Fatalln(err)
					}
					fmt.Println(string(b))
					//os.Exit(0)*/

					envelope := &Envelope{}
					bodyByte, errIOread := io.ReadAll(res.Body)
					if errIOread == nil {
						erXmlUnmarshal := xml.Unmarshal(bodyByte, envelope)
						if erXmlUnmarshal == nil {
							if envelope.Body.GetStatusResponse.Code == 0 {
								statusSlice = append(statusSlice, envelope.Body.GetStatusResponse.StatisId)
								statusSlice = append(statusSlice, envelope.Body.GetStatusResponse.Status)
								myError = 0
							} else {
								fmt.Println("Попытка получения Статуса обращения оборвалась на ПОСЛЕДНЕМ этапе")
								fmt.Println(envelope.Body.GetStatusResponse.Description)
								fmt.Println("Проверь корректность:")
								fmt.Println("SOAP-сервер: " + soapServer)
								fmt.Println("SR id: " + srID)
								fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
								fmt.Println("")
								time.Sleep(60 * time.Second)
								fmt.Println("")
								myError = 1
							}
						} else {
							//log.Fatalln(erXmlUnmarshal)
							fmt.Println("Ошибка перекодировки ответа в xml")
							fmt.Println(erXmlUnmarshal.Error())
							fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
							time.Sleep(60 * time.Second)
							myError++
						}
					} else {
						//log.Fatalln(err)
						fmt.Println("Ошибка чтения байтов из ответа")
						fmt.Println(errIOread.Error())
						fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
						time.Sleep(60 * time.Second)
						myError++
					}
				} else {
					//log.Fatal("Error on dispatching request. ", err.Error())
					//return
					fmt.Println("Ошибка отправки запроса")
					fmt.Println(errClientDo.Error())
					fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
					time.Sleep(60 * time.Second)
					myError++
					//Если ночью нет доступа к SOAP = в ЦОДЕ коллапс. Могу подождать 5 часов
					//if myError == 300 { 					myError = 0				}
				}
			} else {
				//log.Fatal("Error on creating request object. ", err.Error())
				//return
				fmt.Println("Ошибка создания объекта запроса")
				fmt.Println(errHttpReq.Error())
				fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
				time.Sleep(60 * time.Second)
				myError++
				//Если ночью нет доступа к SOAP = в ЦОДЕ коллапс. Могу подождать 5 часов
				//if myError == 300 { 					myError = 0				}
			}
		}
	} else {
		//если передаётся пустая строка, не зная, существует ли заявка
		statusSlice = append(statusSlice, "0")
		statusSlice = append(statusSlice, "Тикет введён не корректно")
	}
	return statusSlice
}

func CheckTicketStatus(bpmServer string, srID string) (statusSlice []string) {
	if len(srID) == 36 {
		//url := "http://10.246.37.15:8060/specs/aoi/tele2/bpm/bpmPortType"
		//url := "http://10.12.15.148/specs/aoi/tele2/bpm/bpmPortType"   //PROD
		url := bpmServer

		//srID := "f0074e96-1ab9-4f63-af29-0acd933b49e8"
		//Убрать из строки \n
		strBefore := "<soapenv:Envelope xmlns:soapenv=\"http://schemas.xmlsoap.org/soap/envelope/\" xmlns:bpm=\"http://www.bercut.com/specs/aoi/tele2/bpm\"><soapenv:Header/><soapenv:Body><bpm:getStatusRequest><CaseID>SRid</CaseID></bpm:getStatusRequest></soapenv:Body></soapenv:Envelope>"
		replacer := strings.NewReplacer("SRid", srID)
		strAfter := replacer.Replace(strBefore)
		payload := []byte(strAfter)

		httpMethod := "POST" // GET запрос не срабатывает
		req, err :=
			http.NewRequest(httpMethod, url, bytes.NewReader(payload))
		if err != nil {
			log.Fatal("Error on creating request object. ", err.Error())
			return
		}
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

		/*Посмотреть response Body, если понадобится
		defer res.Body.Close() //ОСТОРОЖНЕЕ с этой штукой. Дальше могут данные не пойти
		b, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(string(b))
		//os.Exit(0)*/

		type Envelope struct {
			XMLName xml.Name `xml:"Envelope"`
			Text    string   `xml:",chardata"`
			SOAPENV string   `xml:"SOAP-ENV,attr"`
			Body    struct {
				Text              string `xml:",chardata"`
				BerNs0            string `xml:"ber-ns0,attr"`
				GetStatusResponse struct {
					Text     string `xml:",chardata"`
					Code     string `xml:"Code"`
					Status   string `xml:"Status"`
					StatisId string `xml:"StatisId"`
				} `xml:"getStatusResponse"`
			} `xml:"Body"`
		}
		envelope := &Envelope{}
		bodyByte, err := io.ReadAll(res.Body)
		er := xml.Unmarshal(bodyByte, envelope)
		if er != nil {
			log.Fatalln(err)
		}

		statusSlice = append(statusSlice, envelope.Body.GetStatusResponse.StatisId)
		statusSlice = append(statusSlice, envelope.Body.GetStatusResponse.Status)
		//srStatusID = envelope.Body.GetStatusResponse.StatisId
		//srStatus = envelope.Body.GetStatusResponse.Status

	} else {
		statusSlice = append(statusSlice, "0")
		statusSlice = append(statusSlice, "Тикет введён не корректно")
	}
	return statusSlice
}

func ChangeStatusErr(soapServer string, srID string, NewStatus string) (srNewStatus string) {
	UserLogin := "denis.tirskikh"
	//Убрать из строки \n
	strBefore := "<Envelope xmlns=\"http://schemas.xmlsoap.org/soap/envelope/\"><Body><changeCaseStatusRequest xmlns=\"http://www.bercut.com/specs/aoi/tele2/bpm\"><CaseId xmlns=\"\">SRid</CaseId><Status xmlns=\"\">NewStatus</Status><User xmlns=\"\">UserLogin</User></changeCaseStatusRequest></Body></Envelope>"
	replacer := strings.NewReplacer("SRid", srID, "NewStatus", NewStatus, "UserLogin", UserLogin)
	strAfter := replacer.Replace(strBefore)
	payload := []byte(strAfter)
	httpMethod := "POST" // GET запрос не срабатывает

	//https://blog.kowalczyk.info/tools/xmltogo/
	type Envelope struct {
		XMLName xml.Name `xml:"Envelope"`
		Text    string   `xml:",chardata"`
		SOAPENV string   `xml:"SOAP-ENV,attr"`
		Body    struct {
			Text                     string `xml:",chardata"`
			BerNs0                   string `xml:"ber-ns0,attr"`
			ChangeCaseStatusResponse struct {
				Text        string `xml:",chardata"`
				Code        int    `xml:"Code"`
				ModifyOn    string `xml:"ModifyOn"`
				NewStatusId string `xml:"NewStatusId"`
				Description string `xml:"Description"`
			} `xml:"changeCaseStatusResponse"`
		} `xml:"Body"`
	}

	myError := 1
	for myError != 0 {
		//req, err := http.NewRequest(httpMethod, url, bytes.NewReader(payload))
		req, errHttpReq := http.NewRequest(httpMethod, soapServer, bytes.NewReader(payload))
		if errHttpReq == nil {
			client := &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{
						InsecureSkipVerify: true,
					},
				},
			}
			res, errClientDo := client.Do(req)
			if errClientDo == nil {
				/*Посмотреть response Body, если понадобится
				defer res.Body.Close() //ОСТОРОЖНЕЕ с этой штукой. Дальше могут данные не пойти
				b, err := io.ReadAll(res.Body)
				if err != nil {
					log.Fatalln(err)
				}
				fmt.Println(string(b))
				//os.Exit(0)*/

				envelope := &Envelope{}
				bodyByte, errIOread := io.ReadAll(res.Body)
				if errIOread == nil {
					erXmlUnmarshal := xml.Unmarshal(bodyByte, envelope)
					if erXmlUnmarshal == nil {
						if envelope.Body.ChangeCaseStatusResponse.Code == 0 {
							srDateChange := envelope.Body.ChangeCaseStatusResponse.ModifyOn
							srNewStatus = envelope.Body.ChangeCaseStatusResponse.NewStatusId
							fmt.Println("Статус обращения изменён на " + NewStatus + " в: " + srDateChange)
							myError = 0
						} else {
							fmt.Println("НЕ УДАЛОСЬ изменить статус обращения на " + NewStatus)
							fmt.Println(envelope.Body.ChangeCaseStatusResponse.Description)
							fmt.Println("SOAP-сервер: " + soapServer)
							fmt.Println("SR id: " + srID)
							fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
							time.Sleep(60 * time.Second)
							fmt.Println("")
							myError = 1
						}
					} else {
						//log.Fatalln(erXmlUnmarshal)
						fmt.Println("Ошибка перекодировки ответа в xml")
						fmt.Println(erXmlUnmarshal.Error())
						fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
						time.Sleep(60 * time.Second)
						myError++
					}
				} else {
					fmt.Println("Ошибка чтения байтов из ответа")
					fmt.Println(errIOread.Error())
					fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
					time.Sleep(60 * time.Second)
					myError++
				}
			} else {
				//log.Fatal("Error on dispatching request. ", err.Error())
				//return "Error on dispatching request. "
				fmt.Println("Ошибка отправки запроса")
				fmt.Println(errClientDo.Error())
				fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
				time.Sleep(60 * time.Second)
				myError++
			}
		} else {
			//log.Fatal("Error on creating request object. ", err.Error())
			//return "Error on creating request object. "
			fmt.Println("Ошибка создания объекта запроса")
			fmt.Println(errHttpReq.Error())
			fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
			time.Sleep(60 * time.Second)
			myError++
		}
	}
	return srNewStatus
}

func ChangeStatus(soapServer string, srID string, NewStatus string) (srNewStatus string) {
	url := soapServer
	//srID := "fc0d1340-2ccd-4772-a48f-0f60f5ba753e"
	UserLogin := "denis.tirskikh"
	//NewStatus := "На уточнении"
	//NewStatus := "Отменено"

	//Убрать из строки \n
	strBefore := "<Envelope xmlns=\"http://schemas.xmlsoap.org/soap/envelope/\"><Body><changeCaseStatusRequest xmlns=\"http://www.bercut.com/specs/aoi/tele2/bpm\"><CaseId xmlns=\"\">SRid</CaseId><Status xmlns=\"\">NewStatus</Status><User xmlns=\"\">UserLogin</User></changeCaseStatusRequest></Body></Envelope>"
	replacer := strings.NewReplacer("SRid", srID, "NewStatus", NewStatus, "UserLogin", UserLogin)
	strAfter := replacer.Replace(strBefore)
	payload := []byte(strAfter)

	httpMethod := "POST" // GET запрос не срабатывает
	req, err :=
		http.NewRequest(httpMethod, url, bytes.NewReader(payload))
	if err != nil {
		log.Fatal("Error on creating request object. ", err.Error())
		return "Error on creating request object. "
	}
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
		return "Error on dispatching request. "
	}

	/*Посмотреть response Body, если понадобится
	defer res.Body.Close() //ОСТОРОЖНЕЕ с этой штукой. Дальше могут данные не пойти
	b, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(b))
	//os.Exit(0)*/

	//https://blog.kowalczyk.info/tools/xmltogo/
	type Envelope struct {
		XMLName xml.Name `xml:"Envelope"`
		Text    string   `xml:",chardata"`
		SOAPENV string   `xml:"SOAP-ENV,attr"`
		Body    struct {
			Text                     string `xml:",chardata"`
			BerNs0                   string `xml:"ber-ns0,attr"`
			ChangeCaseStatusResponse struct {
				Text        string `xml:",chardata"`
				Code        string `xml:"Code"`
				ModifyOn    string `xml:"ModifyOn"`
				NewStatusId string `xml:"NewStatusId"`
			} `xml:"changeCaseStatusResponse"`
		} `xml:"Body"`
	}
	envelope := &Envelope{}
	bodyByte, err := io.ReadAll(res.Body)
	er := xml.Unmarshal(bodyByte, envelope)
	if er != nil {
		log.Fatalln(err)
	}

	srDateChange := envelope.Body.ChangeCaseStatusResponse.ModifyOn
	srNewStatus = envelope.Body.ChangeCaseStatusResponse.NewStatusId

	if srDateChange != "" && srNewStatus != "" {
		fmt.Println("Статус обращения изменён на " + NewStatus + " в: " + srDateChange)
	} else {
		fmt.Println("НЕ УДАЛОСЬ изменить статус обращения на " + NewStatus)
	}
	return srNewStatus
}

func AddCommentErr(soapServer string, srID string, myComment string, bpmUrl string) {
	userLogin := "denis.tirskikh"
	//Убрать из строки \n
	//strBefore := "<Envelope xmlns=\"http://schemas.xmlsoap.org/soap/envelope/\"><Body><createCommentRequest xmlns=\"http://www.bercut.com/specs/aoi/tele2/bpm\"><CaseId xmlns=\"\">srID</CaseId><Message xmlns=\"\">myComment</Message><Author xmlns=\"\">userLogin</Author></createCommentRequest></Body></Envelope>"
	strBefore := "<Envelope xmlns=\"http://schemas.xmlsoap.org/soap/envelope/\"><Body><createCommentRequest xmlns=\"http://www.bercut.com/specs/aoi/tele2/bpm\"><CaseId>srID</CaseId><Message>myComment</Message><Author>userLogin</Author></createCommentRequest></Body></Envelope>"
	replacer := strings.NewReplacer("srID", srID, "myComment", myComment, "userLogin", userLogin)
	strAfter := replacer.Replace(strBefore)
	//fmt.Println(strAfter)
	payload := []byte(strAfter)
	httpMethod := "POST" // GET запрос не срабатывает

	//https://blog.kowalczyk.info/tools/xmltogo/
	type Envelope struct {
		XMLName xml.Name `xml:"Envelope"`
		Text    string   `xml:",chardata"`
		SOAPENV string   `xml:"SOAP-ENV,attr"`
		Body    struct {
			Text                  string `xml:",chardata"`
			BerNs0                string `xml:"ber-ns0,attr"`
			CreateCommentResponse struct {
				Text        string `xml:",chardata"`
				Code        int    `xml:"Code"`
				CreatedOn   string `xml:"CreatedOn"`
				ID          string `xml:"Id"`
				Description string `xml:"Description"`
			} `xml:"createCommentResponse"`
		} `xml:"Body"`
	}

	myError := 1
	for myError != 0 {
		//req, err :=	http.NewRequest(httpMethod, url, bytes.NewReader(payload))
		req, errHttpReq := http.NewRequest(httpMethod, soapServer, bytes.NewReader(payload))
		if errHttpReq == nil {
			client := &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{
						InsecureSkipVerify: true,
					},
				},
			}
			res, errClientDo := client.Do(req)
			if errClientDo == nil {
				/*Посмотреть response Body, если понадобится
				defer res.Body.Close() //ОСТОРОЖНЕЕ с этой штукой. Дальше могут данные не пойти
				b, err := io.ReadAll(res.Body)
				if err != nil {
					log.Fatalln(err)
				}
				fmt.Println(string(b))
				os.Exit(0)*/

				envelope := &Envelope{}
				bodyByte, errIOread := io.ReadAll(res.Body)
				if errIOread == nil {
					erXmlUnmarshal := xml.Unmarshal(bodyByte, envelope)
					if erXmlUnmarshal == nil {
						if envelope.Body.CreateCommentResponse.Code == 0 {
							srDateComment := envelope.Body.CreateCommentResponse.CreatedOn
							//srNewStatus := envelope.Body.ChangeCaseStatusResponse.NewStatusId
							fmt.Println("Оставлен комментарий в ")
							fmt.Println(bpmUrl + srID)
							fmt.Println(srDateComment)
							myError = 0
						} else {
							fmt.Println("НЕ УДАЛОСЬ оставить комментарий")
							fmt.Println(envelope.Body.CreateCommentResponse.Description)
							fmt.Println("Проверь корректность:")
							fmt.Println("SOAP-сервер: " + soapServer)
							fmt.Println("SR id: " + srID)
							fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
							fmt.Println("")
							time.Sleep(60 * time.Second)
							myError = 1
						}
					} else {
						//log.Fatalln(erXmlUnmarshal)
						fmt.Println("Ошибка перекодировки ответа в xml")
						fmt.Println(erXmlUnmarshal.Error())
						fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
						time.Sleep(60 * time.Second)
						myError++
					}
				} else {
					fmt.Println("Ошибка чтения байтов из ответа")
					fmt.Println(errIOread.Error())
					fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
					time.Sleep(60 * time.Second)
					myError++
				}
			} else {
				//log.Fatal("Error on dispatching request. ", errClientDo.Error())
				//return
				fmt.Println("Ошибка отправки запроса")
				fmt.Println(errClientDo.Error())
				fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
				time.Sleep(60 * time.Second)
				myError++
			}
		} else {
			//log.Fatal("Error on creating request object. ", errHttpReq.Error())
			//return
			fmt.Println("Ошибка создания объекта запроса")
			fmt.Println(errHttpReq.Error())
			fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
			time.Sleep(60 * time.Second)
			myError++
		}
	}
}

func AddComment(bpmServer string, srID string, myComment string, bpmUrl string) {
	url := bpmServer
	//srID := "fc0d1340-2ccd-4772-a48f-0f60f5ba753e"
	userLogin := "denis.tirskikh"
	//myComment := "Моё первое сервисное сообщение!"

	//Убрать из строки \n
	//strBefore := "<Envelope xmlns=\"http://schemas.xmlsoap.org/soap/envelope/\"><Body><createCommentRequest xmlns=\"http://www.bercut.com/specs/aoi/tele2/bpm\"><CaseId xmlns=\"\">srID</CaseId><Message xmlns=\"\">myComment</Message><Author xmlns=\"\">userLogin</Author></createCommentRequest></Body></Envelope>"
	strBefore := "<Envelope xmlns=\"http://schemas.xmlsoap.org/soap/envelope/\"><Body><createCommentRequest xmlns=\"http://www.bercut.com/specs/aoi/tele2/bpm\"><CaseId>srID</CaseId><Message>myComment</Message><Author>userLogin</Author></createCommentRequest></Body></Envelope>"
	replacer := strings.NewReplacer("srID", srID, "myComment", myComment, "userLogin", userLogin)
	strAfter := replacer.Replace(strBefore)
	//fmt.Println(strAfter)
	payload := []byte(strAfter)

	httpMethod := "POST" // GET запрос не срабатывает
	req, err :=
		http.NewRequest(httpMethod, url, bytes.NewReader(payload))
	if err != nil {
		log.Fatal("Error on creating request object. ", err.Error())
		return
	}
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

	/*Посмотреть response Body, если понадобится
	defer res.Body.Close() //ОСТОРОЖНЕЕ с этой штукой. Дальше могут данные не пойти
	b, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(b))
	os.Exit(0)*/

	//https://blog.kowalczyk.info/tools/xmltogo/
	type Envelope struct {
		XMLName xml.Name `xml:"Envelope"`
		Text    string   `xml:",chardata"`
		SOAPENV string   `xml:"SOAP-ENV,attr"`
		Body    struct {
			Text                  string `xml:",chardata"`
			BerNs0                string `xml:"ber-ns0,attr"`
			CreateCommentResponse struct {
				Text      string `xml:",chardata"`
				Code      string `xml:"Code"`
				CreatedOn string `xml:"CreatedOn"`
				ID        string `xml:"Id"`
			} `xml:"createCommentResponse"`
		} `xml:"Body"`
	}
	envelope := &Envelope{}
	bodyByte, err := io.ReadAll(res.Body)
	er := xml.Unmarshal(bodyByte, envelope)
	if er != nil {
		log.Fatalln(err)
	}

	srDateComment := envelope.Body.CreateCommentResponse.CreatedOn
	//srNewStatus := envelope.Body.ChangeCaseStatusResponse.NewStatusId

	if srDateComment != "" {
		fmt.Println("Оставлен комментарий в ")
		fmt.Println(bpmUrl + srID)
		fmt.Println(srDateComment)
	} else {
		fmt.Println("НЕ УДАЛОСЬ оставить комментарий")
	}
}
