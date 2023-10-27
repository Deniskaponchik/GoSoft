package usecase

import (
	"bytes"
	"fmt"
	"github.com/deniskaponchik/GoSoft/Unifi/internal/entity"
	"strings"
	"time"
)

type UnifiUseCase struct {
	repo         UnifiRepo //interface
	soap         UnifiSoap //interface
	ui           Ui        //interface
	everyCodeMap map[int]bool
	//restartHour  int
}

// реализуем Инъекцию зависимостей DI. Используется в app
func NewUnifiUC(r UnifiRepo, s UnifiSoap, ui Ui, everyCode map[int]bool) *UnifiUseCase {
	return &UnifiUseCase{
		//Мы можем передать сюда ЛЮБОЙ репозиторий (pg, s3 и т.д.) НО КОД НЕ ПОМЕНЯЕТСЯ! В этом смысл DI
		repo:         r,
		soap:         s,
		ui:           ui,
		everyCodeMap: everyCode,
		//restartHour:  restartHour,
	}
}

// Переменные, которые используются во всех методах ниже
var mapAp map[string]*entity.Ap

// var siteapName_ForTickets map[string][]ForApsTicket //НЕ должна создаваться новая раз в 12 минут
var siteNameApCutName_Ap map[string][]*entity.Ap //По новой логике должна создаваться новая раз в 12 минут
var siteApCutName_Login map[string]string        //мапа ответственных сотрудников по офису

var mapClient map[string]*entity.Client             //string = client.mac. client = machine
var mac_OneHourAnomalies map[string]*entity.Anomaly //string = client.mac. обнуляется каждый час

var srStatusCodesForNewTicket map[string]bool
var srStatusCodesForCancelTicket map[string]bool

var timeNow time.Time
var err error

func (puc *UnifiUseCase) InfinityProcessingUnifi() error {

	count12minute := 0
	//count20minute := 0
	countHourAnom := 0
	countHourDBap := 0
	//countHourDBmachine := 0
	//countDayAnom := 0
	countDayDBmachine := 0
	countDay := time.Now().Day()

	srStatusCodesForNewTicket = map[string]bool{
		"Отменено":     true, //Cancel  6e5f4218-f46b-1410-fe9a-0050ba5d6c38
		"Решено":       true, //Resolve  ae7f411e-f46b-1410-009b-0050ba5d6c38
		"Закрыто":      true, //Closed  3e7f420c-f46b-1410-fc9a-0050ba5d6c38
		"На уточнении": true, //Clarification 81e6a1ee-16c1-4661-953e-dde140624fb
		"Тикет введён не корректно": true,
		//"": true,
	}
	srStatusCodesForCancelTicket = map[string]bool{
		"Визирование":  true,
		"Назначено":    true,
		"На уточнении": true, //Clarification 81e6a1ee-16c1-4661-953e-dde140624fb
	}

	//apMyMap := DownloadMapFromDBapsErr(wifiConf.GlpiConnectStringITsupport, bdController)
	mapAp, err = puc.repo.DownloadMapFromDBapsErr()
	if err != nil {
		fmt.Println("мапа точек доступа не смогла загрузиться из БД")
		return err //прекращаем работу скрипта
	}
	//machineMyMap := DownloadMapFromDBmachinesErr(wifiConf.GlpiConnectStringITsupport, bdController)
	mapClient, err = puc.repo.DownloadMapFromDBmachinesErr()
	if err != nil {
		fmt.Println("мапа машин не смогла загрузиться из БД")
		return err //прекращаем работу скрипта
	}
	//siteApCutNameLogin := DownloadMapFromDBerr(wifiConf.GlpiConnectStringITsupport)
	siteApCutName_Login, err = puc.repo.DownloadMapFromDBerr()
	if err != nil {
		fmt.Println("мапа соответствия сайта и логина ответственного сотрудника не загрузилась")
		return err //прекращаем работу скрипта
	}

	for true {
		timeNow := time.Now()
		//Снятие показаний с контроллера раз в 12 минут. Промежутки разные для контроллеров
		//if timeNow.Minute() != 0 && every12start[timeNow.Minute()] && timeNow.Minute() != count12minute {
		if timeNow.Minute() != 0 && puc.everyCodeMap[timeNow.Minute()] && timeNow.Minute() != count12minute {
			count12minute = timeNow.Minute()
			fmt.Println(timeNow.Format("02 January, 15:04:05"))

			err = puc.ui.GetSites() //в puc *UnifiUseCase подгружаются Sites
			if err == nil {
				err = puc.ui.AddAps(mapAp)
				if err == nil {
					siteNameApCutName_Ap = map[string][]*entity.Ap{}

					mac_OneHourAnomalies, err = puc.ui.GetHourAnomalies(mapClient)

				} else {
					fmt.Println(err.Error())
					fmt.Println("точки доступа НЕ загрузились с контроллера")
				}
			} else {
				fmt.Println(err.Error())
				fmt.Println("sites НЕ загрузились с контроллера")
			}
		}
		fmt.Println("Sleep 58s")
		fmt.Println("")
		time.Sleep(58 * time.Second)
	} // while TRUE
	return nil
}

