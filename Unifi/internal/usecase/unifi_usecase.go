package usecase

import (
	"bytes"
	//"fmt"
	"github.com/deniskaponchik/GoSoft/Unifi/internal/entity"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
	//"context"
)

type UnifiUseCase struct {
	repo UnifiRepo //interface. НЕ ИСПОЛЬЗОВАТЬ *
	//repoRostov  UnifiRepo //interface
	//repoNovosib UnifiRepo //interface

	soap UnifiSoap //interface. НЕ ИСПОЛЬЗОВАТЬ *
	//soapTest    UnifiSoap //interface
	//soapProd    UnifiSoap //interface

	//uint		    Ui
	uiRostov  Ui //interface. НЕ ИСПОЛЬЗОВАТЬ *
	uiNovosib Ui //interface. НЕ ИСПОЛЬЗОВАТЬ *

	everyCodeMap             map[int]int //map[int]bool
	countDayTicketCreateAnom int
	countHourAnom            [3]int

	controllerInt  int
	timezone       int
	httpUrl        string
	mx             sync.RWMutex
	hostnameClient map[string]*entity.Client
	hostnameAp     map[string]*entity.Ap
}

// реализуем Инъекцию зависимостей DI. Используется в app
// rr UnifiRepo, rn UnifiRepo, st UnifiSoap, sp UnifiSoap, uiRostov Ui, uiNovosib Ui,
func NewUnifiUC(r UnifiRepo, s UnifiSoap, uiRostov Ui, uiNovosib Ui, everyCodeInt map[int]int, timezone int, httpUrl string,
	countDayTicketCreateAnom int, h1 int, h2 int) *UnifiUseCase {
	return &UnifiUseCase{
		//Мы можем передать сюда ЛЮБОЙ репозиторий (pg, s3 и т.д.) НО КОД НЕ ПОМЕНЯЕТСЯ! В этом смысл DI
		repo: r, //interface
		//repoRostov:   rr,
		//repoNovosib:	rn,
		soap: s, //interface
		//soapTest:     st,
		//soapProd: 	sp,
		uiRostov:                 uiRostov,
		uiNovosib:                uiNovosib,
		everyCodeMap:             everyCodeInt,
		timezone:                 timezone,
		httpUrl:                  httpUrl,
		hostnameClient:           make(map[string]*entity.Client),
		hostnameAp:               make(map[string]*entity.Ap),
		countDayTicketCreateAnom: countDayTicketCreateAnom,
		countHourAnom:            [3]int{0, h1, h2},
	}
}

func (uuc *UnifiUseCase) GetClientForRest(hostName string) *entity.Client { //c context.Context
	uuc.mx.RLock()
	defer uuc.mx.RUnlock()
	client, exisHost := uuc.hostnameClient[hostName]
	if exisHost {
		return client
	} else {
		return nil
	}
}

func (uuc *UnifiUseCase) GetApForRest(hostName string) *entity.Ap { //c context.Context
	uuc.mx.RLock()
	defer uuc.mx.RUnlock()
	ap, exisHost := uuc.hostnameAp[hostName]
	if exisHost {
		return ap
	} else {
		return nil
	}
}

// Переменные, которые используются во всех методах ниже
var mac_Ap map[string]*entity.Ap
var siteApCutName_Office map[string]*entity.Office //мапа ответственных сотрудников по офису. Нужно ТОЛЬКО заявкам по точкам
var mac_Client map[string]*entity.Client           //string = client.mac. client = machine. Не обнуляется + передаётся между функциями

var srStatusCodesForNewTicket map[string]bool
var srStatusCodesForCancelTicket map[string]bool
var sleepHoursUnifi map[int]bool

var timeNowU time.Time
var ui Ui
var before30days string
var err error
var exis bool

