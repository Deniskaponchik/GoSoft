package main

import (
	"fmt"
	"github.com/unpoller/unifi"
	"io"
	"log"
	"strings"
	"time"
)

func main() {
	//bpmServer := "http://10.12.15.148/specs/aoi/tele2/bpm/bpmPortType"   //PROD
	bpmServer := "http://10.246.37.15:8060/specs/aoi/tele2/bpm/bpmPortType" //TEST
	//unifiController :=  "https://localhost:8443/"
	//unifiController := "https://10.78.221.142:8443/"  //Rostov
	unifiController := "https://10.8.176.8:8443/" //Novosib

	countMinute := 0
	count3minute := 0
	//count5minute := 5
	countHourAnom := 0
	countHourDB := 0
	countDay := time.Now().Day()

	//Download MAPs from DB
	//noutnameLogin :=map[string]string{}     //clientHostName - > userLogin
	noutnameLogin := DownloadMapFromDB("glpi_db", "name", "contact", "glpi_db.glpi_computers", "date_mod")
	siteApCutNameLogin := DownloadMapFromDB("wifi_db", "site_apcut", "login", "wifi_db.site_apcut_login", "site_apcut")
	//maschineMacName := map[string]string{}   // clientMAC -> clientHostName  // maschineMAC -> maschineHostName
	maschineMacName := DownloadMapFromDB("wifi_db", "mac", "hostname", "wifi_db.maschine_mac_name", "hostname")
	//apMacName := map[string]string{}      // apMac -> apName
	apMacName := DownloadMapFromDB("wifi_db", "mac", "name", "wifi_db.ap_mac_name", "name")
	//namesClientAps := map[string]string{} // clientName -> apName
	namesClientAp := DownloadMapFromDB("wifi_db", "mascine_name", "ap_name", "wifi_db.names_mascine_ap", "mascine_name")
	//maschineMacSRid := DownloadMapFromDB("wifi_db", "hostname", "srid", "wifi_db.mascine_name_srid", "hostname")
	maschineMacSRid := DownloadMapFromDB("wifi_db", "mac", "srid", "wifi_db.maschine_mac_srid", "mac")
	//apMacSRid := DownloadMapFromDB("wifi_db", "apname", "srid", "wifi_db.ap_name_srid", "apname")
	apMacSRid := DownloadMapFromDB("wifi_db", "mac", "srid", "wifi_db.ap_mac_srid", "mac")
	/*
		for k, v := range apnameSRid {
			//fmt.Printf("key: %d, value: %t\n", k, v)
			fmt.Println("newMap "+k, v)
		}*/
	//os.Exit(0)
	fmt.Println("")

	// Пока без них:
	type clientT2Corp struct {
		mac            string
		hostname       string
		region         string
		apName         string
		anomls         []string
		srNumber       string
		srID           string
		srLink         string
		countAnomalies int  //сколько раз бы создалась заявка, если бы не стояло ограничение в 1 заявку на пользователя
		exception      bool //исключить создание обращений по нему, ели нет. 1 -исключить
	}
	type ap struct {
	}

	c := unifi.Config{
		//c := *unifi.Config{  //ORIGINAL
		User: "unifi",
		Pass: "FORCEpower23",
		//URL:  "https://localhost:8443/"
		//URL:  "https://10.78.221.142:8443/", //ROSTOV
		//URL: "https://10.8.176.8:8443/", //NOVOSIB
		URL: unifiController,
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

		//Снятие показаний с контрллера каждую МИНУТУ. Изменить на 3 минуты на ПРОДе
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

			clients, err := uni.GetClients(sites) //client = Notebook or Mobile = maschine
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
					maschineMacName[client.Mac] = client.Hostname //Добавить КОРП клиентов в map
					namesClientAp[client.Name] = apName           //Добавить Соответсвие имён клиентов и точек
				}
			}
			/*Вывести maschineMacName мапу на экран
			for k, v := range maschineMacName {
				//fmt.Printf("key: %d, value: %t\n", k, v)
				fmt.Println(k, v)
			}*/
			/*Вывести соответсвие имён клиентов и имён точек на экран
			for k, v := range namesClientAps {
				//fmt.Printf("key: %d, value: %t\n", k, v)
				fmt.Println(k, v)
			}*/

			countMinute = time.Now().Minute()

			// блок кода про точки, а уже потом аномалии
			//if time.Now().Minute()%5 == 0 && time.Now().Minute() != count5minute { //запускается раз в 5 минут
			if time.Now().Minute()%3 == 0 && time.Now().Minute() != count3minute { //запускается раз в 5 минут
				fmt.Println("Обработка точек доступа...")
				type ForApsTicket struct {
					site string
					//apName string
					apMac string
					//userLogin     string //не помещаю, чтобы не делать лишних запросов к БД, если заявка НЕ будет создаваться
					apNames []string //сделано для массовых отключений точек при отключении света в офисе
				}
				//создаётся локально в блоке раз в 5 минут. Резервировать в БД НЕ нужно
				siteapNameForTickets := map[string]ForApsTicket{}
				//sliceForTicket := []ForApsTicket{}

				srStatusCodesForNewTicket := map[string]bool{
					"Отменено":     true, //Cancel  6e5f4218-f46b-1410-fe9a-0050ba5d6c38
					"Решено":       true, //Resolve  ae7f411e-f46b-1410-009b-0050ba5d6c38
					"Закрыто":      true, //Closed  3e7f420c-f46b-1410-fc9a-0050ba5d6c38
					"На уточнении": true, //Clarification 81e6a1ee-16c1-4661-953e-dde140624fb
				}

				for _, ap := range devices.UAPs {
					//if ap.SiteName[:len(ap.SiteName)-11] != "Резерв/Склад" {
					if ap.SiteID != "5f2285f3a1a7693ae6139c00" { //NOVOSIB

						fmt.Println(ap.Name)
						fmt.Println(ap.SiteID)
						apLastSeen := ap.LastSeen.Int()
						_, exisApMacSRid := apMacSRid[ap.Mac]

						//Точка доступна. Заявки нет.
						if apLastSeen != 0 && !exisApMacSRid {
							//Идём дальше

							//Точка доступна. Заявка есть
						} else if apLastSeen != 0 && exisApMacSRid {
							//ОЧИЩАЕМ мапу, оставляем коммент, ПЫТАЕМСЯ закрыть тикет, если на визировании
							//оставить комментарий, что точка стала доступна
							comment := "Точка появилась в сети: " + ap.Name
							AddComment(bpmServer, apMacSRid[ap.Mac], comment)

							//удалить запись из мапы, предварительно сохранив Srid
							srID := apMacSRid[ap.Mac]
							delete(apMacSRid, ap.Mac)

							//проверить, не последняя ли это запись была в мапе в массиве
							countOfIncident := 0
							for _, v := range apMacSRid {
								if v == srID {
									countOfIncident++
								}
							}
							if countOfIncident == 0 {
								//Пробуем закрыть тикет, только ЕСЛИ он на Визировании
								sliceTicketStatus := CheckTicketStatus(bpmServer, apMacSRid[ap.Mac]) //получаем статус
								if sliceTicketStatus[1] == "На визировании" {
									//Если статус заявки по-прежнему на Визировании
									ChangeStatus(bpmServer, srID, "На уточнении")
									AddComment(bpmServer, srID, "Обращение отменено, т.к. все точки из него появились в сети")
									ChangeStatus(bpmServer, srID, "Отменено")
								}
							}

							/*Точка Недоступна. Заявки нет
							} else if apLastSeen == 0 && !exisApMacSRid {
								//Заполняем переменные, которые понадобятся дальше
								siteName := ap.SiteName[:len(ap.SiteName)-11]
								//apCutName := ap.Name[:len(ap.Name)3]
								apCutName := strings.Split(ap.Name, "-")[0]
								siteApCutName := siteName + apCutName

								//Проверяем и Вносим во временную мапу. Заявка на данном этапе никакая ещё НЕ создаётся
								_, exisSiteName := siteapNameForTickets[siteName] //проверяем, есть ли siteName в мапе ДЛЯтикетов
								//for _, ticket := range sliceForTicket {

								//если в мапе дляТикета сайта ещё НЕТ
								if !exisSiteName {
									siteapNameForTickets[siteApCutName] = ForApsTicket{
										//siteName,
										//ap.Name,
										ap.Mac,
										[]string{ap.Name},
									}

									//если в мапе дляТикета сайт уже есть, добавляем в массив точку
								} else {
									//в мапе нельзя просто изменить значение.
									for k, v := range siteapNameForTickets {
										if k == siteApCutName {
											//https://stackoverflow.com/questions/42716852/how-to-update-map-values-in-go
											//1.Using pointers. не смог победить указатели...
											//v2 := v
											//v2.corpAnomalies = append(v2.corpAnomalies, anomaly.Anomaly)
											//mapNoutnameFortickets[k] = v2

											//2.Reassigning the modified struct.
											v.apNames = append(v.apNames, ap.Name)
											siteapNameForTickets[k] = v
										}
									}
								}*/

							//Точка недоступна.
						} else if apLastSeen == 0 {
							//Проверяем заявку на НЕ закрытость. если заявки нет - ничего страшного
							checkSlice := CheckTicketStatus(bpmServer, apMacSRid[ap.Mac])
							if srStatusCodesForNewTicket[checkSlice[1]] || !exisApMacSRid {
								delete(apMacSRid, ap.Mac) //удаляем заявку. если заявки нет - ничего страшного

								//Заполняем переменные, которые понадобятся дальше
								fmt.Println(ap.SiteID)
								var siteName string
								if ap.SiteID != "6360a823a1a769286dc707f2" {
									siteName = "Урал"
								} else {
									siteName = ap.SiteName[:len(ap.SiteName)-11]
								}
								//apCutName := ap.Name[:len(ap.Name)3]
								apCutName := strings.Split(ap.Name, "-")[0]
								siteApCutName := siteName + "_" + apCutName

								//Проверяем и Вносим во временную мапу. Заявка на данном этапе никакая ещё НЕ создаётся
								_, exisSiteName := siteapNameForTickets[siteName] //проверяем, есть ли siteName в мапе ДЛЯтикетов
								//for _, ticket := range sliceForTicket {

								//если в мапе дляТикета сайта ещё НЕТ
								if !exisSiteName {
									siteapNameForTickets[siteApCutName] = ForApsTicket{
										siteName,
										//ap.Name,
										ap.Mac,
										[]string{ap.Name},
									}

									//если в мапе дляТикета сайт уже есть, добавляем в массив точку
								} else {
									//в мапе нельзя просто изменить значение.
									for k, v := range siteapNameForTickets {
										if k == siteApCutName {
											//https://stackoverflow.com/questions/42716852/how-to-update-map-values-in-go
											/*1.Using pointers. не смог победить указатели...
											v2 := v
											v2.corpAnomalies = append(v2.corpAnomalies, anomaly.Anomaly)
											mapNoutnameFortickets[k] = v2 */

											//2.Reassigning the modified struct.
											v.apNames = append(v.apNames, ap.Name)
											siteapNameForTickets[k] = v
										}
									}
								}
							}
						}
					}
				}
				//Пробежались по всем точкам. Заводим заявки
				fmt.Println("")
				fmt.Println("Создание заявок по точкам:")
				for k, v := range siteapNameForTickets {
					fmt.Println(k)
					for _, s := range v.apNames {
						fmt.Println(s)
					}
					//usrLogin := noutnameLogin[v.clientName]
					usrLogin := siteApCutNameLogin[k]
					fmt.Println(usrLogin)

					desAps := strings.Join(v.apNames, "\n")
					description := "Зафиксировано отключение точек:" + "\n" +
						desAps + "\n" +
						"" + "\n" +
						"Рекомендации по выполнению таких инцидентов собраны на страничке корпоративной wiki" + "\n" +
						"https://wiki.tele2.ru/display/ITKB/%5BHelpdesk+IT%5D+System+Monitoring" + "\n" +
						""
					incidentType := "Недоступна точка доступа"

					//srTicketSlice := CreateApTicket(bpmServer, usrLogin, description, v.site, incidentType)
					srTicketSlice := CreateSmacWiFiTicket(bpmServer, usrLogin, description, v.site, incidentType)
					fmt.Println(srTicketSlice[2])
					apMacSRid[v.apMac] = srTicketSlice[0] //добавить в мапу apMac - ID Тикета
					fmt.Println("")
				}
				fmt.Println("")

				//count5minute = time.Now().Minute()
				count3minute = time.Now().Minute()
			}
			//
			//

			//АНОМАЛИИ. Блок кода запустится, если в этот ЧАС он ещё НЕ выполнялся
			//if time.Now().Minute() == 47 { // Если время 3 минуты от начала часа то блок для аномаоий
			if time.Now().Hour() != countHourAnom {
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

				type ForAnomalyTicket struct {
					site       string
					apName     string
					clientName string
					noutMac    string
					//userLogin     string //не помещаю, чтобы не делать лишних запросов к БД, если заявка НЕ будет создаваться
					corpAnomalies []string
				}

				//srStatusCodesForNewTicket := []string{
				srStatusCodesForNewTicket := map[string]bool{
					"Отменено":     true, //Cancel   6e5f4218-f46b-1410-fe9a-0050ba5d6c38
					"Решено":       true, //Resolve  ae7f411e-f46b-1410-009b-0050ba5d6c38
					"Закрыто":      true, //Closed  3e7f420c-f46b-1410-fc9a-0050ba5d6c38
					"На уточнении": true, //Clarification  81e6a1ee-16c1-4661-953e-dde140624fb
				}

				//mapNoutnameFortickets создаётся локально в блоке аномалий каждый час. Резервировать в БД НЕ нужно
				mapNoutnameForTickets := map[string]ForAnomalyTicket{} //https://stackoverflow.com/questions/42716852/how-to-update-map-values-in-go
				//
				for _, anomaly := range anomalies {
					_, existence := maschineMacName[anomaly.DeviceMAC] //проверяем, соответствует ли мак мапе corp клиентов
					//fmt.Println("Аномалии Tele2Corp клиентов:")
					if existence { //блок кода для Tele2Corp
						//если есть, выводим на экран с именем ПК, взятым из мапы
						//siteName := anomaly.SiteName[:len(anomaly.SiteName)-11]
						clientHostName := maschineMacName[anomaly.DeviceMAC]
						apName := namesClientAp[clientHostName]
						//usrLogin := GetLogin(clientHostName) //чтобы не делать лишних запросов к БД
						//fmt.Println(siteName, clientHostName, usrLogin, apName, anomaly.Datetime, anomaly.Anomaly)
						//fmt.Println(siteName, clientHostName, apName, anomaly.Datetime, anomaly.Anomaly) //без usrLogin

						_, exisClHostName := mapNoutnameForTickets[clientHostName] //проверяем, есть ли client hostname в мапе ДЛЯтикетов
						if !exisClHostName {                                       //если нет, добавляем новый
							mapNoutnameForTickets[clientHostName] = ForAnomalyTicket{ //https://stackoverflow.com/questions/42716852/how-to-update-map-values-in-go
								//siteName,
								anomaly.SiteName[:len(anomaly.SiteName)-11],
								apName,
								clientHostName,
								anomaly.DeviceMAC,
								[]string{anomaly.Anomaly},
							}
						} else { //если есть, добавляем данные в мапу
							for k, v := range mapNoutnameForTickets {
								if k == clientHostName {
									//https://stackoverflow.com/questions/42716852/how-to-update-map-values-in-go
									/*1.Using pointers. не смог победить указатели...
									v2 := v
									v2.corpAnomalies = append(v2.corpAnomalies, anomaly.Anomaly)
									mapNoutnameFortickets[k] = v2 */

									//2.Reassigning the modified struct.
									v.corpAnomalies = append(v.corpAnomalies, anomaly.Anomaly)
									mapNoutnameForTickets[k] = v
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
				for _, v := range mapNoutnameForTickets {
					if len(v.corpAnomalies) > 2 {
						fmt.Println(v.clientName)
						for _, s := range v.corpAnomalies {
							fmt.Println(s)
						}
						//SoapCreateTicket(clientHostName, v.clientName, v.corpAnomalies, siteName)
						//usrLogin := GetLogin(v.clientName)
						usrLogin := noutnameLogin[v.clientName]
						fmt.Println(usrLogin)
						// Проверяет, есть ли заявка в мапе ClientMacName - ID Тикета
						srID, existence := maschineMacSRid[v.noutMac]
						//Проверяем заявку на НЕ закрытость. если заявки нет - ничего страшного
						checkSlice := CheckTicketStatus(bpmServer, srID)

						desAnomalies := strings.Join(v.corpAnomalies, "\n")

						if srStatusCodesForNewTicket[checkSlice[1]] || !existence {
							//Если статус заявки Отменено, Закрыто, Решено, На уточнении
							//или заявки нет в мапе maschineMacSRid
							//То создаём новую
							delete(maschineMacSRid, v.noutMac) //удаляем заявку. если заявки нет - ничего страшного

							description := "На ноутбуке:" + "\n" +
								v.clientName + "\n" + "" + "\n" +
								"зафиксированы следующие Аномалии:" + "\n" +
								desAnomalies + "\n" +
								"" + "\n" +
								"Предполагаемое, но не на 100% точное имя точки:" + "\n" +
								v.apName + "\n" +
								"" + "\n" +
								"Рекомендации по выполнению таких инцидентов собраны на страничке корпоративной wiki" + "\n" +
								"https://wiki.tele2.ru/display/ITKB/%5BHelpdesk+IT%5D+System+Monitoring" + "\n" +
								""
							//fmt.Println(description)
							incidentType := "Плохое качество соединения клиента"

							//srTicketSlice := CreateAnomalyTicket(bpmServer, usrLogin, v.clientName, v.corpAnomalies, v.apName, v.site)
							srTicketSlice := CreateSmacWiFiTicket(bpmServer, usrLogin, description, v.site, incidentType)
							fmt.Println(srTicketSlice[2])
							maschineMacSRid[v.noutMac] = srTicketSlice[0] //добавить в мапу ClientMac - ID Тикета
							fmt.Println("")

						} else {
							//Если заявка уже есть, то добавить комментарий с новыми аномалтями
							comment := "Возникли новые аномалии за последний час:" + "\n" + desAnomalies
							AddComment(bpmServer, srID, comment)
						}

						/*старый блок кода
						if existence {
							//1. Если заявка В МАПЕ есть, проверить её статус
							statusSlice := CheckTicketStatus(bpmServer, srID)
							if srStatusCodesForNewTicket[statusSlice[1]] {
								//2. Если Статус закрыто, решено, отменено завести новую
								srTicketSlice := CreateAnomalyTicket(bpmServer, usrLogin, v.clientName, v.corpAnomalies, v.apName, v.site)
								fmt.Println(srTicketSlice[2])
								maschineMacSRid[v.noutMac] = srTicketSlice[0] //добавить в мапу ClientMac - ID Тикета
								fmt.Println("")
							} else {
								//2. Если статус НЕ закрыто, решено, отменено
								//В случае с аномалиями не делаем ничего
								//В ветке с точками будет проверка по каждой точке, не поднялась ли + комментарий
							}
						} else {
							//1. Если заявки В МАПЕ НЕТ
							srTicketSlice := CreateAnomalyTicket(bpmServer, usrLogin, v.clientName, v.corpAnomalies, v.apName, v.site)
							fmt.Println(srTicketSlice[2])
							maschineMacSRid[v.noutMac] = srTicketSlice[0] //добавить в мапу ClientMac - ID Тикета
							fmt.Println("")
						}*/
					}
				}
				fmt.Println("")
				countHourAnom = time.Now().Hour()
			} // END of ANOMALIES block

			//Обновление мап и БД. Блок кода запустится, если в этот ЧАС он ещё НЕ выполнялся
			if time.Now().Hour() != countHourDB {
				//noutnameLogin выгружать НЕ нужно
				UploadsMapsToDB(maschineMacName, "wifi_db", "wifi_db.maschine_mac_name", "TRUNCATE")
				UploadsMapsToDB(apMacName, "wifi_db", "wifi_db.ap_mac_name", "TRUNCATE")
				UploadsMapsToDB(namesClientAp, "wifi_db", "wifi_db.names_mascine_ap", "TRUNCATE")
				UploadsMapsToDB(maschineMacSRid, "wifi_db", "wifi_db.maschine_mac_srid", "DELETE")
				UploadsMapsToDB(apMacSRid, "wifi_db", "wifi_db.ap_mac_srid", "DELETE")
				countHourDB = time.Now().Hour()
			}
			//Обновление мап раз в сутки
			if time.Now().Day() != countDay {
				//noutnameLogin :=map[string]string{}     //clientHostName - > userLogin
				noutnameLogin = DownloadMapFromDB("glpi_db", "name", "contact", "glpi_db.glpi_computers", "date_mod")
				countDay = time.Now().Day()
			}
		} // Поминутный if
		fmt.Println("Time to Sleep")
		time.Sleep(60 * time.Second) //Изменить на 5 секунд на ПРОДе
	} // while TRUE
} //main func
