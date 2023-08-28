package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

//func main() {}

func apiLineInfo(ip string) (status string) {

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
	client := http.Client{Timeout: 5 * time.Second}

	myError := 1
	for myError != 0 {
		url := "http://" + ip + "/api/v1/mgmt/lineInfo"
		//fmt.Println(url)

		req, errNewRequest := http.NewRequest(http.MethodGet, url, http.NoBody)
		if errNewRequest == nil {
			req.SetBasicAuth("Polycom", "3214")
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
						}
					} else {
						fmt.Println(status)
						fmt.Println("От устройства получен Статус НЕ 2000")
						fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
						time.Sleep(30 * time.Second)
						myError++
					}
				} else {
					fmt.Println(errDecode.Error())
					fmt.Println("Ошибка перекодировки ответа")
					fmt.Println("Скорее всего, API недоступен")
					fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
					time.Sleep(30 * time.Second)
					myError++
				}
			} else {
				fmt.Println(errClientDo.Error())
				fmt.Println("Ошибка отправки запроса")
				fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
				time.Sleep(30 * time.Second)
				myError++
			}
		} else {
			fmt.Println(errNewRequest.Error())
			fmt.Println("Ошибка создания ОБЪЕКТА запроса")
			fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
			time.Sleep(30 * time.Second)
			myError++
		}
		if myError == 4 {
			myError = 0
			fmt.Println("После 3 неудачных попыток идём дальше. Получить статус работы skype не удалось")
			status = ""
			//statuses = append(statuses, 0)
			//statuses = append(statuses, 0)
		}
	}
	//fmt.Printf("Status: %d\n", res.StatusCode)
	//fmt.Printf("Body: %s\n", string(resBody))
	return status
}

func apiSafeRestart2(ip string) (status string) {

	type Envelope struct {
		//status string `json:"Status"`
		Status string `json:"Status"`
	}
	client := http.Client{Timeout: 5 * time.Second}

	myError := 1
	for myError != 0 {
		url := "http://" + ip + "/api/v1/mgmt/safeRestart"
		//url := ip + "/api/v1/mgmt/safeRestart"
		req, errNewRequest := http.NewRequest(http.MethodPost, url, http.NoBody)
		if errNewRequest == nil {
			req.SetBasicAuth("Polycom", "3214")
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
					status = envelope.Status
					if status == "2000" {
						fmt.Println("Запрос на перезагрузку прошёл успешно. Устройство перезагрузится в течение 5 минут")
						myError = 0
					} else {
						fmt.Println(status)
						fmt.Println("От устройства получен Статус НЕ 2000")
						fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
						time.Sleep(30 * time.Second)
						myError++
					}
				} else {
					fmt.Println(errDecode.Error())
					fmt.Println("Ошибка перекодировки ответа")
					fmt.Println("Скорее всего, API недоступен")
					fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
					time.Sleep(30 * time.Second)
					myError++
				}
			} else {
				fmt.Println(errClientDo.Error())
				fmt.Println("Ошибка отправки запроса")
				fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
				time.Sleep(30 * time.Second)
				myError++
			}
		} else {
			fmt.Println(errNewRequest.Error())
			fmt.Println("Ошибка создания ОБЪЕКТА запроса")
			fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
			time.Sleep(30 * time.Second)
			myError++
		}
		if myError == 4 {
			myError = 0
			fmt.Println("После 3 неудачных попыток идём дальше. Перезагрузка не была осуществлена")
			status = ""
			//statuses = append(statuses, 0)
			//statuses = append(statuses, 0)
		}
	}
	//fmt.Printf("Status: %d\n", res.StatusCode)
	//fmt.Printf("Body: %s\n", string(resBody))
	return status
}