func (uuc *UnifiUseCase) InfinityProcessingUnifi() {

	//удалить префикс времени в логах
	//https://stackoverflow.com/questions/48629988/remove-timestamp-prefix-from-go-logger
	//log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	//log.SetFlags(0)

	//count12minute := 0
	//count20minute := 0
	countHourDBap := 0
	//countHourAnom := 0 //Здесь заложены 2 процесса, объединённых в 1 счётчик: получение аномалий с контроллера и выгрузка их в БД
	//countHourAnom := [3]int{} //перенёс в uuc
	countDayDownloadMapsWithAnomalies := 0
	//countDayTicketCreateAnom := 0 //перенёс в uuc

	countDayUploadMachineToDB := 0
	countDayDownlSiteApCutName := time.Now().Day()

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

	mac_Ap, uuc.hostnameAp, err = uuc.repo.Download2MapFromDBaps()
	if err != nil {
		//log.Println("мапа точек доступа не смогла загрузиться из БД")
		//return err //прекращаем работу скрипта
		log.Fatalf("мапа точек доступа не смогла загрузиться из БД")
	}

	siteApCutName_Office, err = uuc.repo.DownloadMapOffice()
	if err != nil {
		//log.Println("мапа соответствия сайта и логина ответственного сотрудника не загрузилась")
		//return err //прекращаем работу скрипта
		log.Fatalf("мапа ответсвенных за офис не смогла загрузиться из БД")
	}

	uuc.mx.Lock() //блокируем на всю загрузку из БД мютекс у hostnameClient
	mac_Client, uuc.hostnameClient, err = uuc.repo.Download2MapFromDBclient()
	if err != nil {
		//log.Println("мапа машин не смогла загрузиться из БД")
		//return err //прекращаем работу скрипта
		log.Fatalf("мапа машин не смогла загрузиться из БД")
	}
	uuc.mx.Unlock()
	//for k, v := range mac_Client {log.Println(k, v.Mac, v.Controller, v.Exception, v.ApMac, v.Modified, v.Hostname, v.SrID)}

	timeNowU = time.Now()
	before30days = timeNowU.Add(time.Duration(-720) * time.Hour).Format("2006-01-02 15:04:05")
	timeNowU = time.Now()
	//errDownClwithAnom := uuc.repo.DownloadClientsWithAnomalySlice(mac_Client, before30days, timeNowU)
	errDownClwithAnom := uuc.repo.DownloadMacMapsClientApWithAnomaly(mac_Client, mac_Ap, before30days, timeNowU)
	if errDownClwithAnom != nil {
		log.Fatalf("мапа соответсвия hostname и клиентов не смогла загрузиться из БД")
	} else {
		countDayDownloadMapsWithAnomalies = timeNowU.Day()
	}

	for true {
		timeNowU = time.Now()

		//intCodeController, exisCodeRun := uuc.everyCodeMap[timeNowU.Minute()]
		uuc.controllerInt, exis = uuc.everyCodeMap[timeNowU.Minute()]
		if exis {
			if uuc.controllerInt == 1 {
				ui = uuc.uiRostov
				uuc.repo.ChangeCntrlNumber(1)
			} else {
				ui = uuc.uiNovosib
				uuc.repo.ChangeCntrlNumber(2)
			}

			//if uuc.everyCodeMap[timeNowU.Minute()] { //актуально для ИЛИ условия: timeNowU.Minute() != 0 || timeNowU.Minute() != count12minute
			//count12minute = timeNowU.Minute()
			log.Println(timeNowU.Format("02 January, 15:04:05"))

			//err = uuc.ui.GetSites() //в uuc *UnifiUseCase подгружаются Sites
			err = ui.GetSites() //в uuc *UnifiUseCase подгружаются Sites
			if err == nil {

				//обработка точек
				uuc.mx.Lock()
				err = ui.AddAps2Maps(mac_Ap, uuc.hostnameAp) //для загрузки требуются Sites. Берутся из ui
				uuc.mx.Unlock()
				if err == nil {
					siteNameApCutName_Ap, errHandlAps := uuc.HandlingAps()
					if errHandlAps == nil {
						err = uuc.TicketsCreatingAps(siteNameApCutName_Ap)
						if err != nil {
							log.Println(err.Error())
							log.Println("функция создания заявок по точкам завершилась ошибкой")
						}
					}

					//Обновление БД ap раз в час.
					if timeNowU.Hour() != countHourDBap {
						log.Println("Ежечасовая выгрузка точек в БД")
						countHourDBap = timeNowU.Hour()
						err = uuc.repo.UpdateDbAp(mac_Ap)
						if err != nil {
							log.Println(err.Error())
							log.Println("выгрузка точек в БД завершилось ошибкой")
						}
					}

					//Загрузка Клиентов с контроллера и обновление двух мап Клиентов
					uuc.mx.Lock() //блокируем на всю загрузку из БД мютекс у hostnameClient
					err = ui.UpdateClients2MapWithoutApMap(mac_Client, uuc.hostnameClient, timeNowU.Format("2006-01-02"))
					if err != nil {
						log.Println(err.Error())
						log.Println("Клиенты НЕ загрузились с контроллера")
					}
					uuc.mx.Unlock()

				} else {
					log.Println(err.Error())
					log.Println("точки доступа НЕ загрузились с контроллера")
				}

				//if timeNowU.Hour() != countHourAnom {
				if timeNowU.Hour() != uuc.countHourAnom[uuc.controllerInt] {
					log.Println("")
					log.Println("Ежечасовое получение и занесение аномалий в БД")

					//macClient_OneHourAnomalies, errGetHourAnom := uuc.ui.GetHourAnomalies(mac_Client, mac_Ap)
					//macClient_OneHourAnomalies, errGetHourAnom := ui.GetHourAnomalies(mac_Client, mac_Ap)
					macClient_OneHourAnomalies, errGetHourAnom := ui.GetHourAnomaliesAddSlice(mac_Client, mac_Ap)
					if errGetHourAnom == nil {
						err = uuc.repo.UpdateDbAnomaly(macClient_OneHourAnomalies)
						if err != nil {
							log.Println("Ежечасовое занесение аномалий в БД завершилось ошибкой")
							log.Println(err.Error())
						} else {
							//Успешное прохождение Получения аномалий с контроллера + выгрузка их БД
							//countHourAnom = timeNowU.Hour()
							uuc.countHourAnom[uuc.controllerInt] = timeNowU.Hour()
						}
					}
				}

			} else {
				log.Println(err.Error())
				log.Println("sites НЕ загрузились с контроллера")
			}

			if timeNowU.Day() != uuc.countDayTicketCreateAnom {
				log.Println("")
				log.Println("Ежесуточное создание заявок по аномалиям")

				var errDayDownMapsWithAnomalies error
				if timeNowU.Day() != countDayDownloadMapsWithAnomalies {
					//если выгрузка сегодня ещё не производилась
					log.Println("Ежесуточная выгрузка аномалий за предыдущие 30дн.")

					before30days = timeNowU.Add(time.Duration(-720) * time.Hour).Format("2006-01-02 15:04:05")
					//before30days := timeNowU.Add(time.Duration(-3) * time.Hour).Format("2006-01-02 15:04:05")
					//errDownClwithAnom := uuc.repo.DownloadMacClientsWithAnomalies(mac_Client, before30days, timeNowU)
					//errDownClwithAnom := uuc.repo.DownloadClientsWithAnomalySlice(mac_Client, before30days, timeNowU)
					errDayDownMapsWithAnomalies = uuc.repo.DownloadMacMapsClientApWithAnomaly(mac_Client, mac_Ap, before30days, timeNowU)

					countDayDownloadMapsWithAnomalies = timeNowU.Day()
				}

				//если ошибок при загрузке из БД за 30 дн. нет, то пробуем создавать заявки
				if errDayDownMapsWithAnomalies == nil {
					err = uuc.TicketsCreatingClientsWithAnomalySlice(mac_Client)
					if err != nil {
						log.Println("Создание заявок на основании аномалий за 30 дней завершилось ошибкой")
						log.Println(err.Error())
					} else {
						uuc.countDayTicketCreateAnom = timeNowU.Day()
					}
				}
			}

			if timeNowU.Day() != countDayUploadMachineToDB {
				log.Println("")
				log.Println("Ежесуточная выгрузка мапы клиентов в БД")

				err = uuc.repo.UpdateDbClient(mac_Client)
				if err == nil {
					countDayUploadMachineToDB = timeNowU.Day()
				} else {
					log.Println(err.Error())
					log.Println("Ежесуточная выгрузка мапы клиентов в БД завершилось ошибкой")
				}
			}

			if timeNowU.Day() != countDayDownlSiteApCutName {
				log.Println("")
				log.Println("Ежесуточное обновление мапы контактных лиц в офисах по точкам")

				//siteApCutName_Login, err = uuc.repo.DownloadMapFromDBerr()
				siteApCutName_Office, err = uuc.repo.DownloadMapOffice()
				if err == nil {
					countDayDownlSiteApCutName = timeNowU.Day()
				} else {
					log.Println(err.Error())
					log.Println("Ежесуточное обновление мапы контактных лиц в офисах по точкам завершилось ошибкой")
				}
			}

			//} //every 12 minutes

		} //if exis in every code map

		//log.Println("Sleep 58s")
		//log.Println("")
		log.Println("Sleep 58s")
		log.Println("")
		time.Sleep(58 * time.Second)
	} // while TRUE
	//return nil
}

