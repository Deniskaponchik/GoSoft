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
		Name   string `json:"name"`
		Status string `json:"status"`

		Data []struct {
			Monitor1_Model   string `json:"Monitor1_Model"`
			Monitor1_SN      string `json:"Monitor1_SN"`
			Monitor1_Vendor  string `json:"Monitor1_Vendor"`
			Monitor2_Model   string `json:"Monitor2_Model"`
			Monitor2_SN      string `json:"Monitor2_SN"`
			Monitor2_Vendor  string `json:"Monitor2_Vendor"`
			OZU              int    `json:"OZU"`
			disk1_model      string `json:"disk1_model"`
			disk2_model      string `json:"disk2_model"`
			disk3_model      string `json:"disk3_model"`
			last_scan        string `json:"last_scan"`
			os_build         string `json:"os_build"`
			os_version       string `json:"os_version"`
			pc_cpu           string `json:"pc_cpu"`
			pc_manufacturer  string `json:"pc_manufacturer"`
			pc_model         string `json:"pc_model"`
			pc_name          string `json:"pc_name"`
			pc_serial_number string `json:"pc_serial_number"`
			ram1             int    `json:"ram1"`
			ram2             int    `json:"ram2"`
			ram3             int    `json:"ram3"`
			ram4             int    `json:"ram4"`
			samaccountname   string `json:"samaccountname"`
		} `json:"data"`
	}
	//client := http.Client{Timeout: 5 * time.Second}
	client := uc3po.client

	//var err error
	myError := 1
	for myError != 0 {
		//url := "http://" + polyStruct.IP + "/api/v1/mgmt/lineInfo"
		url := uc3po.url + "pc/" + notebook.Hostname
		log.Println(url)

		req, errNewRequest := http.NewRequest(http.MethodGet, url, http.NoBody)
		if errNewRequest == nil {
			//req.SetBasicAuth(polyLogin, polyPassword)
			//req.Header.Add("Content-Type", "application/json")

			res, errClientDo := client.Do(req)
			if errClientDo == nil {
				defer res.Body.Close()
				//body = string(resBody)
				//statusHttp = res.StatusCode
				//statuses[0] = res.StatusCode
				envelope := &Envelope{}

				//if errDecode := json.NewDecoder(res.Body).Decode(&envelope); errDecode == nil {
				if errDecode := json.NewDecoder(res.Body).Decode(envelope); errDecode == nil {
					if envelope.Status == "ok" {
						//log.Println("Запрос статуса прошёл.")
						if len(envelope.Data) > 0 {
							//Успешное выполнение функции
							notebook.UserLogin = envelope.Data[0].samaccountname
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
			return err
		}
	}
	return nil
}
