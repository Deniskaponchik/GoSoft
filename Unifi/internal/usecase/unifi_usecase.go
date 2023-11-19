package usecase

import (
	"bytes"
	"fmt"
	"github.com/deniskaponchik/GoSoft/Unifi/internal/entity"
	"strconv"
	"strings"
	"time"
)

type UnifiUseCase struct {
	repo         UnifiRepo //interface
	soap         UnifiSoap //interface
	ui           Ui        //interface
	everyCodeMap map[int]bool
	timezone     int
}

// реализуем Инъекцию зависимостей DI. Используется в app
func NewUnifiUC(r UnifiRepo, s UnifiSoap, ui Ui, everyCode map[int]bool, timezone int) *UnifiUseCase {
	return &UnifiUseCase{
		//Мы можем передать сюда ЛЮБОЙ репозиторий (pg, s3 и т.д.) НО КОД НЕ ПОМЕНЯЕТСЯ! В этом смысл DI
		repo:         r,
		soap:         s,
		ui:           ui,
		everyCodeMap: everyCode,
		timezone:     timezone,
	}
}

// Переменные, которые используются во всех методах ниже
var mac_Ap map[string]*entity.Ap

// var siteApCutName_Login map[string]string //мапа ответственных сотрудников по офису. Нужно ТОЛЬКО заявкам по точкам
var siteApCutName_Office map[string]*entity.Office

var mac_Client map[string]*entity.Client //string = client.mac. client = machine. Не обнуляется + передаётся между функциями
//var mac_HourAnomalies map[string]*entity.Anomaly //string = client.mac. обнуляется каждый час. Создаётся и умирает внутри одной функции InfinityProcessingUnifi
//var clientsWith30daysAnomalies map[string]*entity.Client  //string = client.mac. обнуляется каждый день + живёт внутри одной функции

var srStatusCodesForNewTicket map[string]bool
var srStatusCodesForCancelTicket map[string]bool
var sleepHoursUnifi map[int]bool

var timeNowU time.Time
var err error