// Обработка точек доступа
func (uuc *UnifiUseCase) HandlingAps() (siteNameApCutName_Ap map[string][]*entity.Ap, err error) {
	var apCutName string
	var siteApCutName string
	siteNameApCutName_Ap = make(map[string][]*entity.Ap)

	for _, ap := range mac_Ap {
		if ap.Controller == uuc.controllerInt {

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
						log.Println(ap.Name)
						log.Println(ap.Mac)
						log.Println("Точка доступна. Заявка есть")
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
								log.Println(ticket.Status)

								if srStatusCodesForCancelTicket[ticket.Status] {
									//Если статус заявки на Уточнении, Визирование, Назначено
									if ap.CommentCount < 2 {
										ticket.Comment = "Будет предпринята попытка отмены обращения, т.к. все точки из него появились в сети"
										err = uuc.soap.AddCommentErr(ticket)
										if err == nil {
											ap.CommentCount = 2
										}
									}

									log.Println("Попытка изменить статус в На уточнении")
									ticket.Status = "На уточнении"
									err = uuc.soap.ChangeStatusErr(ticket)
									//if error не делаю, т.к. лишним не будет при любом раскладе попытаться вернуть на уточнение

									log.Println("Попытка изменить статус в Отменено")
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
								log.Println("Статус заявки получить не удалось.Никакие действия с заявкой не будут производиться")
							}
						} else {
							//Если запись НЕ последняя, только удалить из мапы sr и comment, заодно и имя обновим
							ap.SrID = ""
							ap.CommentCount = 0
						}
						log.Println("")
					}
				} else { //Точка НЕ доступна
					log.Println(ap.Name)
					log.Println(ap.Mac)
					log.Println("Точка НЕ доступна")

					ticket := &entity.Ticket{}
					//Проверяем заявку на НЕ закрытость.
					if ap.SrID != "" {
						//status = CheckTicketStatusErr(soapServer, srID)
						ticket.ID = ap.SrID
						err = uuc.soap.CheckTicketStatusErr(ticket)

						log.Println("Созданное обращение:")
						log.Println(ticket.BpmServer + ap.SrID) //bpmUrl + srID)
						log.Println(ticket.Status)              //checkSlice[1])
					}

					//if srStatusCodesForNewTicket[checkSlice[1]] || srID == "" {
					if srStatusCodesForNewTicket[ticket.Status] || ap.SrID == "" {
						//Заявки нет
						log.Println("Заявка Закрыта, Отменена, Отклонена или заявки НЕТ вовсе")

						ap.SrID = "" //удаляем заявку
						//ap.CountAttempts = 0

						//Заполняем переменные, которые понадобятся дальше
						//log.Println("Site ID: " + ap.SiteID)
						log.Println(siteApCutName)

						//Проверяем и вносим во временную мапу. Заявка на данном этапе никакая ещё НЕ создаётся
						//_, exisSiteName := siteapNameForTickets[siteApCutName] //проверяем, есть ли в мапе ДЛЯтикетов
						k, exisSiteName := siteNameApCutName_Ap[siteApCutName] //проверяем, есть ли в мапе ДЛЯтикетов
						//k - Ap slice
						if !exisSiteName {
							log.Println("в мапе для Тикета записи ещё НЕТ")
							//apSlice := []*entity.Ap{ap}
							//создаём массив и вставляем в мапу ДляТикета
							siteNameApCutName_Ap[siteApCutName] = []*entity.Ap{ap}
						} else {
							log.Println("в мапе для Тикета запись ЕСТЬ")
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
					log.Println("")
				}

			} // if != Резерв/Склад
		} //if uuc.controllerInt == ap.controller
	} //for
	return siteNameApCutName_Ap, nil
}

