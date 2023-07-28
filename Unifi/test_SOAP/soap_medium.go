package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)
import (
	"bytes"
	"crypto/tls"
	"log"
	"net/http"
)

var (
	//Server = "http://10.12.15.148/specs/aoi/tele2/bpm/bpmPortType"   //Prod
	Server = "http://10.246.37.15:8060/specs/aoi/tele2/bpm/bpmPortType" //TEST
	/*
		StatusIdVising := "b32b613a-0282-4e8a-b831-1027e7c7972f"
		StatusIdCancel := "6e5f4218-f46b-1410-fe9a-0050ba5d6c38"
		StatusIdResolve := "ae7f411e-f46b-1410-009b-0050ba5d6c38"
		StatusIdClarification := "81e6a1ee-16c1-4661-953e-dde140624fb3"
		CloseCode_FullSolution := 200
	*/
)

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

func mainCOMMENT() {
	url := Server
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

func mainCHANGE(bpmServer string, srID string, NewStatus string) (srNewStatus string) {

	url := bpmServer
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

func mainCHECK() {
	//url := "http://10.246.37.15:8060/specs/aoi/tele2/bpm/bpmPortType"
	url := Server
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