// Удаление из массива точек. не используется в новой логике, где мапа точек ДляТикета обнуляется каждые 12 минут
func removeFromSliceAp(s []*entity.Ap, i int) []*entity.Ap {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

// Обработка точек доступа
func (puc *UnifiUseCase) HandlingAps() error {
	var apCutName string
	var siteApCutName string

	for _, ap := range mapAp {
		if ap.SiteName != "Резерв/Склад" {

			apCutName = strings.Split(ap.Name, "-")[0]    //берём первые 3 буквы от имени точки
			siteApCutName = ap.SiteName + "_" + apCutName //приклеиваем к имени сайта

			if ap.StateInt != 0 { //Точка доступна.
				if ap.SrID == "" { //Заявки нет
					//Новая логика, где мапа ДляТикета обновляется каждые 12 минут
					if ap.CountAttempts != 0 {
						ap.CountAttempts = 0
					}

					/*Старая логика, где мапа ДляТикета статична всегда и за ней нужно следить
					//Пытаемся удалить запись и в мапе ДляТикета, если она там начала создаваться
					for k, v := range siteNameApCutName_Ap {
						// v = массив структур
						if k == siteApCutName {
							ap.CountAttempts = 0
							if len(v) > 1 {
								//если элементов с массиве больше 1, то удаляю лишний
								for i, eap := range v {
									if ap.Mac == eap.Mac {
										v = removeFromSliceAp(v, i)
									}
								}
							} else {
								//если всего 1 элемент в массиве, удаляю бакет в мапе
								delete(siteNameApCutName_Ap, k)
								break
							}
						}
					}*/
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
						err = puc.soap.AddCommentErr(ticket)
						if err == nil {
							ap.CommentCount = 1
						}
					}

					//проверить, не последняя ли это запись в мапе в массиве
					countOfIncident := 0
					for _, v := range mapAp {
						if v.SrID == ap.SrID {
							countOfIncident++
							//BREAK здесь НЕ нужен. Пробежаться нужно по всем
						}
					}

					if countOfIncident == 1 {
						//если последняя запись, пробуем закрыть тикет
						//status := CheckTicketStatusErr(soapServer, srID)
						err = puc.soap.CheckTicketStatusErr(ticket)
						if err == nil {
							fmt.Println(ticket.Status)

							if srStatusCodesForCancelTicket[ticket.Status] {
								//Если статус заявки на Уточнении, Визирование, Назначено
								if ap.CommentCount < 2 {
									ticket.Comment = "Будет предпринята попытка отмены обращения, т.к. все точки из него появились в сети"
									err = puc.soap.AddCommentErr(ticket)
									if err == nil {
										ap.CommentCount = 2
									}
								}

								fmt.Println("Попытка изменить статус в На уточнении")
								ticket.Status = "На уточнении"
								err = puc.soap.ChangeStatusErr(ticket)
								//if error не делаю, т.к. лишним не будет при любом раскладе попытаться вернуть на уточнение

								fmt.Println("Попытка изменить статус в Отменено")
								ticket.Status = "Отменено"
								err = puc.soap.ChangeStatusErr(ticket)
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
					err = puc.soap.CheckTicketStatusErr(ticket)
				}

				fmt.Println("Созданное обращение:")
				fmt.Println(ticket.BpmServer + ap.SrID) //bpmUrl + srID)
				fmt.Println(ticket.Status)              //checkSlice[1])

				//if srStatusCodesForNewTicket[checkSlice[1]] || srID == "" {
				if srStatusCodesForNewTicket[ticket.Status] || ap.SrID == "" {
					//Заявки нет
					fmt.Println("Заявка Закрыта, Отменена, Отклонена или заявки НЕТ вовсе")

					ap.SrID = "" //удаляем заявку
					//ap.CountAttempts = 0

					//Заполняем переменные, которые понадобятся дальше
					fmt.Println("Site ID: " + ap.SiteID)
					fmt.Println(siteApCutName)

					//Проверяем и вносим во временную мапу. Заявка на данном этапе никакая ещё НЕ создаётся
					//_, exisSiteName := siteapNameForTickets[siteApCutName] //проверяем, есть ли в мапе ДЛЯтикетов
					k, exisSiteName := siteNameApCutName_Ap[siteApCutName] //проверяем, есть ли в мапе ДЛЯтикетов

					if !exisSiteName {
						fmt.Println("в мапе для Тикета записи ещё НЕТ")
						//apSlice := []*entity.Ap{ap}
						//создаём массив и вставляем в мапу ДляТикета
						siteNameApCutName_Ap[siteApCutName] = []*entity.Ap{ap}
					} else {
						fmt.Println("в мапе для Тикета запись ЕСТЬ")
						// k - slice
						k = append(k, ap) //просто добавляем точку в уже созданный массив в мапе ДляТикета
					}
				} else {
					//Заявка создана и её статус позволяет её оставить в таком виде
					//ничего не делаем
				}
				fmt.Println("")
			}
		} // if != Резерв/Склад
	} //for
	return nil
}

// Создание заявок
func (puc *UnifiUseCase) TicketsCreatingAps() error {
	fmt.Println("")
	fmt.Println("Создание заявок по точкам:")

	var countAttempts int
	var apsNames []string

	for k, v := range siteNameApCutName_Ap {
		// k - siteNameApCutName
		for _, ap := range v {
			//пробегаемся по массиву точек
			ap.CountAttempts++
			countAttempts = ap.CountAttempts
			apsNames = append(apsNames, ap.Name)
		}

		if countAttempts >= 2 {
			//create ticket
			desAps := strings.Join(apsNames, "\n")

			ticket := &entity.Ticket{
				UserLogin:    siteApCutName_Login[k],
				IncidentType: "Недоступна точка доступа",
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
			err = puc.soap.CreateTicketSmacWifi(ticket)
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
	}
	fmt.Println("")
	return nil
}

func (puc *UnifiUseCase) TicketsCreatingAnomalies() error {

	before30days := timeNow.Add(time.Duration(-720) * time.Hour).Format("2006-01-02 15:04:05")
	//before30days := timeNow.Add(time.Duration(-3) * time.Hour).Format("2006-01-02 15:04:05")

	//Загружаем из БД аномалии за последние 30 дн. в массив структур DateSiteAnom
	//macDay_DateSiteAnom := map[string]DateSiteAnom{}
	//macDay_DateSiteAnom := DownloadMapFromDBanomaliesErr(wifiConf.GlpiConnectStringITsupport, bdController, before30days)
	dayMac_Anomaly, errDownAnomFromDB := puc.repo.DownloadMapFromDBanomaliesErr(before30days)
	if errDownAnomFromDB == nil {

		mac_DateSiteAnomSlice := map[string][]DateSiteAnom{}
		dateSiteAnomSlice := []DateSiteAnom{}

		//Переделываем в мапу типа мак-структура DateSiteAnom за каждый день
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
		//var b2 bytes.Buffer

		for k, v := range mac_DateSiteAnomSlice {
			// TODO: TEST. Change count of caught anomalies
			if len(v) > 9 {
				_, exis := machineMyMap[k]
				if exis {
					for ke, va := range machineMyMap {
						if k == ke {
							noutName = va.Hostname
							fmt.Println(noutName)
							fmt.Println(k) //mac
							apName = va.ApName
							fmt.Println(apName)
							srID = va.SrID
							//fmt.Println(srID)
							exceptionInt = va.Exception

							var statusTicket string
							if srID != "" {
								statusTicket = CheckTicketStatusErr(soapServer, srID)
							}
							if exceptionInt == 0 && (srStatusCodesForNewTicket[statusTicket] || srID == "") {
								//Если заявки ещё нет, либо закрыта отменена
								usrLogin = GetLoginPCerr(wifiConf.GlpiConnectStringGlpi, va.Hostname)
								fmt.Println(usrLogin)

								var b2 bytes.Buffer
								for _, val := range v {
									region = val.SiteName //А если сотрудник был в разных регионах?
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
								//Либо заявка уже есть.
								fmt.Println("Созданное обращение:")
								fmt.Println(bpmUrl + srID)
								fmt.Println(statusTicket)
								//Добавить коммент с аномалиями за последние сутки
								yesterday := timeNow.Add(time.Duration(-22) * time.Hour).Format("2006-01-02")
								dayMac := yesterday + k
								_, exisDayMac := macDay_DateSiteAnom[dayMac]
								if exisDayMac {
									fmt.Println("Есть аномалии за прошедшие сутки. Попытка добавить комментарий...")
									anomSlice := macDay_DateSiteAnom[dayMac].AnomSlice
									var b2 bytes.Buffer
									for _, val := range anomSlice {
										b2.WriteString(val + "\n")
									}
									comment := "За последние сутки появились новые аномалии:" + "\n" +
										b2.String() +
										""
									//fmt.Println(comment)
									AddCommentErr(soapServer, srID, comment, bpmUrl)
								}

							}
							fmt.Println("")
							break
						}
					}
				} else {
					fmt.Println("Не удалось найти запись по маку в мапе машин. Создать заявку невозможно")
				}

			}
		}
	}

	return nil
}