func (uuc *UnifiUseCase) TicketsCreatingAps(siteNameApCutName_Ap map[string][]*entity.Ap) error {
	log.Println("")
	log.Println("Создание заявок по точкам:")

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
					log.Println(ticket.UserLogin)

					//srTicketSlice := CreateSmacWiFiTicketErr(soapServer, bpmUrl, usrLogin, description, v.site, incidentType)
					err = uuc.soap.CreateTicketSmacWifi(ticket)
					if err == nil {
						log.Println(ticket.Url) //srTicketSlice[2])
						//После создания снова пробегаемся по всему массиву точек и прописываем SrID
						for _, ap := range v {
							ap.SrID = ticket.ID
							ap.CountAttempts = 0
						}
						//Удаляем запись в мапе. По новой логике, где мапа ДляТикета обновляется каждые 12 минут это не нужно
						//delete(siteNameApCutName_Ap, k)
					} else {
						log.Println("тикет НЕ был создан. В точках srID НЕ был прописан")
					}
				} else {
					//do nothing. Не создаём тикет. Переходим к следующему бакету мапы ДляТикета
				}
			} else {
				log.Println(k)
				log.Println("Аларм попадает в спящие часы")
				log.Println("Текущий час на сервере: " + strconv.Itoa(timeNowU.Hour()))
				log.Println("Временная зона сервера: " + strconv.Itoa(uuc.timezone))
				log.Println("Временная зона региона: " + strconv.Itoa(office.TimeZone))
				log.Println("Час в регионе: " + strconv.Itoa(trueHour))
			}

		} else {
			log.Println("в мапе siteApCutName_Office нет соответствующего бакета офиса:")
			log.Println(k)

			ticket := &entity.Ticket{
				//UserLogin:    siteApCutName_Login[k],
				UserLogin:    "denis.tirskikh",
				IncidentType: "Недоступна точка доступа",
				Region:       "БиДВ",
				Description: "Не создано соответствие Сайт_ИмяТочки и ответственного сотрудника по офису:" + "\n" +
					k + "\n" +
					"" + "\n",
			}

			err = uuc.soap.CreateTicketSmacWifi(ticket)
			if err == nil {
				log.Println(ticket.Url) //srTicketSlice[2])
				//После создания снова пробегаемся по всему массиву точек и прописываем SrID
				for _, ap := range v {
					ap.SrID = ticket.ID
					ap.CountAttempts = 0
				}
				//Удаляем запись в мапе. По новой логике, где мапа ДляТикета обновляется каждые 12 минут это не нужно
				//delete(siteNameApCutName_Ap, k)
			} else {
				log.Println("тикет НЕ был создан. В точках srID НЕ был прописан")
			}

		}
		log.Println("")
	}
	log.Println("")
	return nil
}

