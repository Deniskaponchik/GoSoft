package main

import (
	"bytes"
	"fmt"
	"github.com/unpoller/unifi"
	"log"
	"strconv"
	"strings"
	"time"
)

func main() {
	fmt.Println("")

	unifiController := 21 //10-Rostov Local; 11-Rostov ip; 20-Novosib Local; 21-Novosib ip
	var urlController string
	var bdController int8 //Да string, потому что значение пойдёт в replace для БД
	everyStartCode := map[int]bool{}
	//ROSTOV
	if unifiController == 10 || unifiController == 11 {
		bdController = 1
		if unifiController == 10 {
			urlController = "https://localhost:8443/"
		} else {
			urlController = "https://10.78.221.142:8443/"
		}
		//everyStartCode := [10] int8 {3, 9, 15, 21, 27, 33, 39, 45, 51, 57}
		everyStartCode = map[int]bool{
			3:  true,
			9:  true,
			15: true,
			21: true,
			27: true,
			33: true,
			39: true,
			45: true,
			51: true,
			57: true,
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
		//everyStartCode := [10] int8 {6, 12, 18, 24, 30, 36, 42, 48, 54, 59}
		everyStartCode = map[int]bool{
			6:  true,
			12: true,
			18: true,
			24: true,
			30: true,
			36: true,
			42: true,
			48: true,
			54: true,
			59: true,
		}
	}
	fmt.Println("Unifi controller")
	fmt.Println(urlController)

	var soapServer string
	soapServerProd := "http://10.12.15.148/specs/aoi/tele2/bpm/bpmPortType"      //PROD
	soapServerTest := "http://10.246.37.15:8060/specs/aoi/tele2/bpm/bpmPortType" //TEST
	var bpmUrl string
	bpmUrlProd := "https://bpm.tele2.ru/0/Nui/ViewModule.aspx#CardModuleV2/CasePage/edit/"
	bpmUrlTest := "https://t2ru-tr-tst-01.corp.tele2.ru/0/Nui/ViewModule.aspx#CardModuleV2/CasePage/edit/"
	/*
		bpm := 1 // 0-PROD; 1-TEST
		var soapServer string
		var bpmUrl string
		if bpm == 0 {
			soapServer = "http://10.12.15.148/specs/aoi/tele2/bpm/bpmPortType" //PROD
			bpmUrl = "https://bpm.tele2.ru/0/Nui/ViewModule.aspx#CardModuleV2/CasePage/edit/"
		} else {
			soapServer = "http://10.246.37.15:8060/specs/aoi/tele2/bpm/bpmPortType" //TEST
			bpmUrl = "https://t2ru-tr-tst-01.corp.tele2.ru/0/Nui/ViewModule.aspx#CardModuleV2/CasePage/edit/"
		}
		fmt.Println("SOAP")
		fmt.Println(soapServer)
		fmt.Println("BPM")
		fmt.Println(bpmUrl)
		fmt.Println("")
	*/

	count6minute := 0
	countHourAnom := 0
	countHourDBap := 0
	countHourDBmachine := 0
	countDay := time.Now().Day()

	srStatusCodesForNewTicket := map[string]bool{
		"Отменено":     true, //Cancel  6e5f4218-f46b-1410-fe9a-0050ba5d6c38
		"Решено":       true, //Resolve  ae7f411e-f46b-1410-009b-0050ba5d6c38
		"Закрыто":      true, //Closed  3e7f420c-f46b-1410-fc9a-0050ba5d6c38
		"На уточнении": true, //Clarification 81e6a1ee-16c1-4661-953e-dde140624fb
		"Тикет введён не корректно": true,
		//"": true,
	}
	srStatusCodesForCancelTicket := map[string]bool{
		"Визирование":  true,
		"Назначено":    true,
		"На уточнении": true, //Clarification 81e6a1ee-16c1-4661-953e-dde140624fb
	}
	sitesException := map[string]bool{
		"5f2285f3a1a7693ae6139c00": true, //Novosi. Резерв/Склад
		"5f5b49d1a9f6167b55119c9b": true, //Ростов. Резерв/Склад
		//"Закрыто":      true, //Closed  3e7f420c-f46b-1410-fc9a-0050ba5d6c38
		//"На уточнении": true, //Clarification 81e6a1ee-16c1-4661-953e-dde140624fb
	}

	//Download MAPs from DB
	//apMyMap := map[string]ApMyStruct{}
	//apMyMap := DownloadMapFromDBaps(bdController)
	apMyMap := DownloadMapFromDBapsErr(bdController)

	//НЕ должна создаваться новая раз в 5 минут, поэтому здесь в отличие от аномальной
	siteapNameForTickets := map[string]ForApsTicket{}

	//machineMyMap := map[string]MachineMyStruct{}
	//machineMyMap := DownloadMapFromDBmachines(bdController)
	machineMyMap := DownloadMapFromDBmachinesErr(bdController)

	//siteApCutNameLogin := map[string]string{}
	//siteApCutNameLogin := DownloadMapFromDB("it_support_db", "site_apcut", "login", "it_support_db.site_apcut_login", 0, "site_apcut")
	siteApCutNameLogin := DownloadMapFromDBerr()

	//fmt.Println("Вывод мапы СНАРУЖИ функции")
	/*
		for k, v := range siteApCutNameLogin {
			//fmt.Printf("key: %d, value: %t\n", k, v)
			fmt.Println("newMap "+k, v)
		}
		os.Exit(0)
	*/
	fmt.Println("")
	//
	//

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

	//log.SetOutput(io.Discard) //Отключить вывод лога

	//var uni unifi.Unifi
	//sites := []*unifi.Site()

	for true { //зацикливаем навечно
		//currentMinute := time.Now().Minute()
		timeNow := time.Now()

		//Снятие показаний с контрллера раз в 6 минут. промежутки разные для контроллеров
		//if time.Now().Minute() != 0 && time.Now().Minute()%3 == 0 && time.Now().Minute() != count3minute {
		//if currentMinute != 0 && everyStartCode[currentMinute] && currentMinute != count6minute {
		if timeNow.Minute() != 0 && everyStartCode[timeNow.Minute()] && timeNow.Minute() != count6minute {
			//count6minute = time.Now().Minute()
			count6minute = timeNow.Minute()
			//fmt.Println(time.Now().String())
			fmt.Println(timeNow.Format("02 January, 15:04:05"))

			//uni, err := unifi.NewUnifi(c)
			uni, errNewUnifi := unifi.NewUnifi(&c)
			if errNewUnifi == nil {
				fmt.Println("uni загрузился")

				sites, errGetSites := uni.GetSites()
				if errGetSites == nil {
					fmt.Println("sites загрузились")

					devices, errGetDevices := uni.GetDevices(sites) //devices = APs
					if errGetDevices == nil {
						fmt.Println("devices загрузились")
						fmt.Println("")

						//
						//
						//ТОЧКИ
						soapServer = soapServerProd
						fmt.Println("SOAP")
						fmt.Println(soapServer)
						bpmUrl = bpmUrlProd
						fmt.Println("BPM")
						fmt.Println(bpmUrl)
						fmt.Println("")

						fmt.Println("Обработка точек доступа...")
						fmt.Println("")
						for _, ap := range devices.UAPs {
							siteID := ap.SiteID
							if !sitesException[siteID] { // НЕ Резерв/Склад
								apMac := ap.Mac
								apName := ap.Name
								apLastSeen := ap.State.Int()
								//apState := ap.State.Int()

								//fmt.Println(ap.Name)	fmt.Println(ap.Uptime.Int())  fmt.Println(ap.Uptime.String()) fmt.Println(ap.Uptime.Val)
								//_, exisApMacSRid := apMacSRid[ap.Mac]
								_, exisApMyMap := apMyMap[apMac]
								if !exisApMyMap { //если в мапе нет записи, создаём
									apMyMap[apMac] = ApMyStruct{
										apName,
										0,
										"",
									}
								}
								for keyAp, valueAp := range apMyMap {
									if keyAp == apMac {
										srID := valueAp.SrID

										apSiteName := ap.SiteName
										var siteName string
										if siteID == "5e74aaa6a1a76964e770815c" {
											siteName = "Урал" //именно с дефолтными сайтами так почему-то
										} else if siteID == "5e758bdca9f6163bb0c3c962" {
											siteName = "Волга" //именно с дефолтными сайтами так почему-то
										} else {
											siteName = apSiteName[:len(apSiteName)-11]
										}
										apCutName := strings.Split(apName, "-")[0]
										siteApCutName := siteName + "_" + apCutName

										//
										//Точка доступна. Заявки нет.
										if apLastSeen != 0 && srID == "" {
											//for site_apCut
											for k, v := range siteapNameForTickets {
												if k == siteApCutName {
													for ke, _ := range v.apsMacName {
														if ke == apMac {
															delete(v.apsMacName, ke)
														}
													}
												}
											}

											//Точка доступна. Заявка есть.   +Имя точки обновляю
										} else if apLastSeen != 0 && srID != "" {
											fmt.Println(apName)
											fmt.Println(apMac)
											fmt.Println("Точка доступна. Заявка есть")
											//Оставляем коммент, Очищаем запись в мапе, ПЫТАЕМСЯ закрыть тикет, если на визировании

											//comment := "Точка появилась в сети: " + ap.Name
											comment := "Точка появилась в сети: " + apName
											//AddComment(soapServer, srID, comment, bpmUrl)
											AddCommentErr(soapServer, srID, comment, bpmUrl)

											/*удалить запись из мапы, предварительно сохранив Srid
											srID := apMacSRid[ap.Mac]
											delete(apMacSRid, ap.Mac)
											//сложной мапы здесь уже нет. И удалять её не нужно и нечего
											*/
											//удалить запись из мапы, заодно и имя обновим
											valueAp.Name = apName
											valueAp.SrID = ""
											apMyMap[keyAp] = valueAp

											//проверить, не последняя ли это запись была в мапе в массиве
											countOfIncident := 0
											for _, v := range apMyMap {
												if v.SrID == srID {
													countOfIncident++
													//fmt.Println(countOfIncident)
												}
											}
											if countOfIncident == 0 {
												//Пробуем закрыть тикет, только ЕСЛИ он на Визировании
												//fmt.Println("Попали в блок изменения статусов заявок")
												//sliceTicketStatus := CheckTicketStatus(soapServer, srID) //получаем статус
												sliceTicketStatus := CheckTicketStatusErr(soapServer, srID)
												fmt.Println(sliceTicketStatus[1])
												//if sliceTicketStatus[1] == "Визирование" || sliceTicketStatus[1] == "Назначено" {
												if srStatusCodesForCancelTicket[sliceTicketStatus[1]] {
													//Если статус заявки на Уточнении, Визирование, Назначено
													//ChangeStatus(soapServer, srID, "На уточнении")
													ChangeStatusErr(soapServer, srID, "На уточнении")
													//AddComment(soapServer, srID, "Обращение отменено, т.к. все точки из него появились в сети", bpmUrl)
													AddCommentErr(soapServer, srID, "Обращение отменено, т.к. все точки из него появились в сети", bpmUrl)
													//ChangeStatus(soapServer, srID, "Отменено")
													ChangeStatusErr(soapServer, srID, "Отменено")
												}
											}
											fmt.Println("")

											//
											//Точка недоступна
										} else if apLastSeen == 0 {
											fmt.Println(apName)
											fmt.Println(apMac)
											fmt.Println("Точка НЕ доступна")

											//Проверяем заявку на НЕ закрытость. если заявки нет - ничего страшного
											//checkSlice := CheckTicketStatus(soapServer, srID)
											checkSlice := CheckTicketStatusErr(soapServer, srID)

											//if srStatusCodesForNewTicket[checkSlice[1]] || !exisApMacSRid {
											if srStatusCodesForNewTicket[checkSlice[1]] || srID == "" {
												fmt.Println(bpmUrl + srID)
												fmt.Println("Статус: " + checkSlice[1])
												fmt.Println("Заявка Закрыта, Отменена, Отклонена ИЛИ заявки нет вовсе")

												//delete(apMacSRid, ap.Mac) //удаляем заявку. если заявки нет - ничего страшного
												//удаляем заявку + обновить имя
												valueAp.Name = apName
												valueAp.SrID = ""
												apMyMap[keyAp] = valueAp

												//Заполняем переменные, которые понадобятся дальше
												fmt.Println("Site ID: " + siteID)
												fmt.Println(siteApCutName)

												//Проверяем и вносим во временную мапу. Заявка на данном этапе никакая ещё НЕ создаётся
												_, exisSiteName := siteapNameForTickets[siteApCutName] //проверяем, есть ли в мапе ДЛЯтикетов
												//если в мапе дляТикета сайта ещё НЕТ
												if !exisSiteName {
													fmt.Println("в мапе для Тикета записи ещё НЕТ")
													//aps.mac := [string]
													siteapNameForTickets[siteApCutName] = ForApsTicket{
														siteName,
														0,
														map[string]string{apMac: apName},
													}

													//если в мапе дляТикета сайт уже есть, добавляем в массив точку
												} else {
													fmt.Println("в мапе для Тикета запись ЕСТЬ")
													//в мапе нельзя просто изменить значение.
													for k, v := range siteapNameForTickets {
														if k == siteApCutName {
															_, exisApsMacName := v.apsMacName[apMac]
															//for _, apMac := range v.apsMac {
															//if !cointains(v.apsMac, ap.Name) { //своя функция contains
															if !exisApsMacName {
																//https://stackoverflow.com/questions/42716852/how-to-update-map-values-in-go
																/*1.Using pointers. не смог победить указатели...
																v2 := v
																v2.corpAnomalies = append(v2.corpAnomalies, anomaly.Anomaly)
																mapNoutnameFortickets[k] = v2 */

																//2.Reassigning the modified struct.
																/*первичное решение через другую мапу
																//v.apNames = append(v.apNames, ap.Name)
																//v.apsMac = append(v.apsMac, ap.Mac)
																//прошлое решение через массив
																//v.apsMac = append(v.apsMac, apMac)
																//siteapNameForTickets[k] = v
																*/
																v.apsMacName[apMac] = apName
																siteapNameForTickets[k] = v

																break // ЗДЕСЬ break НЕ НУЖЕН! да вроде, нужен
															}
														}
													}
												}
											} else {
												fmt.Println("Созданное обращение:")
												fmt.Println(bpmUrl + srID)
												fmt.Println(checkSlice[1])
											}
											fmt.Println("")
										}
										break
									}
								}
							} //fmt.Println("")
						}
						//Пробежались по всем точкам. Заводим заявки
						fmt.Println("")

						fmt.Println("Создание заявок по точкам:")
						for k, v := range siteapNameForTickets {
							if len(v.apsMacName) > 0 {
								//vCountIncident := v.countIncident
								fmt.Println(k)
								fmt.Println(v.countIncident) //"Число циклов захода на создание заявки: " +
								v.countIncident++

								//Если v.count < 10
								if v.countIncident < 5 {
									//обновляем мапу и инкрементируем count
									siteapNameForTickets[k] = v
								} else {
									//Если count == 10, Создаём заявку
									var apsNames []string
									//for _, s := range v.apNames {	fmt.Println(s)	}
									//for _, mac := range v.apsMac {
									for _, name := range v.apsMacName {
										//apName := apMacName[mac] //сходить в другую мапу
										//apName := apMacName[name]
										//apsNames = append(apsNames, apName)
										apsNames = append(apsNames, name)
										fmt.Println(name)
									}

									//usrLogin := noutnameLogin[v.clientName]
									usrLogin := siteApCutNameLogin[k]
									if usrLogin == "" {
										usrLogin = "denis.tirskikh"
									}
									fmt.Println(usrLogin)

									//desAps := strings.Join(v.apNames, "\n")
									desAps := strings.Join(apsNames, "\n")
									description := "Зафиксировано отключение Wi-Fi точек доступа:" + "\n" +
										desAps + "\n" +
										"" + "\n" +
										"Рекомендации по выполнению таких инцидентов собраны на страничке корпоративной wiki" + "\n" +
										"https://wiki.tele2.ru/display/ITKB/%5BHelpdesk+IT%5D+System+Monitoring" + "\n" +
										"" + "\n" +
										"!!! Не нужно решать/отменять/отклонять/возвращать/закрывать заявку, пока работа точек не будет восстановлена - автоматически создастся новый тикет !!!" + "\n" +
										""
									incidentType := "Недоступна точка доступа"

									//srTicketSlice := CreateSmacWiFiTicket(soapServer, usrLogin, description, v.site, incidentType)
									srTicketSlice := CreateSmacWiFiTicketErr(soapServer, bpmUrl, usrLogin, description, v.site, incidentType)
									fmt.Println(srTicketSlice[2])

									//apMacSRid[v.apMac] = srTicketSlice[0] //добавить в мапу apMac - ID Тикета
									//for _, mac := range v.apsMac {
									for mac, _ := range v.apsMacName {
										//apMacSRid[mac] = srTicketSlice[0]
										for key, value := range apMyMap {
											if key == mac {
												value.SrID = srTicketSlice[0]
												apMyMap[key] = value
												break
											}
										}
									}
									fmt.Println("")

									//Удаляем запись в мапе
									delete(siteapNameForTickets, k)
								}
							} else {
								delete(siteapNameForTickets, k)
							}
						}
						fmt.Println("")

						//Обновление БД ap раз в час
						if timeNow.Hour() != countHourDBap {
							countHourDBap = timeNow.Hour()

							bdCntrl := strconv.Itoa(int(bdController))
							var lenMap int
							var count int
							var exception string
							var b1 bytes.Buffer
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
								//UploadMapsToDBstring("it_support_db", query)
								UploadMapsToDBerr(query)
							} else {
								fmt.Println("Передана пустая карта. Запрос не выполнен")
							}
							fmt.Println("")
						}

						//
						//
						clients, errGetClients := uni.GetClients(sites) //client = Notebook or Mobile = machine
						if errGetClients == nil {
							fmt.Println("clients загрузились")
							fmt.Println("")
							//var apName string
							for _, client := range clients {
								if !client.IsGuest.Val {

									apName := client.ApName //НИЧЕГО не выводит и не содержит...
									clientMac := client.Mac
									clientName := client.Name

									var clExInt int
									if client.Noted.Val {
										clientExceptionStr := strings.Split(client.Note, " ")[0]
										if clientExceptionStr == "Exception" {
											clExInt = 1
										} else {
											clExInt = 0
										}
									}
									/*1. Если разработчик исправит скрипт, и будет возможность получать имя точек у клиентов
									_, exisNoutMyMap := machineMyMap[clientMac]
									if !exisNoutMyMap { //если записи клиента НЕТ
										machineMyMap[clientMac] = MachineMyStruct{
											clientName,
											0,
											"",
											apName,
										}
									} else {
										for ke, va := range machineMyMap {
											if ke == clientMac {
												va.ApName = apName
												va.Hostname = clientName
												va.Exception = clExInt
												machineMyMap[ke] = va
												break //прекращаем цикл, когда найден клиент и имя точки присвоено ему
											}
										}
									}*/
									//2. Если разработчик НЕ исправит: https://github.com/unpoller/unifi/issues/90
									//пробегаемся по всей мапе точек и получаем имя соответствию мака
									for k, v := range apMyMap {
										if k == clientMac {
											apName = v.Name
											apException := v.Exception
											//пробегаемся по всей мапе клиентов и назначаем имя точки клиенту
											_, exisNoutMyMap := machineMyMap[clientMac]
											if !exisNoutMyMap { //если записи клиента НЕТ
												machineMyMap[clientMac] = MachineMyStruct{
													clientName,
													clExInt + apException,
													"",
													apName,
												}
											} else { //если запись клиента создана, обновляем её
												for ke, va := range machineMyMap {
													if ke == client.Mac {
														va.Hostname = clientName
														va.ApName = apName
														va.Exception = clExInt + apException
														machineMyMap[ke] = va
														break //прекращаем цикл, когда найден клиент и имя точки присвоено ему
													}
												}
											}
											break //прекращаем цикл, когда найден мак точки
										}
									}
								}
							}

							//
							//
							//АНОМАЛИИ
							if timeNow.Hour() != countHourAnom {
								//countHourAnom = time.Now().Hour()
								countHourAnom = timeNow.Hour()

								soapServer = soapServerTest
								fmt.Println("SOAP")
								fmt.Println(soapServer)
								bpmUrl = bpmUrlTest
								fmt.Println("BPM")
								fmt.Println(bpmUrl)
								fmt.Println("")

								count := 60 //минус 70 минут
								//then := now.Add(time.Duration(-count) * time.Minute)
								then := timeNow.Add(time.Duration(-count) * time.Minute)

								anomalies, errGetAnomalies := uni.GetAnomalies(sites,
									//time.Date(2023, 07, 11, 7, 0, 0, 0, time.Local), time.Now()
									then,
								)
								if errGetAnomalies == nil {
									fmt.Println("anomalies загрузились")
									fmt.Println("")

									//mapNoutnameFortickets создаётся локально в блоке аномалий каждый час. Резервировать в БД НЕ нужно
									mapNoutnameForTickets := map[string]ForAnomalyTicket{}
									//https://stackoverflow.com/questions/42716852/how-to-update-map-values-in-go

									for _, anomaly := range anomalies {
										noutMac := anomaly.DeviceMAC
										siteName := anomaly.SiteName
										anomalyStr := anomaly.Anomaly

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
																/*1.Using pointers. не смог победить указатели...
																v2 := v
																v2.corpAnomalies = append(v2.corpAnomalies, anomaly.Anomaly)
																mapNoutnameFortickets[k] = v2 */

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
											//usrLogin := GetLoginPC(k)
											usrLogin := GetLoginPCerr(k)
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
													checkSlice := CheckTicketStatusErr(soapServer, srID)

													desAnomalies := strings.Join(v.corpAnomalies, "\n")

													//if srStatusCodesForNewTicket[checkSlice[1]] || !existence {
													if srStatusCodesForNewTicket[checkSlice[1]] {
														fmt.Println("Заявка Закрыта, Отменена, Отклонена ИЛИ в мапе нет записи")

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

														//srTicketSlice := CreateSmacWiFiTicket(soapServer, usrLogin, description, v.site, incidentType)
														srTicketSlice := CreateSmacWiFiTicketErr(soapServer, bpmUrl, usrLogin, description, v.site, incidentType)
														fmt.Println(srTicketSlice[2])

														//machineMacSRid[v.noutMac] = srTicketSlice[0] //добавить в мапу ClientMac - ID Тикета
														va.SrID = srTicketSlice[0]
														machineMyMap[ke] = va

													} else {
														//Если заявка уже есть, то добавить комментарий с новыми аномалиями
														comment := "Возникли новые аномалии за последний час:" + "\n" + desAnomalies
														//AddComment(soapServer, srID, comment, bpmUrl)
														AddCommentErr(soapServer, srID, comment, bpmUrl)
													}
													break
												}
											}
											fmt.Println("")
										}
									}

									//UploadMapsToDBreplace(machineMacSRid, "wifi_db", "wifi_db.machine_mac_srid", "srid", bdController)
									fmt.Println("")

									//
									//
									//Обновление БД machine раз в час
									if timeNow.Hour() != countHourDBmachine {
										countHourDBmachine = timeNow.Hour()
										/* OLD
										UploadMapsToDBreplace(machineMacName, "wifi_db", "wifi_db.machine_mac_name", "", bdController)
										UploadMapsToDBreplace(apMacName, "wifi_db", "wifi_db.ap_mac_name", "", bdController)
										UploadMapsToDBreplace(namesClientAp, "wifi_db", "wifi_db.names_machine_ap", "", bdController)
										UploadMapsToDBreplace(apMacSRid, "wifi_db", "wifi_db.ap_mac_srid", "", bdController)
										//UploadsMapsToDB(machineMacSRid, "wifi_db", "wifi_db.machine_mac_srid", "DELETE")
										*/

										bdCntrl := strconv.Itoa(int(bdController))
										var lenMap int
										var count int
										var exception string
										var b2 bytes.Buffer
										var query string

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
											//UploadMapsToDBstring("it_support_db", query)
											UploadMapsToDBerr(query)
										} else {
											fmt.Println("Передана пустая карта. Запрос не выполнен")
										}
										fmt.Println("")
									}
									//
									//
								} else {
									//panic(errGetAnomalies.Error())
									//log.Fatalln("Error:", errGetAnomalies)
									fmt.Println(errGetAnomalies.Error())
									fmt.Println("anomalies НЕ загрузились")
								}
							} // END of ANOMALIES block
						} else {
							//panic(errGetClients.Error())
							//log.Fatalln("Error:", errGetClients)
							fmt.Println(errGetClients.Error())
							fmt.Println("clients НЕ загрузились")
						}
					} else {
						//panic(errGetDevices.Error())
						//log.Fatalln("Error:", errGetDevices)
						fmt.Println(errGetDevices.Error())
						fmt.Println("devices НЕ загрузились")
					}
				} else {
					//panic(errGetSites.Error())
					//log.Fatalln("Error:", errGetSites)
					fmt.Println(errGetSites.Error())
					fmt.Println("sites НЕ загрузились")
				}
			} else {
				//panic(errNewUnifi.Error())
				fmt.Println(errNewUnifi.Error())
				fmt.Println("NewUnifi не загрузился")
				//log.Fatalln("Error:", errNewUnifi)
			}
		} //Снятие показаний раз в 6 минут

		//
		//Обновление мапы site_apcut_login раз в сутки (первичное обновление происходит при старте кода вначале)
		if timeNow.Day() != countDay {
			countDay = timeNow.Day()

			//siteApCutNameLogin = DownloadMapFromDB("wifi_db", "site_apcut", "login", "wifi_db.site_apcut_login", 0, "site_apcut")
			//siteApCutNameLogin = DownloadMapFromDB("it_support_db", "site_apcut", "login", "it_support_db.site_apcut_login", 0, "site_apcut")
			siteApCutNameLogin = DownloadMapFromDBerr()
		}
		//
		//

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

type ApMyStruct struct {
	Name      string
	Exception int //Это исключение НЕ для заявок по Точкам, а для Аномалий!!!
	SrID      string
}
type ForApsTicket struct {
	site          string
	countIncident int
	//apsMac        []string
	apsMacName map[string]string
	//apNames []string //сделано для массовых отключений точек при отключении света в офисе
}

type MachineMyStruct struct {
	Hostname  string
	Exception int
	SrID      string
	ApName    string
}
type ForAnomalyTicket struct {
	site   string
	apName string
	//clientName string  //имя ноутбука будет в ключе мапы, в которую будет встроена эта структура
	noutMac       string //нужен для проверки тикета на открытость
	corpAnomalies []string
}
