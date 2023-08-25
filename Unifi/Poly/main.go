package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	//"github.com/go-ping/ping"
)

func main() {
	fmt.Println("")
	/*
		every20Code := map[int]bool{
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
		every20Code := map[int]bool{ //6 minutes
			3:  true,
			9:  true,
			15: true,
			21: true,
			33: true,
			39: true,
			45: true,
			51: true,
			57: true,
		}*/
	every20Code := map[int]bool{
		5:  true,
		25: true,
		45: true,
	}

	var soapServer string
	soapServerProd := "http://10.12.15.148/specs/aoi/tele2/bpm/bpmPortType" //PROD
	//soapServerTest := "http://10.246.37.15:8060/specs/aoi/tele2/bpm/bpmPortType" //TEST
	var bpmUrl string
	bpmUrlProd := "https://bpm.tele2.ru/0/Nui/ViewModule.aspx#CardModuleV2/CasePage/edit/"
	//bpmUrlTest := "https://t2ru-tr-tst-01.corp.tele2.ru/0/Nui/ViewModule.aspx#CardModuleV2/CasePage/edit/"

	count20minute := 0
	countHourFromDB := 0
	countHourToDB := 0
	reboot := 0

	srStatusCodesForNewTicket := map[string]bool{
		"Отменено":     true, //Cancel  6e5f4218-f46b-1410-fe9a-0050ba5d6c38
		"Решено":       true, //Resolve  ae7f411e-f46b-1410-009b-0050ba5d6c38
		"Закрыто":      true, //Closed  3e7f420c-f46b-1410-fc9a-0050ba5d6c38
		"На уточнении": true, //Clarification 81e6a1ee-16c1-4661-953e-dde140624fb
		"Тикет введён не корректно": true,
		"": true,
	}
	srStatusCodesForCancelTicket := map[string]bool{
		"Визирование":  true,
		"Назначено":    true,
		"На уточнении": true, //Clarification 81e6a1ee-16c1-4661-953e-dde140624fb
	}

	//Download MAPs from DB
	polyMap := map[string]PolyStruct{} //просто создаю пустую
	//polyMap := DownloadMapFromDBvcsErr()

	//fmt.Println("Вывод мапы СНАРУЖИ функции")
	/*
		for k, v := range siteApCutNameLogin {
			//fmt.Printf("key: %d, value: %t\n", k, v)
			fmt.Println("newMap "+k, v)
		}
		os.Exit(0)
	*/

	fmt.Println("")

	//log.SetOutput(io.Discard) //Отключить вывод лога
	//

	for true { //зацикливаем навечно
		timeNow := time.Now()
		//fmt.Println(timeNow)

		//
		//
		//Обновление мап раз в час. для контроля корректности ip-адресов
		if timeNow.Hour() != countHourFromDB {
			countHourFromDB = timeNow.Hour()

			//polyMap = make(map[string]PolyStruct{})
			polyMap = map[string]PolyStruct{}
			//clear(polyMyMap)
			polyMap = DownloadMapFromDBvcsErr()
		}

		//
		//
		if timeNow.Minute() != 0 && every20Code[timeNow.Minute()] && timeNow.Minute() != count20minute {
			count20minute = timeNow.Minute()
			fmt.Println(timeNow.Format("02 January, 15:04:05"))

			//Опрос каждые 20 минут
			//soapServer = soapServerTest
			soapServer = soapServerProd
			fmt.Println("SOAP")
			fmt.Println(soapServer)
			//bpmUrl = bpmUrlTest
			bpmUrl = bpmUrlProd
			fmt.Println("BPM")
			fmt.Println(bpmUrl)
			fmt.Println("")

			//siteapNameForTickets := map[string]ForPolyTicket{}
			regionVcsSlice := map[string][]PolyStruct{}

			//fmt.Println("Обработка codec устройств")
			//fmt.Println("")
			for k, v := range polyMap {
				if v.Exception == 0 {
					ip := v.IP
					region := v.Region
					roomName := v.RoomName
					login := v.Login
					srID := v.SrID

					/*
						fmt.Println(ip)
						fmt.Println(region)
						fmt.Println(roomName)
						fmt.Println(v.PolyType)
						fmt.Println(srID)
					*/
					//var commentUnreach string
					var statusReach string
					var vcsType string

					if v.PolyType == 1 {
						vcsType = "Codec"
						//commentUnreach = "Codec не отвечает на API-запросы"
						statusReach = apiLineInfo(ip)
					} else {
						vcsType = "Visual"
						//commentUnreach = "Visual не доступен по http"
						statusReach = netDialTmtErr(ip)
					}

					//ВКС доступно. Заявки нет.
					if statusReach != "" && srID == "" {
						//Идём дальше

						//ВКС доступно. Заявка есть
					} else if statusReach != "" && srID != "" {
						fmt.Println(region)
						fmt.Println(roomName)
						fmt.Println(vcsType)
						fmt.Println("ВКС доступно. Заявка есть")
						//Оставляем коммент, ПЫТАЕМСЯ закрыть тикет, если на визировании, Очищаем запись в мапе,

						commentForUpdate := v.Comment
						comment := vcsType + " появился в сети: " + roomName
						if v.Comment < 1 {
							if AddCommentErr(soapServer, srID, comment, bpmUrl) != "" {
								commentForUpdate = 1
							}
						}

						//проверить, не последняя ли это запись в мапе в массиве
						countOfIncident := 0
						for _, va := range polyMap {
							if va.SrID == srID {
								countOfIncident++
								//BREAK здесь НЕ нужен. Пробежаться нужно по всем
							}
						}

						if countOfIncident == 1 {
							//если последняя запись, пробуем закрыть тикет
							statusTicket := CheckTicketStatusErr(soapServer, srID)
							fmt.Println(statusTicket)

							if srStatusCodesForCancelTicket[statusTicket] {
								//Если статус заявки на Уточнении, Визирование, Назначено
								if v.Comment < 2 {
									comment = "Будет предпринята попытка по отмене обращения, т.к. все точки из него появились в сети"
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
									//valueAp.Name = apName
									v.SrID = ""
									v.Comment = 0 //также обнулить параметр COMMENT
									polyMap[k] = v
								} else {
									//Если НЕ удалось отменить заявку
									//valueAp.Name = apName
									//valueAp.SrID не зануляем, т.к. будет второй заход через 12 минут
									v.Comment = commentForUpdate
									polyMap[k] = v
								}
							} else {
								//Если статус заявки В работе, Решено, Закрыто и т.д.
								//valueAp.Name = apName
								v.SrID = ""
								v.Comment = 0
								polyMap[k] = v
							}
						} else {
							//Если запись НЕ последняя, только удалить из мапы sr и comment, заодно и имя обновим
							//valueAp.Name = apName
							v.SrID = ""
							v.Comment = 0
							polyMap[k] = v
						}

						fmt.Println("")

						//ВКС недоступна
					} else if statusReach == "" {
						fmt.Println(region)
						fmt.Println(roomName)
						fmt.Println(vcsType)
						fmt.Println("ВКС Недоступно")

						//Проверяем заявку на НЕ закрытость. если заявки нет - ничего страшного
						var statusTicket string
						if srID != "" {
							statusTicket = CheckTicketStatusErr(soapServer, srID)
						}

						if srStatusCodesForNewTicket[statusTicket] || srID == "" {
							fmt.Println(bpmUrl + srID)
							fmt.Println("Статус: " + statusTicket) //checkSlice[1])
							fmt.Println("Заявка Закрыта, Отменена, Отклонена ИЛИ её нет вовсе")

							//удаляем заявку
							//valueAp.Name = apName
							v.SrID = ""
							polyMap[k] = v

							//Заполняем переменные, которые понадобятся дальше
							fmt.Println(k)
							fmt.Println(login)

							//Проверяем и вносим во временную мапу. Заявка на данном этапе никакая ещё НЕ создаётся
							//_, exisSiteName := siteapNameForTickets[siteApCutName] //проверяем, есть ли в мапе ДЛЯтикетов
							_, exisRegion := regionVcsSlice[region] //проверяем, есть ли в мапе ДЛЯтикетов

							//если в мапе дляТикета сайта ещё НЕТ
							if !exisRegion {
								fmt.Println("в мапе для Тикета записи ещё НЕТ")
								newPolySlice := []PolyStruct{}
								newPolySlice = append(newPolySlice, v)
								regionVcsSlice[region] = newPolySlice

								//если в мапе дляТикета сайт уже есть, добавляем в массив точку
							} else {
								fmt.Println("в мапе для Тикета запись ЕСТЬ")
								//в мапе нельзя просто изменить значение.
								for ke, va := range regionVcsSlice {
									if ke == region {
										//https://stackoverflow.com/questions/42716852/how-to-update-map-values-in-go
										//2.Reassigning the modified struct.
										va = append(va, v)
										regionVcsSlice[ke] = va
										break
									}
								}
							}
						} else {
							fmt.Println("Созданное обращение:")
							fmt.Println(bpmUrl + srID)
							fmt.Println(statusTicket) //checkSlice[1])
						}
						fmt.Println("")
					}
				}
				fmt.Println("")
			} //for

			//
			//
			//Пробежались по всем vcs. Заводим заявки
			fmt.Println("")
			fmt.Println("Создание заявок по ВКС:")
			for k, v := range regionVcsSlice {
				fmt.Println(k)

				var vcsInfo []string
				var usrLogin string

				for _, vcs := range v {
					//apsNames = append(apsNames, name)
					vcsInfo = append(vcsInfo, vcs.RoomName)
					vcsInfo = append(vcsInfo, vcs.IP)
					if vcs.PolyType == 1 {
						vcsInfo = append(vcsInfo, "Codec не отвечает на API-запросы")
					} else {
						vcsInfo = append(vcsInfo, "Visual недоступен по http")
					}
					vcsInfo = append(vcsInfo, "")
					usrLogin = vcs.Login

					fmt.Println(vcs.RoomName)
					fmt.Println(vcs.IP)
					fmt.Println(vcs.PolyType)
				}

				if usrLogin == "" {
					usrLogin = "denis.tirskikh"
				}
				fmt.Println(usrLogin)

				//desAps := strings.Join(apsNames, "\n")
				desVcs := strings.Join(vcsInfo, "\n")
				description := "Зафиксировано отключение устройств ВКС Poly:" + "\n" +
					desVcs + "\n" +
					"" + "\n" +
					"Рекомендации по выполнению таких инцидентов собраны на страничке корпоративной wiki" + "\n" +
					"https://wiki.tele2.ru/display/ITKB/%5BHelpdesk+IT%5D+System+Monitoring" + "\n" +
					"" + "\n" +
					"!!! Не нужно решать/отменять/отклонять/возвращать/закрывать заявку, пока работа всех ВКС устройств не будет восстановлена - автоматически создастся новый тикет !!!" + "\n" +
					""
				incidentType := "Устройство недоступно"
				//monitoring := "https://monitoring.tele2.ru/zabbix1/zabbix.php?show=1&application=&name=&inventory%5B0%5D%5Bfield%5D=type&inventory%5B0%5D%5Bvalue%5D=&evaltype=0&tags%5B0%5D%5Btag%5D=&tags%5B0%5D%5Boperator%5D=0&tags%5B0%5D%5Bvalue%5D=&show_tags=3&tag_name_format=0&tag_priority=&show_opdata=0&show_timeline=1&filter_name=&filter_show_counter=0&filter_custom_time=0&sort=clock&sortorder=DESC&age_state=0&show_suppressed=0&unacknowledged=0&compact_view=0&details=0&highlight_row=0&action=problem.view&groupids%5B%5D=163&hostids%5B%5D=11224&hostids%5B%5D=11381"
				monitoring := "https://r.tele2.ru/aV4MBGZ"

				fmt.Println("Попытка создания заявки")
				srTicketSlice := CreatePolyTicketErr(soapServer, bpmUrl, usrLogin, description, "", k, monitoring, incidentType)
				if srTicketSlice[0] != "" {
					fmt.Println(srTicketSlice[2])
					//delete(regionVcsSlice, k)  //думаю, что удалять не стоит, т.к. будет каждый раз новая мапа создаваться

					//обновляем в мапе srid
					for _, va := range v {
						for key, val := range polyMap {
							if va.IP == val.IP {
								val.SrID = srTicketSlice[0]
								polyMap[key] = val
								break
							}
						}
					}
				}
				fmt.Println("")
			}
			fmt.Println("")
			//
			//

			//
			//
			//Обновление БД. Запустится, если в этот ЧАС он ещё НЕ выполнялся
			if timeNow.Hour() != countHourToDB {
				countHourToDB = timeNow.Hour()

				var queries []string
				for k, v := range polyMap {
					queries = append(queries, "UPDATE it_support_db.poly SET srid = '"+v.SrID+"', comment = "+strconv.Itoa(int(v.Comment))+" WHERE mac = '"+k+"';")
				}
				UpdateMapsToDBerr(queries)
				fmt.Println("")
			}

			//
			//
			//Перезагрузка
			if timeNow.Hour() == 7 && reboot == 0 {
				for _, v := range polyMap {
					if v.PolyType == 1 {
						fmt.Println(v.RoomName)
						apiSafeRestart2(v.IP)
					}
				}
				reboot = 1
			}
			if timeNow.Hour() == 8 {
				reboot = 0
			}

		} //Снятие показаний раз в 20 минут
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

type PolyStruct struct {
	IP        string
	Region    string
	RoomName  string
	Login     string
	SrID      string
	PolyType  int
	Comment   int
	Exception int
}

/*
type ForPolyTicket struct {
	site          string
	countIncident int
	apsMacName    map[string]string
}*/