// 2 раза проверяю наличие тикета
func (uuc *UnifiUseCase) TicketsCreatingClientsWithAnomalySlice(mac_Client map[string]*entity.Client) error {

	var lenAnomStructSlice int
	var anomalyStruct *entity.Anomaly
	var anomalyTempMap map[string]string
	var date string
	//var webView string

	for _, client := range mac_Client {

		if client.SrID != "" {
			log.Println(client.Hostname)
			log.Println("Созданное обращение:")

			ticket := &entity.Ticket{
				ID: client.SrID,
			}
			//var errCheckStatus error
			errCheckStatus := uuc.soap.CheckTicketStatusErr(ticket)
			if errCheckStatus == nil {
				log.Println(ticket.Url)
				log.Println(ticket.Status)

				if srStatusCodesForNewTicket[ticket.Status] {
					//Если заявка закрыта, отменена, удаляем запись srid
					log.Println("Удаляем запись о заявке у клиента")
					client.SrID = ""

				} else {
					//Если заявка в работе, визирование, назначено, добавляем комментарий
					yesterday := timeNowU.Add(time.Duration(-22) * time.Hour).Format("2006-01-02")
					var b1 bytes.Buffer

					//беру последнюю добавленную аномалию в массив
					//anomalyStruct = client.SliceAnomalies[lenAnomStructSlice-1]
					anomalyStruct = client.SliceAnomalies[len(client.SliceAnomalies)-1]
					date = strings.Split(anomalyStruct.DateHour, " ")[0] //обрезаю только Date

					//log.Println("yesterday date    = " + yesterday)
					log.Println("last anomaly date = " + date)

					if yesterday == date {
						//если за прошедшие сутки были аномалии
						b1.WriteString(anomalyStruct.ApName + "\n")
						for _, oneAnomaly := range anomalyStruct.SliceAnomStr {
							b1.WriteString(oneAnomaly + "\n")
						}
						b1.WriteString("\n")
					}

					if b1.Len() != 0 {
						ticket.Comment = "За последние сутки появились новые аномалии:" + "\n" +
							b1.String() +
							""
						err = uuc.soap.AddCommentErr(ticket)
						if err == nil {
							log.Println("оставлен комментарий, что за последние сутки были новые аномалии")
						} else {
							log.Println("Комментарий не смог добавиться в обращение")
							log.Println(err.Error())
						}
					}
				}
			}
			log.Println("")
		}

		if client.SrID == "" {

			//У каждого клиента проверить длину массива Аномалий. из бд взяты записи с Exception = 0
			//if len(client.Date_Anomaly) > 9 {
			lenAnomStructSlice = len(client.SliceAnomalies)
			if lenAnomStructSlice > 9 {

				anomalyTempMap = make(map[string]string)
				var b2 bytes.Buffer
				ticket := &entity.Ticket{} //

				//пробегаемся по всем элементам массива аномалий
				for _, anomalyStruct = range client.SliceAnomalies {

					date = strings.Split(anomalyStruct.DateHour, " ")[0]
					anomalyTempMap[date] = date

					b2.WriteString(anomalyStruct.ApName + "\n")
					b2.WriteString(anomalyStruct.DateHour + "\n")
					for _, oneAnomaly := range anomalyStruct.SliceAnomStr {
						b2.WriteString(oneAnomaly + "\n")
					}
					b2.WriteString("\n")

					//имя точки уже получено в каждой аномалии
					ticket.Region = anomalyStruct.SiteName //у клиентов не получаю SiteName. Беру из Аномалий
				}

				if len(anomalyTempMap) > 9 { //если больше 9 дней с аномалями
					log.Println(client.Hostname)

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

					webView := "http://" + uuc.httpUrl + "/client/view/" + client.Hostname

					//https://wiki.tele2.ru/display/ITKB/%5BHelpdesk+IT%5D+System+Monitoring
					//https://wiki.tele2.ru/pages/viewpage.action?pageId=168680976#id-[HelpdeskIT]SystemMonitoring-Аномалии
					ticket.Description = "На ноутбуке:" + "\n" +
						client.Hostname + "\n" + "" + "\n" +
						"За последние 30 дней зафиксировано более 10 дней с Аномалиями качества работы Wi-Fi сети Tele2Corp" + "\n" +
						"" + "\n" +
						"Рекомендации по выполнению таких инцидентов собраны на страничке корпоративной wiki" + "\n" +
						"https://wiki.tele2.ru/pages/viewpage.action?pageId=168680976#id-[HelpdeskIT]SystemMonitoring-Аномалии" + "\n" +
						"" + "\n" +
						"Не нужно закрывать обращение, если кол-во дней с аномалиями за последние 30 дн. больше 10" + "\n" +
						"!!! Создастся новый тикет !!!" + "\n" +
						"" + "\n" +
						"Ресурс для просмотра актуальных аномалий на клиенте:" + "\n" +
						webView + "\n" +
						"" + "\n" +
						"Время аномалий:" + "\n" +
						"для Урала, Сибири и ДВ - Новосибирское" + "\n" +
						"для всей остальной западной России - Московское" + "\n" +
						"" + "\n" +
						"Аномалии обновляются в начале каждого часа" + "\n" +
						//b2.String() +
						""

					log.Println("Попытка создания заявки")
					errCreateTicket := uuc.soap.CreateTicketSmacWifi(ticket)
					if errCreateTicket == nil {
						log.Println(ticket.Url)
						client.SrID = ticket.ID
					} else {
						log.Println("Ошибка создания обращения")
						log.Println(errCreateTicket.Error())
					}

					//log.Println("Sleep 120s")
					//time.Sleep(120 * time.Second)
					log.Println("")
					//log.Println("")

				} //if len(anomalyTempMap) > 9 {

			} //if len(anom.TimeStr_sliceAnomStr) > 9 {

		} //if client.SrID == ""{

	} //for _, client := range mac_Client

	return nil
}

