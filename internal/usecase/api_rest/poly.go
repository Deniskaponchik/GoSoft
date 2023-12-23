package api_rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/deniskaponchik/GoSoft/internal/entity"
	"net/http"
	"time"
)

type PolyWebAPI struct {
	client       http.Client
	polyUserName string
	polyPassword string
}

func NewPolyWebApi(u string, p string) *PolyWebAPI {
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	return &PolyWebAPI{
		client:       client,
		polyUserName: u, //cfg.PolyUsername
		polyPassword: p, //cfg.PolyPassword
	}
}

func (pwa *PolyWebAPI) ApiLineInfoErr(polyStruct *entity.PolyStruct) (err error) {
	//https://www.poly.com/content/dam/www/products/support/voice/trio/other/rest-api-ref-trio-5-9-5-en.pdf

	type Envelope struct {
		//status string `json:"Status"`
		Status string `json:"Status"`
		Data   []struct {
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
	client := pwa.client

	//var err error
	myError := 1
	for myError != 0 {
		//url := "http://" + ip + "/api/v1/mgmt/lineInfo"
		url := "http://" + polyStruct.IP + "/api/v1/mgmt/lineInfo"
		//fmt.Println(url)

		req, errNewRequest := http.NewRequest(http.MethodGet, url, http.NoBody)
		if errNewRequest == nil {
			//req.SetBasicAuth(polyLogin, polyPassword)
			req.SetBasicAuth(pwa.polyUserName, pwa.polyPassword)
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
						//fmt.Println("Запрос статуса прошёл.")
						if len(envelope.Data) > 0 {
							//Успешное выполнение функции
							//fmt.Println("Получен статус skype")
							polyStruct.Status = envelope.Data[0].RegistrationStatus //Registered - правильный ответ
							myError = 0
							return nil //polyStruct, nil
						} else {
							fmt.Println("Получен ответ 2000 от устройства, но тело ответа пустое")
							fmt.Println("Будет предпринята новая попытка отправки запроса через 30 сек.")
							time.Sleep(30 * time.Second)
							myError++
							err = errors.New("получен ответ 2000, но тело ответа пустое")
						}
					} else {
						//fmt.Println(status)
						fmt.Println(polyStruct.Status)
						fmt.Println("От устройства получен Статус НЕ 2000")
						fmt.Println("Будет предпринята новая попытка отправки запроса через 30 сек.")
						time.Sleep(30 * time.Second)
						myError++
						err = errors.New("от устройства получен статус не 2000")
					}
				} else {
					fmt.Println(errDecode.Error())
					fmt.Println("Ошибка перекодировки ответа")
					fmt.Println("Скорее всего, API недоступен")
					fmt.Println("Будет предпринята новая попытка отправки запроса через 30 сек.")
					time.Sleep(30 * time.Second)
					myError++
					err = errDecode
				}
			} else {
				fmt.Println(errClientDo.Error())
				fmt.Println("Ошибка отправки запроса")
				fmt.Println("Будет предпринята новая попытка отправки запроса через 30 сек.")
				time.Sleep(30 * time.Second)
				myError++
				err = errClientDo
			}
		} else {
			fmt.Println(errNewRequest.Error())
			fmt.Println("Ошибка создания ОБЪЕКТА запроса")
			fmt.Println("Будет предпринята новая попытка отправки запроса через 30 сек.")
			time.Sleep(30 * time.Second)
			myError++
			err = errNewRequest
		}
		if myError == 4 {
			myError = 0
			fmt.Println("После 3 неудачных попыток идём дальше. Получить статус работы skype не удалось")
			//status = ""
			//return "", err
			return err //polyStruct, err
		}
	}
	//return status, nil
	return nil //polyStruct, nil
}

