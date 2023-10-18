package soap

import (
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/deniskaponchik/GoSoft/Unifi/internal/entity"
	"io"
	"net/http"
	"strings"
	"time"
)

type PolySoap struct {
	soapUrl string
	bpmUrl  string
	//2 вариант
	//usecase.PolyTicket
	//3 вариант
	//srStatusCodesForNewTicket    map[string]bool
	//srStatusCodesForCancelTicket map[string]bool
}

func New(s string, b string) *PolySoap {
	return &PolySoap{
		soapUrl: s,
		bpmUrl:  b,
	}
}

func (ps *PolySoap) CreatePolyTicketErr(ticket entity.Ticket) (entity.Ticket, error) { //srSlice []string, err error) {

	ticket.BpmServer = ps.bpmUrl //оставь в таком виде. не нужно при успешном выполнении формировать полную ссылку.

	if ticket.UserLogin != "" {
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
		//fmt.Println(strAfter)
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

		var err error
		myError := 1
		for myError != 0 {
			//req, err :=	http.NewRequest(httpMethod, url, bytes.NewReader(payload))
			req, errHttpReq := http.NewRequest("POST", ps.soapUrl, bytes.NewReader(payload))
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
						errXmlUnmarshal := xml.Unmarshal(bodyByte, envelope)
						if errXmlUnmarshal == nil {
							if envelope.Body.CreateRequestResponse.Code != 0 || envelope.Body.CreateRequestResponse.ID == "" {
								fmt.Println(envelope.Body.CreateRequestResponse.Description)
								//fmt.Println(envelope.Body.CreateRequestResponse.Code)
								//fmt.Println(envelope.Body.CreateRequestResponse.Text)
								fmt.Println("Заявка НЕ создалась на ФИНАЛЬНОМ этапе")
								fmt.Println("Проверь доступность SOAP-сервера и корректность входных данных:")
								fmt.Println("SOAP-сервер: " + ps.soapUrl)
								fmt.Println("User login: " + ticket.UserLogin)
								fmt.Println("Регион: " + ticket.Region)
								fmt.Println("Тип инцидента: " + ticket.IncidentType)
								fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
								fmt.Println("")
								time.Sleep(30 * time.Second)
								myError++
								err = errors.New("заявка не создалась на финальном этапе")
							} else {
								//Успешное завершение функции
								ticket.ID = envelope.Body.CreateRequestResponse.ID
								ticket.Number = envelope.Body.CreateRequestResponse.Number
								ticket.Url = ps.bpmUrl + ticket.ID
								myError = 0
								return ticket, nil
							}
						} else {
							fmt.Println(errXmlUnmarshal.Error())
							fmt.Println("Ошибка перекодировки ответа в xml")
							fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
							time.Sleep(30 * time.Second)
							myError++
							err = errXmlUnmarshal
						}
					} else {
						fmt.Println(errIOread.Error())
						fmt.Println("Ошибка чтения байтов из ответа")
						fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
						time.Sleep(30 * time.Second)
						myError++
						err = errIOread
					}
				} else {
					fmt.Println(errClientDo.Error())
					fmt.Println("Ошибка отправки запроса")
					fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
					time.Sleep(30 * time.Second)
					myError++
					err = errClientDo
				}
			} else {
				fmt.Println(errHttpReq.Error())
				fmt.Println("Ошибка создания объекта запроса")
				fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
				time.Sleep(30 * time.Second)
				myError++
				err = errHttpReq
			}
			if myError == 6 {
				myError = 0
				fmt.Println("После 6 неудачных попыток идём дальше. Заявка не была создана")
				//nil в ticket использовать не рекомендую, потому что значения теоретически потом пойдут в БД
				//ticket.ID = ""
				//ticket.Number = ""
				//ticket.Url = ""
				return ticket, err
			}
		}
	} else {
		//Для аномальных заявок
		fmt.Println("Заявка не была создана. User Login пустой")
		//nil в ticket использовать не рекомендую, потому что значения теоретически потом пойдут в БД
		ticket.ID = ""
		ticket.Number = ""
		ticket.Url = ""
		return ticket, errors.New("userlogin is empty")
	}
	return ticket, nil
}

