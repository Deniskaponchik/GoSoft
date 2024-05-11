package api_rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/deniskaponchik/GoSoft/internal/entity"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type C3po struct {
	client http.Client
	//serverC3po	 string
	url string
}

func NewC3po(url string) *C3po {
	client := http.Client{
		Timeout: 240 * time.Second,
	}
	return &C3po{
		client: client,
		url:    url,
	}
}

func (c3po *C3po) GetUserLogin(notebook *entity.Client) (err error) {

	type Envelope struct {
		Name   string `json:"name"`
		Status string `json:"status"`

		Data []struct {
			//Monitor1_Model  string `json:"Monitor1_Model"`
			//Monitor1_SN     string `json:"Monitor1_SN"`
			//Monitor1_Vendor string `json:"Monitor1_Vendor"`
			//Monitor2_Model  string `json:"Monitor2_Model"`
			//Monitor2_SN     string `json:"Monitor2_SN"`
			//Monitor2_Vendor string `json:"Monitor2_Vendor"`
			//OZU              int    `json:"OZU"`
			//Disk1_model      string `json:"disk1_model"`
			//Disk2_model      string `json:"disk2_model"`
			//Disk3_model      string `json:"disk3_model"`
			//Last_scan        string `json:"last_scan"`
			//Os_build         string `json:"os_build"`
			//Os_version       string `json:"os_version"`
			//Pc_cpu           string `json:"pc_cpu"`
			//Pc_manufacturer  string `json:"pc_manufacturer"`
			//Pc_model         string `json:"pc_model"`
			//Pc_name          string `json:"pc_name"`
			//Pc_serial_number string `json:"pc_serial_number"`
			//Ram1             int    `json:"ram1"`
			//Ram2             int    `json:"ram2"`
			//Ram3             int    `json:"ram3"`
			//Ram4             int    `json:"ram4"`
			Samaccountname string `json:"samaccountname"`
		} `json:"data"`
	}
	//client := http.Client{Timeout: 5 * time.Second}
	client := c3po.client

	//var err error
	myError := 1
	for myError != 0 {
		//url := "http://" + polyStruct.IP + "/api/v1/mgmt/lineInfo"
		url := c3po.url + "pc/" + notebook.Hostname
		log.Println(url)

		req, errNewRequest := http.NewRequest(http.MethodGet, url, http.NoBody)
		if errNewRequest == nil {
			//req.SetBasicAuth(polyLogin, polyPassword)
			//req.Header.Add("Content-Type", "application/json")

			res, errClientDo := client.Do(req)
			if errClientDo == nil {
				defer res.Body.Close()

				envelope := &Envelope{}

				//https://forum.golangbridge.org/t/why-getting-err-eof-while-decoding-responce-body-into-struct/27444
				/*Посмотреть res.Body
				b, errIoRead := io.ReadAll(res.Body)
				if errIoRead != nil {
					log.Fatalln(errIoRead)
				}
				fmt.Println("Вывод из ioReadAll")
				fmt.Println(string(b))
				*/

				//
				body, errReadAll := ioutil.ReadAll(res.Body)
				if errReadAll != nil {
					log.Fatal(errReadAll)
				}
				buf := bytes.NewBuffer(body)
				//err = json.NewDecoder(buf).Decode(envelope)
				/*
					errUnmarshal := json.Unmarshal(body, envelope)
					if errUnmarshal == nil {
						log.Fatal(errUnmarshal)
					}*/

				//if errDecode := json.NewDecoder(res.Body).Decode(envelope); errDecode == nil {
				if errDecode := json.NewDecoder(buf).Decode(envelope); errDecode == nil {
					log.Println(envelope.Status)
					if envelope.Status == "ok" {
						//log.Println("Запрос статуса прошёл.")
						if len(envelope.Data) > 0 {
							//Успешное выполнение функции
							notebook.UserLogin = strings.ReplaceAll(envelope.Data[0].Samaccountname, " ", "") //cut spaces
							myError = 0
						} else {
							log.Println("Получен ответ OK от c3po, но информация о машине не была найдена")
							myError = 0
						}
						return nil
					} else {
						log.Println("От c3po получен Статус НЕ OK")
						log.Println("Будет предпринята новая попытка отправки запроса через 30 сек.")
						time.Sleep(30 * time.Second)
						myError++
						err = errors.New("от устройства получен статус не ok")
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
			return err
		}
	}
	return nil
}
