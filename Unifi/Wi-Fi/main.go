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

	unifiController := 11 //10-Rostov Local; 11-Rostov ip; 20-Novosib Local; 21-Novosib ip
	var urlController string
	var bdController int8 //Да string, потому что значение пойдёт в replace для БД
	//everyStartCode := map[int]bool{}
	every12start := map[int]bool{}
	//ROSTOV
	if unifiController == 10 || unifiController == 11 {
		bdController = 1
		if unifiController == 10 {
			urlController = "https://localhost:8443/"
		} else {
			urlController = "https://10.78.221.142:8443/"
		}
		/*everyStartCode := [10] int8 {3, 9, 15, 21, 27, 33, 39, 45, 51, 57}
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
		}*/
		every12start = map[int]bool{
			9:  true,
			21: true,
			33: true,
			45: true,
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
		/*everyStartCode := [10] int8 {6, 12, 18, 24, 30, 36, 42, 48, 54, 59}
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
		}*/
		every12start = map[int]bool{
			3:  true,
			15: true,
			27: true,
			39: true,
			51: true,
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

	//count6minute := 0
	count12minute := 0
	//countHourAnom := 0
	countHourDBap := 0
	//countHourDBmachine := 0
	//countDayAnom := 0
	countDayDBmachine := 0
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
		//if currentMinute != 0 && everyStartCode[currentMinute] && currentMinute != count6minute {
		//if timeNow.Minute() != 0 && everyStartCode[timeNow.Minute()] && timeNow.Minute() != count6minute {
		if timeNow.Minute() != 0 && every12start[timeNow.Minute()] && timeNow.Minute() != count12minute {
			//count6minute = timeNow.Minute()
			count12minute = timeNow.Minute()
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

						//Теперь заявки на точки заводятся online сразу в каждый 6 минутный заход. Переносить мапу не нужно
						siteapNameForTickets := map[string]ForApsTicket{}

						for _, ap := range devices.UAPs {
							siteID := ap.SiteID
							if !sitesException[siteID] { // НЕ Резерв/Склад
								apMac := ap.Mac
								apName := ap.Name
								apLastSeen := ap.State.Int()
								//apState := ap.State.Int()

								//_, exisApMacSRid := apMacSRid[ap.Mac]
								_, exisApMyMap := apMyMap[apMac]
								if !exisApMyMap { //если в мапе нет записи, создаём
									apMyMap[apMac] = ApMyStruct{
										apName,
										0,
										"",
										0,
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
											//Идём дальше
											//

											//Точка доступна. Заявка есть. +Имя точки обновляю
										} else if apLastSeen != 0 && srID != "" {
											fmt.Println(apName)
											fmt.Println(apMac)
											fmt.Println("Точка доступна. Заявка есть")
											//Оставляем коммент, ПЫТАЕМСЯ закрыть тикет, если на визировании, Очищаем запись в мапе,

											comment := "Точка появилась в сети: " + apName
											//AddComment(soapServer, srID, comment, bpmUrl)
											//createdOn := AddCommentErr(soapServer, srID, comment, bpmUrl)
											if valueAp.Comment < 1 {
												if AddCommentErr(soapServer, srID, comment, bpmUrl) != "" {
													valueAp.Comment = 1
												}
											}

											//проверить, не последняя ли это запись в мапе в массиве
											countOfIncident := 0
											for _, v := range apMyMap {
												if v.SrID == srID {
													countOfIncident++
													//BREAK здесь НЕ нужен. Пробежаться нужно по всем
												}
											}

											if countOfIncident == 1 {
												//если последняя запись, пробуем закрыть тикет
												//sliceTicketStatus := CheckTicketStatusErr(soapServer, srID)
												//fmt.Println(sliceTicketStatus[1])
												status := CheckTicketStatusErr(soapServer, srID)
												fmt.Println(status)

												//if srStatusCodesForCancelTicket[sliceTicketStatus[1]] {
												if srStatusCodesForCancelTicket[status] {
													//Если статус заявки на Уточнении, Визирование, Назначено

													if valueAp.Comment < 2 {
														comment = "Будет предпринята попытка по отмене обращения, т.к. все точки из него появились в сети"
														if AddCommentErr(soapServer, srID, comment, bpmUrl) != "" {
															valueAp.Comment = 2
														}
													}

													fmt.Println("Попытка изменить статус в На уточнении")
													ChangeStatusErr(soapServer, srID, "На уточнении")

													fmt.Println("Попытка изменить статус в Отменено")
													//ChangeStatusErr(soapServer, srID, "Отменено")
													if ChangeStatusErr(soapServer, srID, "Отменено") != "" {
														//Если отмена заявки прошла успешно, удалить запись из мапы, заодно и имя обновим
														valueAp.Name = apName
														valueAp.SrID = ""
														valueAp.Comment = 0 //также обнулить параметр COMMENT
														apMyMap[keyAp] = valueAp
													}
												}
											} else {
												//Если запись НЕ последняя, только удалить из мапы, заодно и имя обновим
												valueAp.Name = apName
												valueAp.SrID = ""
												//valueAp.comment уже обновлён выше ДО if
												apMyMap[keyAp] = valueAp
											}

											fmt.Println("")

											//
											//Точка недоступна
										} else if apLastSeen == 0 {
											fmt.Println(apName)
											fmt.Println(apMac)
											fmt.Println("Точка НЕ доступна")

											//Проверяем заявку на НЕ закрытость. если заявки нет - ничего страшного
											var status string
											if srID != "" {
												//checkSlice := CheckTicketStatus(soapServer, srID)
												status = CheckTicketStatusErr(soapServer, srID)
											}
											//if srStatusCodesForNewTicket[checkSlice[1]] || !exisApMacSRid {
											//if srStatusCodesForNewTicket[checkSlice[1]] || srID == "" {
											if srStatusCodesForNewTicket[status] || srID == "" {
												fmt.Println(bpmUrl + srID)
												fmt.Println("Статус: " + status) //checkSlice[1])
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
												fmt.Println(status) //checkSlice[1])
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

							var apsNames []string
							//for _, mac := range v.apsMac {
							for _, name := range v.apsMacName {
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

							fmt.Println("Попытка создания заявки по точке")
							//srTicketSlice := CreateSmacWiFiTicketErr(soapServer, bpmUrl, usrLogin, description, v.site, incidentType)
							srTicketSlice := CreateWiFiTicketErr(soapServer, bpmUrl, usrLogin, description, "", v.site, "", incidentType)
							if srTicketSlice[0] != "" {
								fmt.Println(srTicketSlice[2])

								for mac, _ := range v.apsMacName {
									for key, value := range apMyMap {
										if key == mac {
											value.SrID = srTicketSlice[0]
											apMyMap[key] = value
											break
										}
									}
								}
							}
							fmt.Println("")
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
							//if timeNow.Hour() != countHourAnom {
							//if timeNow.Day() != countDayAnom {
							if false {
								//countHourAnom = timeNow.Hour()
								//countDayAnom = timeNow.Day()

								soapServer = soapServerTest
								fmt.Println("SOAP")
								fmt.Println(soapServer)
								bpmUrl = bpmUrlTest
								fmt.Println("BPM")
								fmt.Println(bpmUrl)
								fmt.Println("")

								count := 720 //минус 30 день
								//then := timeNow.Add(time.Duration(-count) * time.Minute)
								then := timeNow.Add(time.Duration(-count) * time.Hour)

								anomalies, errGetAnomalies := uni.GetAnomalies(sites,
									//time.Date(2023, 07, 11, 7, 0, 0, 0, time.Local), time.Now()
									//time.Date(2023, 07, 01, 0, 0, 0, 0, time.Local), //time.Now(),
									then,
								)
								if errGetAnomalies == nil {
									fmt.Println("anomalies загрузились")
									fmt.Println("")

									//mac_dateSite_anom := map[string]DateSite_anom{}
									mac_dateSite_anom := map[string]map[string][]string{}

									var siteName string
									var noutMac string
									var anomalyStr string
									var anomalyDatetime time.Time

									for _, anomaly := range anomalies {
										siteName = anomaly.SiteName
										noutMac = anomaly.DeviceMAC
										anomalyStr = anomaly.Anomaly
										anomalyDatetime = anomaly.Datetime
										//fmt.Println(anomalyDatetime, siteName, noutMac, anomalyStr)

										anomalyDatetime.String()
										dateSite := anomalyDatetime.Format("2006-01-02") + "_" + siteName

										dateSite_anom := map[string][]string{}

										_, exisMac := mac_dateSite_anom[noutMac]
										if !exisMac {
											//Если мака вообще нет, создаём новую мапу внутри мапы
											dateSite_anom[dateSite] = []string{anomalyStr}
											mac_dateSite_anom[noutMac] = dateSite_anom

										} else {
											//Если мак есть
											dateSite_anom = mac_dateSite_anom[noutMac]
											_, exisDateSite := dateSite_anom[dateSite]
											if !exisDateSite {
												//если НЕТ записи с датой дня
												dateSite_anom[dateSite] = []string{anomalyStr}
											} else {
												//если запись дня есть
												sliceAnom := dateSite_anom[dateSite]
												sliceAnom = append(sliceAnom, anomalyStr)
												dateSite_anom[dateSite] = sliceAnom
											}
											mac_dateSite_anom[noutMac] = dateSite_anom
										}
									}

									fmt.Println("")
									fmt.Println("Tele2Corp клиенты с более чем 2 аномалиями:")

									var countIncident int
									var noutName string
									var usrLogin string
									var srID string
									var exception int
									var apName string
									var region string
									incidentType := "Плохое качество соединения клиента"

									var b bytes.Buffer

									for k, v := range mac_dateSite_anom {
										countIncident = len(v)
										if countIncident > 10 {
											for ke, va := range machineMyMap {
												if k == ke {
													noutName = va.Hostname
													fmt.Println(noutName)
													srID = va.SrID
													exception = va.Exception
													apName = va.ApName
													fmt.Println(apName)

													//checkSlice := CheckTicketStatusErr(soapServer, srID)
													//fmt.Println(checkSlice[1])  //статус обращения
													status := CheckTicketStatusErr(soapServer, srID)

													if exception == 0 && srStatusCodesForNewTicket[status] { //checkSlice[1]] { //srID != "" {
														//завести заявку
														usrLogin = GetLoginPCerr(va.Hostname)
														fmt.Println(usrLogin)

														for key, val := range v {
															region = strings.Split(key, "_")[1]
															b.WriteString(key + "\n")
															for _, valu := range val {
																b.WriteString(valu + "\n")
															}
															b.WriteString("\n")
														}

														description := "На ноутбуке:" + "\n" +
															noutName + "\n" + "" + "\n" +
															"За последние 30 дней зафиксировано более 10 Аномалий" + "\n" +
															"" + "\n" +
															"Предполагаемое, но не на 100% точное имя точки:" + "\n" +
															apName + "\n" +
															"" + "\n" +
															"Рекомендации по выполнению таких инцидентов собраны на страничке корпоративной wiki" + "\n" +
															"https://wiki.tele2.ru/display/ITKB/%5BHelpdesk+IT%5D+System+Monitoring" + "\n" +
															"" + "\n" +
															b.String() +
															""

														srTicketSlice := CreateWiFiTicketErr(soapServer, bpmUrl, usrLogin, description, noutName, region, apName, incidentType)
														fmt.Println(srTicketSlice[2])

														va.SrID = srTicketSlice[0]
														machineMyMap[ke] = va

													} else if exception != 0 && srStatusCodesForNewTicket[status] { //checkSlice[1]] {
														fmt.Println("По пользователю или точке выставлено исключение")
														//fmt.Println("Exception = " + )
													} else {
														fmt.Println("Обращение по пользователю уже создано")
														//Как бы добавить аномалии, появившиеся за последние сутки?
														fmt.Println(bpmUrl + srID)
														fmt.Println(status) //checkSlice[1])
													}
													break
												}
											}

										}
									}

									//UploadMapsToDBreplace(machineMacSRid, "wifi_db", "wifi_db.machine_mac_srid", "srid", bdController)
									fmt.Println("")

									//
									//
									//Обновление БД machine раз в час  countDayDBmachine
									//if timeNow.Hour() != countHourDBmachine {
									if timeNow.Day() != countDayDBmachine {
										//countHourDBmachine = timeNow.Hour()
										countDayDBmachine = timeNow.Day()

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
	Comment   int8
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
type DateSite_anom struct {
	dateSite   string
	anom_slice []string
}
type ForAnomalyTicket struct {
	site   string
	apName string
	//clientName string  //имя ноутбука будет в ключе мапы, в которую будет встроена эта структура
	noutMac       string //нужен для проверки тикета на открытость
	corpAnomalies []string
}
