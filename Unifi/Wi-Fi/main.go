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

func main34564() {
	fmt.Println("")

	unifiController := 21 //10-Rostov Local; 11-Rostov ip; 20-Novosib Local; 21-Novosib ip
	var urlController string
	var bdController int8 //Да string, потому что значение пойдёт в replace для БД
	every12start := map[int]bool{}
	//every20start := map[int]bool{}

	//ROSTOV
	if unifiController == 10 || unifiController == 11 {
		bdController = 1
		if unifiController == 10 {
			urlController = "https://localhost:8443/"
		} else {
			urlController = "https://10.78.221.142:8443/"
		}
		//
		every12start = map[int]bool{
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
		} /*
			every12start = map[int]bool{
				9:  true,
				21: true,
				33: true,
				45: true,
				57: true,
			}
				every20start = map[int]bool{
					5:  true,
					25: true,
					45: true,
				}*/

		//NOVOSIB
	} else if unifiController == 20 || unifiController == 21 {
		//else{
		bdController = 2
		if unifiController == 20 {
			urlController = "https://localhost:8443/"
		} else {
			urlController = "https://10.8.176.8:8443/"
		}
		//
		every12start = map[int]bool{
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
		} /*
			every12start = map[int]bool{
				3:  true,
				15: true,
				27: true,
				39: true,
				51: true,
			}
				every20start = map[int]bool{
					15: true,
					35: true,
					55: true,
				}*/
	}
	fmt.Println("Unifi controller")
	fmt.Println(urlController)

	var soapServer string
	soapServerProd := "http://10.12.15.148/specs/aoi/tele2/bpm/bpmPortType"      //PROD
	soapServerTest := "http://10.246.37.15:8060/specs/aoi/tele2/bpm/bpmPortType" //TEST
	var bpmUrl string
	bpmUrlProd := "https://bpm.tele2.ru/0/Nui/ViewModule.aspx#CardModuleV2/CasePage/edit/"
	bpmUrlTest := "https://t2ru-tr-tst-01.corp.tele2.ru/0/Nui/ViewModule.aspx#CardModuleV2/CasePage/edit/"

	count12minute := 0
	//count20minute := 0
	countHourAnom := 0
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
		"5f2285f3a1a7693ae6139c00": true, //Novosib. Резерв/Склад
		"5f5b49d1a9f6167b55119c9b": true, //Ростов. Резерв/Склад
		//"Закрыто":      true, //Closed  3e7f420c-f46b-1410-fc9a-0050ba5d6c38
		//"На уточнении": true, //Clarification 81e6a1ee-16c1-4661-953e-dde140624fb
	}

	//Download MAPs from DB
	//apMyMap := map[string]ApMyStruct{}
	//apMyMap := DownloadMapFromDBaps(bdController)
	apMyMap := DownloadMapFromDBapsErr(bdController)

	//НЕ должна создаваться новая раз в 12 минут
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

	for true { //зацикливаем навечно
		//currentMinute := time.Now().Minute()
		timeNow := time.Now()

		//Снятие показаний с контроллера раз в 12 минут. Промежутки разные для контроллеров
		if timeNow.Minute() != 0 && every12start[timeNow.Minute()] && timeNow.Minute() != count12minute {
			//if timeNow.Minute() != 0 && every20start[timeNow.Minute()] && timeNow.Minute() != count20minute {
			count12minute = timeNow.Minute()
			//count20minute = timeNow.Minute()

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
											//Пытаемся удалить запись и в мапе ДляТикета, если она там начала создаваться
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
											//Оставляем коммент, ПЫТАЕМСЯ закрыть тикет, если на визировании, Очищаем запись в мапе,

											commentForUpdate := valueAp.Comment
											comment := "Точка появилась в сети: " + apName
											if valueAp.Comment < 1 {
												if AddCommentErr(soapServer, srID, comment, bpmUrl) != "" {
													commentForUpdate = 1
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
												status := CheckTicketStatusErr(soapServer, srID)
												fmt.Println(status)

												if srStatusCodesForCancelTicket[status] {
													//Если статус заявки на Уточнении, Визирование, Назначено
													if valueAp.Comment < 2 {
														comment = "Будет предпринята попытка отмены обращения, т.к. все точки из него появились в сети"
														if AddCommentErr(soapServer, srID, comment, bpmUrl) != "" {
															commentForUpdate = 2
														}
													}

													fmt.Println("Попытка изменить статус в На уточнении")
													ChangeStatusErr(soapServer, srID, "На уточнении")
													//if error не делаю, т.к. лишним не будет при любом раскладе попытаться вернуть на уточнение

													fmt.Println("Попытка изменить статус в Отменено")
													if ChangeStatusErr(soapServer, srID, "Отменено") != "" {
														//Если отмена заявки прошла успешно, удалить запись из мапы, заодно и имя обновим
														valueAp.Name = apName
														valueAp.SrID = ""
														valueAp.Comment = 0 //также обнулить параметр COMMENT
														apMyMap[keyAp] = valueAp
													} else {
														//Если НЕ удалось отменить заявку
														valueAp.Name = apName
														//valueAp.SrID не зануляем, т.к. будет второй заход через 12 минут
														valueAp.Comment = commentForUpdate
														apMyMap[keyAp] = valueAp
													}
												} else {
													//Если статус заявки В работе, Решено, Закрыто и т.д.
													valueAp.Name = apName
													valueAp.SrID = ""
													valueAp.Comment = 0
													apMyMap[keyAp] = valueAp
												}
											} else {
												//Если запись НЕ последняя, только удалить из мапы sr и comment, заодно и имя обновим
												valueAp.Name = apName
												valueAp.SrID = ""
												valueAp.Comment = 0
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
											//checkSlice := CheckTicketStatus(soapServer, srID)
											//checkSlice := CheckTicketStatusErr(soapServer, srID)
											var status string
											if srID != "" {
												//checkSlice := CheckTicketStatus(soapServer, srID)
												status = CheckTicketStatusErr(soapServer, srID)
											}

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
							if len(v.apsMacName) > 0 {
								fmt.Println(k)
								fmt.Println(v.countIncident) //"Число циклов захода на создание заявки: " +
								v.countIncident++

								if v.countIncident < 2 {
									//обновляем мапу и инкрементируем count
									siteapNameForTickets[k] = v
								} else {
									//Если count == 2, Создаём заявку
									var apsNames []string
									for _, name := range v.apsMacName {
										apsNames = append(apsNames, name)
										fmt.Println(name)
									}

									usrLogin := siteApCutNameLogin[k]
									if usrLogin == "" {
										usrLogin = "denis.tirskikh"
									}
									fmt.Println(usrLogin)

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

									srTicketSlice := CreateSmacWiFiTicketErr(soapServer, bpmUrl, usrLogin, description, v.site, incidentType)
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

							var apName string
							var clientMac string
							var clientName string
							//var clientIP string
							//var siteName string
							//region_guestClients := map[string][]ForGuestClientTicket{}

							for _, client := range clients {
								apName = client.ApName //НИЧЕГО не выводит и не содержит. Имя точки берётся ниже на основании сравнения мапой точек
								clientMac = client.Mac
								clientName = client.Name
								//clientIP = client.IP
								//siteName = client.SiteName

								if !client.IsGuest.Val {
									var clExInt int
									if client.Noted.Val {
										clientExceptionStr := strings.Split(client.Note, " ")[0]
										if clientExceptionStr == "Exception" {
											clExInt = 1
										} else {
											clExInt = 0
										}
									}
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
								} /* До будущих времён, когда буду обрабатывать Клиентов
								else {
									//Если клиент Guest
									splitIP := strings.Split(clientIP, ".")[0]
									if splitIP == "169" {
										forGuestClientTicket := ForGuestClientTicket{
											clientMac,
											clientName,
											clientIP,
										}

										//Заносим в мапу для заявки
										_, exisRegion := region_guestClients[region]
										if exisRegion {
											for k, v := range region_guestClients {
												if k == region {
													v = append(v, forGuestClientTicket)
													region_guestClients[k] = v
													break
												}
											}
										} else {
											forGuestClientTicketSlice := []ForGuestClientTicket{
												forGuestClientTicket,
											}
											region_guestClients[region] = forGuestClientTicketSlice
										}
									}
								}*/
							}
							//Пробежались по всем клиентам. Заводим заявки по Guest
							//
							//

							//
							//
							//АНОМАЛИИ. Corp
							if timeNow.Hour() != countHourAnom {
								//if timeNow.Day() != countDayAnom {
								//if false {
								countHourAnom = timeNow.Hour()
								//countDayAnom = timeNow.Day()

								soapServer = soapServerTest
								fmt.Println("SOAP")
								fmt.Println(soapServer)
								bpmUrl = bpmUrlTest
								fmt.Println("BPM")
								fmt.Println(bpmUrl)
								fmt.Println("")

								//count := 720 //минус 30 день
								count := 1 //минус 1 час
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
									//mac_dateSite_anom := map[string]map[string][]string{}
									mac_DateSiteAnom := map[string]DateSiteAnom{}

									var siteName string
									var noutMac string
									var anomalyStr string
									//var anomalyDatetime time.Time
									var anomalyDatetimeMySQL string
									var dateSiteAnom DateSiteAnom

									for _, v := range anomalies {
										noutMac = v.DeviceMAC
										_, exisMac1 := machineMyMap[noutMac]
										if exisMac1 {
											siteName = v.SiteName
											anomalyStr = v.Anomaly
											anomalyDatetimeMySQL = v.Datetime.Format("2006-01-02 15:04:05")

											_, exisMac2 := mac_DateSiteAnom[noutMac]
											if !exisMac2 {
												//Если мака вообще нет, создаём новую мапу внутри мапы
												anomSlice := []string{anomalyStr}
												mac_DateSiteAnom[noutMac] = DateSiteAnom{
													noutMac,
													anomalyDatetimeMySQL,
													siteName,
													anomSlice,
												}
											} else {
												//Если мак есть
												dateSiteAnom = mac_DateSiteAnom[noutMac]
												dateSiteAnom.AnomSlice = append(dateSiteAnom.AnomSlice, anomalyStr)
												mac_DateSiteAnom[noutMac] = dateSiteAnom
											}
										}
									} // k,v

									//
									//
									fmt.Println("")
									fmt.Println("Ежечасовое занесение аномалий в БД")

									bdCntrl := strconv.Itoa(int(bdController))
									var anomSliceString string
									var query string
									var b1 bytes.Buffer
									b1.WriteString("INSERT INTO it_support_db.anomalies VALUES ")
									lenMap := len(mac_DateSiteAnom)
									var lenSlice int
									var siteNameCut string
									countB1 := 0

									for k, v := range mac_DateSiteAnom {
										lenSlice = len(v.AnomSlice)
										if lenSlice > 1 {
											countB1++
											siteNameCut = v.SiteName[:len(v.SiteName)-11]
											anomSliceString = strings.Join(v.AnomSlice, ";")
											b1.WriteString("('" + v.DateTime + "','" + k + "','" + bdCntrl + "','" + siteNameCut + "','" + anomSliceString + "'),")
											/*
												if countB1 != lenMap {
													b1.WriteString("('" + v.datetime + "','" + k + "','" + bdCntrl + "','" + siteNameCut + "','" + anomSliceString + "'),")
												} else {
													b1.WriteString("('" + v.datetime + "','" + k + "','" + bdCntrl + "','" + siteNameCut + "','" + anomSliceString + "')")
													//в конце НЕ ставим запятую
												}*/
										}
									}
									query = b1.String()
									//strings.TrimSuffix(query, ",")
									//Возможно, не самый энергоэффективный метод обрезать строку с конца, но рабочий
									if last := len(query) - 1; last >= 0 && query[last] == ',' {
										query = query[:last]
									}
									fmt.Println(query)
									if countB1 != 0 {
										//UploadMapsToDBerr(query)
									} else {
										fmt.Println("Передана пустая карта. Запрос не выполнен")
									}
									fmt.Println("")
									//
									//

									//
									//
									//Создание заявок по машинам раз в сутки
									//if timeNow.Hour() != countHourDBmachine {
									if timeNow.Day() != countDayDBmachine {
										//countHourDBmachine = timeNow.Hour()
										countDayDBmachine = timeNow.Day()

										//Загружаем из БД аномалии за последние 30 дн. в массив структур DateSiteAnom
										before30days := timeNow.Add(time.Duration(-720) * time.Hour).Format("2006-01-02 15:04:05")
										//before30days := timeNow.Add(time.Duration(-3) * time.Hour).Format("2006-01-02 15:04:05")

										//macDay_DateSiteAnom := map[string]DateSiteAnom{}
										macDay_DateSiteAnom := DownloadMapFromDBanomaliesErr(bdController, before30days)

										mac_DateSiteAnomSlice := map[string][]DateSiteAnom{}
										dateSiteAnomSlice := []DateSiteAnom{}

										//Обрабатываем
										for _, v := range macDay_DateSiteAnom {
											//fmt.Println(k, v.DateTime, v.Mac)
											_, exis := mac_DateSiteAnomSlice[v.Mac]
											if !exis {
												dateSiteAnomSlice = []DateSiteAnom{v}
												mac_DateSiteAnomSlice[v.Mac] = dateSiteAnomSlice
											} else {
												dateSiteAnomSlice = mac_DateSiteAnomSlice[v.Mac]
												dateSiteAnomSlice = append(dateSiteAnomSlice, v)
												mac_DateSiteAnomSlice[v.Mac] = dateSiteAnomSlice
											}
										}

										//Создаём заявки
										//var countIncident int
										var noutName string
										var usrLogin string
										var srID string
										var exceptionInt int
										//var apName string
										var region string
										incidentType := "Плохое качество соединения клиента"
										var b2 bytes.Buffer

										for k, v := range mac_DateSiteAnomSlice {
											if len(v) > 9 {
												_, exis := machineMyMap[k]
												if exis {
													for ke, va := range machineMyMap {
														if k == ke {
															noutName = va.Hostname
															fmt.Println(noutName)
															srID = va.SrID
															exceptionInt = va.Exception
															apName = va.ApName
															fmt.Println(apName)

															var statusTicket string
															if srID != "" {
																statusTicket = CheckTicketStatusErr(soapServer, srID)
															}
															if exceptionInt == 0 && (srStatusCodesForNewTicket[statusTicket] || srID == "") {
																//Если заявки ещё нет, либо закрыта отменена
																usrLogin = GetLoginPCerr(va.Hostname)
																fmt.Println(usrLogin)

																for _, val := range v {
																	//region = val.SiteName
																	b2.WriteString(val.SiteName + "\n")
																	b2.WriteString(val.DateTime + "\n")
																	for _, valu := range val.AnomSlice {
																		b2.WriteString(valu + "\n")
																	}
																	b2.WriteString("\n")
																}

																description := "На ноутбуке:" + "\n" +
																	noutName + "\n" + "" + "\n" +
																	"За последние 30 дней зафиксировано более 10 дней с Аномалиями качества работы Wi-Fi сети Tele2Corp" + "\n" +
																	"" + "\n" +
																	"Предполагаемое, но не на 100% точное имя точки:" + "\n" +
																	apName + "\n" +
																	"" + "\n" +
																	"Рекомендации по выполнению таких инцидентов собраны на страничке корпоративной wiki" + "\n" +
																	"https://wiki.tele2.ru/display/ITKB/%5BHelpdesk+IT%5D+System+Monitoring" + "\n" +
																	"" + "\n" +
																	b2.String() +
																	""

																fmt.Println("Попытка создания заявки")
																srTicketSlice := CreateWiFiTicketErr(soapServer, bpmUrl, usrLogin, description, noutName, region, apName, incidentType)

																if srTicketSlice[0] != "" {
																	fmt.Println(srTicketSlice[2])

																	va.SrID = srTicketSlice[0]
																	machineMyMap[ke] = va
																}
															} else if exceptionInt > 0 {
																fmt.Println("Клиент добавлен в исключение")
															} else {
																//либо заявка уже есть. добавить коммент?
																fmt.Println("Созданное обращение:")
																fmt.Println(bpmUrl + srID)
																fmt.Println(statusTicket)
															}

															break
														}
													}
												} else {
													fmt.Println("Не удалось найти запись по маку в мапе машин. Создать заявку невозможно")
												}

											}
										}

										//
										//
										fmt.Println("")
										fmt.Println("Ежесуточное занесение машин в БД")
										bdCntrl = strconv.Itoa(int(bdController))
										//var lenMap int
										//var count int
										var exceptionStr string
										var b3 bytes.Buffer
										//var query string

										//b.WriteString("REPLACE INTO " + tableName + " VALUES ")
										b3.WriteString("REPLACE INTO " + "it_support_db.machine" + " VALUES ")
										//lenMap := len(uploadMap)
										lenMap = len(machineMyMap)
										count = 0
										//for k, v := range uploadMap {
										for k, v := range machineMyMap {
											exceptionStr = strconv.Itoa(int(v.Exception))
											count++
											if count != lenMap {
												// mac, hostname, controller, exception, srid, apname
												b3.WriteString("('" + k + "','" + v.Hostname + "','" + bdCntrl + "','" + exceptionStr + "','" + v.SrID + "','" + v.ApName + "'),")
											} else {
												b3.WriteString("('" + k + "','" + v.Hostname + "','" + bdCntrl + "','" + exceptionStr + "','" + v.SrID + "','" + v.ApName + "')")
												//в конце НЕ ставим запятую
											}
										}
										query = b3.String()
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
				//log.Fatalln("Error:", errNewUnifi)
				fmt.Println(errNewUnifi.Error())
				fmt.Println("NewUnifi не загрузился")
			}
		} //Снятие показаний раз в 12 минут

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

		fmt.Println("Sleep 58s")
		fmt.Println("")
		time.Sleep(58 * time.Second)

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

type GuestClient struct {
}
type ForGuestClientTicket struct {
	mac      string
	hostname string
	ip       string
}

type MachineMyStruct struct {
	Hostname  string
	Exception int
	SrID      string
	ApName    string
}
type DateSiteAnom struct {
	Mac       string
	DateTime  string
	SiteName  string
	AnomSlice []string
}
type ForAnomalyTicket struct {
	site   string
	apName string
	//clientName string  //имя ноутбука будет в ключе мапы, в которую будет встроена эта структура
	noutMac       string //нужен для проверки тикета на открытость
	corpAnomalies []string
}
