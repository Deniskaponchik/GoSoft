package main

import (
	"bytes"
	"fmt"
	"github.com/unpoller/unifi"
	"io"
	"log"
	"strconv"
	"strings"
	"time"
)

func main() {
	fmt.Println("")

	var url string

	everyStartCode := map[int]bool{}
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

	var soapServer string
	soapServerProd := "http://10.12.15.148/specs/aoi/tele2/bpm/bpmPortType"      //PROD
	soapServerTest := "http://10.246.37.15:8060/specs/aoi/tele2/bpm/bpmPortType" //TEST
	var bpmUrl string
	bpmUrlProd := "https://bpm.tele2.ru/0/Nui/ViewModule.aspx#CardModuleV2/CasePage/edit/"
	bpmUrlTest := "https://t2ru-tr-tst-01.corp.tele2.ru/0/Nui/ViewModule.aspx#CardModuleV2/CasePage/edit/"

	count6minute := 0
	countHourAnom := 0
	countHourDB := 0
	countDay := time.Now().Day()

	srStatusCodesForNewTicket := map[string]bool{
		"Отменено":     true, //Cancel  6e5f4218-f46b-1410-fe9a-0050ba5d6c38
		"Решено":       true, //Resolve  ae7f411e-f46b-1410-009b-0050ba5d6c38
		"Закрыто":      true, //Closed  3e7f420c-f46b-1410-fc9a-0050ba5d6c38
		"На уточнении": true, //Clarification 81e6a1ee-16c1-4661-953e-dde140624fb
		"Тикет введён не корректно": true,
		"": true,
	}

	//Download MAPs from DB
	//polyMyMap := map[string]PolyStruct{}
	poly1map := DownloadMapFromDBaps(1)
	poly2map := DownloadMapFromDBaps(2)

	siteapNameForTickets := map[string]ForApsTicket{} //
	//fmt.Println("Вывод мапы СНАРУЖИ функции")
	/*
		for k, v := range siteApCutNameLogin {
			//fmt.Printf("key: %d, value: %t\n", k, v)
			fmt.Println("newMap "+k, v)
		}
		os.Exit(0)
	*/

	fmt.Println("")

	log.SetOutput(io.Discard) //Отключить вывод лога
	//
	//

	for true { //зацикливаем навечно
		timeNow := time.Now()

		if timeNow.Minute() != 0 && everyStartCode[timeNow.Minute()] && timeNow.Minute() != count6minute {
			count6minute = timeNow.Minute()
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

			devices, err := uni.GetDevices(sites) //devices = APs
			if err != nil {
				log.Fatalln("Error:", err)
			} else {
				fmt.Println("devices загрузились")
			}
			fmt.Println("")
			//
			//

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

			//if time.Now().Minute()%3 == 0 && time.Now().Minute() != count3minute { //запускается раз в 3 минуты
			fmt.Println("Обработка точек доступа...")
			fmt.Println("")
			for _, ap := range devices.UAPs {
				siteID := ap.SiteID
				//fmt.Println(ap.Name)	fmt.Println(ap.SiteID)
				//if ap.SiteName[:len(ap.SiteName)-11] != "Резерв/Склад" {
				//if ap.SiteID != "5f2285f3a1a7693ae6139c00" { //NOVOSIB
				if !sitesException[siteID] { // НЕ Резерв/Склад

					apMac := ap.Mac
					apName := ap.Name
					apLastSeen := ap.Uptime.Int()

					//fmt.Println(ap.Name)	fmt.Println(ap.Uptime.Int())  fmt.Println(ap.Uptime.String()) fmt.Println(ap.Uptime.Val)

					//_, exisApMacSRid := apMacSRid[ap.Mac]
					_, exisApMyMap := apMyMap[apMac]
					//если в мапе нет записи, создаём
					if !exisApMyMap {
						apMyMap[apMac] = ApMyStruct{
							apName,
							0,
							"",
						}
					}
					/*else {  //случай, когда запись в мапе уже создана, думаю, можно упустить в пользу того,
								//чтобы обновить данные с чем-нибудь другим далее ниже. Только ВСПОМНИТЬ ОБНОВИТь!!!
						for k, v := range apMyMap {
							if k == ap.Mac {
								v.name = ap.Name
								apMyMap[k] = v
								break
							}
						}
					}*/

					for keyAp, valueAp := range apMyMap {
						if keyAp == apMac {
							srID := valueAp.SrID

							//Точка доступна. Заявки нет.
							//if apLastSeen != 0 && !exisApMy {
							if apLastSeen != 0 && srID == "" {
								//fmt.Println("Точка доступна. Заявки нет. Идём дальше")
								//

								//Точка доступна. Заявка есть.   +Имя точки обновляю
								//} else if apLastSeen != 0 && exisApMacSRid {
							} else if apLastSeen != 0 && srID != "" {
								fmt.Println(apName)
								fmt.Println(apMac)
								fmt.Println("Точка доступна. Заявка есть")
								//Оставляем коммент, Очищаем запись в мапе, ПЫТАЕМСЯ закрыть тикет, если на визировании

								//comment := "Точка появилась в сети: " + ap.Name
								comment := "Точка появилась в сети: " + apName
								//AddComment(soapServer, apMacSRid[ap.Mac], comment, bpmUrl)
								AddComment(soapServer, srID, comment, bpmUrl)

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
								/*Старый блок
								for _, v := range apMacSRid {
									if v == srID {
										countOfIncident++
									}
								}*/
								for _, v := range apMyMap {
									if v.SrID == srID {
										countOfIncident++
										//fmt.Println(countOfIncident)
									}
								}
								if countOfIncident == 0 {
									//Пробуем закрыть тикет, только ЕСЛИ он на Визировании
									//fmt.Println("Попали в блок изменения статусов заявок")
									//sliceTicketStatus := CheckTicketStatus(soapServer, apMacSRid[ap.Mac]) //получаем статус
									sliceTicketStatus := CheckTicketStatus(soapServer, srID) //получаем статус
									fmt.Println(sliceTicketStatus[1])
									if sliceTicketStatus[1] == "Визирование" {
										//Если статус заявки по-прежнему на Визировании
										ChangeStatus(soapServer, srID, "На уточнении")
										AddComment(soapServer, srID, "Обращение отменено, т.к. все точки из него появились в сети", bpmUrl)
										ChangeStatus(soapServer, srID, "Отменено")
									}
								}
								fmt.Println("")

								//
								//Точка недоступна.
							} else if apLastSeen == 0 {
								apSiteName := ap.SiteName
								fmt.Println(apName)
								fmt.Println(apMac)
								fmt.Println("Точка НЕ доступна")

								//Проверяем заявку на НЕ закрытость. если заявки нет - ничего страшного
								//checkSlice := CheckTicketStatus(soapServer, apMacSRid[ap.Mac])
								checkSlice := CheckTicketStatus(soapServer, srID)

								//if srStatusCodesForNewTicket[checkSlice[1]] || !exisApMacSRid {
								if srStatusCodesForNewTicket[checkSlice[1]] || srID == "" {
									fmt.Println("Заявка Закрыта, Отменена, Отклонена ИЛИ в мапе нет записи")

									//delete(apMacSRid, ap.Mac) //удаляем заявку. если заявки нет - ничего страшного
									//удаляем заявку + обновить имя
									valueAp.Name = apName
									valueAp.SrID = ""
									apMyMap[keyAp] = valueAp

									//Заполняем переменные, которые понадобятся дальше
									//fmt.Println(ap.SiteID)
									fmt.Println("Site ID: " + siteID)
									var siteName string
									//if ap.SiteID == "5e74aaa6a1a76964e770815c" { //6360a823a1a769286dc707f2
									if siteID == "5e74aaa6a1a76964e770815c" { //6360a823a1a769286dc707f2
										siteName = "Урал"
									} else {
										//siteName = ap.SiteName[:len(ap.SiteName)-11]
										siteName = apSiteName[:len(apSiteName)-11]
									}
									//apCutName := ap.Name[:len(ap.Name)3]
									//apCutName := strings.Split(ap.Name, "-")[0]
									apCutName := strings.Split(apName, "-")[0]
									siteApCutName := siteName + "_" + apCutName
									fmt.Println(siteApCutName)

									//Проверяем и Вносим во временную мапу. Заявка на данном этапе никакая ещё НЕ создаётся
									_, exisSiteName := siteapNameForTickets[siteApCutName] //проверяем, есть ли в мапе ДЛЯтикетов
									//for _, ticket := range sliceForTicket {

									//если в мапе дляТикета сайта ещё НЕТ
									if !exisSiteName {
										fmt.Println("в мапе для Тикета записи ещё НЕТ")
										//aps.mac := [string]
										siteapNameForTickets[siteApCutName] = ForApsTicket{
											siteName,
											0,
											//[]string{ap.Mac},
											//apMac[apName],
											map[string]string{apMac: apName},
											//map[apMac]apName,
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
									//fmt.Println(bpmUrl + apMacSRid[ap.Mac])
									fmt.Println(bpmUrl + srID)
									fmt.Println(checkSlice[1])
								}
								fmt.Println("")
							}
						}
					}
				} //fmt.Println("")
			}
			//Пробежались по всем точкам. Заводим заявки
			fmt.Println("")
			fmt.Println("Создание заявок по точкам:")
			for k, v := range siteapNameForTickets {
				//vCountIncident := v.countIncident
				fmt.Println(k)
				fmt.Println(v.countIncident) //"Число циклов захода на создание заявки: " +
				//fmt.Println(vСountIncident) //"Число циклов захода на создание заявки: " +
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
			}
			fmt.Println("")
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

			//
			//
		} //Снятие показаний раз в 6 минут
		fmt.Println("Sleep 58s")
		fmt.Println("")
		time.Sleep(58 * time.Second) //Изменить на 5 секунд на ПРОДе

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

type PolyStruct struct {
	Region   string
	Room     string
	login    string
	srid     string
	type	 int
	countInc int
}
type ForPolyTicket struct {
	site          string
	countIncident int
	apsMacName    map[string]string
}