func (ps *PolySoap) CheckTicketStatusErr(ticket entity.Ticket) (entity.Ticket, error) {

	ticket.BpmServer = ps.bpmUrl

	//if len(srID) == 36 {
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

	var err error
	myError := 1
	for myError != 0 {
		//req, errHttpReq := http.NewRequest(httpMethod, url, bytes.NewReader(payload))
		req, errHttpReq := http.NewRequest(httpMethod, ps.soapUrl, bytes.NewReader(payload))
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
					errXmlUnmarshal := xml.Unmarshal(bodyByte, envelope)
					if errXmlUnmarshal == nil {
						//if envelope.Body.GetStatusResponse.Code == 0 || envelope.Body.GetStatusResponse.StatisId != "" { //не решился пока что поменять if и else местами
						if envelope.Body.GetStatusResponse.Code != 0 || envelope.Body.GetStatusResponse.StatisId == "" {
							fmt.Println(envelope.Body.GetStatusResponse.Description)
							fmt.Println("Попытка получения Статуса обращения оборвалась на ПОСЛЕДНЕМ этапе")
							fmt.Println("Проверь доступность SOAP-сервера и корректность входных данных:")
							fmt.Println("SOAP-сервер: " + ps.soapUrl)
							fmt.Println("SR id: " + ticket.ID)
							fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
							fmt.Println("")
							time.Sleep(30 * time.Second)
							myError++
							err = errors.New("не удалось проверить статус на финальном этапе")
						} else {
							//Успешное завершение функции
							ticket.Status = envelope.Body.GetStatusResponse.Status
							myError = 0
							return ticket, nil
						}
					} else {
						fmt.Println("Ошибка перекодировки ответа в xml")
						fmt.Println(errXmlUnmarshal.Error())
						fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
						time.Sleep(30 * time.Second)
						myError++
						err = errXmlUnmarshal
					}
				} else {
					fmt.Println("Ошибка чтения байтов из ответа")
					fmt.Println(errIOread.Error())
					fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
					time.Sleep(30 * time.Second)
					myError++
					err = errIOread
				}
			} else {
				fmt.Println("Ошибка отправки запроса")
				fmt.Println(errClientDo.Error())
				fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
				time.Sleep(30 * time.Second)
				myError++
				err = errClientDo
				//Если ночью нет доступа к SOAP = в ЦОДЕ коллапс. Могу подождать 5 часов
				//if myError == 300 { 					myError = 0				}
			}
		} else {
			fmt.Println("Ошибка создания объекта запроса")
			fmt.Println(errHttpReq.Error())
			fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
			time.Sleep(30 * time.Second)
			myError++
			err = errHttpReq
			//Если ночью нет доступа к SOAP = в ЦОДЕ коллапс. Могу подождать 5 часов
			//if myError == 300 { 					myError = 0				}
		}
		if myError == 6 {
			myError = 0
			fmt.Println("После 6 неудачных попыток идём дальше. Статус заявки НЕ был уточнён")
			//ticketOut.Status = ""
			return ticket, err
		}
	}
	/* }else {
		//если передаётся пустая строка, не зная, существует ли заявка
		statusSlice = append(statusSlice, "0")
		statusSlice = append(statusSlice, "Тикет введён не корректно") //МЕНЯТЬ НЕ НУЖНО!!!
	}*/
	return ticket, nil
}

func (ps *PolySoap) ChangeStatusErr(ticket entity.Ticket) error {
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

	var err error
	myError := 1
	for myError != 0 {
		//req, err := http.NewRequest(httpMethod, url, bytes.NewReader(payload))
		req, errHttpReq := http.NewRequest(httpMethod, ps.soapUrl, bytes.NewReader(payload))
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
					errXmlUnmarshal := xml.Unmarshal(bodyByte, envelope)
					if errXmlUnmarshal == nil {
						if envelope.Body.ChangeCaseStatusResponse.Code != 0 || envelope.Body.ChangeCaseStatusResponse.NewStatusId == "" {
							fmt.Println(envelope.Body.ChangeCaseStatusResponse.Description)
							fmt.Println("НЕ УДАЛОСЬ изменить статус обращения на " + ticket.Status)
							fmt.Println("Проверь доступность SOAP-сервера и корректность входных данных:")
							fmt.Println("SOAP-сервер: " + ps.soapUrl)
							fmt.Println("SR id: " + ticket.ID)
							fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
							time.Sleep(30 * time.Second)
							fmt.Println("")
							myError++
							err = errors.New("не удалось изменить статус на финальном этапе")
						} else {
							//Успешное завершение функции
							srDateChange := envelope.Body.ChangeCaseStatusResponse.ModifyOn
							ticket.Status = envelope.Body.ChangeCaseStatusResponse.NewStatusId
							fmt.Println("Статус обращения изменён на " + ticket.Status + " в: " + srDateChange)
							myError = 0
							return nil
						}
					} else {
						fmt.Println("Ошибка перекодировки ответа в xml")
						fmt.Println(errXmlUnmarshal.Error())
						fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
						time.Sleep(30 * time.Second)
						myError++
						err = errXmlUnmarshal
					}
				} else {
					fmt.Println("Ошибка чтения байтов из ответа")
					fmt.Println(errIOread.Error())
					fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
					time.Sleep(30 * time.Second)
					myError++
					err = errIOread
				}
			} else {
				fmt.Println("Ошибка отправки запроса")
				fmt.Println(errClientDo.Error())
				fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
				time.Sleep(30 * time.Second)
				myError++
				err = errClientDo
			}
		} else {
			fmt.Println("Ошибка создания объекта запроса")
			fmt.Println(errHttpReq.Error())
			fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
			time.Sleep(30 * time.Second)
			myError++
			err = errHttpReq
		}
		if myError == 6 {
			myError = 0
			fmt.Println("После 6 неудачных попыток идём дальше. Статус заявки НЕ был изменён")
			//srNewStatus = ""
			return err
		}
	}
	return nil
}

