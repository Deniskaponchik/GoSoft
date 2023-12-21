package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {
	login, err := GetUserLogin("http://login:password@c3po.corp.tele2.ru/sccm/api/info/pc/WSNS-TROFIMOV2")
	//login, err := GetUserLogin("http://c3po.corp.tele2.ru/sccm/api/info/pc/WSNS-TROFIMOV2")
	//login, err := GetUserLogin("http://10.57.188.15/api/v1/mgmt/lineInfo")
	if err == nil {
		fmt.Print(login)
	} else {
		fmt.Printf(err.Error())
	}
}

func GetUserLogin(url string) (login string, err error) {

	type Envelope struct {
		Data []struct {
			//Monitor1_Model   string `json:"Monitor1_Model"`
			//Monitor1_SN      string `json:"Monitor1_SN"`
			Monitor1_Vendor  string `json:"Monitor1_Vendor"`
			Monitor2_Model   string `json:"Monitor2_Model"`
			Monitor2_SN      string `json:"Monitor2_SN"`
			Monitor2_Vendor  string `json:"Monitor2_Vendor"`
			OZU              int    `json:"OZU"`
			Disk1_model      string `json:"disk1_model;omitempty"`
			Disk2_model      string `json:"disk2_model;omitempty"`
			Disk3_model      string `json:"disk3_model;omitempty"`
			Last_scan        string `json:"last_scan"`
			Os_build         string `json:"os_build"`
			Os_version       string `json:"os_version"`
			Pc_cpu           string `json:"pc_cpu"`
			Pc_manufacturer  string `json:"pc_manufacturer"`
			Pc_model         string `json:"pc_model"`
			Pc_name          string `json:"pc_name"`
			Pc_serial_number string `json:"pc_serial_number"`
			Ram1             int    `json:"ram1;omitempty"`
			Ram2             int    `json:"ram2;omitempty"`
			Ram3             int    `json:"ram3;omitempty"`
			Ram4             int    `json:"ram4;omitempty"`
			Samaccountname   string `json:"samaccountname"`
		} `json:"data"`
		Name   string `json:"name"`
		Status string `json:"status"`
	}
	//client := http.Client{Timeout: 5 * time.Second}
	client := http.Client{Timeout: 240 * time.Second}

	//var err error
	myError := 1
	for myError != 0 {
		//url := "http://" + polyStruct.IP + "/api/v1/mgmt/lineInfo"
		//url := uc3po.url + "pc/" + notebook.Hostname
		log.Println(url)

		req, errNewRequest := http.NewRequest(http.MethodGet, url, http.NoBody)
		if errNewRequest == nil {
			//req.SetBasicAuth("c3po login", "c3po password")
			//req.Header.Add("Content-Type", "application/json")

			res, errClientDo := client.Do(req)
			if errClientDo == nil {
				envelope := &Envelope{}

				//https://forum.golangbridge.org/t/why-getting-err-eof-while-decoding-responce-body-into-struct/27444

				//defer res.Body.Close()

				/*Посмотреть res.Body
				b, errIoRead := io.ReadAll(res.Body)
				if errIoRead != nil {
					log.Fatalln(errIoRead)
				}
				fmt.Println(string(b))
				*/

				//
				body, errReadAll := ioutil.ReadAll(res.Body)
				if errReadAll != nil {
					log.Fatal(errReadAll)
				}
				buf := bytes.NewBuffer(body)
				/*
					errUnmarshal := json.Unmarshal(body, envelope)
					if errUnmarshal == nil {
						log.Fatal(errUnmarshal)
					}*/

				//if errDecode := json.NewDecoder(res.Body).Decode(&envelope); errDecode == nil {
				//errDecode := json.NewDecoder(res.Body).Decode(envelope)
				errDecode := json.NewDecoder(buf).Decode(envelope)
				if errDecode == nil { //errDecode.Error() == "EOF"
					if envelope.Status == "ok" { //envelope.Status == "ok" envelope.Status != ""
						//log.Println("Запрос статуса прошёл.")
						if len(envelope.Data) > 0 {
							//Успешное выполнение функции
							//notebook.UserLogin = envelope.Data[0].Samaccountname
							login = strings.ReplaceAll(envelope.Data[0].Samaccountname, " ", "")
							fmt.Println(envelope.Status)
							fmt.Println(envelope.Name)
							return login, nil
						} else {
							log.Println("Получен ответ OK от c3po, но информация о машине не была найдена")
							return "", nil
						}

					} else {
						log.Println("От c3po получен Статус НЕ OK")
						log.Println("Будет предпринята новая попытка отправки запроса через 30 сек.")
						time.Sleep(30 * time.Second)
						myError++
						err = errors.New("от устройства получен статус не 2000")
					}
				} else {
					log.Println(errDecode.Error())
					log.Println("Ошибка перекодировки ответа")
					log.Println("Скорее всего, API недоступен")
					log.Println("Будет предпринята новая попытка отправки запроса через 30 сек.")
					time.Sleep(30 * time.Second)
					myError++
					err = errDecode
				}
			} else {
				log.Println(errClientDo.Error())
				log.Println("Ошибка отправки запроса")
				log.Println("Будет предпринята новая попытка отправки запроса через 30 сек.")
				time.Sleep(30 * time.Second)
				myError++
				err = errClientDo
			}
		} else {
			log.Println(errNewRequest.Error())
			log.Println("Ошибка создания ОБЪЕКТА запроса")
			log.Println("Будет предпринята новая попытка отправки запроса через 30 сек.")
			time.Sleep(30 * time.Second)
			myError++
			err = errNewRequest
		}
		if myError == 4 {
			myError = 0
			log.Println("После 3 неудачных попыток идём дальше. Получить логин от c3po не удалось")
			return "", err
		}
	}
	return "", nil
}