/* Заявки создаём всё по той же mac_Client. Клиенты содержат мапу Аномалий
func (uuc *UnifiUseCase) TicketsCreatingClientsWithAnomalySlice(mac_Client map[string]*entity.Client) error {

	var lenAnomStructSlice int
	var anomalyStruct *entity.Anomaly
	var anomalyTempMap map[string]string
	var date string
	//var webView string

	for _, client := range mac_Client {

		if client.SrID != "" {

			ticket := &entity.Ticket{
				ID: client.SrID,
			}

			//var errCheckStatus error
			errCheckStatus := uuc.soap.CheckTicketStatusErr(ticket)
			if errCheckStatus == nil {
				if srStatusCodesForNewTicket[ticket.Status]{
					//Если заявка закрыта,отменена, удаляем запись srid
					client.SrID = ""
				}else{
					//Если заявка в работе, визирование, назначено, добавляем комментарий
					log.Println("Созданное обращение:")
					log.Println(ticket.Url)
					log.Println(ticket.Status)

					//Добавить коммент с аномалиями за последние сутки
					yesterday := timeNowU.Add(time.Duration(-22) * time.Hour).Format("2006-01-02")
					var b1 bytes.Buffer

					//беру последнюю добавленную аномалию в массив
					anomalyStruct = client.SliceAnomalies[lenAnomStructSlice-1]
					date = strings.Split(anomalyStruct.DateHour, " ")[0] //обрезаю только Date

					log.Println("yesterday date    = " + yesterday)
					log.Println("last anomaly date = " + date)

					if yesterday == date {
						//если за прошедшие сутки были аномалии
						b1.WriteString(anomalyStruct.ApName + "\n")
						for _, oneAnomaly := range anomalyStruct.SliceAnomStr {
							b1.WriteString(oneAnomaly + "\n")
						}
						b1.WriteString("\n")
					}

					if b1.Len() != 0 {
						ticket.Comment = "За последние сутки появились новые аномалии:" + "\n" +
							b1.String() +
							""
						err = uuc.soap.AddCommentErr(ticket)
						if err == nil {
							log.Println("оставлен комментарий, что за последние сутки были новые аномалии")
						} else {
							log.Println("Комментарий не смог добавиться в обращение")
							log.Println(err.Error())
						}
					}
				}
			}
		}

		if client.SrID == ""{

		}



		//У каждого клиента проверить длину массива Аномалий. из бд взяты записи с Exception = 0
		//if len(client.Date_Anomaly) > 9 {
		lenAnomStructSlice = len(client.SliceAnomalies)
		if lenAnomStructSlice > 9 {
			anomalyTempMap = make(map[string]string)
			var b2 bytes.Buffer
			ticket := &entity.Ticket{}

			//пробегаемся по всем элементам массива аномалий
			for _, anomalyStruct = range client.SliceAnomalies {

				date = strings.Split(anomalyStruct.DateHour, " ")[0]
				anomalyTempMap[date] = date

				b2.WriteString(anomalyStruct.ApName + "\n")
				b2.WriteString(anomalyStruct.DateHour + "\n")
				for _, oneAnomaly := range anomalyStruct.SliceAnomStr {
					b2.WriteString(oneAnomaly + "\n")
				}
				b2.WriteString("\n")

				//имя точки уже получено в каждой аномалии
				ticket.Region = anomalyStruct.SiteName //у клиентов не получаю SiteName. Беру из Аномалий
			}

			if len(anomalyTempMap) > 9 { //если больше 9 дней с аномалями
				log.Println(client.Hostname)

				var errCheckStatus error
				if client.SrID != "" {
					ticket.ID = client.SrID
					errCheckStatus = uuc.soap.CheckTicketStatusErr(ticket)
					//log.Println("Заведённое ранее обращение:")
					//log.Println(ticket.BpmServer + client1.SrID)
					//log.Println(ticket.Status)
				} else {
					log.Println("заявка ещё не была создана")
					ticket.Status = ""
				}
				//errCheckStatus := uuc.soap.CheckTicketStatusErr(ticket)
				if errCheckStatus == nil { //&& (srStatusCodesForNewTicket[ticket.Status] || client.SrID == "") {

					if srStatusCodesForNewTicket[ticket.Status] || client.SrID == "" {
						//Если заявки ещё нет, либо закрыта отменена

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

						webView := "http://" + uuc.httpUrl + "/client/view/" + client.Hostname

						ticket.Description = "На ноутбуке:" + "\n" +
							client.Hostname + "\n" + "" + "\n" +
							"За последние 30 дней зафиксировано более 10 дней с Аномалиями качества работы Wi-Fi сети Tele2Corp" + "\n" +
							"" + "\n" +
							"Рекомендации по выполнению таких инцидентов собраны на страничке корпоративной wiki" + "\n" +
							"https://wiki.tele2.ru/display/ITKB/%5BHelpdesk+IT%5D+System+Monitoring" + "\n" +
							"" + "\n" +
							"Не нужно закрывать обращение, если кол-во дней с аномалиями за последние 30 дн. больше 10" + "\n" +
							"!!! Создастся новый тикет !!!" + "\n" +
							"" + "\n" +
							"Ресурс для просмотра актуальных аномалий на клиенте:" + "\n" +
							webView + "\n" +
							"" + "\n" +
							"Время аномалий:" + "\n" +
							"для Урала, Сибири и ДВ - Новосибирское" + "\n" +
							"для всей остальной западной России - Московское" + "\n" +
							"" + "\n" +
							"Аномалии обновляются в начале каждого часа" + "\n" +
							//b2.String() +
							""

						log.Println("Попытка создания заявки")
						errCreateTicket := uuc.soap.CreateTicketSmacWifi(ticket)
						if errCreateTicket == nil {
							log.Println(ticket.Url)
							client.SrID = ticket.ID
						} else {
							log.Println("Ошибка создания обращения")
							log.Println(errCreateTicket.Error())
						}
					} else {
						log.Println("Созданное обращение:")
						log.Println(ticket.Url)
						log.Println(ticket.Status)

						//Добавить коммент с аномалиями за последние сутки
						yesterday := timeNowU.Add(time.Duration(-22) * time.Hour).Format("2006-01-02")
						var b1 bytes.Buffer

						//беру последнюю добавленную аномалию в массив
						anomalyStruct = client.SliceAnomalies[lenAnomStructSlice-1]
						date = strings.Split(anomalyStruct.DateHour, " ")[0] //обрезаю только Date

						log.Println("yesterday date    = " + yesterday)
						log.Println("last anomaly date = " + date)

						if yesterday == date {
							//если за прошедшие сутки были аномалии
							b1.WriteString(anomalyStruct.ApName + "\n")
							for _, oneAnomaly := range anomalyStruct.SliceAnomStr {
								b1.WriteString(oneAnomaly + "\n")
							}
							b1.WriteString("\n")
						}

						if b1.Len() != 0 {
							ticket.Comment = "За последние сутки появились новые аномалии:" + "\n" +
								b1.String() +
								""
							err = uuc.soap.AddCommentErr(ticket)
							if err == nil {
								log.Println("оставлен комментарий, что за последние сутки были новые аномалии")
							} else {
								log.Println("Комментарий не смог добавиться в обращение")
								log.Println(err.Error())
							}
						}
					}
				} else {
					log.Println("Ошибка при получении статуса обращения")
					log.Println("Дальнейшее создание обращения прекращено")
				}

				//log.Println("Sleep 120s")
				//time.Sleep(120 * time.Second)
				log.Println("")
				//log.Println("")

			} //if len(anomalyTempMap) > 9 {

		} //if len(anom.TimeStr_sliceAnomStr) > 9 {

	} //for _, client := range mac_Client

	return nil
}
*/

