package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

//https://blog.logrocket.com/making-http-requests-in-go/

func main() {
	//ip := ""
	//getPass()
	statuses := safeRestart2("10.57.178.42")
	fmt.Println(statuses)
	//fmt.Println(statuses[1])
}

func getPass() {
	client := http.Client{Timeout: 5 * time.Second}

	req, err := http.NewRequest(http.MethodGet, "http://10.57.178.41/api/v1/mgmt/network/info", http.NoBody)
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth("Polycom", "3214")

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Status: %d\n", res.StatusCode)
	fmt.Printf("Body: %s\n", string(resBody))
}

// func safeRestart2(ip string) (statusHttp int, statusPoly int) {
func safeRestart2(ip string) (status string) {

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
						myError = 0
					} else {
						fmt.Println(status)
						fmt.Println("От устройства получен Статус НЕ 2000")
						fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
						time.Sleep(60 * time.Second)
						myError++
					}
				} else {
					fmt.Println(errDecode.Error())
					fmt.Println("Ошибка перекодировки ответа")
					fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
					time.Sleep(60 * time.Second)
					myError++
				}
			} else {
				fmt.Println(errClientDo.Error())
				fmt.Println("Ошибка отправки запроса")
				fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
				time.Sleep(60 * time.Second)
				myError++
			}
		} else {
			fmt.Println(errNewRequest.Error())
			fmt.Println("Ошибка создания ОБЪЕКТА запроса")
			fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
			time.Sleep(60 * time.Second)
			myError++
		}
		if myError == 6 {
			myError = 0
			fmt.Println("После 6 неудачных попыток идём дальше. Заявка не была создана")
			//statuses = append(statuses, 0)
			//statuses = append(statuses, 0)
		}
	}
	//fmt.Printf("Status: %d\n", res.StatusCode)
	//fmt.Printf("Body: %s\n", string(resBody))
	return status
}

func safeRestart1(ip string) (statuses []int) {

	type Envelope struct {
		status int `json:"Status"`
	}
	url := "http://" + ip + "/api/v1/mgmt/safeRestart"
	client := http.Client{Timeout: 5 * time.Second}

	myError := 1
	for myError != 0 {
		req, errNewRequest := http.NewRequest(http.MethodPost, url, http.NoBody)
		if errNewRequest == nil {
			req.SetBasicAuth("Polycom", "3214")
			req.Header.Add("Content-Type", "application/json")
			//req.Header.Set("Content-Type", "application/json")

			res, errClientDo := client.Do(req)
			if errClientDo == nil {
				defer res.Body.Close()

				resBody, errIOread := io.ReadAll(res.Body)
				if errIOread == nil {
					statuses[0] = res.StatusCode
					//body = string(resBody)
					envelope := &Envelope{}

					/*if err := json.NewDecoder(resBody).Decode(&cResp); err != nil {
					log.Fatal("ooopsss! an error occurred, please try again")	}*/
					errJsonUnmarshal := json.Unmarshal(resBody, envelope)
					if errJsonUnmarshal == nil {
						statuses[1] = envelope.status
						myError = 0
					} else {
						fmt.Println(errJsonUnmarshal.Error())
						fmt.Println("Ошибка перекодировки ответа в json")
						fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
						time.Sleep(60 * time.Second)
						myError++
					}
				} else {
					fmt.Println(errIOread.Error())
					fmt.Println("Ошибка чтения байтов из ответа")
					fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
					time.Sleep(60 * time.Second)
					myError++
				}
			} else {
				fmt.Println(errClientDo.Error())
				fmt.Println("Ошибка отправки запроса")
				fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
				time.Sleep(60 * time.Second)
				myError++
			}
		} else {
			fmt.Println(errNewRequest.Error())
			fmt.Println("Ошибка создания ОБЪЕКТА запроса")
			fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
			time.Sleep(60 * time.Second)
			myError++
		}
		if myError == 6 {
			myError = 0
			fmt.Println("После 6 неудачных попыток идём дальше. Заявка не была создана")
			statuses = append(statuses, 0)
			statuses = append(statuses, 0)
		}
	}

	//fmt.Printf("Status: %d\n", res.StatusCode)
	//fmt.Printf("Body: %s\n", string(resBody))
	return statuses
}

func post() {
	postBody, _ := json.Marshal(map[string]string{
		"name":  "Toby",
		"email": "Toby@example.com",
	})
	responseBody := bytes.NewBuffer(postBody)

	//Leverage Go's HTTP Post function to make request
	resp, err := http.Post("https://postman-echo.com/post", "application/json", responseBody)
	//Handle Error
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}

	defer resp.Body.Close()
	//Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	sb := string(body)
	log.Printf(sb)
}
