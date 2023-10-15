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

func (puc *PolyUseCase) InfinityPolyProcessing(bpmUrl string, soapUrl string) error {

	/*
		every66Code := map[int]bool{
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
		every63Code := map[int]bool{ //6 minutes
			3:  true,
			9:  true,
			15: true,
			21: true,
			33: true,
			39: true,
			45: true,
			51: true,
			57: true,
		}
	*/
	every20Code := map[int]bool{
		5:  true,
		25: true,
		45: true,
	}

	everyCode = every20Code
	count20minute = 0
	countHourFromDB = 0
	countHourToDB = 0
	reboot = 0

	polyMap, errDownMapFromDB := puc.repo.DownloadMapFromDBvcsErr(0)
	if errDownMapFromDB != nil {
		return errDownMapFromDB
	}

	return nil
}

// Получение списка устройств
func (puc *PolyUseCase) GetEntityMap(int) (map[string]entity.PolyStruct, error) {
	return puc.repo.DownloadMapFromDBvcsErr(0)
}

// Опрос устройств
func (puc *PolyUseCase) Survey(polyMap map[string]entity.PolyStruct, bpmUrl string) (
	region_VcsSlice map[string][]entity.PolyStruct, err error) {

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
			polyTicket := entity.PolyTicket{
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
						polyTicket, errCheckStatus = puc.soap.CheckTicketStatusErr(polyTicket)
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
								errChangeStatus := puc.soap.ChangeStatusErr(polyTicket)
								//ChangeStatusErr(soapServer, srID, "На уточнении")
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
				var errCheckStatus error
				if polyTicket.ID != "" { //srID != "" {
					//statusTicket = CheckTicketStatusErr(soapServer, srID)
					polyTicket, errCheckStatus = puc.soap.CheckTicketStatusErr(polyTicket)
				}

				if srStatusCodesForNewTicket[polyTicket.Status] || polyTicket.ID == "" || errCheckStatus != nil {
					//if srStatusCodesForNewTicket[statusTicket] || srID == "" {
					fmt.Println(bpmUrl + polyTicket.ID)         //srID)
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
					fmt.Println(bpmUrl + polyTicket.ID) //srID)
					fmt.Println(polyTicket.Status)      //statusTicket)
				}
				fmt.Println("")
			}
		}
		//fmt.Println("")
	} //for
	return region_VcsSlice, nil
}

// Создание заявок
func (puc *PolyUseCase) TicketsCreating(region_VcsSlice map[string][]entity.PolyStruct) (polyTicket entity.PolyTicket, err error) {

}

//Перезагрузка устройств
