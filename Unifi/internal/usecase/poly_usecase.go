package usecase

import (
	//Not have package imports from the outer layer.
	//UseCase ВСЕГДА остаётся чистым. Он ничего не знает про тех, кто его вызывает. Чистота архитектуры - это про UseCase
	"fmt"
	"github.com/deniskaponchik/GoSoft/Unifi/internal/entity"
	"strconv"
	"strings"
	"time"
)

type PolyUseCase struct {
	repo    PolyRepo    //interface
	webAPI  PolyWebApi  //interface
	netDial PolyNetDial //interface
	soap    PolySoap    //interface

	everyCodeMap map[int]bool
	restartHour  int
	timezone     int
}

// реализуем Инъекцию зависимостей DI. Используется в app
func NewPoly(r PolyRepo, a PolyWebApi, n PolyNetDial, s PolySoap, everyCode map[int]bool, restartHour int, timezone int) *PolyUseCase {
	return &PolyUseCase{
		//Мы можем передать сюда ЛЮБОЙ репозиторий (pg, s3 и т.д.) НО КОД НЕ ПОМЕНЯЕТСЯ! В этом смысл DI
		repo:         r,
		webAPI:       a,
		netDial:      n,
		soap:         s,
		everyCodeMap: everyCode,
		restartHour:  restartHour,
		timezone:     timezone,
	}
}

// Переменные, которые используются во всех методах ниже
var polyMap map[string]entity.PolyStruct
var region_VcsSlice map[string][]entity.PolyStruct
var timeNowP time.Time
var errp error
var sleepHoursPoly map[int]bool

func (puc *PolyUseCase) InfinityPolyProcessing() error {

	everyCode := puc.everyCodeMap
	count20minute := 0
	countHourFromDB := 0
	countHourToDB := 0
	reboot := 0

	sleepHoursPoly = map[int]bool{
		20: true,
		21: true,
		22: true,
		23: true,
		0:  true,
		1:  true,
		2:  true,
		3:  true,
		4:  true,
		5:  true,
		6:  true,
	}

	//polyMap = make(map[string]entity.PolyStruct)
	polyMap, errp = puc.repo.DownloadMapFromDBvcsErr(0) // 0 - при старте приложения. код не пойдёт дальше приошибках
	if errp != nil {
		return errp
	}

	fmt.Println("")
	//log.SetOutput(io.Discard) //Отключить вывод лога. Not for Zerolog

	for true { //зацикливаем навечно
		timeNowP = time.Now()

		//Запуск каждые 20 минут
		if timeNowP.Minute() != 0 && everyCode[timeNowP.Minute()] && timeNowP.Minute() != count20minute {
			//if timeNowP.Minute() != 0 && every20Code[timeNowP.Minute()] && timeNowP.Minute() != count20minute {
			count20minute = timeNowP.Minute()
			fmt.Println(timeNowP.Format("02 January, 15:04:05"))

			//Опрос устройств
			errp = puc.Survey() //polyMap) //, cfg.BpmUrl)  region_VcsSlice,
			if errp != nil {
				//l.Info(fmt.Errorf("app - Run - Download polyMap from DB: %w", errSurvey))
				fmt.Errorf("app - Run - Error in Survey: %w", errp)
			} else {
				//Если опрос устройств не завершился ошибкой, то заводим заявки
				errp = puc.TicketsCreating() //(region_VcsSlice)
				if errp != nil {
					fmt.Errorf("app - Run - Error in Ticketing: %w", errp)
				}
			}

			//Обновление БД. Запустится, если в этот ЧАС он ещё НЕ выполнялся
			if timeNowP.Hour() != countHourToDB {
				countHourToDB = timeNowP.Hour()

				//UpdateMapsToDBerr(polyConf.GlpiConnectStringITsupport, queries)
				errp = puc.repo.UpdateMapsToDBerr(polyMap)
				fmt.Println("")
			}

			//Обновление мап раз в час. для контроля корректности ip-адресов
			if timeNowP.Hour() != countHourFromDB {
				countHourFromDB = timeNowP.Hour()

				//polyMap = DownloadMapFromDBvcsErr(polyConf.GlpiConnectStringITsupport)
				polyMap, errp = puc.repo.DownloadMapFromDBvcsErr(1) //1 - код пойдёт дальше при ошибках
				fmt.Println("")
			}

			//Перезагрузка
			if timeNowP.Hour() == puc.restartHour && reboot == 0 { //Было 7 часов
				for _, v := range polyMap {
					if v.PolyType == 1 {
						fmt.Println(v.RoomName)
						//apiSafeRestart2(v.IP, polyConf.PolyUsername, polyConf.PolyPassword)
						errp = puc.webAPI.ApiSafeRestart(v)
					}
				}
				reboot = 1
				time.Sleep(2400 * time.Second) //40 minutes
			}
			if timeNowP.Hour() == puc.restartHour+1 { //стояло 8
				reboot = 0
			}
		} //Снятие показаний раз в 20 минут
		fmt.Println("Sleep 58s")
		fmt.Println("")
		time.Sleep(58 * time.Second)
	} //while true

	return nil
}

