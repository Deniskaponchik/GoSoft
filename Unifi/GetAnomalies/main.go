package main

import (
	"fmt"
	"github.com/unpoller/unifi"
	"log"
	"time"
)

func main567567657657() {
	//c := *unifi.Config{
	c := unifi.Config{
		User: "unifi",
		Pass: "FORCEpower23",
		//URL:  "https://localhost:8443/"
		//URL:  "https://10.78.221.142:8443/", //ROSTOV
		URL: "https://10.8.176.8:8443/", //NOVOSIB
		// Log with log.Printf or make your own interface that accepts (msg, test_SOAP)
		ErrorLog: log.Printf,
		DebugLog: log.Printf,
	}

	clientMacName := map[string]string{}  // clientMAC  -> clientName
	apMacName := map[string]string{}      // apMac      -> apName
	namesClientAps := map[string]string{} // clientName -> apName

	for true { //зацикливаем
		//uni, err := unifi.NewUnifi(c)
		uni, err := unifi.NewUnifi(&c) //в аргументах функций обычно всегда используется &. вставляем переменную из этой функции
		if err != nil {
			log.Fatalln("Error:", err)
		}
		//
		sites, err := uni.GetSites()
		if err != nil {
			log.Fatalln("Error:", err)
		}
		log.Println(len(sites), "Unifi Sites Found: ", sites)

		//
		//ORIGINAL
		devices, err := uni.GetDevices(sites) //devices = APs
		if err != nil {
			log.Fatalln("Error:", err)
		}
		/* ORIGINAL
		log.Println(len(devices.UAPs), "Unifi Wireless APs Found:")
		for i, uap := range devices.UAPs {
			log.Println(i+1, uap.Name, uap.IP, uap.Mac)
		}*/
		// Добавляем маки и имена точек в map
		for _, uap := range devices.UAPs {
			_, existence := apMacName[uap.Mac] //проверяем, есть ли мак в мапе
			if !existence {
				apMacName[uap.Mac] = uap.Name
			}
		}
		/*Вывести AP мапу на экран
		for k, v := range apMacName {
			//fmt.Printf("key: %d, value: %t\n", k, v)
			fmt.Println(k, v)
		}*/

		//
		//ORIGINAL
		clients, err := uni.GetClients(sites)
		if err != nil {
			log.Fatalln("Error:", err)
		}
		/* ORIGINAL
		log.Println(len(clients), "Clients connected:")
		for i, client := range clients {
			log.Println(i+1, client.SiteName, client.IsGuest.Val, client.Mac, client.Hostname, client.IP, client.LastSeen, client.Anomalies) //i+1
		}*/
		for _, client := range clients {
			if !client.IsGuest.Val {
				//Вывод на экран
				//siteName := client.SiteName[:len(client.SiteName)-11]
				apHostName := apMacName[client.ApMac]
				//fmt.Println(siteName, apHostName, client.Hostname, client.Mac, client.IP)

				//Обновление мапы
				clientMacName[client.Mac] = client.Hostname //Добавить КОРП клиентов в map
				namesClientAps[client.Name] = apHostName    //Добавить Соответсвие имён клиентов и точек
			}
		}
		/*Вывести CLIENT мапу на экран
		for k, v := range clientMacName {
			//fmt.Printf("key: %d, value: %t\n", k, v)
			fmt.Println(k, v)
		}*/
		/*Вывести соответсвие имён клиентов и имён точек на экран
		for k, v := range namesClientAps {
			//fmt.Printf("key: %d, value: %t\n", k, v)
			fmt.Println(k, v)
		}*/

		// Если время НЕ 1 минута от начала часа
		if time.Now().Minute() == 47 {
			now := time.Now()
			count := 61 //минус 70 минут
			then := now.Add(time.Duration(-count) * time.Minute)
			//ORIGINAL
			anomalies, err := uni.GetAnomalies(sites,
				//time.Date(2023, 07, 11, 7, 0, 0, 0, time.Local), time.Now()
				then,
			)
			if err != nil {
				log.Fatalln("Error:", err)
			}
			/* ORIGINAL
			log.Println(len(anomalies), "Anomalies:")
			for i, anomaly := range anomalies {
				log.Println(i+1, anomaly.Datetime, anomaly.DeviceMAC, anomaly.Anomaly) //i+1
			}*/
			//bpmTickets := []BpmTicket{}
			bpmTickets := map[string]BpmTicket{} //https://stackoverflow.com/questions/42716852/how-to-update-map-values-in-go
			//
			for _, anomaly := range anomalies {
				_, existence := clientMacName[anomaly.DeviceMAC] //проверяем, есть ли мак в мапе corp clients
				//fmt.Println("Аномалии Tele2Corp клиентов:")
				if existence { //блок кода для Tele2Corp
					//если есть, выводим на экран с именем ПК, взятым из мапы
					siteName := anomaly.SiteName[:len(anomaly.SiteName)-11]
					clientHostName := clientMacName[anomaly.DeviceMAC]
					apHostName := namesClientAps[clientHostName]
					//usrLogin := GetLogin(clientHostName)
					//fmt.Println(siteName, clientHostName, usrLogin, apHostName, anomaly.Datetime, anomaly.Anomaly)
					fmt.Println(siteName, clientHostName, apHostName, anomaly.Datetime, anomaly.Anomaly) //без usrLogin

					_, exisClHostName := bpmTickets[clientHostName] //проверяем, есть ли client hostname в мапе тикетов
					if !exisClHostName {                            //если нет, создаём
						bpmTickets[clientHostName] = BpmTicket{ //https://stackoverflow.com/questions/42716852/how-to-update-map-values-in-go
							//site:
							siteName,
							//apName:
							apHostName,
							clientHostName,
							//	usrLogin,
							[]string{anomaly.Anomaly},
							//"за последний час у пользователя возникли следующие аномалии на Wi-Fi сети Tele2Corp:",
							//"",
						}
					} else {
						for k, v := range bpmTickets {
							if k == clientHostName {
								//https://stackoverflow.com/questions/42716852/how-to-update-map-values-in-go
								/*1.Using pointers. не смог победить указатели...
								v2 := v
								v2.corpAnomalies = append(v2.corpAnomalies, anomaly.Anomaly)
								bpmTickets[k] = v2 */

								//2.Reassigning the modified struct
								v.corpAnomalies = append(v.corpAnomalies, anomaly.Anomaly)
								bpmTickets[k] = v
							}
						}

					}
				} else {
					//Обработка аномалий для Tele2Guest.
					//Пока просто шапка
				}
			}

			fmt.Println("Tele2Corp клиенты с больше чем 1 аномалией:")
			for _, v := range bpmTickets {
				if len(v.corpAnomalies) > 2 {
					fmt.Println(v.clientName)
					for _, s := range v.corpAnomalies {
						fmt.Println(s)
					}
					//SoapCreateTicket(clientHostName, v.clientName, v.corpAnomalies, siteName)
					usrLogin := GetLogin(v.clientName)
					SoapCreateTicket(usrLogin, v.clientName, v.corpAnomalies, v.apName, v.site)
					fmt.Println("")
				}
			}
			//fmt.Println(bpmTickets)
			//jsonStr, err := json.Marshal(bpmTickets)
			//fmt.Println(string(jsonStr))

		} //else

		time.Sleep(60 * time.Second) //Ставим на паузу на 1 минуту
	} // Do while
} //main func

/*
func GetClientsCorpWithAnomalies(anoms []*Anomaly) ([]*ClientCorp) {
	return
}*/

type BpmTicket struct { //структура ДОЛЖНА находиться ВНЕ main
	site       string
	apName     string
	clientName string
	//userLogin     string
	corpAnomalies []string
	//description    string
	//recomendations string
}
