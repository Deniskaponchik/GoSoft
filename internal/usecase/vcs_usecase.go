package usecase

import (
	//Not have package imports from the outer layer.
	//UseCase ВСЕГДА остаётся чистым. Он ничего не знает про тех, кто его вызывает. Чистота архитектуры - это про UseCase
	"fmt"
	"github.com/deniskaponchik/GoSoft/internal/entity"
	"strconv"
	"strings"
	"time"
)

type VcsUseCase struct {
	repo    PolyRepo    //interface
	webAPI  PolyWebApi  //interface
	netDial PolyNetDial //interface
	soap    PolySoap    //interface

	everyCodeMap map[int]bool
	restartHour  int
	timezone     int
}

// реализуем Инъекцию зависимостей DI. Используется в app
func NewVcs(r PolyRepo, a PolyWebApi, n PolyNetDial, s PolySoap, everyCode map[int]bool, restartHour int, timezone int) *PolyUseCase {
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

// Создание заявок
func (puc *PolyUseCase) TicketsCreatingVCS() error {
	//region_VcsSlice map[string][]entity.PolyStruct,	polyMap map[string]entity.PolyStruct) error {

	var usrLogin string
	var trueHour int

	fmt.Println("")
	fmt.Println("Создание заявок по ВКС:")
	for k, v := range region_VcsSlice {
		// k - region
		fmt.Println(k)

		trueHour = timeNowP.Add(time.Duration(v[0].TimeZone-puc.timezone) * time.Hour).Hour()
		if !sleepHoursPoly[trueHour] || puc.timezone == 100 {

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