// Опрос устройств
func (puc *PolyUseCase) Survey() error {
	//polyMap map[string]entity.PolyStruct) ( //, bpmUrl string) (
	//map[string]entity.PolyStruct, region_VcsSlice map[string][]entity.PolyStruct, errp error) {

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

	//siteapNameForTickets := map[string]ForPolyTicket{}
	//Создавать новую мапу каждые 20 минут:
	region_VcsSlice = map[string][]entity.PolyStruct{}

	//fmt.Println("")
	for k, v := range polyMap { // v == polyStruct
		//polyStruct := &v
		if v.Exception == 0 {

			/*теперь передаю структуру в сервисы, а не текст
			ip := v.IP
			region := v.Region
			roomName := v.RoomName
			login := v.Login
			srID := v.SrID
			*/
			polyTicket := &entity.Ticket{
				ID:        v.SrID,
				UserLogin: v.Login,
				Region:    v.Region,
			}
			//statusReach := v.Status  //нигде не использую пока убираю

			/*
				fmt.Println(ip)
				fmt.Println(region)
				fmt.Println(roomName)
				fmt.Println(v.PolyType)
				fmt.Println(srID)
			*/
			//var commentUnreach string

			//var statusReach string   //нигде не использую пока убираю
			var errGetStatus error
			var vcsType string
			if v.PolyType == 1 {
				vcsType = "Codec"
				//commentUnreach = "Codec не отвечает на API-запросы"
				//statusReach = webapi.apiLineInfo(ip, polyConf.PolyUsername, polyConf.PolyPassword)
				//statusReach, errGetStatus = puc.webAPI.ApiLineInfo(v) //возвращает строку
				//v, errGetStatus = puc.webAPI.ApiLineInfoErr(v) //возвращает структуру
				errGetStatus = puc.webAPI.ApiLineInfoErr(&v) //возвращает структуру
			} else {
				vcsType = "Visual"
				//commentUnreach = "Visual не доступен по netdial"
				//statusReach = netDialTmtErr(ip)
				//v, errGetStatus = puc.netDial.NetDialTmtErr(v) //возвращает структуру
				errGetStatus = puc.netDial.NetDialTmtErr(&v) //возвращает структуру
			}

			//ВКС доступно.
			if errGetStatus == nil { //&& srID == "" {  //&& statusReach == "Registered"  оставляю на будущее
				//if statusReach != "" && srID == "" {
				if v.SrID == "" {
					//Заявки нет. Просто Идём дальше

				} else {
					//Заявка есть
					//} else if errGetStatus == nil && srID != "" { 	//} else if statusReach != "" && srID != "" {

					fmt.Println(v.Region)
					fmt.Println(v.RoomName)
					fmt.Println(vcsType)
					fmt.Println("ВКС доступно. Заявка есть")
					//Оставляем коммент, ПЫТАЕМСЯ закрыть тикет, если на визировании, Очищаем запись в мапе,

					commentForUpdate := v.Comment
					//comment := vcsType + " появился в сети: " + v.RoomName
					polyTicket.Comment = vcsType + " появился в сети: " + v.RoomName
					if v.Comment < 1 {
						//если по данному устройству коммент раньше НЕ оставлялся
						errAddComment := puc.soap.AddCommentErr(polyTicket)
						if errAddComment == nil {
							//if AddCommentErr(soapServer, srID, comment, bpmUrl) != "" {
							commentForUpdate = 1
						}
					}

					//проверить, не последняя ли это запись в мапе в массиве
					countOfIncident := 0
					for _, va := range polyMap {
						if va.SrID == v.SrID { // srID {
							countOfIncident++
							//BREAK здесь НЕ нужен. Пробежаться нужно по всем
						}
					}

					if countOfIncident == 1 {
						//если последняя запись, пробуем закрыть тикет
						var errCheckStatus error
						//statusTicket := CheckTicketStatusErr(soapServer, srID)
						//polyTicket, errCheckStatus = puc.soap.CheckTicketStatusErr(polyTicket)
						errCheckStatus = puc.soap.CheckTicketStatusErr(polyTicket)
						fmt.Println(polyTicket.Status) //statusTicket)

						//Статус заявки удалось получить
						if errCheckStatus == nil {
							//Если статус заявки на Уточнении, Визирование, Назначено
							if srStatusCodesForCancelTicket[polyTicket.Status] { //statusTicket] {
								if v.Comment < 2 {
									polyTicket.Comment = "Будет предпринята попытка по отмене обращения, т.к. все устройства из него появились в сети"
									//var errAddComment error
									//если добавление коммента прошло без ошибок
									if puc.soap.AddCommentErr(polyTicket) == nil {
										//if AddCommentErr(soapServer, srID, comment, bpmUrl) != "" {
										commentForUpdate = 2
									}
								}

								//var errChangeStatus error
								fmt.Println("Попытка изменить статус в На уточнении")
								polyTicket.Status = "На уточнении"

								//ChangeStatusErr(soapServer, srID, "На уточнении")
								errChangeStatus := puc.soap.ChangeStatusErr(polyTicket)
								//if error не делаю, т.к. лишним не будет при любом раскладе попытаться вернуть на уточнение

								fmt.Println("Попытка изменить статус в Отменено")
								polyTicket.Status = "Отменено"
								errChangeStatus = puc.soap.ChangeStatusErr(polyTicket)
								//Если отмена заявки прошла успешно, удалить запись из мапы, заодно и имя обновим
								if errChangeStatus == nil {
									//if ChangeStatusErr(soapServer, srID, "Отменено") != "" {
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
								//Если статус заявки В работе, на 3 линии, Решено, Закрыто
								v.SrID = ""
								v.Comment = 0
								polyMap[k] = v
							}
						} else {
							//Получение статуса обращения завершилось ошибкой
							v.Comment = commentForUpdate
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
				}

			} else {
				//ВКС недоступна			//} else if statusReach == "" {
				fmt.Println(v.Region)
				fmt.Println(v.RoomName)
				fmt.Println(vcsType)
				fmt.Println("ВКС Недоступно")

				//Проверяем заявку на НЕ закрытость. если заявки нет - ничего страшного
				//var errCheckStatus error
				if polyTicket.ID != "" { //srID != "" {
					//statusTicket = CheckTicketStatusErr(soapServer, srID)
					//polyTicket, errCheckStatus = puc.soap.CheckTicketStatusErr(polyTicket)
					errCheckStatus := puc.soap.CheckTicketStatusErr(polyTicket)
					if errCheckStatus != nil {
						fmt.Println("Статус обращения выяснить не удалось. Никаких действий не предпринимаем")
					}
				}

				if srStatusCodesForNewTicket[polyTicket.Status] || polyTicket.ID == "" { //} || errCheckStatus != nil {
					//if srStatusCodesForNewTicket[statusTicket] || srID == "" {
					//fmt.Println(bpmUrl + polyTicket.ID)         //srID)
					fmt.Println(polyTicket.BpmServer + polyTicket.ID)
					fmt.Println("Статус: " + polyTicket.Status) //statusTicket)
					fmt.Println("Заявка Закрыта, Отменена, Отклонена ИЛИ её нет вовсе")

					//удаляем заявку
					//valueAp.Name = apName
					v.SrID = ""
					polyMap[k] = v

					//Заполняем переменные, которые понадобятся дальше
					fmt.Println(k)
					fmt.Println(v.Login) //login)

					//Проверяем и вносим во временную мапу. Заявка на данном этапе никакая ещё НЕ создаётся
					//_, exisSiteName := siteapNameForTickets[siteApCutName] //проверяем, есть ли в мапе ДЛЯтикетов
					_, exisRegion := region_VcsSlice[v.Region] //проверяем, есть ли в мапе ДЛЯтикетов

					//если в мапе дляТикета сайта ещё НЕТ
					if !exisRegion {
						fmt.Println("в мапе для Тикета записи ещё НЕТ")
						newPolySlice := []entity.PolyStruct{}
						newPolySlice = append(newPolySlice, v)
						region_VcsSlice[v.Region] = newPolySlice

						//если в мапе дляТикета сайт уже есть, добавляем в массив точку
					} else {
						fmt.Println("в мапе для Тикета запись ЕСТЬ")
						//в мапе нельзя просто изменить значение.
						for ke, va := range region_VcsSlice {
							if ke == v.Region {
								//https://stackoverflow.com/questions/42716852/how-to-update-map-values-in-go
								//2.Reassigning the modified struct.
								va = append(va, v)
								region_VcsSlice[ke] = va
								break
							}
						}
					}
				} else {
					//Если ticketID не пустое ИЛИ Не было ошибок при получении статуса ИЛИ статуса нет в мапе ДляНовогоСтатуса
					fmt.Println("Созданное обращение:")
					//fmt.Println(bpmUrl + polyTicket.ID) //srID)
					fmt.Println(polyTicket.BpmServer + polyTicket.ID)
					fmt.Println(polyTicket.Status) //statusTicket)
				}
				fmt.Println("")
			}
		}
		//fmt.Println("")
	} //for
	return nil //polyMap, region_VcsSlice, nil
}

// Создание заявок
func (puc *PolyUseCase) TicketsCreating() error {
	//region_VcsSlice map[string][]entity.PolyStruct,	polyMap map[string]entity.PolyStruct) error {

	var usrLogin string
	var trueHour int

	fmt.Println("")
	fmt.Println("Создание заявок по ВКС:")
	for k, v := range region_VcsSlice {
		// k - region
		fmt.Println(k)

		trueHour = timeNowP.Add(time.Duration(v[0].TimeZone-puc.timezone) * time.Hour).Hour()
		if !sleepHoursUnifi[trueHour] || puc.timezone == 100 {

			var vcsInfo []string

			for _, vcs := range v {
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
			description := "Зафиксированы сбои в работе устройств ВидеоКонференцСвязи Poly Trio 8800:" + "\n" +
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

			polyTicket := &entity.Ticket{
				UserLogin:    usrLogin,
				Description:  description,
				Region:       k,
				IncidentType: incidentType,
				Reason:       "",
				Monitoring:   monitoring,
			}

			fmt.Println("Попытка создания заявки")
			//srTicketSlice := CreatePolyTicketErr(soapServer, bpmUrl, usrLogin, description, "", k, monitoring, incidentType)
			//polyTicket, errCreateTicket := puc.soap.CreatePolyTicketErr(polyTicket)
			errCreateTicket := puc.soap.CreatePolyTicketErr(polyTicket)
			if errCreateTicket == nil {
				fmt.Println(polyTicket.BpmServer + polyTicket.ID)
				//delete(regionVcsSlice, k)  //думаю, что удалять не стоит, т.к. будет каждый раз новая мапа создаваться

				//обновляем в мапе srid
				for _, va := range v {
					for key, val := range polyMap {
						if va.IP == val.IP {
							val.SrID = polyTicket.ID //srTicketSlice[0]
							polyMap[key] = val
							break
						}
					}
				}
			} else {
				//если создание заявки прошло с ошибкой, то у меня внутри функции, итак, заложены уведомления
				//в мапу ничего не обновляю. Вернёмся через полчаса, если устройство по-прежнему будет недоступно
			}
			fmt.Println("")

		} else {
			fmt.Println(k)
			fmt.Println("Аларм попадает в спящие часы")
			fmt.Println("Текущий час на сервере: " + strconv.Itoa(timeNowP.Hour()))
			fmt.Println("Временная зона сервера: " + strconv.Itoa(puc.timezone))
			fmt.Println("Временная зона региона: " + strconv.Itoa(v[0].TimeZone))
			fmt.Println("Час в регионе: " + strconv.Itoa(trueHour))
		}

	}
	fmt.Println("")

	return nil
}

//Перезагрузка устройств
