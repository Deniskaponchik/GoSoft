package usecase

import (
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
var siteNameApCutName_Ap map[string][]*entity.Ap //НЕ должна создаваться новая раз в 12 минут
var siteApCutName_Login map[string]string        //мапа ответственных сотрудников по офису

var mapClient map[string]*entity.Client //client = machine
//var mapAnomaly map[string]*entity.Anomaly

var srStatusCodesForNewTicket map[string]bool
var srStatusCodesForCancelTicket map[string]bool

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

// Опрос устройств
func (puc *UnifiUseCase) ApsHandling() error {
	var apMac string
	var apName string
	//var apLastSeen int

	var apCutName string
	var siteApCutName string

	for _, ap := range mapAp {
		if ap.SiteName != "Резерв/Склад" {

			apMac = ap.Mac
			apName = ap.Name

			apCutName = strings.Split(ap.Name, "-")[0]    //берём первые 3 буквы от имени точки
			siteApCutName = ap.SiteName + "_" + apCutName //приклеиваем к имени сайта

			if ap.StateInt != 0 { //Точка доступна.
				if ap.SrID == "" { //Заявки нет
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
					ap.CountAttempts = 0

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
func (puc *UnifiUseCase) ApsTicketsCreating() error {

	return nil
}

//Перезагрузка устройств

func removeFromSliceAp(s []*entity.Ap, i int) []*entity.Ap {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
