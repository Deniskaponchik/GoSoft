package api_soap

import (
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"errors"
	"github.com/deniskaponchik/GoSoft/Unifi/internal/entity"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type Soap struct {
	soapUrl    string
	bpmUrl     string
	httpClient *http.Client
	//2 вариант
	//usecase.PolyTicket
	//3 вариант
	//srStatusCodesForNewTicket    map[string]bool
	//srStatusCodesForCancelTicket map[string]bool
}

func NewSoap(s string, b string) *Soap {
	log.Println(s)
	log.Println(b)

	return &Soap{
		soapUrl: s,
		bpmUrl:  b,
		httpClient: &http.Client{
			Timeout: 240 * time.Second, //bpm часто лагает. при 120 выдаёт: (context deadline exceeded (Client.Timeout exceeded while awaiting headers)

			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
	}
}

func (ss *Soap) CreateTicketSmacWifi(ticket *entity.Ticket) (err error) { //(entity.Ticket, error) { //srSlice []string, err error) {

	ticket.BpmServer = ss.bpmUrl //оставь в таком виде. не нужно при успешном выполнении формировать полную ссылку.

	if ticket.UserLogin == "" {
		ticket.UserLogin = "denis.tirskikh"
	}

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
	replacer := strings.NewReplacer("Description", ticket.Description, "UserLogin", ticket.UserLogin, "Reason", ticket.Reason,
		"Region", ticket.Region, "Monitoring", ticket.Monitoring, "incidentType", ticket.IncidentType)
	strAfter := replacer.Replace(strBefore)
	//log.Println(strAfter)
	payload := []byte(strAfter)
	//os.Exit(0)

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

	//var err error
	myError := 1
	for myError != 0 {
		//req, err :=	http.NewRequest(httpMethod, url, bytes.NewReader(payload))
		req, errHttpReq := http.NewRequest("POST", ss.soapUrl, bytes.NewReader(payload))
		if errHttpReq == nil {
			/*
				client := &http.Client{
					Timeout: 120 * time.Second, //bpm часто лагает
					Transport: &http.Transport{
						TLSClientConfig: &tls.Config{
							InsecureSkipVerify: true,
						},
					},
				}
				res, errClientDo := client.Do(req) */
			res, errClientDo := ss.httpClient.Do(req)
			if errClientDo == nil {
				/*Посмотреть response Body, если понадобится
				defer res.Body.Close()
				b, err := io.ReadAll(res.Body)
				if err != nil {
					log.Fatalln(err)
				}
				log.Println(string(b))
				//os.Exit(0)
				*/

				// Смог победить только через unmarshal. Кривенько косо, но работает и куча времени угрохано даже на это
				envelope := &Envelope{}
				bodyByte, errIOread := io.ReadAll(res.Body)
				if errIOread == nil {
					errXmlUnmarshal := xml.Unmarshal(bodyByte, envelope)
					if errXmlUnmarshal == nil {
						if envelope.Body.CreateRequestResponse.Code != 0 || envelope.Body.CreateRequestResponse.ID == "" {
							log.Println(envelope.Body.CreateRequestResponse.Description)
							//log.Println(envelope.Body.CreateRequestResponse.Code)
							//log.Println(envelope.Body.CreateRequestResponse.Text)
							log.Println("Заявка НЕ создалась на ФИНАЛЬНОМ этапе")
							log.Println("Проверь доступность SOAP-сервера и корректность входных данных:")
							log.Println("SOAP-сервер: " + ss.soapUrl)
							log.Println("User login: " + ticket.UserLogin)
							log.Println("Регион: " + ticket.Region)
							log.Println("Тип инцидента: " + ticket.IncidentType)
							log.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
							log.Println("")
							time.Sleep(30 * time.Second)
							myError = 6
							err = errors.New("заявка не создалась на финальном этапе")
							return err
						} else {
							//Успешное завершение функции
							ticket.ID = envelope.Body.CreateRequestResponse.ID
							ticket.Number = envelope.Body.CreateRequestResponse.Number
							ticket.Url = ss.bpmUrl + ticket.ID
							myError = 0
							return nil //ticket, nil
						}
					} else {
						log.Println(errXmlUnmarshal.Error())
						log.Println("Ошибка перекодировки ответа в xml")
						log.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
						time.Sleep(30 * time.Second)
						myError++
						err = errXmlUnmarshal
					}
				} else {
					log.Println(errIOread.Error())
					log.Println("Ошибка чтения байтов из ответа")
					log.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
					time.Sleep(30 * time.Second)
					myError++
					err = errIOread
				}
			} else {
				log.Println(errClientDo.Error())
				log.Println("Ошибка отправки запроса")
				log.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
				time.Sleep(30 * time.Second)
				myError++
				err = errClientDo
			}
		} else {
			log.Println(errHttpReq.Error())
			log.Println("Ошибка создания объекта запроса")
			log.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
			time.Sleep(30 * time.Second)
			myError++
			err = errHttpReq
		}
		if myError == 6 {
			myError = 0
			log.Println("После 6 неудачных попыток идём дальше. Заявка не была создана")
			//nil в ticket использовать не рекомендую, потому что значения теоретически потом пойдут в БД
			return err
		}
	}
	return nil
}

func (ss *Soap) CreateTicketSmacVcs(ticket *entity.Ticket) (err error) { //(entity.Ticket, error) { //srSlice []string, err error) {

	ticket.BpmServer = ss.bpmUrl //оставь в таком виде. не нужно при успешном выполнении формировать полную ссылку.

	if ticket.UserLogin == "" {
		ticket.UserLogin = "denis.tirskikh"
	}

	strBefore :=
		"<soapenv:Envelope " +
			"xmlns:soapenv=\"http://schemas.xmlsoap.org/soap/envelope/\" " +
			"xmlns:bpm=\"http://www.bercut.com/specs/aoi/tele2/bpm\">" +
			"<soapenv:Header/>" +
			"<soapenv:Body>" +
			"<bpm:createRequestRequest>" +
			"<SystemId>5594b877-3bb7-46db-99f5-3c75b3e46556</SystemId>" +
			"<ServiceId>8ec1af6d-c717-449b-837b-7bd443fab97a</ServiceId>" +
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
			"<ID>a83e9532-8951-4dde-b1bf-fe5dcf26c50e</ID>" +
			"<Value>incidentType</Value>" +
			"</Filds>" +
			"</bpm:createRequestRequest>" +
			"</soapenv:Body>" +
			"</soapenv:Envelope>"
	//replacer := strings.NewReplacer("Description", "My des", "UserLogin", "denis.tirskikh", "Region", "Москва ЦФ")
	//replacer := strings.NewReplacer("Description", description, "UserLogin", userLogin, "incidentType", incidentType, "Region", region)
	replacer := strings.NewReplacer("Description", ticket.Description, "UserLogin", ticket.UserLogin, "Reason", ticket.Reason,
		"Region", ticket.Region, "Monitoring", ticket.Monitoring, "incidentType", ticket.IncidentType)
	strAfter := replacer.Replace(strBefore)
	//log.Println(strAfter)
	payload := []byte(strAfter)
	//os.Exit(0)

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

	//var err error
	myError := 1
	for myError != 0 {
		//req, err :=	http.NewRequest(httpMethod, url, bytes.NewReader(payload))
		req, errHttpReq := http.NewRequest("POST", ss.soapUrl, bytes.NewReader(payload))
		if errHttpReq == nil {
			/*
				client := &http.Client{
					Timeout: 120 * time.Second, //bpm часто лагает
					Transport: &http.Transport{
						TLSClientConfig: &tls.Config{
							InsecureSkipVerify: true,
						},
					},
				}
				res, errClientDo := client.Do(req) */
			res, errClientDo := ss.httpClient.Do(req)
			if errClientDo == nil {
				/*Посмотреть response Body, если понадобится
				defer res.Body.Close()
				b, err := io.ReadAll(res.Body)
				if err != nil {
					log.Fatalln(err)
				}
				log.Println(string(b))
				//os.Exit(0)
				*/

				// Смог победить только через unmarshal. Кривенько косо, но работает и куча времени угрохано даже на это
				envelope := &Envelope{}
				bodyByte, errIOread := io.ReadAll(res.Body)
				if errIOread == nil {
					errXmlUnmarshal := xml.Unmarshal(bodyByte, envelope)
					if errXmlUnmarshal == nil {
						if envelope.Body.CreateRequestResponse.Code != 0 || envelope.Body.CreateRequestResponse.ID == "" {
							log.Println(envelope.Body.CreateRequestResponse.Description)
							//log.Println(envelope.Body.CreateRequestResponse.Code)
							//log.Println(envelope.Body.CreateRequestResponse.Text)
							log.Println("Заявка НЕ создалась на ФИНАЛЬНОМ этапе")
							log.Println("Проверь доступность SOAP-сервера и корректность входных данных:")
							log.Println("SOAP-сервер: " + ss.soapUrl)
							log.Println("User login: " + ticket.UserLogin)
							log.Println("Регион: " + ticket.Region)
							log.Println("Тип инцидента: " + ticket.IncidentType)
							log.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
							log.Println("")
							time.Sleep(30 * time.Second)
							myError = 6
							err = errors.New("заявка не создалась на финальном этапе")
							return err
						} else {
							//Успешное завершение функции
							ticket.ID = envelope.Body.CreateRequestResponse.ID
							ticket.Number = envelope.Body.CreateRequestResponse.Number
							ticket.Url = ss.bpmUrl + ticket.ID
							myError = 0
							return nil
						}
					} else {
						log.Println(errXmlUnmarshal.Error())
						log.Println("Ошибка перекодировки ответа в xml")
						log.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
						time.Sleep(30 * time.Second)
						myError++
						err = errXmlUnmarshal
					}
				} else {
					log.Println(errIOread.Error())
					log.Println("Ошибка чтения байтов из ответа")
					log.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
					time.Sleep(30 * time.Second)
					myError++
					err = errIOread
				}
			} else {
				log.Println(errClientDo.Error())
				log.Println("Ошибка отправки запроса")
				log.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
				time.Sleep(30 * time.Second)
				myError++
				err = errClientDo
			}
		} else {
			log.Println(errHttpReq.Error())
			log.Println("Ошибка создания объекта запроса")
			log.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
			time.Sleep(30 * time.Second)
			myError++
			err = errHttpReq
		}
		if myError == 6 {
			myError = 0
			log.Println("После 6 неудачных попыток идём дальше. Заявка не была создана")
			return err
		}
	}
	return nil
}

func (ss *Soap) CheckTicketStatusErr(ticket *entity.Ticket) (err error) {

	ticket.BpmServer = ss.bpmUrl //не убирать. Используется вне функции

	//Убрать из строки \n
	strBefore := "<soapenv:Envelope xmlns:soapenv=\"http://schemas.xmlsoap.org/soap/envelope/\" xmlns:bpm=\"http://www.bercut.com/specs/aoi/tele2/bpm\"><soapenv:Header/><soapenv:Body><bpm:getStatusRequest><CaseID>SRid</CaseID></bpm:getStatusRequest></soapenv:Body></soapenv:Envelope>"
	replacer := strings.NewReplacer("SRid", ticket.ID)
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

	//var err error
	myError := 1
	for myError != 0 {
		//req, errHttpReq := http.NewRequest(httpMethod, url, bytes.NewReader(payload))
		req, errHttpReq := http.NewRequest(httpMethod, ss.soapUrl, bytes.NewReader(payload))
		if errHttpReq == nil {
			/*
				client := &http.Client{
					Timeout: 120 * time.Second, //bpm часто лагает
					Transport: &http.Transport{
						TLSClientConfig: &tls.Config{
							InsecureSkipVerify: true,
						},
					},
				}
				res, errClientDo := client.Do(req) */
			res, errClientDo := ss.httpClient.Do(req)
			if errClientDo == nil {
				/*Посмотреть response Body, если понадобится
				defer res.Body.Close() //ОСТОРОЖНЕЕ с этой штукой. Дальше могут данные не пойти
				b, err := io.ReadAll(res.Body)
				if err != nil {
					log.Fatalln(err)
				}
				log.Println(string(b))
				//os.Exit(0)*/

				envelope := &Envelope{}
				bodyByte, errIOread := io.ReadAll(res.Body)
				if errIOread == nil {
					errXmlUnmarshal := xml.Unmarshal(bodyByte, envelope)
					if errXmlUnmarshal == nil {
						//if envelope.Body.GetStatusResponse.Code == 0 || envelope.Body.GetStatusResponse.StatisId != "" { //не решился пока что поменять if и else местами
						if envelope.Body.GetStatusResponse.Code != 0 || envelope.Body.GetStatusResponse.StatisId == "" {
							log.Println(envelope.Body.GetStatusResponse.Description)
							log.Println("Попытка получения Статуса обращения оборвалась на ПОСЛЕДНЕМ этапе")
							log.Println("Проверь доступность SOAP-сервера и корректность входных данных:")
							log.Println("SOAP-сервер: " + ss.soapUrl)
							log.Println("SR id: " + ticket.ID)
							log.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
							log.Println("")
							time.Sleep(30 * time.Second)
							myError++
							err = errors.New("не удалось проверить статус на финальном этапе")
						} else {
							//Успешное завершение функции
							ticket.Status = envelope.Body.GetStatusResponse.Status
							ticket.Url = ticket.BpmServer + ticket.ID
							myError = 0
							return nil
						}
					} else {
						log.Println("Ошибка перекодировки ответа в xml")
						log.Println(errXmlUnmarshal.Error())
						log.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
						time.Sleep(30 * time.Second)
						myError++
						err = errXmlUnmarshal
					}
				} else {
					log.Println("Ошибка чтения байтов из ответа")
					log.Println(errIOread.Error())
					log.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
					time.Sleep(30 * time.Second)
					myError++
					err = errIOread
				}
			} else {
				log.Println("Ошибка отправки запроса")
				log.Println(errClientDo.Error())
				log.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
				time.Sleep(30 * time.Second)
				myError++
				err = errClientDo
				//Если ночью нет доступа к SOAP = в ЦОДЕ коллапс. Могу подождать 5 часов
				//if myError == 300 { 					myError = 0				}
			}
		} else {
			log.Println("Ошибка создания объекта запроса")
			log.Println(errHttpReq.Error())
			log.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
			time.Sleep(30 * time.Second)
			myError++
			err = errHttpReq
			//Если ночью нет доступа к SOAP = в ЦОДЕ коллапс. Могу подождать 5 часов
			//if myError == 300 { 					myError = 0				}
		}
		if myError == 6 {
			myError = 0
			log.Println("После 6 неудачных попыток идём дальше. Статус заявки НЕ был уточнён")
			//ticketOut.Status = ""
			return err
		}
	}
	return nil
}

func (ss *Soap) ChangeStatusErr(ticket *entity.Ticket) (err error) {
	UserLogin := "denis.tirskikh"
	//Убрать из строки \n
	strBefore := "<Envelope xmlns=\"http://schemas.xmlsoap.org/soap/envelope/\"><Body><changeCaseStatusRequest xmlns=\"http://www.bercut.com/specs/aoi/tele2/bpm\"><CaseId xmlns=\"\">SRid</CaseId><Status xmlns=\"\">NewStatus</Status><User xmlns=\"\">UserLogin</User></changeCaseStatusRequest></Body></Envelope>"
	replacer := strings.NewReplacer("SRid", ticket.ID, "NewStatus", ticket.Status, "UserLogin", UserLogin)
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

	//var err error
	myError := 1
	for myError != 0 {
		//req, err := http.NewRequest(httpMethod, url, bytes.NewReader(payload))
		req, errHttpReq := http.NewRequest(httpMethod, ss.soapUrl, bytes.NewReader(payload))
		if errHttpReq == nil {
			/*
				client := &http.Client{
					Timeout: 120 * time.Second, //bpm часто лагает
					Transport: &http.Transport{
						TLSClientConfig: &tls.Config{
							InsecureSkipVerify: true,
						},
					},
				}
				res, errClientDo := client.Do(req) */
			res, errClientDo := ss.httpClient.Do(req)
			if errClientDo == nil {

				envelope := &Envelope{}
				bodyByte, errIOread := io.ReadAll(res.Body)
				if errIOread == nil {
					errXmlUnmarshal := xml.Unmarshal(bodyByte, envelope)
					if errXmlUnmarshal == nil {
						if envelope.Body.ChangeCaseStatusResponse.Code != 0 || envelope.Body.ChangeCaseStatusResponse.NewStatusId == "" {
							log.Println(envelope.Body.ChangeCaseStatusResponse.Description)
							log.Println("НЕ УДАЛОСЬ изменить статус обращения на " + ticket.Status)
							log.Println("Проверь доступность SOAP-сервера и корректность входных данных:")
							log.Println("SOAP-сервер: " + ss.soapUrl)
							log.Println("SR id: " + ticket.ID)
							log.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
							time.Sleep(30 * time.Second)
							log.Println("")
							myError++
							err = errors.New("не удалось изменить статус на финальном этапе")
						} else {
							//Успешное завершение функции
							srDateChange := envelope.Body.ChangeCaseStatusResponse.ModifyOn
							//ticket.Status = envelope.Body.ChangeCaseStatusResponse.NewStatusId
							log.Println("Статус обращения изменён на " + ticket.Status + " в: " + srDateChange)
							myError = 0
							return nil
						}
					} else {
						log.Println("Ошибка перекодировки ответа в xml")
						log.Println(errXmlUnmarshal.Error())
						log.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
						time.Sleep(30 * time.Second)
						myError++
						err = errXmlUnmarshal
					}
				} else {
					log.Println("Ошибка чтения байтов из ответа")
					log.Println(errIOread.Error())
					log.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
					time.Sleep(30 * time.Second)
					myError++
					err = errIOread
				}
			} else {
				log.Println("Ошибка отправки запроса")
				log.Println(errClientDo.Error())
				log.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
				time.Sleep(30 * time.Second)
				myError++
				err = errClientDo
			}
		} else {
			log.Println("Ошибка создания объекта запроса")
			log.Println(errHttpReq.Error())
			log.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
			time.Sleep(30 * time.Second)
			myError++
			err = errHttpReq
		}
		if myError == 6 {
			myError = 0
			log.Println("После 6 неудачных попыток идём дальше. Статус заявки НЕ был изменён")
			return err
		}
	}
	return nil
}

func (ss *Soap) AddCommentErr(ticket *entity.Ticket) (err error) {
	userLogin := "denis.tirskikh"
	//Убрать из строки \n
	//strBefore := "<Envelope xmlns=\"http://schemas.xmlsoap.org/soap/envelope/\"><Body><createCommentRequest xmlns=\"http://www.bercut.com/specs/aoi/tele2/bpm\"><CaseId xmlns=\"\">srID</CaseId><Message xmlns=\"\">myComment</Message><Author xmlns=\"\">userLogin</Author></createCommentRequest></Body></Envelope>"
	strBefore := "<Envelope xmlns=\"http://schemas.xmlsoap.org/soap/envelope/\"><Body><createCommentRequest xmlns=\"http://www.bercut.com/specs/aoi/tele2/bpm\"><CaseId>srID</CaseId><Message>myComment</Message><Author>userLogin</Author></createCommentRequest></Body></Envelope>"
	replacer := strings.NewReplacer("srID", ticket.ID, "myComment", ticket.Comment, "userLogin", userLogin)
	strAfter := replacer.Replace(strBefore)
	//log.Println(strAfter)
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
		req, errHttpReq := http.NewRequest(httpMethod, ss.soapUrl, bytes.NewReader(payload))
		if errHttpReq == nil {
			/*
				client := &http.Client{
					Timeout: 120 * time.Second, //bpm часто лагает
					Transport: &http.Transport{
						TLSClientConfig: &tls.Config{
							InsecureSkipVerify: true,
						},
					},
				}
				res, errClientDo := client.Do(req) */
			res, errClientDo := ss.httpClient.Do(req)
			if errClientDo == nil {
				/*Посмотреть response Body, если понадобится
				defer res.Body.Close() //ОСТОРОЖНЕЕ с этой штукой. Дальше могут данные не пойти
				b, err := io.ReadAll(res.Body)
				if err != nil {
					log.Fatalln(err)
				}
				log.Println(string(b))
				os.Exit(0)*/

				envelope := &Envelope{}
				bodyByte, errIOread := io.ReadAll(res.Body)
				if errIOread == nil {
					errXmlUnmarshal := xml.Unmarshal(bodyByte, envelope)
					if errXmlUnmarshal == nil {
						if envelope.Body.CreateCommentResponse.Code != 0 || envelope.Body.CreateCommentResponse.CreatedOn == "" {
							log.Println(envelope.Body.CreateCommentResponse.Description)
							log.Println("Попытка оставить комментарий ОБОРВАЛАСЬ на ПОСЛЕДНЕМ этапе")
							log.Println("Проверь доступность SOAP-сервера и корректность входных данных:")
							log.Println("SOAP-сервер: " + ss.soapUrl)
							log.Println("SR id: " + ticket.ID)
							log.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
							log.Println("")
							time.Sleep(30 * time.Second)
							myError++
							err = errors.New("не удалось добавить комментарий на финальном этапе")
						} else {
							//srDateComment := envelope.Body.CreateCommentResponse.CreatedOn
							//createdOn = envelope.Body.CreateCommentResponse.CreatedOn
							log.Println("Оставлен комментарий в ")
							log.Println(ss.bpmUrl + ticket.ID)
							log.Println(envelope.Body.CreateCommentResponse.CreatedOn)
							myError = 0
							return nil
						}
					} else {
						log.Println("Ошибка перекодировки ответа в xml")
						log.Println(errXmlUnmarshal.Error())
						log.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
						time.Sleep(30 * time.Second)
						myError++
						err = errXmlUnmarshal
					}
				} else {
					log.Println("Ошибка чтения байтов из ответа")
					log.Println(errIOread.Error())
					log.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
					time.Sleep(30 * time.Second)
					myError++
					err = errIOread
				}
			} else {
				log.Println("Ошибка отправки запроса")
				log.Println(errClientDo.Error())
				log.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
				time.Sleep(30 * time.Second)
				myError++
				err = errClientDo
			}
		} else {
			log.Println("Ошибка создания объекта запроса")
			log.Println(errHttpReq.Error())
			log.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
			time.Sleep(30 * time.Second)
			myError++
			err = errHttpReq
		}
		if myError == 6 {
			myError = 0
			log.Println("После 6 неудачных попыток идём дальше. Комментарий не был оставлен")
			return err
		}
	}
	return nil
}
