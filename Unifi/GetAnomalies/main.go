package main

import (
	"fmt"
	"github.com/unpoller/unifi"
	"io"
	"log"
	"os"
	"time"
)

func main() {

	//MAPs
	//clientHostName - > userLogin
	noutnameLogin := DownloadMapFromDB("glpi_db", "name", "contact", "glpi_db.glpi_computers", "date_mod")
	clientMacName := map[string]string{}  // clientMAC  -> clientHostName
	apMacName := map[string]string{}      // apMac      -> apName
	namesClientAps := map[string]string{} // clientName -> apName
	//clientnameTicketid
	//apnameTicketid

	for k, v := range noutnameLogin {
		//fmt.Printf("key: %d, value: %t\n", k, v)
		fmt.Println("newMap "+k, v)
	}
	os.Exit(0)

	countMinute := 0
	count5minute := 5
	countHourAnom := 0
	countHourDB := 0

	//c := *unifi.Config{  //ORIGINAL
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

	log.SetOutput(io.Discard) //Отключить вывод лога

	//uni, err := unifi.NewUnifi(c)
	uni, err := unifi.NewUnifi(&c)
	if err != nil {
		log.Fatalln("Error:", err)
	}

	for true { //зацикливаем навечно

		//Снятие показаний с контрллера каждую МИНУТУ
		if time.Now().Minute() != countMinute { //Блок кода запустится, если в эту минуту он ещё НЕ выполнялся
			sites, err := uni.GetSites()
			if err != nil {
				log.Fatalln("Error:", err)
			}
			//log.Println(len(sites), "Unifi Sites Found: ", sites)

			devices, err := uni.GetDevices(sites) //devices = APs
			if err != nil {
				log.Fatalln("Error:", err)
			}
			/* ORIGINAL
			log.Println(len(devices.UAPs), "Unifi Wireless APs Found:")
			for i, uap := range devices.UAPs {
				log.Println(i+1, uap.Name, uap.IP, uap.Mac)
			}*/
			// Добавляем маки и имена точек в apMacName map
			for _, uap := range devices.UAPs {
				_, existence := apMacName[uap.Mac] //проверяем, есть ли мак в мапе
				if !existence {
					apMacName[uap.Mac] = uap.Name
				}
			}
			/*Вывести apMacName мапу на экран
			for k, v := range apMacName {
				//fmt.Printf("key: %d, value: %t\n", k, v)
				fmt.Println(k, v)
			}*/

			clients, err := uni.GetClients(sites) //clients = Notebooks, Mobiles
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
					//siteName := client.SiteName[:len(client.SiteName)-11]
					apName := apMacName[client.ApMac]
					//fmt.Println(siteName, apName, client.Hostname, client.Mac, client.IP)

					//Обновление мапы clientMAC-clientHOST
					clientMacName[client.Mac] = client.Hostname //Добавить КОРП клиентов в map
					namesClientAps[client.Name] = apName        //Добавить Соответсвие имён клиентов и точек
				}
			}
			/*Вывести clientMacName мапу на экран
			for k, v := range clientMacName {
				//fmt.Printf("key: %d, value: %t\n", k, v)
				fmt.Println(k, v)
			}*/
			/*Вывести соответсвие имён клиентов и имён точек на экран
			for k, v := range namesClientAps {
				//fmt.Printf("key: %d, value: %t\n", k, v)
				fmt.Println(k, v)
			}*/

			countMinute = time.Now().Minute()

			// Тут должен быть блок кода про точки, а уже потом аномалии
			if time.Now().Minute()%5 == 0 && time.Now().Minute() != count5minute { //запускается раз в 5 минут

				count5minute = time.Now().Minute()
			}

			//АНОМАЛИИ
			//if time.Now().Minute() == 47 { // Если время 3 минуты от начала часа то блок для аномаоий
			if time.Now().Hour() != countHourAnom { //Блок кода запустится, если в этот ЧАС он ещё НЕ выполнялся
				now := time.Now()
				count := 60 //минус 70 минут
				then := now.Add(time.Duration(-count) * time.Minute)

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

				type ForTicket struct {
					site       string
					apName     string
					clientName string
					//userLogin     string //не помещаю, чтобы не делать лишних запросов к БД, если заявка НЕ будет создаваться
					corpAnomalies []string
				}

				//mapNoutnameFortickets создаётся локально в блоке аномалий каждый час. Резервировать в БД НЕ нужно
				mapNoutnameFortickets := map[string]ForTicket{} //https://stackoverflow.com/questions/42716852/how-to-update-map-values-in-go
				//
				for _, anomaly := range anomalies {
					_, existence := clientMacName[anomaly.DeviceMAC] //проверяем, соответсвует ли мак мапе corp клиентов
					//fmt.Println("Аномалии Tele2Corp клиентов:")
					if existence { //блок кода для Tele2Corp
						//если есть, выводим на экран с именем ПК, взятым из мапы
						siteName := anomaly.SiteName[:len(anomaly.SiteName)-11]
						clientHostName := clientMacName[anomaly.DeviceMAC]
						apName := namesClientAps[clientHostName]
						//usrLogin := GetLogin(clientHostName)
						//fmt.Println(siteName, clientHostName, usrLogin, apName, anomaly.Datetime, anomaly.Anomaly)
						fmt.Println(siteName, clientHostName, apName, anomaly.Datetime, anomaly.Anomaly) //без usrLogin

						_, exisClHostName := mapNoutnameFortickets[clientHostName] //проверяем, есть ли client hostname в мапе тикетов
						if !exisClHostName {                                       //если нет, добавляем новый
							mapNoutnameFortickets[clientHostName] = ForTicket{ //https://stackoverflow.com/questions/42716852/how-to-update-map-values-in-go
								//site:
								siteName,
								//apName:
								apName,
								clientHostName,
								//	usrLogin,
								[]string{anomaly.Anomaly},
								//"за последний час у пользователя возникли следующие аномалии на Wi-Fi сети Tele2Corp:",
								//"",
							}
						} else {
							for k, v := range mapNoutnameFortickets {
								if k == clientHostName {
									//https://stackoverflow.com/questions/42716852/how-to-update-map-values-in-go
									/*1.Using pointers. не смог победить указатели...
									v2 := v
									v2.corpAnomalies = append(v2.corpAnomalies, anomaly.Anomaly)
									mapNoutnameFortickets[k] = v2 */

									//2.Reassigning the modified struct
									v.corpAnomalies = append(v.corpAnomalies, anomaly.Anomaly)
									mapNoutnameFortickets[k] = v
								}
							}

						}
					} else {
						//Обработка аномалий для Tele2Guest.
						//Пока просто шапка
					}
				}

				fmt.Println("")
				fmt.Println("Tele2Corp клиенты с более чем 2 аномалиями:")
				for _, v := range mapNoutnameFortickets {
					if len(v.corpAnomalies) > 2 {
						fmt.Println(v.clientName)
						for _, s := range v.corpAnomalies {
							fmt.Println(s)
						}
						//SoapCreateTicket(clientHostName, v.clientName, v.corpAnomalies, siteName)
						usrLogin := GetLogin(v.clientName)

						//1. Проверяет, есть ли заявка в мапе ClientHostName - ID Тикета
						//2. Если заявка В МАПЕ есть, проверить её статус
						//3. Если Статус закрыто, решено, завести новую
						srTicketSlice := CreateSmacWiFiTicket(usrLogin, v.clientName, v.corpAnomalies, v.apName, v.site)
						fmt.Println(srTicketSlice[2])
						fmt.Println("")
						//после создания заявки добавить в мапу ClientHostName - ID Тикета

					}
				}
				//fmt.Println(mapNoutnameFortickets)
				//jsonStr, err := json.Marshal(mapNoutnameFortickets)
				//fmt.Println(string(jsonStr))

				countHourAnom = time.Now().Hour()
			} // END of ANOMALIES block

			//Занесение мап в БД
			if time.Now().Hour() != countHourDB { //Блок кода запустится, если в этот ЧАС он ещё НЕ выполнялся

				countHourDB = time.Now().Hour()
			}
		} // Поминутный if
		time.Sleep(60 * time.Second) //Изменить на 5 секунд на ПРОДе
	} // while TRUE
} //main func