func (pwa *PolyWebAPI) ApiLineInfo(polyStruct entity.PolyStruct) (status string, err error) {
	//Метод возвращает строку

	type Envelope struct {
		//status string `json:"Status"`
		Status string `json:"Status"`
		Data   []struct {
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
	client := pwa.client

	myError := 1
	for myError != 0 {
		//url := "http://" + ip + "/api/v1/mgmt/lineInfo"
		url := "http://" + polyStruct.IP + "/api/v1/mgmt/lineInfo"
		//fmt.Println(url)

		req, errNewRequest := http.NewRequest(http.MethodGet, url, http.NoBody)
		if errNewRequest == nil {
			//req.SetBasicAuth(polyLogin, polyPassword)
			req.SetBasicAuth(pwa.polyUserName, pwa.polyPassword)
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
					if envelope.Status == "2000" {
						//fmt.Println("Запрос статуса прошёл.")
						if len(envelope.Data) > 0 {
							//fmt.Println("Получен статус skype")
							status = envelope.Data[0].RegistrationStatus
							myError = 0
						} else {
							fmt.Println("Получен ответ 2000 от устройства, но тело ответа пустое")
							fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
							time.Sleep(30 * time.Second)
							myError++
							err = errors.New("получен ответ 2000 от устройства, но тело ответа пустое")
						}
					} else {
						fmt.Println(status)
						fmt.Println("От устройства получен Статус НЕ 2000")
						fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
						time.Sleep(30 * time.Second)
						myError++
						err = errors.New("от устройства получен статус не 2000")
					}
				} else {
					fmt.Println(errDecode.Error())
					fmt.Println("Ошибка перекодировки ответа")
					fmt.Println("Скорее всего, API недоступен")
					fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
					time.Sleep(30 * time.Second)
					myError++
					err = errDecode
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
			fmt.Println(errNewRequest.Error())
			fmt.Println("Ошибка создания ОБЪЕКТА запроса")
			fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
			time.Sleep(30 * time.Second)
			myError++
			err = errNewRequest
		}
		if myError == 4 {
			myError = 0
			fmt.Println("После 3 неудачных попыток идём дальше. Получить статус работы skype не удалось")
			//status = ""
			return "", err
		}
	}
	//fmt.Printf("Status: %d\n", res.StatusCode)
	//fmt.Printf("Body: %s\n", string(resBody))
	return status, nil
}

func (pwa *PolyWebAPI) ApiSafeRestart(polyStruct entity.PolyStruct) (err error) {

	type Envelope struct {
		//status string `json:"Status"`
		Status string `json:"Status"`
	}
	//client := http.Client{Timeout: 5 * time.Second}
	client := pwa.client
	ip := polyStruct.IP //работает только в таком виде

	myError := 1
	for myError != 0 {
		//url := "http://10.21.178.78/api/v1/mgmt/safeRestart"
		url := "http://" + ip + "/api/v1/mgmt/safeRestart" //работает только в таком виде
		//url := "http://" + polyStruct.IP + "/api/v1/mgmt/lineInfo" //в таком виде, почему-то не работает
		fmt.Println(url)

		req, errNewRequest := http.NewRequest(http.MethodPost, url, http.NoBody)
		if errNewRequest == nil {
			//req.SetBasicAuth(polyLogin, polyPassword)
			req.SetBasicAuth(pwa.polyUserName, pwa.polyPassword)
			req.Header.Add("Content-Type", "application/json")

			res, errClientDo := client.Do(req)
			if errClientDo == nil {
				defer res.Body.Close()

				//body = string(resBody)
				//statusHttp = res.StatusCode
				//statuses[0] = res.StatusCode
				envelope := &Envelope{}

				//if errDecode := json.NewDecoder(res.Body).Decode(&envelope); errDecode == nil {
				if errDecode := json.NewDecoder(res.Body).Decode(envelope); errDecode == nil {
					//log.Fatal("ooopsss! an error occurred, please try again")
					//statusPoly = envelope.status
					//statuses[1] = envelope.status
					status := envelope.Status
					if status == "2000" {
						fmt.Println("Запрос на перезагрузку прошёл успешно. Устройство перезагрузится в течение 5 минут")
						myError = 0
						return nil
					} else {
						fmt.Println(status)
						fmt.Println("От устройства получен Статус НЕ 2000")
						fmt.Println("Будет предпринята новая попытка отправки запроса через 30 сек.")
						time.Sleep(30 * time.Second)
						myError++
						err = errors.New("от устройства получен статус не 2000")
					}
				} else {
					fmt.Println(errDecode.Error())
					fmt.Println("Ошибка перекодировки ответа")
					fmt.Println("Скорее всего, API недоступен")
					fmt.Println("Будет предпринята новая попытка отправки запроса через 30 сек.")
					time.Sleep(30 * time.Second)
					myError++
					err = errDecode
				}
			} else {
				fmt.Println(errClientDo.Error())
				fmt.Println("Ошибка отправки запроса")
				fmt.Println("Будет предпринята новая попытка отправки запроса через 30 сек.")
				time.Sleep(30 * time.Second)
				myError++
				err = errClientDo
			}
		} else {
			fmt.Println(errNewRequest.Error())
			fmt.Println("Ошибка создания ОБЪЕКТА запроса")
			fmt.Println("Будет предпринята новая попытка отправки запроса через 30 сек.")
			time.Sleep(30 * time.Second)
			myError++
			err = errNewRequest
		}
		if myError == 4 {
			myError = 0
			fmt.Println("После 3 неудачных попыток идём дальше. Перезагрузка не была осуществлена")
			//status = ""
			return err
		}
	}
	//fmt.Printf("Status: %d\n", res.StatusCode)
	//fmt.Printf("Body: %s\n", string(resBody))
	return nil
}
