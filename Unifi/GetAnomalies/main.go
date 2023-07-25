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
	fmt.Println("")

	bpm := 1 // 0 -PROD; 1 -TEST
	var soapServer string
	var bpmUrl string
	if bpm == 0 {
		soapServer = "http://10.12.15.148/specs/aoi/tele2/bpm/bpmPortType" //PROD
		bpmUrl = "https://bpm.tele2.ru/0/Nui/ViewModule.aspx#CardModuleV2/CasePage/edit/"
	} else {
		soapServer = "http://10.246.37.15:8060/specs/aoi/tele2/bpm/bpmPortType" //TEST
		bpmUrl = "https://t2ru-tr-tst-01.corp.tele2.ru/0/Nui/ViewModule.aspx#CardModuleV2/CasePage/edit/"
	}

	unifiController := 11 //10-Rostov Local; 11-Rostov ip; 20-Novosib Local; 21-Novosib ip
	var urlController string
	var bdController int8 //Да string, потому что значение пойдёт в replace для БД
	//ROSTOV
	if unifiController == 10 || unifiController == 11 {
		bdController = 1
		if unifiController == 10 {
			urlController = "https://localhost:8443/"
		} else {
			urlController = "https://10.78.221.142:8443/"
		}
		//NOVOSIB
	} else if unifiController == 20 || unifiController == 21 {
		//else{
		bdController = 2
		if unifiController == 20 {
			urlController = "https://localhost:8443/"
		} else {
			urlController = "https://10.8.176.8:8443/"
		}
	}

	//countMinute := 0
	count3minute := 0
	//count5minute := 5
	countHourAnom := 0
	countHourDB := 0
	countDay := time.Now().Day()

	srStatusCodesForNewTicket := map[string]bool{
		"Отменено":     true, //Cancel  6e5f4218-f46b-1410-fe9a-0050ba5d6c38
		"Решено":       true, //Resolve  ae7f411e-f46b-1410-009b-0050ba5d6c38
		"Закрыто":      true, //Closed  3e7f420c-f46b-1410-fc9a-0050ba5d6c38
		"На уточнении": true, //Clarification 81e6a1ee-16c1-4661-953e-dde140624fb
		"":             true,
		"Тикет введён не корректно": true,
	}
	sitesException := map[string]bool{
		"5f2285f3a1a7693ae6139c00": true, //Novosi. Резерв/Склад
		"5f5b49d1a9f6167b55119c9b": true, //Ростов. Резерв/Склад
		//"Закрыто":      true, //Closed  3e7f420c-f46b-1410-fc9a-0050ba5d6c38
		//"На уточнении": true, //Clarification 81e6a1ee-16c1-4661-953e-dde140624fb
	}

	type ForApsTicket struct {
		site          string
		countIncident int
		apsMac        []string
		//apMac         string
		//apNames []string //сделано для массовых отключений точек при отключении света в офисе
	}
	type ForAnomalyTicket struct {
		site       string
		apName     string
		clientName string
		noutMac    string
		//userLogin     string //не помещаю, чтобы не делать лишних запросов к БД, если заявка НЕ будет создаваться
		corpAnomalies []string
	}

	//Download MAPs from DB
	//Выгружает только 4500 записей из 10000. отключаю. Проще делать разовый запрос по паре клиентов раз в час.
	//noutnameLogin :=map[string]string{}     //clientHostName - > userLogin
	//noutnameLogin := DownloadMapFromDB("glpi_db", "name", "contact", "glpi_db.glpi_computers", 0, "date_mod")
	siteApCutNameLogin := DownloadMapFromDB("wifi_db", "site_apcut", "login", "wifi_db.site_apcut_login", 0, "site_apcut")

	//machineMacName := map[string]string{}   // clientMAC -> clientHostName  // machineMAC -> machineHostName
	machineMacName := DownloadMapFromDB("wifi_db", "mac", "hostname", "wifi_db.machine_mac_name", bdController, "hostname")
	//apMacName := map[string]string{}      // apMac -> apName
	apMacName := DownloadMapFromDB("wifi_db", "mac", "name", "wifi_db.ap_mac_name", bdController, "name")
	//namesClientAps := map[string]string{} // clientName -> apName
	namesClientAp := DownloadMapFromDB("wifi_db", "machine_name", "ap_name", "wifi_db.names_machine_ap", bdController, "machine_name")
	//machineMacSRid := DownloadMapFromDB("wifi_db", "hostname", "srid", "wifi_db.mascine_name_srid", "hostname")
	machineMacSRid := DownloadMapFromDB("wifi_db", "mac", "srid", "wifi_db.machine_mac_srid", bdController, "mac")
	//apMacSRid := DownloadMapFromDB("wifi_db", "apname", "srid", "wifi_db.ap_name_srid", "apname")
	apMacSRid := DownloadMapFromDB("wifi_db", "mac", "srid", "wifi_db.ap_mac_srid", bdController, "mac")
	siteapNameForTickets := map[string]ForApsTicket{} //НЕ должна создаваться новая раз в 5 минут, поэтому здесь в отличие от аномальной
	//siteapNameForTickets := DownloadHardMapFromDB  //НЕ нужно резервировать, не делает погоду
	/*
		for k, v := range noutnameLogin {
			//fmt.Printf("key: %d, value: %t\n", k, v)
			fmt.Println("newMap "+k, v)
		}*/
	//os.Exit(0)
	fmt.Println("")

	c := unifi.Config{
		//c := *unifi.Config{  //ORIGINAL
		User: "unifi",
		Pass: "FORCEpower23",
		//URL:  "https://localhost:8443/"
		//URL:  "https://10.78.221.142:8443/", //ROSTOV
		//URL: "https://10.8.176.8:8443/", //NOVOSIB
		URL: urlController,
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
		//if time.Now().Minute() != countMinute { //Блок кода запустится, если в эту минуту он ещё НЕ выполнялся
		if time.Now().Minute() != 0 && time.Now().Minute()%3 == 0 && time.Now().Minute() != count3minute { //запускается раз в 3 минут
			count3minute = time.Now().Minute()

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

			clients, err := uni.GetClients(sites) //client = Notebook or Mobile = machine
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
					machineMacName[client.Mac] = client.Hostname //Добавить КОРП клиентов в map
					namesClientAp[client.Name] = apName          //Добавить Соответсвие имён клиентов и точек
				}
			}
			/*Вывести machineMacName мапу на экран
			for k, v := range machineMacName {
				//fmt.Printf("key: %d, value: %t\n", k, v)
				fmt.Println(k, v)
			}*/
			/*Вывести соответсвие имён клиентов и имён точек на экран
			for k, v := range namesClientAps {
				//fmt.Printf("key: %d, value: %t\n", k, v)
				fmt.Println(k, v)
			}*/
			//
			//countMinute = time.Now().Minute()
			//

			//
			//
			// блок кода про точки, а уже потом аномалии
			//if time.Now().Minute()%5 == 0 && time.Now().Minute() != count5minute { //запускается раз в 5 минут
			//if time.Now().Minute()%3 == 0 && time.Now().Minute() != count3minute { //запускается раз в 3 минуты
			fmt.Println("Обработка точек доступа...")

			for _, ap := range devices.UAPs {
				//fmt.Println(ap.Name)	fmt.Println(ap.SiteID)
				//if ap.SiteName[:len(ap.SiteName)-11] != "Резерв/Склад" {
				//if ap.SiteID != "5f2285f3a1a7693ae6139c00" { //NOVOSIB
				if !sitesException[ap.SiteID] {

					//fmt.Println(ap.Name)  fmt.Println(ap.SiteName)  fmt.Println(ap.SiteID)
					apLastSeen := ap.LastSeen.Int()
					_, exisApMacSRid := apMacSRid[ap.Mac]

					//Точка доступна. Заявки нет.
					if apLastSeen != 0 && !exisApMacSRid {
						//Идём дальше
						//fmt.Println("Точка доступна. Заявки нет")

						//Точка доступна. Заявка есть
					} else if apLastSeen != 0 && exisApMacSRid {
						fmt.Println("Точка доступна. Заявка есть")
						//ОЧИЩАЕМ мапу, оставляем коммент, ПЫТАЕМСЯ закрыть тикет, если на визировании
						//оставить комментарий, что точка стала доступна
						comment := "Точка появилась в сети: " + ap.Name
						AddComment(soapServer, apMacSRid[ap.Mac], comment, bpmUrl)

						//удалить запись из мапы, предварительно сохранив Srid
						srID := apMacSRid[ap.Mac]
						delete(apMacSRid, ap.Mac)
						//сложной мапы здесь уже нет. И удалять её не нужно и нечего

						//проверить, не последняя ли это запись была в мапе в массиве
						countOfIncident := 0
						for _, v := range apMacSRid {
							if v == srID {
								countOfIncident++
							}
						}
						if countOfIncident == 0 {
							//Пробуем закрыть тикет, только ЕСЛИ он на Визировании
							sliceTicketStatus := CheckTicketStatus(soapServer, apMacSRid[ap.Mac]) //получаем статус
							if sliceTicketStatus[1] == "На визировании" {
								//Если статус заявки по-прежнему на Визировании
								ChangeStatus(soapServer, srID, "На уточнении")
								AddComment(soapServer, srID, "Обращение отменено, т.к. все точки из него появились в сети", bpmUrl)
								ChangeStatus(soapServer, srID, "Отменено")
							}
						}

						//Точка недоступна.
					} else if apLastSeen == 0 {
						fmt.Println(ap.Name)
						fmt.Println(ap.Mac)
						fmt.Println("Точка НЕ доступна")
						//Проверяем заявку на НЕ закрытость. если заявки нет - ничего страшного
						checkSlice := CheckTicketStatus(soapServer, apMacSRid[ap.Mac])
						if srStatusCodesForNewTicket[checkSlice[1]] || !exisApMacSRid {
							fmt.Println("Заявка Закрыта, Отменена, Отклонена ИЛИ в мапе нет записи")
							delete(apMacSRid, ap.Mac) //удаляем заявку. если заявки нет - ничего страшного

							//Заполняем переменные, которые понадобятся дальше
							fmt.Println(ap.SiteID)
							var siteName string
							if ap.SiteID == "5e74aaa6a1a76964e770815c" { //6360a823a1a769286dc707f2
								siteName = "Урал"
							} else {
								siteName = ap.SiteName[:len(ap.SiteName)-11]
							}
							//apCutName := ap.Name[:len(ap.Name)3]
							apCutName := strings.Split(ap.Name, "-")[0]
							siteApCutName := siteName + "_" + apCutName
							fmt.Println(siteApCutName)

							//Проверяем и Вносим во временную мапу. Заявка на данном этапе никакая ещё НЕ создаётся
							_, exisSiteName := siteapNameForTickets[siteApCutName] //проверяем, есть ли siteName в мапе ДЛЯтикетов
							//for _, ticket := range sliceForTicket {

							//если в мапе дляТикета сайта ещё НЕТ
							if !exisSiteName {
								fmt.Println("в мапе для Тикета записи ещё НЕТ")
								//aps.mac := [string]
								siteapNameForTickets[siteApCutName] = ForApsTicket{
									siteName,
									0,
									[]string{ap.Mac},
									//ap.Mac,
									//[]string{ap.Name},
								}

								//если в мапе дляТикета сайт уже есть, добавляем в массив точку
							} else {
								fmt.Println("в мапе для Тикета запись ЕСТЬ")
								//в мапе нельзя просто изменить значение.
								for k, v := range siteapNameForTickets {
									if k == siteApCutName {
										//for _, apMac := range v.apsMac {
										if !cointains(v.apsMac, ap.Name) { //своя функция contains
											//https://stackoverflow.com/questions/42716852/how-to-update-map-values-in-go
											/*1.Using pointers. не смог победить указатели...
											v2 := v
											v2.corpAnomalies = append(v2.corpAnomalies, anomaly.Anomaly)
											mapNoutnameFortickets[k] = v2 */

											//2.Reassigning the modified struct.
											//v.apNames = append(v.apNames, ap.Name)
											v.apsMac = append(v.apsMac, ap.Mac)
											//Инкрементировать countIncident ? Вроде, нет
											siteapNameForTickets[k] = v
										}
									}
								}
							}
						} else {
							fmt.Println("Созданное обращение:")
							fmt.Println(bpmUrl + apMacSRid[ap.Mac])
							fmt.Println(checkSlice[1])
						}
						fmt.Println("")
					}
				} //fmt.Println("")
			}
			//Пробежались по всем точкам. Заводим заявки
			fmt.Println("")
			fmt.Println("Создание заявок по точкам:")
			for k, v := range siteapNameForTickets {
				fmt.Println(k)
				fmt.Println(v.countIncident) //"Число циклов захода на создание заявки: " +
				v.countIncident++

				//Если v.count < 10
				if v.countIncident < 10 {
					//обновляем мапу и инкрементируем count
					siteapNameForTickets[k] = v
				} else {
					//Если count == 10, Создаём заявку
					var apsNames []string
					//for _, s := range v.apNames {	fmt.Println(s)	}
					for _, mac := range v.apsMac {
						apName := apMacName[mac]
						apsNames = append(apsNames, apName)
						fmt.Println(apName)
					}

					//usrLogin := noutnameLogin[v.clientName]
					usrLogin := siteApCutNameLogin[k]
					fmt.Println(usrLogin)

					//desAps := strings.Join(v.apNames, "\n")
					desAps := strings.Join(apsNames, "\n")
					description := "Зафиксировано отключение точек:" + "\n" +
						desAps + "\n" +
						"" + "\n" +
						"Рекомендации по выполнению таких инцидентов собраны на страничке корпоративной wiki" + "\n" +
						"https://wiki.tele2.ru/display/ITKB/%5BHelpdesk+IT%5D+System+Monitoring" + "\n" +
						""
					incidentType := "Недоступна точка доступа"

					//srTicketSlice := CreateApTicket(soapServer, usrLogin, description, v.site, incidentType)
					srTicketSlice := CreateSmacWiFiTicket(soapServer, usrLogin, description, v.site, incidentType)
					fmt.Println(srTicketSlice[2])

					//apMacSRid[v.apMac] = srTicketSlice[0] //добавить в мапу apMac - ID Тикета
					for _, mac := range v.apsMac {
						apMacSRid[mac] = srTicketSlice[0]
					}
					fmt.Println("")

					//Удаляем запись в мапе
					delete(siteapNameForTickets, k)
				}
				//fmt.Println("")
				//count5minute = time.Now().Minute()
				//count3minute = time.Now().Minute()
			}
			//
			//

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

				//mapNoutnameFortickets создаётся локально в блоке аномалий каждый час. Резервировать в БД НЕ нужно
				mapNoutnameForTickets := map[string]ForAnomalyTicket{} //https://stackoverflow.com/questions/42716852/how-to-update-map-values-in-go
				//
				for _, anomaly := range anomalies {
					_, existence := machineMacName[anomaly.DeviceMAC] //проверяем, соответствует ли мак мапе corp клиентов
					//fmt.Println("Аномалии Tele2Corp клиентов:")
					if existence {
						//блок кода для Tele2Corp
						//если есть, выводим на экран с именем ПК, взятым из мапы
						//siteName := anomaly.SiteName[:len(anomaly.SiteName)-11]
						clientHostName := machineMacName[anomaly.DeviceMAC]
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
						usrLogin := GetLoginPC(v.clientName)
						//usrLogin := noutnameLogin[v.clientName]
						fmt.Println(usrLogin)
						// Проверяет, есть ли заявка в мапе ClientMacName - ID Тикета
						srID, existence := machineMacSRid[v.noutMac]
						//Проверяем заявку на НЕ закрытость. если заявки нет - ничего страшного
						checkSlice := CheckTicketStatus(soapServer, srID)

						desAnomalies := strings.Join(v.corpAnomalies, "\n")

						if srStatusCodesForNewTicket[checkSlice[1]] || !existence {
							fmt.Println("Заявка закрыта, Отменена, Отклонена ИЛИ в мапе нет записи")
							//Если статус заявки Отменено, Закрыто, Решено, На уточнении
							//или заявки нет в мапе machineMacSRid
							//То создаём новую
							delete(machineMacSRid, v.noutMac) //удаляем заявку. если заявки нет - ничего страшного

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

							//srTicketSlice := CreateAnomalyTicket(soapServer, usrLogin, v.clientName, v.corpAnomalies, v.apName, v.site)
							srTicketSlice := CreateSmacWiFiTicket(soapServer, usrLogin, description, v.site, incidentType)
							fmt.Println(srTicketSlice[2])
							machineMacSRid[v.noutMac] = srTicketSlice[0] //добавить в мапу ClientMac - ID Тикета
						} else {
							//Если заявка уже есть, то добавить комментарий с новыми аномалтями
							comment := "Возникли новые аномалии за последний час:" + "\n" + desAnomalies
							AddComment(soapServer, srID, comment, bpmUrl)
							//fmt.Println(comment)
						}
						fmt.Println("")
					}
				}
				//раз в час выполняет код по аномалиям. И БД обновляется в то же время
				UploadMapsToDBreplace(machineMacSRid, "wifi_db", "wifi_db.machine_mac_srid", "srid", bdController)
				fmt.Println("")
				countHourAnom = time.Now().Hour()
			} // END of ANOMALIES block

			//
			//
			//Обновление мап и БД. Блок кода запустится, если в этот ЧАС он ещё НЕ выполнялся
			if time.Now().Hour() != countHourDB {
				UploadMapsToDBreplace(machineMacName, "wifi_db", "wifi_db.machine_mac_name", "", bdController)
				UploadMapsToDBreplace(apMacName, "wifi_db", "wifi_db.ap_mac_name", "", bdController)
				UploadMapsToDBreplace(namesClientAp, "wifi_db", "wifi_db.names_machine_ap", "", bdController)

				///Закрытие заявок - это ДОП.функционал. Главное - СОЗДАТЬ заявку
				UploadMapsToDBreplace(apMacSRid, "wifi_db", "wifi_db.ap_mac_srid", "", bdController)
				//UploadsMapsToDB(machineMacSRid, "wifi_db", "wifi_db.machine_mac_srid", "DELETE")

				countHourDB = time.Now().Hour()
			}
			//Обновление мап раз в сутки
			if time.Now().Day() != countDay {
				//noutnameLogin :=map[string]string{}     //clientHostName - > userLogin
				//noutnameLogin = DownloadMapFromDB("glpi_db", "name", "contact", "glpi_db.glpi_computers", 0, "date_mod")
				siteApCutNameLogin = DownloadMapFromDB("wifi_db", "site_apcut", "login", "wifi_db.site_apcut_login", 0, "site_apcut")
				countDay = time.Now().Day()
			}
			//} // 3 минутный if про точки

		} //Снятие показаний раз в 3 минуты
		fmt.Println("Sleep 45s")
		fmt.Println("")
		time.Sleep(45 * time.Second) //Изменить на 5 секунд на ПРОДе

	} // while TRUE

} //main func

func cointains(slice []string, compareString string) bool {
	for _, v := range slice {
		if v == compareString {
			return true
		}
	}
	return false
}