/* Заявки создаём всё по той же mac_Client. Клиенты содержат мапу Аномалий
func (uuc *UnifiUseCase) TicketsCreatingMacClients(mac_Client map[string]*entity.Client) error {

	before30days = timeNowU.Add(time.Duration(-720) * time.Hour).Format("2006-01-02 15:04:05")
	//before30days := timeNowU.Add(time.Duration(-3) * time.Hour).Format("2006-01-02 15:04:05")

	//mac_Anomaly, errDownAnomFromDB := uuc.repo.DownloadMapFromDBanomaliesErr(before30days)
	//clientsWith30daysAnomalies, errDownClwithAnom := uuc.repo.DownloadClientsWithAnomalies(before30days)
	errDownClwithAnom := uuc.repo.DownloadMacClientsWithAnomalies(mac_Client, before30days, timeNowU)
	if errDownClwithAnom == nil {

		for _, client := range mac_Client {

			//У каждого клиента проверить длину мапы Аномалий. Если длина 10 и более, то пробуем заводить заявку
			if len(client.Date_Anomaly) > 9 {
				log.Println(client.Mac)

				//Проверяем, есть ли hostname. Без него всё бессмысленно
				//if client.Hostname != "" {

				if client.Exception == 0 { //из бд взяты записи с Exception = 0
					log.Println(client.Hostname)

					ticket := &entity.Ticket{}
					var errCheckStatus error
					if client.SrID != "" {
						ticket.ID = client.SrID
						errCheckStatus = uuc.soap.CheckTicketStatusErr(ticket)
						//log.Println("Заведённое ранее обращение:")
						//log.Println(ticket.BpmServer + client1.SrID)
						//log.Println(ticket.Status)
					} else {
						log.Println("заявка ещё не была создана")
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

							log.Println("Попытка создания заявки")
							errCreateTicket := uuc.soap.CreateTicketSmacWifi(ticket)
							if errCreateTicket == nil {
								log.Println(ticket.Url)
								client.SrID = ticket.ID
							} else {
								log.Println("Ошибка создания обращения")
								log.Println(errCreateTicket.Error())
							}
						} else {
							log.Println("Созданное обращение:")
							log.Println(ticket.Url)
							log.Println(ticket.Status)

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
									log.Println("оставлен комментарий, что за последние сутки были новые аномалии")
								} else {
									log.Println("Комментарий не смог добавиться в обращение")
									log.Println(err.Error())
								}
							}
						}
					} else {
						log.Println("Ошибка при получении статуса обращения")
						log.Println("Дальнейшее создание обращения прекращено")
					}
				} else {
					log.Println("Клиент или точка добавлены в исключение")
				}

				log.Println("")
			} //if len(anom.TimeStr_sliceAnomStr) > 9 {
		} //for _, anom := range mac_Anomaly
	} else {
		log.Println("ошибка загрузки мапы аномалий за последние 30 дн. из БД")
		log.Println(errDownClwithAnom.Error())
	}
	return nil
}*/
