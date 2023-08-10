package main

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

func main() {
	//bpmUrl := "https://bpm.tele2.ru/0/Nui/ViewModule.aspx#CardModuleV2/CasePage/edit/"
	//soapServer := "http://10.12.15.148/specs/aoi/tele2/bpm/bpmPortType" //RIGHT
	//soapServer := "http://10.12.15.149/specs/aoi/tele2/bpm/bpmPortType" //WRONG

	//bpmUrl := "https://t2ru-tr-tst-01.corp.tele2.ru/0/Nui/ViewModule.aspx#CardModuleV2/CasePage/edit/"
	//soapServer := "http://10.246.37.15:8060/specs/aoi/tele2/bpm/bpmPortType"

	//srID := "42255953-46aa-40c3-8df0-65a82e31b1d1"

	//cswt := CheckTicketStatusErr(soapServer, "4b34ea8c-76df-40f5-a617-5d9843f5fc69")
	//cswt := CreateWiFiTicketErr(soapServer, bpmUrl, "denis.tirskikh", "description", "WSIR-BRONER", "БиДВ","IRK-CO-01", "Плохое качество соединения клиента")
	//fmt.Println(cswt[0])
	//fmt.Println(cswt[1])
	//fmt.Println(cswt[2])
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
				Code        string `xml:"Code"`
				ModifyOn    string `xml:"ModifyOn"`
				NewStatusId string `xml:"NewStatusId"`
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
						srDateChange := envelope.Body.ChangeCaseStatusResponse.ModifyOn
						srNewStatus = envelope.Body.ChangeCaseStatusResponse.NewStatusId

						if srDateChange != "" && srNewStatus != "" {
							fmt.Println("Статус обращения изменён на " + NewStatus + " в: " + srDateChange)
						} else {
							fmt.Println("НЕ УДАЛОСЬ изменить статус обращения на " + NewStatus)
						}
						myError = 0
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

func mainCHANGE(soapServer string, srID string, NewStatus string) (srNewStatus string) {

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
						if envelope.Body.CreateCommentResponse.Code != 0 || envelope.Body.CreateCommentResponse.CreatedOn == "" {
							//srDateComment := envelope.Body.CreateCommentResponse.CreatedOn
							fmt.Println("Оставлен комментарий в ")
							fmt.Println(bpmUrl + srID)
							fmt.Println(envelope.Body.CreateCommentResponse.CreatedOn)
							myError = 0
						} else {
							fmt.Println(envelope.Body.CreateCommentResponse.Description)
							fmt.Println("Попытка оставить комментарий ОБОРВАЛАСЬ на ПОСЛЕДНЕМ этапе")
							fmt.Println("Проверь доступность SOAP-сервера и корректность входных данных:")
							fmt.Println("SOAP-сервер: " + soapServer)
							fmt.Println("SR id: " + srID)
							fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
							fmt.Println("")
							time.Sleep(60 * time.Second)
							myError = 1
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
}

func mainCOMMENT(soapServer string) {
	//url := Server
	url := soapServer
	srID := "fc0d1340-2ccd-4772-a48f-0f60f5ba753e"
	userLogin := "denis.tirskikh"
	myComment := "Моё первое сервисное сообщение!"

	//Убрать из строки \n
	strBefore := "<Envelope xmlns=\"http://schemas.xmlsoap.org/soap/envelope/\"><Body><createCommentRequest xmlns=\"http://www.bercut.com/specs/aoi/tele2/bpm\"><CaseId xmlns=\"\">srID</CaseId><Message xmlns=\"\">myComment</Message><Author xmlns=\"\">userLogin</Author></createCommentRequest></Body></Envelope>"
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

	////https://blog.kowalczyk.info/tools/xmltogo/
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
		fmt.Println("Оставлен комментарий в " + srDateComment)
		//fmt.Println(srNewStatus)
	} else {
		fmt.Println("НЕ УДАЛОСЬ оставить комментарий")
	}
}

func CheckTicketStatusErr(soapServer string, srID string) (statusSlice []string) {
	if len(srID) == 36 {
		//url := bpmServer
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
		envelope := &Envelope{}

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

					bodyByte, errIOread := io.ReadAll(res.Body)
					if errIOread == nil {
						erXmlUnmarshal := xml.Unmarshal(bodyByte, envelope)
						if erXmlUnmarshal == nil {
							if envelope.Body.GetStatusResponse.Code != 0 || envelope.Body.GetStatusResponse.StatisId == "" {
								fmt.Println(envelope.Body.GetStatusResponse.Description)
								fmt.Println("Попытка получения Статуса обращения оборвалась на ПОСЛЕДНЕМ этапе")
								fmt.Println("Проверь доступность SOAP-сервера и корректность входных данных:")
								fmt.Println("SOAP-сервер: " + soapServer)
								fmt.Println("SR id: " + srID)
								fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
								fmt.Println("")
								time.Sleep(60 * time.Second)
								myError++
							} else {
								statusSlice = append(statusSlice, envelope.Body.GetStatusResponse.StatisId)
								statusSlice = append(statusSlice, envelope.Body.GetStatusResponse.Status)
								myError = 0
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
					time.Sleep(10 * time.Second)
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
			if myError == 6 {
				myError = 0
				fmt.Println("После 6 неудачных попыток идём дальше. Статус заявки получить не удалось")
				statusSlice = append(statusSlice, "")
				statusSlice = append(statusSlice, "")
			}
		}
	} else {
		statusSlice = append(statusSlice, "0")
		statusSlice = append(statusSlice, "Тикет введён не корректно")
	}
	return statusSlice
}

func mainCHECK(soapServer string) {
	//url := Server
	url := soapServer
	srID := "f0074e96-1ab9-4f63-af29-0acd933b49e8"
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

	srStatus := envelope.Body.GetStatusResponse.Status
	srStatusId := envelope.Body.GetStatusResponse.StatisId
	fmt.Println(srStatus)
	fmt.Println(srStatusId)
}

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
	soapServer string, bpmUrl string, userLogin string, description string, region string, incidentType string) (srSlice []string) {

	if userLogin != "" {
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
					Text       string `xml:",chardata"`
					Code       string `xml:"Code"`
					ID         string `xml:"ID"`
					Number     string `xml:"Number"`
					SystemName string `xml:"SystemName"`
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
							srID := envelope.Body.CreateRequestResponse.ID
							srNumber := envelope.Body.CreateRequestResponse.Number
							bpmLink := bpmUrl + srID
							srSlice = append(srSlice, srID)
							srSlice = append(srSlice, srNumber)
							srSlice = append(srSlice, bpmLink)
							if srSlice[0] == "" {
								fmt.Println("Итог пустой. По каким-то причинам заявка не создалась на ФИНАЛЬНОМ этапе")
								fmt.Println("Проверь корректность:")
								fmt.Println("SOAP-сервер: " + soapServer)
								fmt.Println("User login: " + userLogin)
								fmt.Println("Регион: " + region)
								fmt.Println("Тип инцидента: " + incidentType)
								fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
								fmt.Println("")
								time.Sleep(60 * time.Second)
								myError++
							} else {
								myError = 0
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
		srSlice = append(srSlice, "Заявка не была создана. User Login пустой")
		srSlice = append(srSlice, "Заявка не была создана. User Login пустой")
		srSlice = append(srSlice, "Заявка не была создана. User Login пустой")
	}
	return srSlice
}

func mainCreateAnomalyTicket() {
	//
	// https://novalagung.medium.com/soap-wsdl-request-in-go-language-3861cfb5949e

	//url := Server
	url := "http://10.246.37.15:8060/specs/aoi/tele2/bpm/bpmPortType"
	userlogin := "denis.tirskikh"
	pcName := "wsir-tirskikh"
	anomalies := []string{
		"anomal1",
		"anomaly2",
	}
	apName := "IRK-CO-1FL"
	desAnomalies := strings.Join(anomalies, "\n")

	//description := "У клиента зафиксированы следующие Аномалии:" + "\n" + desAnomalies + "\n" + ""
	description := "На ноутбуке:" + "\n" +
		pcName + "\n" + "" + "\n" +
		"зафиксированы следующие Аномалии:" + "\n" +
		desAnomalies + "\n" + "" + "\n" +
		"Предполагаемое, но не на 100% точное имя точки:" + "\n" +
		apName + "\n" + "" + "\n" +
		"Рекомендации по выполнению таких инцидентов собраны на страничке корпоративной wiki" + "\n" +
		"https://wiki.tele2.ru/display/ITKB/%5BHelpdesk+IT%5D+System+Monitoring" + "\n" +
		""
	region := "БиДВ"
	incidentType := "Плохое качество соединения клиента"

	// Именно двойные кавычки нужны, обрамляющие SOAP запрос. с одинарными не работает
	//strBefore := "<soapenv:Envelope xmlns:soapenv=\"http://schemas.xmlsoap.org/soap/envelope/\" xmlns:bpm=\"http://www.bercut.com/specs/aoi/tele2/bpm\"><soapenv:Header/><soapenv:Body><bpm:createRequestRequest><SystemId>5594b877-3bb7-46db-99f5-3c75b3e46556</SystemId><ServiceId>ed84a37f-4b31-4dab-85fe-ba4fe87325b1</ServiceId><Subject>Description</Subject><UserName>UserLogin</UserName><RequestType>Request</RequestType><Priority>Normal</Priority><Filds><ID>5c8dee23-e48a-45bc-a084-573e1a6cc5ca</ID><Value>Region</Value></Filds></bpm:createRequestRequest></soapenv:Body></soapenv:Envelope>"
	//strBefore := "<soapenv:Envelope xmlns:soapenv=\"http://schemas.xmlsoap.org/soap/envelope/\" xmlns:bpm=\"http://www.bercut.com/specs/aoi/tele2/bpm\"><soapenv:Header/><soapenv:Body><bpm:createRequestRequest><SystemId>5594b877-3bb7-46db-99f5-3c75b3e46556</SystemId><ServiceId>ed84a37f-4b31-4dab-85fe-ba4fe87325b1</ServiceId><Subject>Description</Subject><UserName>UserLogin</UserName><RequestType>Request</RequestType><Priority>Normal</Priority><Filds><ID>5c8dee23-e48a-45bc-a084-573e1a6cc5ca</ID><Value>Region</Value></Filds><Filds><ID>bde054e7-2b91-41c1-abba-2dcbe3a8f3f4</ID><Value>incidentType</Value></Filds></bpm:createRequestRequest></soapenv:Body></soapenv:Envelope>"
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
	/* НЕ РАБОТАЕТ:
	strBefore := `
	<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/"	xmlns:bpm="http://www.bercut.com/specs/aoi/tele2/bpm>"
			<soapenv:Header/>
			<soapenv:Body>
				<bpm:createRequestRequest>
					<SystemId>5594b877-3bb7-46db-99f5-3c75b3e46556</SystemId>
					<ServiceId>ed84a37f-4b31-4dab-85fe-ba4fe87325b1</ServiceId>
					<Subject>Description</Subject>
					<UserName>UserLogin</UserName>
					<RequestType>Request</RequestType>
					<Priority>Normal</Priority>
					<Filds>
						<ID>5c8dee23-e48a-45bc-a084-573e1a6cc5ca</ID>
						<Value>Region</Value>
					</Filds>
					<Filds>
						<ID>bde054e7-2b91-41c1-abba-2dcbe3a8f3f4</ID>
						<Value>incidentType</Value>
					</Filds>
				</bpm:createRequestRequest>
			</soapenv:Body>
	</soapenv:Envelope>`
	*/
	//replacer := strings.NewReplacer("Description", "My des", "UserLogin", "denis.tirskikh", "Region", "Москва ЦФ")
	replacer := strings.NewReplacer("Description", description, "UserLogin", userlogin, "Region", region, "incidentType", incidentType)
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

	/*<bpm:Top></bpm:Top>
	<bpm:Skip></bpm:Skip>
	soapAction := "urn:listUsers"          // The format is `urn:<soap_action>`
	soapAction := "urn:readSystemsRequest" // The format is `urn:<soap_action>`
	username := "admin"	//password := "admin"
	*/
	httpMethod := "POST"
	req, err :=
		http.NewRequest(httpMethod, url, bytes.NewReader(payload))
	if err != nil {
		log.Fatal("Error on creating request object. ", err.Error())
		return
	}

	/*
		req.Header.Set("Content-type", "application/xml")
		req.Header.Set("SOAPAction", soapAction)
		req.SetBasicAuth(username, password)
	*/
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

	//Вбиваем результат из постмана сюда
	//https://tool.hiofd.com/en/xml-to-go/
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
	sr := envelope.Body.CreateRequestResponse.Number
	srID := envelope.Body.CreateRequestResponse.ID
	bpmLink := "https://t2ru-tr-tst-01.corp.tele2.ru/0/Nui/ViewModule.aspx#CardModuleV2/CasePage/edit/" + srID
	//fmt.Println(envelope.Body.CreateRequestResponse.Number, error)
	//fmt.Println(envelope.Body.CreateRequestResponse.Number)
	fmt.Println(sr)
	fmt.Println(srID)
	fmt.Println(bpmLink)

	/* через xml.DECODE НЕ смог обработать результат
	//ORIGINAL
	type UserList struct {
		XMLName xml.Name
		Body    struct {
			XMLName           xml.Name
			ListUsersResponse struct {
				XMLName xml.Name
				Return  []string `xml:"return"`
			} `xml:"listUsersResponse"`
		}
	}
	type TicketNumberID struct {
		XMLName xml.Name
		//XMLNS   xml.Attr
		Body struct {
			XMLName               xml.Name
			createRequestResponse struct {
				XMLName    xml.Name
				Code       int    `xml:"Code,omitempty"`
				ID         string `xml:"ID"`
				Number     string `xml:"Number"`
				SystemName string `xml:"SystemName"`
			} `xml:"createRequestResponse"`
		} `xml:"Body"`
	}
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

	//result := new(UserList)
	result := new(TicketNumberID)
	//result := new(Envelope)
	err = xml.NewDecoder(res.Body).Decode(result)
	if err != nil {
		log.Fatal("Error on unmarshaling xml. ", err.Error())
		return
	}
	//users := result.Body.ListUsersResponse.Return
	//systems := result.Body.readSystemsResponse.Table.row.Cell[0]
	ticket := result.Body.createRequestResponse.Number
	//ticket := result.Body.CreateRequestResponse.Number
	//fmt.Println(strings.Join(users, ", "))
	//fmt.Println(strings.Join(systems, ", "))
	fmt.Println(ticket)
	*/

}

func mainCreateZabbixTicket() {
	// https://novalagung.medium.com/soap-wsdl-request-in-go-language-3861cfb5949e
	// Создание копирующей заявки, как это реализовано сейчас у группы мониторинга БП

	url := "http://10.12.15.148/specs/aoi/tele2/bpm/bpmPortType" //PROD
	//url := "http://10.246.37.15:8060/specs/aoi/tele2/bpm/bpmPortType" //TEST
	userlogin := "service.monbpget"
	/*
		aps := []string{
			"NOV-FL2-CONFROOM-01",
			"NOV-FL2-CONFROOM-02",
			"NOV-FL2-CONFROOM-03",
		}*/
	controller := "t2ru-cntrl-01 "
	//controller := "t2ru-cntrl-01 "
	apName := "NOV-FL2-CONFROOM-01"
	//desAPs := strings.Join(aps, "\n")

	//Problem on host: t2ru-cntrl-01
	//Problem is [Северо-Запад] UAP "NOV-FL2-CONFROOM-01" is disconnected > 30 min: PROBLEM
	//description := "У клиента зафиксированы следующие Аномалии:" + "\n" + desAnomalies + "\n" + ""
	description := "Problem on host: " + controller + "\n" + "Problem is [Северо-Запад] UAP " + apName + " is disconnected > 30 min: PROBLEM" + "\n" + ""
	region := "Северо-Запад"

	strBefore := "<soapenv:Envelope xmlns:soapenv=\"http://schemas.xmlsoap.org/soap/envelope/\" xmlns:bpm=\"http://www.bercut.com/specs/aoi/tele2/bpm\"><soapenv:Header/><soapenv:Body><bpm:createRequestRequest><SystemId>5594b877-3bb7-46db-99f5-3c75b3e46556</SystemId><ServiceId>ed84a37f-4b31-4dab-85fe-ba4fe87325b1</ServiceId><Subject>Description</Subject><UserName>UserLogin</UserName><RequestType>Request</RequestType><Priority>Normal</Priority><Filds><ID>5c8dee23-e48a-45bc-a084-573e1a6cc5ca</ID><Value>Region</Value></Filds></bpm:createRequestRequest></soapenv:Body></soapenv:Envelope>"
	//replacer := strings.NewReplacer("Description", "My des", "UserLogin", "denis.tirskikh", "Region", "Москва ЦФ")
	replacer := strings.NewReplacer("Description", description, "UserLogin", userlogin, "Region", region)
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
	//Посмотреть response Body, если понадобится
	//
	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(b))
	//os.Exit(0)

	//Вбиваем результат из постмана сюда
	//https://tool.hiofd.com/en/xml-to-go/
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
	sr := envelope.Body.CreateRequestResponse.Number
	srID := envelope.Body.CreateRequestResponse.ID
	//bpmLink := "https://t2ru-tr-tst-01.corp.tele2.ru/0/Nui/ViewModule.aspx#CardModuleV2/CasePage/edit/" + srID //TEST
	bpmLink := "https://bpm.tele2.ru/0/Nui/ViewModule.aspx#CardModuleV2/CasePage/edit/" + srID //PROD
	//fmt.Println(envelope.Body.CreateRequestResponse.Number, error)
	//fmt.Println(envelope.Body.CreateRequestResponse.Number)
	fmt.Println(sr)
	fmt.Println(srID)
	fmt.Println(bpmLink)

	/* через xml.DECODE НЕ смог обработать результат
	//ORIGINAL
	type UserList struct {
		XMLName xml.Name
		Body    struct {
			XMLName           xml.Name
			ListUsersResponse struct {
				XMLName xml.Name
				Return  []string `xml:"return"`
			} `xml:"listUsersResponse"`
		}
	}
	type TicketNumberID struct {
		XMLName xml.Name
		//XMLNS   xml.Attr
		Body struct {
			XMLName               xml.Name
			createRequestResponse struct {
				XMLName    xml.Name
				Code       int    `xml:"Code,omitempty"`
				ID         string `xml:"ID"`
				Number     string `xml:"Number"`
				SystemName string `xml:"SystemName"`
			} `xml:"createRequestResponse"`
		} `xml:"Body"`
	}
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

	//result := new(UserList)
	result := new(TicketNumberID)
	//result := new(Envelope)
	err = xml.NewDecoder(res.Body).Decode(result)
	if err != nil {
		log.Fatal("Error on unmarshaling xml. ", err.Error())
		return
	}
	//users := result.Body.ListUsersResponse.Return
	//systems := result.Body.readSystemsResponse.Table.row.Cell[0]
	ticket := result.Body.createRequestResponse.Number
	//ticket := result.Body.CreateRequestResponse.Number
	//fmt.Println(strings.Join(users, ", "))
	//fmt.Println(strings.Join(systems, ", "))
	fmt.Println(ticket)
	*/

}
