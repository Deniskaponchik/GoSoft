package main

//test_SOAP/soap_medium
import (
	"encoding/xml"
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
	//Server = "http://10.12.15.148/specs/aoi/tele2/bpm/bpmPortType"     //PROD
	Server = "http://10.246.37.15:8060/specs/aoi/tele2/bpm/bpmPortType" //TEST
)

func CreateSmacWiFiTicket(
	userLogin string, pcName string, anomalies []string, apName string, region string) (
	srNumber string, srID string, bpmLink string) {

	//url := "http://10.246.37.15:8060/specs/aoi/tele2/bpm/bpmPortType" //TEST
	//url := "http://10.12.15.148/specs/aoi/tele2/bpm/bpmPortType"   //PROD
	url := Server

	//fmt.Println("")
	//fmt.Println("SOAP function:")
	//userlogin := "denis.tirskikh"
	//fmt.Println(userLogin)
	//pcName := "wsir-tirskikh"
	//fmt.Println(pcName)
	/*
		anomaly := []string{
			"anomal1",
			"anomaly2",
			"anomaly3",
		}*/
	//fmt.Println(anomalies)
	desAnomalies := strings.Join(anomalies, "\n")

	//description := "Tootsie roll tiramisu maca" + "\n" + "Danish topping sugar plum tart bonbon "
	description := "На ноутбуке:" + "\n" + pcName + "\n" + "" + "\n" + "зафиксированы следующие Аномалии:" + "\n" + desAnomalies + "\n" + "" + "\n" + "Предполагаемое, но не на 100% точное имя точки:" + "\n" + apName + "\n" + "" + "\n" + "Рекомендации по выполнению таких инцидентов собраны на страничке корпоративной wiki" + "\n" + "https://wiki.tele2.ru/display/ITKB/%5BHelpdesk+IT%5D+System+Monitoring" + "\n" + ""
	//fmt.Println(description)

	//region := "Москва ЦФ"
	//fmt.Println(region)

	strBefore := "<soapenv:Envelope xmlns:soapenv=\"http://schemas.xmlsoap.org/soap/envelope/\" xmlns:bpm=\"http://www.bercut.com/specs/aoi/tele2/bpm\"><soapenv:Header/><soapenv:Body><bpm:createRequestRequest><SystemId>5594b877-3bb7-46db-99f5-3c75b3e46556</SystemId><ServiceId>ed84a37f-4b31-4dab-85fe-ba4fe87325b1</ServiceId><Subject>Description</Subject><UserName>UserLogin</UserName><RequestType>Request</RequestType><Priority>Normal</Priority><Filds><ID>5c8dee23-e48a-45bc-a084-573e1a6cc5ca</ID><Value>Region</Value></Filds></bpm:createRequestRequest></soapenv:Body></soapenv:Envelope>"
	//replacer := strings.NewReplacer("Description", "My des", "UserLogin", "denis.tirskikh", "Region", "Москва ЦФ")
	replacer := strings.NewReplacer("Description", description, "UserLogin", userLogin, "Region", region)
	strAfter := replacer.Replace(strBefore)
	//fmt.Println(strAfter)
	//time.Sleep(60 * time.Second)

	payload := []byte(strAfter)
	//os.Exit(0)

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
	srNumber = envelope.Body.CreateRequestResponse.Number
	srID = envelope.Body.CreateRequestResponse.ID
	bpmLink = "https://t2ru-tr-tst-01.corp.tele2.ru/0/Nui/ViewModule.aspx#CardModuleV2/CasePage/edit/" + srID

	//fmt.Println(sr)	fmt.Println(srID)	fmt.Println(bpmLink)
	return srNumber, srID, bpmLink
}

func CheckTicketStatus(srID string) (srStatus string, srStatusID string) {
	//url := "http://10.246.37.15:8060/specs/aoi/tele2/bpm/bpmPortType"
	//url := "http://10.12.15.148/specs/aoi/tele2/bpm/bpmPortType"   //PROD
	url := Server

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

	srStatus = envelope.Body.GetStatusResponse.Status
	srStatusID = envelope.Body.GetStatusResponse.StatisId
	return srStatus, srStatusID
}

func ClarifyTicket(srID string) {

}
