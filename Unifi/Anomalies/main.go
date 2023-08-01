package main

import (
	//"bytes"
	"fmt"
	"github.com/unpoller/unifi"
	"io"
	"log"
	"strings"

	//"strconv"
	//"strings"
	"time"
)

type MachineMyStruct struct {
	Hostname  string
	Exception int
	SrID      string
	ApName    string
}
type Machine struct {
	site     string
	ApName   string
	Hostname string
	Count    int8
}

func main() {
	fmt.Println("")

	unifiController := 21 //10-Rostov Local; 11-Rostov ip; 20-Novosib Local; 21-Novosib ip
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
	fmt.Println("Unifi controller")
	fmt.Println(urlController)
	fmt.Println(bdController)

	//machineMyMap := map[string]MachineMyStruct{}
	machineMyMap := DownloadMapFromDBmachines(bdController)

	//dateMac_mac := map[string]string{}
	dateMac_site := map[string]string{}

	//mac_count := map[string]int8{}
	//mac_machine := map[string]Machine{}

	//fmt.Println("Вывод мапы СНАРУЖИ функции")
	/*
		for k, v := range siteApCutNameLogin {
			//fmt.Printf("key: %d, value: %t\n", k, v)
			fmt.Println("newMap "+k, v)
		}
		os.Exit(0)
	*/

	fmt.Println("")

	c := unifi.Config{
		//c := *unifi.Config{  //ORIGINAL
		User: "unifi",
		Pass: "FORCEpower23",
		//URL: "https://localhost:8443/"
		//URL: "https://10.78.221.142:8443/", //ROSTOV
		//URL: "https://10.8.176.8:8443/",     //NOVOSIB
		URL: urlController,
		// Log with log.Printf or make your own interface that accepts (msg, test_SOAP)
		ErrorLog: log.Printf,
		DebugLog: log.Printf,
	}

	log.SetOutput(io.Discard) //Отключить вывод лога

	timeNow := time.Now()
	fmt.Println(timeNow.Format("02 January, 15:04:05"))

	//uni, err := unifi.NewUnifi(c)
	uni, err := unifi.NewUnifi(&c)
	if err != nil {
		log.Fatalln("Error:", err)
	} else {
		fmt.Println("uni загрузился")
	}

	sites, err := uni.GetSites()
	if err != nil {
		log.Fatalln("Error:", err)
	} else {
		fmt.Println("sites загрузились")
	}
	/*
		devices, err := uni.GetDevices(sites) //devices = APs
		if err != nil {
			log.Fatalln("Error:", err)
		} else {
			fmt.Println("devices загрузились")
		}

		clients, err := uni.GetClients(sites) //client = Notebook or Mobile = machine
		if err != nil {
			log.Fatalln("Error:", err)
		} else {
			fmt.Println("clients загрузились")
		}
	*/
	fmt.Println("")

	//
	/*
		//count := 60 //минус 70 минут
		//count := 3600
		//count := 36000 //+++
		//count := 86400
		//then := now.Add(time.Duration(-count) * time.Minute)
		//then := timeNow.Add(time.Duration(-count) * time.Minute)
		//then := timeNow.Add(time.Duration(-count) * time.)
	*/
	anomalies, err := uni.GetAnomalies(sites,
		time.Date(2023, 07, 01, 0, 0, 0, 0, time.Local), //time.Now(),
		//then,
		//time.Jul,
	)
	if err != nil {
		log.Fatalln("Error:", err)
	}
	for _, anomaly := range anomalies {
		siteName := anomaly.SiteName
		noutMac := anomaly.DeviceMAC
		anomalyStr := anomaly.Anomaly
		anomalyDatetime := anomaly.Datetime
		fmt.Println(anomalyDatetime, siteName, noutMac, anomalyStr)

		anomalyDatetime.String()
		uniqKey := anomalyDatetime.Format("2006-01-02") + "_" + noutMac
		//dateMac_mac[uniqKey] = noutMac
		dateMac_site[uniqKey] = siteName
	}

	//for k, v := range dateMac_mac {
	for k, v := range dateMac_site {
		kMac := strings.Split(k, "_")[1]
		for ke, va := range machineMyMap {
			if kMac == ke {
				va.Exception++
				va.SrID = v
				machineMyMap[ke] = va
			}
		}
	}
	for _, v := range machineMyMap {
		if v.Exception != 0 {
			login := GetLoginPC(v.Hostname)
			fmt.Println(v.SrID, v.ApName, v.Hostname, login, v.Exception)
		}
	}

	/*

					//mapNoutnameFortickets создаётся локально в блоке аномалий каждый час. Резервировать в БД НЕ нужно
					mapNoutnameForTickets := map[string]ForAnomalyTicket{}
					//https://stackoverflow.com/questions/42716852/how-to-update-map-values-in-go

					for _, anomaly := range anomalies {
						noutMac := anomaly.DeviceMAC
						siteName := anomaly.SiteName
						anomalyStr := anomaly.Anomaly
						anomaly.Datetime

						//_, existence := machineMacName[anomaly.DeviceMAC] //проверяем, соответствует ли мак мапе corp клиентов
						_, exMachMyMap := machineMyMap[noutMac] //проверяем, соответствует ли мак мапе corp клиентов

						//fmt.Println("Аномалии Tele2Corp клиентов:")
						//if existence {
						if exMachMyMap {
							//если есть, пробегаемся по той же мапе machineMyMap
							for ke, va := range machineMyMap {
								if ke == noutMac {
									//siteName := anomaly.SiteName[:len(anomaly.SiteName)-11]
									//clientHostName := machineMacName[anomaly.DeviceMAC]
									clientHostName := va.Hostname
									//apName := namesClientAp[clientHostName]
									apName := va.ApName

									//fmt.Println(siteName, clientHostName, apName, anomaly.Datetime, anomaly.Anomaly) //без usrLogin

									_, exisClHostName := mapNoutnameForTickets[clientHostName] //проверяем, есть ли в мапе ДЛЯтикетов
									if !exisClHostName {
										//если нет, добавляем новый
										mapNoutnameForTickets[clientHostName] = ForAnomalyTicket{ //https://stackoverflow.com/questions/42716852/how-to-update-map-values-in-go
											//anomaly.SiteName[:len(anomaly.SiteName)-11],
											siteName[:len(siteName)-11],
											apName,
											//clientHostName,
											//anomaly.DeviceMAC,
											noutMac,
											//[]string{anomaly.Anomaly},
											[]string{anomalyStr},
										}
									} else { //если есть, добавляем данные в мапу
										for k, v := range mapNoutnameForTickets {
											if k == clientHostName {
												//https://stackoverflow.com/questions/42716852/how-to-update-map-values-in-go

												//2.Reassigning the modified struct.
												//v.corpAnomalies = append(v.corpAnomalies, anomaly.Anomaly)
												v.corpAnomalies = append(v.corpAnomalies, anomalyStr)
												mapNoutnameForTickets[k] = v
											}
										}
									}
									break
								}
							}
						} else {
							//Обработка аномалий для Tele2Guest.
							//Пока просто заглушка
						}
					}

					fmt.Println("")
					fmt.Println("Tele2Corp клиенты с более чем 2 аномалиями:")
					for k, v := range mapNoutnameForTickets {
						corpAnomalies := v.corpAnomalies
						noutMac := v.noutMac
						//if len(v.corpAnomalies) > 2 {
						if len(corpAnomalies) > 2 {
							//fmt.Println(v.clientName)
							fmt.Println(k)
							usrLogin := GetLoginPC(k)
							fmt.Println(usrLogin)
							for _, s := range v.corpAnomalies {
								fmt.Println(s)
							}

							//Проверяет, есть ли заявка в мапе ClientMacName - ID Тикета
							//srID, existence := machineMacSRid[v.noutMac]
							//Выходим на создание заявки
							for ke, va := range machineMyMap {
								if ke == noutMac {
									//Если есть исключение, прерываем for
									if va.Exception > 0 {
										fmt.Println("Точка или Клиент добавлены в исключение")
										break
									}

									srID := va.SrID
									//Проверяем заявку на НЕ закрытость. если заявки нет - ничего страшного
									//checkSlice := CheckTicketStatus(soapServer, srID)
									checkSlice := CheckTicketStatus(soapServer, srID)

									desAnomalies := strings.Join(v.corpAnomalies, "\n")

									//if srStatusCodesForNewTicket[checkSlice[1]] || !existence {
									if srStatusCodesForNewTicket[checkSlice[1]] {
										fmt.Println("Заявка закрыта, Отменена, Отклонена ИЛИ в мапе нет записи")

										//Удалять старую запись необязательно. Обновим позже на другую
										//delete(machineMacSRid, v.noutMac) //удаляем заявку. если заявки нет - ничего страшного

										//То создаём новую
										description := "На ноутбуке:" + "\n" +
											k + "\n" + "" + "\n" +
											"За последний ЧАС зафиксированы следующие Аномалии:" + "\n" +
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

										//machineMacSRid[v.noutMac] = srTicketSlice[0] //добавить в мапу ClientMac - ID Тикета
										va.SrID = srTicketSlice[0]
										machineMyMap[ke] = va

									} else {
										//Если заявка уже есть, то добавить комментарий с новыми аномалиями
										comment := "Возникли новые аномалии за последний час:" + "\n" + desAnomalies
										AddComment(soapServer, srID, comment, bpmUrl)
										//fmt.Println(comment)
									}
									break
								}
							}
							fmt.Println("")
						}
					}

					//раз в час выполняет код по аномалиям. И БД обновляется в то же время.
					//Обновление реализовал ниже в другом блоке
					//UploadMapsToDBreplace(machineMacSRid, "wifi_db", "wifi_db.machine_mac_srid", "srid", bdController)
					fmt.Println("")
				}
				// END of ANOMALIES block
				//
				//

				//
				//
				//Обновление мап и БД
				//запустится, если в этот ЧАС он ещё НЕ выполнялся
				//if time.Now().Hour() != countHourDB {
				if timeNow.Hour() != countHourDB {
					//countHourDB = time.Now().Hour()
					countHourDB = timeNow.Hour()

					bdCntrl := strconv.Itoa(int(bdController))
					var lenMap int
					var count int
					var exception string
					var b1 bytes.Buffer
					var b2 bytes.Buffer
					var query string

					//b.WriteString("REPLACE INTO " + tableName + " VALUES ")
					b1.WriteString("REPLACE INTO " + "it_support_db.ap" + " VALUES ")
					//lenMap := len(uploadMap)
					lenMap = len(apMyMap)
					count = 0
					//for k, v := range uploadMap {
					for k, v := range apMyMap {
						exception = strconv.Itoa(int(v.Exception))
						count++
						if count != lenMap {
							// mac, name, controller, exception, srid
							b1.WriteString("('" + k + "','" + v.Name + "','" + bdCntrl + "','" + exception + "','" + v.SrID + "'),")
						} else {
							b1.WriteString("('" + k + "','" + v.Name + "','" + bdCntrl + "','" + exception + "','" + v.SrID + "')")
							//в конце НЕ ставим запятую
						}
					}
					query = b1.String()
					fmt.Println(query)
					if count != 0 {
						UploadMapsToDBstring("it_support_db", query)
					} else {
						fmt.Println("Передана пустая карта. Запрос не выполнен")
					}
					fmt.Println("")

					//

					//b.WriteString("REPLACE INTO " + tableName + " VALUES ")
					b2.WriteString("REPLACE INTO " + "it_support_db.machine" + " VALUES ")
					//lenMap := len(uploadMap)
					lenMap = len(machineMyMap)
					count = 0
					//for k, v := range uploadMap {
					for k, v := range machineMyMap {
						exception = strconv.Itoa(int(v.Exception))
						count++
						if count != lenMap {
							// mac, hostname, controller, exception, srid, apname
							b2.WriteString("('" + k + "','" + v.Hostname + "','" + bdCntrl + "','" + exception + "','" + v.SrID + "','" + v.ApName + "'),")
						} else {
							b2.WriteString("('" + k + "','" + v.Hostname + "','" + bdCntrl + "','" + exception + "','" + v.SrID + "','" + v.ApName + "')")
							//в конце НЕ ставим запятую
						}
					}
					query = b2.String()
					fmt.Println(query)
					if count != 0 {
						UploadMapsToDBstring("it_support_db", query)
					} else {
						fmt.Println("Передана пустая карта. Запрос не выполнен")
					}
					fmt.Println("")
				}

				//
				//
				//Обновление мап раз в сутки
				//if time.Now().Day() != countDay {
				if timeNow.Day() != countDay {
					//countDay = time.Now().Day()
					countDay = timeNow.Day()

					//siteApCutNameLogin = DownloadMapFromDB("wifi_db", "site_apcut", "login", "wifi_db.site_apcut_login", 0, "site_apcut")
					siteApCutNameLogin = DownloadMapFromDB("it_support_db", "site_apcut", "login", "it_support_db.site_apcut_login", 0, "site_apcut")
				}
			}
		} // while TRUE
	*/
} //main func

func cointains(slice []string, compareString string) bool {
	for _, v := range slice {
		if v == compareString {
			return true
		}
	}
	return false
}