func (uuc *UnifiUseCase) InfinityProcessingUnifi() error {

	count12minute := 0
	//count20minute := 0
	countHourDBap := 0

	countHourAnom := 0 //Здесь заложены 2 процесса, объединённых в 1 счётчик: получение аномалий с контроллера и выгрузка их в БД
	countDayTicketCreateAnom := 0
	countDayUploadMachineToDB := 0
	countDayDownlSiteApCutName := time.Now().Day()

	srStatusCodesForNewTicket = map[string]bool{
		"Отменено":                  true, //Cancel  6e5f4218-f46b-1410-fe9a-0050ba5d6c38
		"Решено":                    true, //Resolve  ae7f411e-f46b-1410-009b-0050ba5d6c38
		"Закрыто":                   true, //Closed  3e7f420c-f46b-1410-fc9a-0050ba5d6c38
		"На уточнении":              true, //Clarification 81e6a1ee-16c1-4661-953e-dde140624fb
		"Тикет введён не корректно": true,
		//"": true,
	}
	srStatusCodesForCancelTicket = map[string]bool{
		"Визирование":  true,
		"Назначено":    true,
		"На уточнении": true, //Clarification 81e6a1ee-16c1-4661-953e-dde140624fb
	}
	sleepHoursUnifi = map[int]bool{
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

	//apMyMap := DownloadMapFromDBapsErr(wifiConf.GlpiConnectStringITsupport, bdController)
	mac_Ap, err = uuc.repo.DownloadMapFromDBapsErr()
	if err != nil {
		fmt.Println("мапа точек доступа не смогла загрузиться из БД")
		return err //прекращаем работу скрипта
	}

	//machineMyMap := DownloadMapFromDBmachinesErr(wifiConf.GlpiConnectStringITsupport, bdController)
	mac_Client, err = uuc.repo.DownloadMapFromDBmachinesErr()
	if err != nil {
		fmt.Println("мапа машин не смогла загрузиться из БД")
		return err //прекращаем работу скрипта
	}
	//for k, v := range mac_Client {fmt.Println(k, v.Mac, v.Controller, v.Exception, v.ApMac, v.Modified, v.Hostname, v.SrID)}

	//siteApCutNameLogin := DownloadMapFromDBerr(wifiConf.GlpiConnectStringITsupport)
	//siteApCutName_Login, err = uuc.repo.DownloadMapFromDBerr()
	siteApCutName_Office, err = uuc.repo.DownloadMapOffice()
	if err != nil {
		fmt.Println("мапа соответствия сайта и логина ответственного сотрудника не загрузилась")
		return err //прекращаем работу скрипта
	}

	for true {
		timeNowU = time.Now()
		//Снятие показаний с контроллера раз в 12 минут. Промежутки разные для контроллеров
		//if timeNowU.Minute() != 0 && every12start[timeNowU.Minute()] && timeNowU.Minute() != count12minute {
		if timeNowU.Minute() != 0 && uuc.everyCodeMap[timeNowU.Minute()] && timeNowU.Minute() != count12minute {
			count12minute = timeNowU.Minute()
			fmt.Println(timeNowU.Format("02 January, 15:04:05"))

			err = uuc.ui.GetSites() //в uuc *UnifiUseCase подгружаются Sites
			if err == nil {
				//обработка точек
				err = uuc.ui.AddAps(mac_Ap) //для загрузки требуются Sites. Берутся из ui
				if err == nil {
					siteNameApCutName_Ap, errHandlAps := uuc.HandlingAps()
					if errHandlAps == nil {
						err = uuc.TicketsCreatingAps(siteNameApCutName_Ap)
						if err != nil {
							fmt.Println(err.Error())
							fmt.Println("функция создания заявок по точкам завершилась ошибкой")
						}
					}

					//Обновление БД ap раз в час.
					if timeNowU.Hour() != countHourDBap {
						fmt.Println("Ежечасовая выгрузка точек в БД")
						countHourDBap = timeNowU.Hour()
						err = uuc.repo.UpdateDbAp(mac_Ap)
						if err != nil {
							fmt.Println(err.Error())
							fmt.Println("выгрузка точек в БД завершилось ошибкой")
						}
					}

					//Загрузка Клиентов с контроллера и обновление мапы Клиентов mac_Client
					err = uuc.ui.UpdateClientsWithoutApMap(mac_Client, timeNowU.Format("2006-01-02"))
					if err != nil {
						fmt.Println(err.Error())
						fmt.Println("Клиенты НЕ загрузились с контроллера")
					}
					/*fmt.Println("вывод мапы после AddClients")
					for k, v := range mac_Client {fmt.Println(k, v.Mac, v.Controller, v.Exception, v.ApMac, v.Modified, v.Hostname, v.SrID)}					}
					time.Sleep(6000000 * time.Second)*/

				} else {
					fmt.Println(err.Error())
					fmt.Println("точки доступа НЕ загрузились с контроллера")
				}

				if timeNowU.Hour() != countHourAnom {
					fmt.Println("")
					fmt.Println("Ежечасовое получение и занесение аномалий в БД")

					macClient_OneHourAnomalies, errGetHourAnom := uuc.ui.GetHourAnomalies(mac_Client, mac_Ap)
					if errGetHourAnom == nil {
						err = uuc.repo.UpdateDbAnomaly(macClient_OneHourAnomalies)
						if err != nil {
							fmt.Println("Ежечасовое занесение аномалий в БД завершилось ошибкой")
							fmt.Println(err.Error())
						} else {
							//Успешное прохождение Получения аномалий с контроллера + выгрузка их БД
							countHourAnom = timeNowU.Hour()
						}
					}
				}

			} else {
				fmt.Println(err.Error())
				fmt.Println("sites НЕ загрузились с контроллера")
			}

			if timeNowU.Day() != countDayTicketCreateAnom {
				fmt.Println("")
				fmt.Println("Ежесуточное создание заявок по аномалиям")

				//err = uuc.TicketsCreatingAnomalies(mac_Client)
				err = uuc.TicketsCreatingMacClients(mac_Client)
				if err != nil {
					fmt.Println("Создание заявок на основании аномалий за 30 дней завершилось ошибкой")
					fmt.Println(err.Error())
				} else {
					countDayTicketCreateAnom = timeNowU.Day()
				}
			}

			if timeNowU.Day() != countDayUploadMachineToDB {
				fmt.Println("")
				fmt.Println("Ежесуточная выгрузка мапы клиентов в БД")

				err = uuc.repo.UpdateDbClient(mac_Client)
				if err == nil {
					countDayUploadMachineToDB = timeNowU.Day()
				} else {
					fmt.Println(err.Error())
					fmt.Println("Ежесуточная выгрузка мапы клиентов в БД завершилось ошибкой")
				}
			}

			if timeNowU.Day() != countDayDownlSiteApCutName {
				fmt.Println("")
				fmt.Println("Ежесуточное обновление мапы контактных лиц в офисах по точкам")

				//siteApCutName_Login, err = uuc.repo.DownloadMapFromDBerr()
				siteApCutName_Office, err = uuc.repo.DownloadMapOffice()
				if err == nil {
					countDayDownlSiteApCutName = timeNowU.Day()
				} else {
					fmt.Println(err.Error())
					fmt.Println("Ежесуточное обновление мапы контактных лиц в офисах по точкам завершилось ошибкой")
				}
			}

		}
		fmt.Println("Sleep 58s")
		fmt.Println("")
		time.Sleep(58 * time.Second)
	} // while TRUE
	return nil
}

// Обработка точек доступа
func (uuc *UnifiUseCase) HandlingAps() (siteNameApCutName_Ap map[string][]*entity.Ap, err error) {
	var apCutName string
	var siteApCutName string
	siteNameApCutName_Ap = make(map[string][]*entity.Ap)

	for _, ap := range mac_Ap {
		if ap.SiteName != "Резерв/Склад" {

			apCutName = strings.Split(ap.Name, "-")[0]    //берём первые 3 буквы от имени точки
			siteApCutName = ap.SiteName + "_" + apCutName //приклеиваем к имени сайта

			if ap.StateInt != 0 { //Точка доступна.
				if ap.SrID == "" { //Заявки нет
					//Новая логика, где мапа ДляТикета обновляется каждые 12 минут. Реализовано за счёт ap.CountAttempts
					if ap.CountAttempts != 0 {
						ap.CountAttempts = 0
					}

				} else { //Заявка есть
					fmt.Println(ap.Name)
					fmt.Println(ap.Mac)
					fmt.Println("Точка доступна. Заявка есть")
					//Оставляем коммент, ПЫТАЕМСЯ закрыть тикет, если на визировании, Очищаем запись в мапе,
					ticket := &entity.Ticket{
						ID: ap.SrID,
					}

					ticket.Comment = "Точка появилась в сети: " + ap.Name
					if ap.CommentCount < 1 {
						err = uuc.soap.AddCommentErr(ticket)
						if err == nil {
							ap.CommentCount = 1
						}
					}

					//проверить, не последняя ли это запись в мапе в массиве
					countOfIncident := 0
					for _, v := range mac_Ap {
						if v.SrID == ap.SrID {
							countOfIncident++
							//BREAK здесь НЕ нужен. Пробежаться нужно по всем
						}
					}

					if countOfIncident == 1 {
						//если последняя запись, пробуем закрыть тикет
						//status := CheckTicketStatusErr(soapServer, srID)
						err = uuc.soap.CheckTicketStatusErr(ticket)
						if err == nil {
							fmt.Println(ticket.Status)

							if srStatusCodesForCancelTicket[ticket.Status] {
								//Если статус заявки на Уточнении, Визирование, Назначено
								if ap.CommentCount < 2 {
									ticket.Comment = "Будет предпринята попытка отмены обращения, т.к. все точки из него появились в сети"
									err = uuc.soap.AddCommentErr(ticket)
									if err == nil {
										ap.CommentCount = 2
									}
								}

								fmt.Println("Попытка изменить статус в На уточнении")
								ticket.Status = "На уточнении"
								err = uuc.soap.ChangeStatusErr(ticket)
								//if error не делаю, т.к. лишним не будет при любом раскладе попытаться вернуть на уточнение

								fmt.Println("Попытка изменить статус в Отменено")
								ticket.Status = "Отменено"
								err = uuc.soap.ChangeStatusErr(ticket)
								if err == nil {
									//Если отмена заявки прошла успешно
									ap.SrID = ""
									ap.CommentCount = 0
								} else {
									//Если НЕ удалось отменить заявку
									//valueAp.SrID не зануляем, т.к. будет второй заход через 12 минут
									//ap.CommentCount остаётся равным 2
								}
							} else {
								//Если статус заявки В работе, Решено, Закрыто и т.д.
								ap.SrID = ""
								ap.CommentCount = 0
							}
						} else {
							fmt.Println("Статус заявки получить не удалось.Никакие действия с заявкой не будут производиться")
						}
					} else {
						//Если запись НЕ последняя, только удалить из мапы sr и comment, заодно и имя обновим
						ap.SrID = ""
						ap.CommentCount = 0
					}
					fmt.Println("")
				}
			} else { //Точка НЕ доступна
				fmt.Println(ap.Name)
				fmt.Println(ap.Mac)
				fmt.Println("Точка НЕ доступна")

				ticket := &entity.Ticket{}
				//Проверяем заявку на НЕ закрытость.
				if ap.SrID != "" {
					//status = CheckTicketStatusErr(soapServer, srID)
					ticket.ID = ap.SrID
					err = uuc.soap.CheckTicketStatusErr(ticket)

					fmt.Println("Созданное обращение:")
					fmt.Println(ticket.BpmServer + ap.SrID) //bpmUrl + srID)
					fmt.Println(ticket.Status)              //checkSlice[1])
				}

				//if srStatusCodesForNewTicket[checkSlice[1]] || srID == "" {
				if srStatusCodesForNewTicket[ticket.Status] || ap.SrID == "" {
					//Заявки нет
					fmt.Println("Заявка Закрыта, Отменена, Отклонена или заявки НЕТ вовсе")

					ap.SrID = "" //удаляем заявку
					//ap.CountAttempts = 0

					//Заполняем переменные, которые понадобятся дальше
					//fmt.Println("Site ID: " + ap.SiteID)
					fmt.Println(siteApCutName)

					//Проверяем и вносим во временную мапу. Заявка на данном этапе никакая ещё НЕ создаётся
					//_, exisSiteName := siteapNameForTickets[siteApCutName] //проверяем, есть ли в мапе ДЛЯтикетов
					k, exisSiteName := siteNameApCutName_Ap[siteApCutName] //проверяем, есть ли в мапе ДЛЯтикетов
					//k - Ap slice
					if !exisSiteName {
						fmt.Println("в мапе для Тикета записи ещё НЕТ")
						//apSlice := []*entity.Ap{ap}
						//создаём массив и вставляем в мапу ДляТикета
						siteNameApCutName_Ap[siteApCutName] = []*entity.Ap{ap}
					} else {
						fmt.Println("в мапе для Тикета запись ЕСТЬ")
						// k - slice
						k = append(k, ap)
						siteNameApCutName_Ap[siteApCutName] = k

						//apSlice := k
						//apSlice = append(apSlice, ap)
						//siteNameApCutName_Ap[siteApCutName] = apSlice
					}
				} else {
					//Заявка создана и её статус позволяет её оставить в таком виде
					//ничего не делаем
				}
				fmt.Println("")
			}
		} // if != Резерв/Склад
	} //for
	return siteNameApCutName_Ap, nil
}

func (uuc *UnifiUseCase) TicketsCreatingAps(siteNameApCutName_Ap map[string][]*entity.Ap) error {
	fmt.Println("")
	fmt.Println("Создание заявок по точкам:")

	var countAttempts int
	var region string
	//var office *entity.Office
	var trueHour int

	for k, v := range siteNameApCutName_Ap {
		// k - siteNameApCutName    v - Ap slice

		office, offExis := siteApCutName_Office[k]
		if offExis {
			trueHour = timeNowU.Add(time.Duration(office.TimeZone-uuc.timezone) * time.Hour).Hour()
			if !sleepHoursUnifi[trueHour] || uuc.timezone == 100 {

				/*если зонаКода < зоныПроблемы{
				if uuc.timezone > office.TimeZone {
					sumTime = timeNowU.Hour() - uuc.timezone - office.TimeZone
				}else{
					sumTime = timeNowU.Hour() + office.TimeZone - uuc.timezone
				}*/

				var apsNames []string

				for _, ap := range v {
					//пробегаемся по массиву точек
					ap.CountAttempts++
					countAttempts = ap.CountAttempts
					apsNames = append(apsNames, ap.Name)
					region = ap.SiteName
				}

				if countAttempts >= 2 {
					//create ticket
					desAps := strings.Join(apsNames, "\n")

					ticket := &entity.Ticket{
						//UserLogin:    siteApCutName_Login[k],
						UserLogin:    office.UserLogin,
						IncidentType: "Недоступна точка доступа",
						Region:       region,
						Description: "Зафиксировано отключение Wi-Fi точек доступа:" + "\n" +
							desAps + "\n" +
							"" + "\n" +
							"Рекомендации по выполнению таких инцидентов собраны на страничке корпоративной wiki" + "\n" +
							"https://wiki.tele2.ru/display/ITKB/%5BHelpdesk+IT%5D+System+Monitoring" + "\n" +
							"" + "\n" +
							"!!! Не нужно решать/отменять/отклонять/возвращать/закрывать заявку, пока работа точек не будет восстановлена - автоматически создастся новый тикет !!!" + "\n" +
							"",
					}
					if ticket.UserLogin == "" {
						ticket.UserLogin = "denis.tirskikh"
					}
					fmt.Println(ticket.UserLogin)

					//srTicketSlice := CreateSmacWiFiTicketErr(soapServer, bpmUrl, usrLogin, description, v.site, incidentType)
					err = uuc.soap.CreateTicketSmacWifi(ticket)
					if err == nil {
						fmt.Println(ticket.Url) //srTicketSlice[2])
						//После создания снова пробегаемся по всему массиву точек и прописываем SrID
						for _, ap := range v {
							ap.SrID = ticket.ID
							ap.CountAttempts = 0
						}
						//Удаляем запись в мапе. По новой логике, где мапа ДляТикета обновляется каждые 12 минут это не нужно
						//delete(siteNameApCutName_Ap, k)
					} else {
						fmt.Println("тикет НЕ был создан. В точках srID НЕ был прописан")
					}
				} else {
					//do nothing. Не создаём тикет. Переходим к следующему бакету мапы ДляТикета
				}
			} else {
				fmt.Println(k)
				fmt.Println("Аларм попадает в спящие часы")
				fmt.Println("Текущий час на сервере: " + strconv.Itoa(timeNowU.Hour()))
				fmt.Println("Временная зона сервера: " + strconv.Itoa(uuc.timezone))
				fmt.Println("Временная зона региона: " + strconv.Itoa(office.TimeZone))
				fmt.Println("Час в регионе: " + strconv.Itoa(trueHour))
			}

		} else {
			fmt.Println("в мапе siteApCutName_Office нет соответствующего бакета офиса:")
			fmt.Println(k)
		}
		fmt.Println("")
	}
	fmt.Println("")
	return nil
}

// Заявки создаём всё по той же mac_Client
func (uuc *UnifiUseCase) TicketsCreatingMacClients(mac_Client map[string]*entity.Client) error {

	before30days := timeNowU.Add(time.Duration(-720) * time.Hour).Format("2006-01-02 15:04:05")
	//before30days := timeNowU.Add(time.Duration(-3) * time.Hour).Format("2006-01-02 15:04:05")

	//mac_Anomaly, errDownAnomFromDB := uuc.repo.DownloadMapFromDBanomaliesErr(before30days)
	//clientsWith30daysAnomalies, errDownClwithAnom := uuc.repo.DownloadClientsWithAnomalies(before30days)
	errDownClwithAnom := uuc.repo.DownloadMacClientsWithAnomalies(mac_Client, before30days, timeNowU)
	if errDownClwithAnom == nil {

		for _, client := range mac_Client {

			//У каждого клиента проверить длину мапы Аномалий. Если длина 10 и более, то пробуем заводить заявку
			if len(client.Date_Anomaly) > 9 {
				fmt.Println(client.Mac)

				//Проверяем, есть ли hostname. Без него всё бессмысленно
				//if client.Hostname != "" {

				if client.Exception == 0 { //из бд взяты записи с Exception = 0
					fmt.Println(client.Hostname)

					ticket := &entity.Ticket{}
					var errCheckStatus error
					if client.SrID != "" {
						ticket.ID = client.SrID
						errCheckStatus = uuc.soap.CheckTicketStatusErr(ticket)
						//fmt.Println("Заведённое ранее обращение:")
						//fmt.Println(ticket.BpmServer + client1.SrID)
						//fmt.Println(ticket.Status)
					} else {
						fmt.Println("заявка ещё не была создана")
						ticket.Status = ""
					}
					//errCheckStatus := uuc.soap.CheckTicketStatusErr(ticket)
					if errCheckStatus == nil { //&& (srStatusCodesForNewTicket[ticket.Status] || client.SrID == "") {

						if srStatusCodesForNewTicket[ticket.Status] || client.SrID == "" {
							//Если заявки ещё нет, либо закрыта отменена
							var b2 bytes.Buffer

							for date, anomalyStruct := range client.Date_Anomaly {
								//имя точки уже получено в каждой аномалии
								ticket.Region = anomalyStruct.SiteName //у клиентов не получаю SiteName. Беру из Аномалий

								b2.WriteString(anomalyStruct.ApName + "\n")
								b2.WriteString(date + "\n")
								for _, oneAnomaly := range anomalyStruct.SliceAnomStr {
									b2.WriteString(oneAnomaly + "\n")
								}
								b2.WriteString("\n")
							}

							//Получение userlogin
							if client.Hostname != "" {
								errGetUserLogin := uuc.repo.GetLoginPCerr(client)
								if errGetUserLogin != nil {
									client.UserLogin = "denis.tirskikh"
								}
							} else {
								//если client.Hostname == "" то создаю информационную заявку на себя, чтобы добавить в БД руками hostname
								client.UserLogin = "denis.tirskikh"
								client.Hostname = client.Mac
							}
							ticket.UserLogin = client.UserLogin

							ticket.IncidentType = "Плохое качество соединения клиента"
							//ticket.Region = anom.SiteName //получаю выше в цикле обработки аномалий

							ticket.Description = "На ноутбуке:" + "\n" +
								client.Hostname + "\n" + "" + "\n" +
								"За последние 30 дней зафиксировано более 10 дней с Аномалиями качества работы Wi-Fi сети Tele2Corp" + "\n" +
								"" + "\n" +
								"Рекомендации по выполнению таких инцидентов собраны на страничке корпоративной wiki" + "\n" +
								"https://wiki.tele2.ru/display/ITKB/%5BHelpdesk+IT%5D+System+Monitoring" + "\n" +
								"" + "\n" +
								b2.String() +
								""

							fmt.Println("Попытка создания заявки")
							errCreateTicket := uuc.soap.CreateTicketSmacWifi(ticket)
							if errCreateTicket == nil {
								fmt.Println(ticket.Url)
								client.SrID = ticket.ID
							} else {
								fmt.Println("Ошибка создания обращения")
								fmt.Println(errCreateTicket.Error())
							}
						} else {
							fmt.Println("Созданное обращение:")
							fmt.Println(ticket.Url)
							fmt.Println(ticket.Status)

							//Добавить коммент с аномалиями за последние сутки
							yesterday := timeNowU.Add(time.Duration(-22) * time.Hour).Format("2006-01-02")
							var b1 bytes.Buffer

							for date, anomalyStruct := range client.Date_Anomaly {
								//имя точки уже получено в каждой аномалии
								//ticket.Region = anomalyStruct.SiteName //у клиентов не получаю SiteName. Беру из Аномалий
								//b1.WriteString(date + "\n")

								if date == yesterday {
									b1.WriteString(anomalyStruct.ApName + "\n")
									for _, oneAnomaly := range anomalyStruct.SliceAnomStr {
										b1.WriteString(oneAnomaly + "\n")
									}
									b1.WriteString("\n")

									break
								}
							}
							if b1.Len() != 0 {
								ticket.Comment = "За последние сутки появились новые аномалии:" + "\n" +
									b1.String() +
									""
								err = uuc.soap.AddCommentErr(ticket)
								if err == nil {
									fmt.Println("оставлен комментарий, что за последние сутки были новые аномалии")
								} else {
									fmt.Println("Комментарий не смог добавиться в обращение")
									fmt.Println(err.Error())
								}
							}
						}
					} else {
						fmt.Println("Ошибка при получении статуса обращения")
						fmt.Println("Дальнейшее создание обращения прекращено")
					}
				} else {
					fmt.Println("Клиент или точка добавлены в исключение")
				}
				/*
					} else {
						fmt.Println("у клиента не прописан hostname")
						fmt.Println("Создание заявки без него невозможно")
					}*/

				fmt.Println("")
			} //if len(anom.TimeStr_sliceAnomStr) > 9 {
		} //for _, anom := range mac_Anomaly
	} else {
		fmt.Println("ошибка загрузки мапы аномалий за последние 30 дн. из БД")
		fmt.Println(errDownClwithAnom.Error())
	}
	return nil
}

// приходит ОТДЕЛЬНАЯ мапа Клиентов со вложенной мапой Аномалий
/*
func (uuc *UnifiUseCase) TicketsCreatingAnomalies(mac_Client map[string]*entity.Client) error {

	before30days := timeNowU.Add(time.Duration(-720) * time.Hour).Format("2006-01-02 15:04:05")
	//before30days := timeNowU.Add(time.Duration(-3) * time.Hour).Format("2006-01-02 15:04:05")

	//mac_Anomaly, errDownAnomFromDB := uuc.repo.DownloadMapFromDBanomaliesErr(before30days)
	clientsWith30daysAnomalies, errDownClwithAnom := uuc.repo.DownloadClientsWithAnomalies(before30days)
	if errDownClwithAnom == nil {

		for client2Mac, client2 := range clientsWith30daysAnomalies { //client2 - клиент с мапой Аномалий

			//у каждого клиента проверить длину мапы Аномалий. если длина 10 и более, то заводить заявку
			if len(client2.Date_Anomaly) > 9 {
				fmt.Println(client2.Mac)

				//идём в мапу Клиентов, чтобы узнать, не заведена ли уже заявка
				client1, exisMacClient1 := mac_Client[client2Mac] //cient1 - клиент без мапы аномалий
				if exisMacClient1 {
					if client1.Exception == 0 { //из бд взяты записи с Exception = 0
						fmt.Println(client1.Hostname)

						ticket := &entity.Ticket{}
						var errCheckStatus error
						if client1.SrID != "" {
							ticket.ID = client1.SrID
							errCheckStatus = uuc.soap.CheckTicketStatusErr(ticket)
							//fmt.Println("Заведённое ранее обращение:")
							//fmt.Println(ticket.BpmServer + client1.SrID)
							//fmt.Println(ticket.Status)
						} else {
							fmt.Println("заявка ещё не была создана")
							ticket.Status = ""
						}
						//errCheckStatus := uuc.soap.CheckTicketStatusErr(ticket)
						if errCheckStatus == nil { //&& (srStatusCodesForNewTicket[ticket.Status] || client.SrID == "") {

							if srStatusCodesForNewTicket[ticket.Status] || client1.SrID == "" {
								//Если заявки ещё нет, либо закрыта отменена
								var b2 bytes.Buffer

								for date, anomalyStruct := range client2.Date_Anomaly {
									//имя точки уже получено в каждой аномалии
									ticket.Region = anomalyStruct.SiteName //у клиентов не получаю SiteName. Беру из Аномалий

									b2.WriteString(anomalyStruct.ApName + "\n")
									b2.WriteString(date + "\n")
									for _, oneAnomaly := range anomalyStruct.SliceAnomStr {
										b2.WriteString(oneAnomaly + "\n")
									}
									b2.WriteString("\n")
								}

								errGetUserLogin := uuc.repo.GetLoginPCerr(client1)
								if errGetUserLogin == nil {
									ticket.UserLogin = client1.UserLogin
								} else {
									//в логику GetLoginPCerr уже заложено назначение client.UserLogin = "denis.tirskikh"
									//ticket.UserLogin = "denis.tirskikh"
								}

								ticket.IncidentType = "Плохое качество соединения клиента"
								//ticket.Region = anom.SiteName //получаю выше в цикле обработки аномалий

								ticket.Description = "На ноутбуке:" + "\n" +
									client1.Hostname + "\n" + "" + "\n" +
									"За последние 30 дней зафиксировано более 10 дней с Аномалиями качества работы Wi-Fi сети Tele2Corp" + "\n" +
									"" + "\n" +
									"Рекомендации по выполнению таких инцидентов собраны на страничке корпоративной wiki" + "\n" +
									"https://wiki.tele2.ru/display/ITKB/%5BHelpdesk+IT%5D+System+Monitoring" + "\n" +
									"" + "\n" +
									b2.String() +
									""

								fmt.Println("Попытка создания заявки")
								errCreateTicket := uuc.soap.CreateTicketSmacWifi(ticket)
								if errCreateTicket == nil {
									fmt.Println(ticket.Url)
									client1.SrID = ticket.ID
								} else {
									fmt.Println("Ошибка создания обращения")
									fmt.Println(errCreateTicket.Error())
								}
							} else {
								fmt.Println("Созданное обращение:")
								fmt.Println(ticket.Url)
								fmt.Println(ticket.Status)

								//Добавить коммент с аномалиями за последние сутки
								yesterday := timeNowU.Add(time.Duration(-22) * time.Hour).Format("2006-01-02")
								var b1 bytes.Buffer

								for date, anomalyStruct := range client2.Date_Anomaly {
									//имя точки уже получено в каждой аномалии
									//ticket.Region = anomalyStruct.SiteName //у клиентов не получаю SiteName. Беру из Аномалий
									//b1.WriteString(date + "\n")

									if date == yesterday {
										b1.WriteString(anomalyStruct.ApName + "\n")
										for _, oneAnomaly := range anomalyStruct.SliceAnomStr {
											b1.WriteString(oneAnomaly + "\n")
										}
										b1.WriteString("\n")

										break
									}
									if b1.Len() != 0 {
										ticket.Comment = "За последние сутки появились новые аномалии:" + "\n" +
											b1.String() +
											""
										err = uuc.soap.AddCommentErr(ticket)
										if err == nil {
											fmt.Println("оставлен комментарий, что за последние сутки были новые аномалии")
										} else {
											fmt.Println("Комментарий не смог добавиться в обращение")
											fmt.Println(err.Error())
										}
									}
								}
							}
						} else {
							fmt.Println("Ошибка при получении статуса обращения")
							fmt.Println("Дальнейшее создание обращения прекращено")
						}
					} else {
						fmt.Println("Клиент или точка добавлены в исключение")
					}
				} else {
					fmt.Println("мак не найден в мапе mac_Client")
					fmt.Println("Создание заявки без этих данных невозможно")
				}
			} //if len(anom.TimeStr_sliceAnomStr) > 9 {
			fmt.Println("")
		} //for _, anom := range mac_Anomaly
	} else {
		fmt.Println("ошибка загрузки мапы аномалий за последние 30 дн. из БД")
		fmt.Println(errDownClwithAnom.Error())
	}
	return nil
}*/
