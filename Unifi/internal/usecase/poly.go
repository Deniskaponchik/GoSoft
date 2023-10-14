package usecase

import (
	//Not have package imports from the outer layer.
	//UseCase ВСЕГДА остаётся чистым. Он ничего не знает про тех, кто его вызывает. Чистота архитектуры - это про UseCase
	"fmt"
	"github.com/deniskaponchik/GoSoft/Unifi/internal/entity"
)

type PolyUseCase struct {
	repo    PolyRepo    //interface
	webAPI  PolyWebApi  //interface
	netDial PolyNetDial //interface
	soap    PolySoap    //interface
}

// реализуем Инъекцию зависимостей DI. Используется в app
func New(r PolyRepo, a PolyWebApi, n PolyNetDial, s PolySoap) *PolyUseCase {
	return &PolyUseCase{
		//Мы можем передать сюда ЛЮБОЙ репозиторий (pg, s3 и т.д.) НО КОД НЕ ПОМЕНЯЕТСЯ! В этом смысл DI
		repo:    r,
		webAPI:  a,
		netDial: n,
		soap:    s,
	}
}

// Получение списка устройств
func (puc *PolyUseCase) GetEntityMap(int) (map[string]entity.PolyStruct, error) {
	return puc.repo.DownloadMapFromDBvcsErr(0)
}

// Опрос устройств
func (puc *PolyUseCase) Survey(polyMap map[string]entity.PolyStruct) (
	regionVcsSlice map[string][]entity.PolyStruct, err error) {

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
	//regionVcsSlice := map[string][]entity.PolyStruct{}

	//fmt.Println("")
	for k, v := range polyMap { // v == polyStruct
		if v.Exception == 0 {

			/*теперь передаю структуру в сервисы, а не текст
			ip := v.IP
			region := v.Region
			roomName := v.RoomName
			login := v.Login
			srID := v.SrID
			*/
			polyTicket := &entity.PolyTicket{
				UserLogin: v.Login,
				ID:        v.SrID,
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
			var vcsType string
			//var statusReach string   //нигде не использую пока убираю
			var errGetStatus error

			if v.PolyType == 1 {
				vcsType = "Codec"
				//commentUnreach = "Codec не отвечает на API-запросы"
				//statusReach = webapi.apiLineInfo(ip, polyConf.PolyUsername, polyConf.PolyPassword)
				//statusReach, errGetStatus = puc.webAPI.ApiLineInfo(v) //возвращает строку
				v, errGetStatus = puc.webAPI.ApiLineInfoErr(v) //возвращает структуру
			} else {
				vcsType = "Visual"
				//commentUnreach = "Visual не доступен по netdial"
				//statusReach = netDialTmtErr(ip)
				v, errGetStatus = puc.netDial.NetDialTmtErr(v) //возвращает структуру
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
						errAddComment := puc.soap.AddCommentErr(*polyTicket)
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
						statusTicket := CheckTicketStatusErr(soapServer, srID)
						fmt.Println(statusTicket)

						if srStatusCodesForCancelTicket[statusTicket] {
							//Если статус заявки на Уточнении, Визирование, Назначено
							if v.Comment < 2 {
								comment = "Будет предпринята попытка по отмене обращения, т.к. все устройства из него появились в сети"
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
						} else if statusTicket == "" {
							// если Не удалось получить статус
							v.Comment = commentForUpdate
							polyMap[k] = v
						} else {
							//Если статус заявки В работе, на 3 линии, Решено, Закрыто
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
				}

			} else {
				//ВКС недоступна			//} else if statusReach == "" {
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
		//fmt.Println("")
	} //for
	return regionVcsSlice, nil
}

// Создание заявок
func (puc *PolyUseCase) Ticketing() (polyTicket entity.PolyTicket, err error)

//Перезагрузка устройств
