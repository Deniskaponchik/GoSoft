package api_rest

import (
	"encoding/json"
	"errors"
	"github.com/deniskaponchik/GoSoft/Unifi/internal/entity"
	"log"
	"net/http"
	"time"
)

type UnifiC3po struct {
	client http.Client
	//serverC3po	 string
	url string
}

func NewUnifiC3po(url string) *UnifiC3po {
	client := http.Client{
		Timeout: 240 * time.Second,
	}
	return &UnifiC3po{
		client: client,
		url:    url,
	}
}

func (uc3po *UnifiC3po) GetUserLogin(notebook *entity.Client) (err error) {

	type Envelope struct {
		//status string `json:"Status"`
		Status string `json:"Status"`

		Data []struct {
			LineNumber         string `json:"LineNumber"`
			SIPAddress         string `json:"SIPAddress"`
			LineType           string `json:"LineType"`
			RegistrationStatus string `json:"RegistrationStatus"`
			Label              string `json:"Label"`
			UserID             string `json:"UserID"`
			ProxyAddress       string `json:"ProxyAddress"`
			Protocol           string `json:"Protocol"`
			Port               string `json:"Port"`
		} `json:"data"`
	}
	//client := http.Client{Timeout: 5 * time.Second}
	//client := pwa.client
	client := uc3po.client

	//var err error
	myError := 1
	for myError != 0 {
		//url := "http://" + ip + "/api/v1/mgmt/lineInfo"
		//url := "http://" + polyStruct.IP + "/api/v1/mgmt/lineInfo"
		url := "http://" + uc3po.url + "/" + notebook.Hostname
		//log.Println(url)

		req, errNewRequest := http.NewRequest(http.MethodGet, url, http.NoBody)
		if errNewRequest == nil {
			//req.SetBasicAuth(polyLogin, polyPassword)
			//req.SetBasicAuth(pwa.polyUserName, pwa.polyPassword)
			//req.Header.Add("Content-Type", "application/json")

			res, errClientDo := client.Do(req)
			if errClientDo == nil {
				defer res.Body.Close()
				//body = string(resBody)
				//statusHttp = res.StatusCode  //200
				//statuses[0] = res.StatusCode
				envelope := &Envelope{}

				//if errDecode := json.NewDecoder(res.Body).Decode(&envelope); errDecode == nil {
				if errDecode := json.NewDecoder(res.Body).Decode(envelope); errDecode == nil {
					//statusPoly = envelope.status
					//statuses[1] = envelope.status
					//status = envelope.Status
					if envelope.Status == "2000" { //или 5000 может возвращать
						//log.Println("Запрос статуса прошёл.")
						if len(envelope.Data) > 0 {
							//Успешное выполнение функции
							//polyStruct.Status = envelope.Data[0].RegistrationStatus //Registered - правильный ответ
							myError = 0
							return nil //polyStruct, nil
						} else {
							log.Println("Получен ответ 2000 от устройства, но тело ответа пустое")
							log.Println("Будет предпринята новая попытка отправки запроса через 30 сек.")
							time.Sleep(30 * time.Second)
							myError++
							err = errors.New("получен ответ 2000, но тело ответа пустое")
						}
					} else {
						//log.Println(status)
						//log.Println(polyStruct.Status)
						log.Println("От устройства получен Статус НЕ 2000")
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
			log.Println("После 3 неудачных попыток идём дальше. Получить статус работы skype не удалось")
			//status = ""
			//return "", err
			return err //polyStruct, err
		}
	}
	//return status, nil
	return nil //polyStruct, nil
}