func (ps *PolySoap) AddCommentErr(ticket entity.Ticket) (err error) {
	userLogin := "denis.tirskikh"
	//Убрать из строки \n
	//strBefore := "<Envelope xmlns=\"http://schemas.xmlsoap.org/soap/envelope/\"><Body><createCommentRequest xmlns=\"http://www.bercut.com/specs/aoi/tele2/bpm\"><CaseId xmlns=\"\">srID</CaseId><Message xmlns=\"\">myComment</Message><Author xmlns=\"\">userLogin</Author></createCommentRequest></Body></Envelope>"
	strBefore := "<Envelope xmlns=\"http://schemas.xmlsoap.org/soap/envelope/\"><Body><createCommentRequest xmlns=\"http://www.bercut.com/specs/aoi/tele2/bpm\"><CaseId>srID</CaseId><Message>myComment</Message><Author>userLogin</Author></createCommentRequest></Body></Envelope>"
	replacer := strings.NewReplacer("srID", ticket.ID, "myComment", ticket.Comment, "userLogin", userLogin)
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
		req, errHttpReq := http.NewRequest(httpMethod, ps.soapUrl, bytes.NewReader(payload))
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
					errXmlUnmarshal := xml.Unmarshal(bodyByte, envelope)
					if errXmlUnmarshal == nil {
						if envelope.Body.CreateCommentResponse.Code != 0 || envelope.Body.CreateCommentResponse.CreatedOn == "" {
							fmt.Println(envelope.Body.CreateCommentResponse.Description)
							fmt.Println("Попытка оставить комментарий ОБОРВАЛАСЬ на ПОСЛЕДНЕМ этапе")
							fmt.Println("Проверь доступность SOAP-сервера и корректность входных данных:")
							fmt.Println("SOAP-сервер: " + ps.soapUrl)
							fmt.Println("SR id: " + ticket.ID)
							fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
							fmt.Println("")
							time.Sleep(30 * time.Second)
							myError++
							err = errors.New("не удалось добавить комментарий на финальном этапе")
						} else {
							//srDateComment := envelope.Body.CreateCommentResponse.CreatedOn
							fmt.Println("Оставлен комментарий в ")
							fmt.Println(ps.bpmUrl + ticket.ID)
							//createdOn = envelope.Body.CreateCommentResponse.CreatedOn
							fmt.Println(envelope.Body.CreateCommentResponse.CreatedOn)
							myError = 0
							return nil
						}
					} else {
						fmt.Println("Ошибка перекодировки ответа в xml")
						fmt.Println(errXmlUnmarshal.Error())
						fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
						time.Sleep(30 * time.Second)
						myError++
						err = errXmlUnmarshal
					}
				} else {
					fmt.Println("Ошибка чтения байтов из ответа")
					fmt.Println(errIOread.Error())
					fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
					time.Sleep(30 * time.Second)
					myError++
					err = errIOread
				}
			} else {
				fmt.Println("Ошибка отправки запроса")
				fmt.Println(errClientDo.Error())
				fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
				time.Sleep(30 * time.Second)
				myError++
				err = errClientDo
			}
		} else {
			fmt.Println("Ошибка создания объекта запроса")
			fmt.Println(errHttpReq.Error())
			fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
			time.Sleep(30 * time.Second)
			myError++
			err = errHttpReq
		}
		if myError == 6 {
			myError = 0
			fmt.Println("После 6 неудачных попыток идём дальше. Комментарий не был оставлен")
			//createdOn = ""
			return err
		}
	}
	return nil
}
