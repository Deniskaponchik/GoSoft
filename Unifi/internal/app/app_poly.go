package app

import (
	"fmt"
	"github.com/deniskaponchik/GoSoft/Unifi/config"
	"github.com/deniskaponchik/GoSoft/Unifi/internal/usecase"
	"github.com/deniskaponchik/GoSoft/Unifi/internal/usecase/netdial"
	_ "github.com/deniskaponchik/GoSoft/Unifi/internal/usecase/ping"
	"github.com/deniskaponchik/GoSoft/Unifi/internal/usecase/repo"
	"github.com/deniskaponchik/GoSoft/Unifi/internal/usecase/soap"
	"github.com/deniskaponchik/GoSoft/Unifi/internal/usecase/webapi"
	"github.com/deniskaponchik/GoSoft/Unifi/pkg/logger"
	"strconv"
	"strings"
	"time"
)

// Run creates objects via constructors.
func PolyRun(cfg *config.Config) {
	fmt.Println("")
	l := logger.New(cfg.Log.Level)

	/* Repository
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()
	*/
	polyUseCase := usecase.New(
		repo.New(cfg.GLPI.GlpiConnectStrITsupport),
		webapi.New(cfg.PolyUsername, cfg.PolyPassword),
		netdial.New(),
		soap.New(cfg.SoapUrl, cfg.BpmUrl),
	)

	//Download MAPs from DB
	//polyMap := DownloadMapFromDBvcsErr(polyConf.GlpiConnectStringITsupport)
	polyMap, errGetEntityMap := polyUseCase.GetEntityMap(0) //Запрос к БД может делать только UseCase.  Не напрямую из какого-либо пакета
	if errGetEntityMap != nil {
		l.Fatal(fmt.Errorf("app - Run - Download polyMap from DB: %w", errGetEntityMap))
	}
	//fmt.Println("Вывод мапы СНАРУЖИ функции")
	/*
		for k, v := range siteApCutNameLogin {
			//fmt.Printf("key: %d, value: %t\n", k, v)
			fmt.Println("newMap "+k, v)
		}
		os.Exit(0)
	*/

	//
	fmt.Println("")
	//log.SetOutput(io.Discard) //Отключить вывод лога

	for true { //зацикливаем навечно
		timeNow := time.Now()

		//
		if timeNow.Minute() != 0 && cfg.EveryCode[timeNow.Minute()] && timeNow.Minute() != cfg.Count20minute {
			//if timeNow.Minute() != 0 && every20Code[timeNow.Minute()] && timeNow.Minute() != count20minute {
			cfg.Count20minute = timeNow.Minute()
			fmt.Println(timeNow.Format("02 January, 15:04:05"))

			//soapServer = soapServerTest	//soapServer = soapServerProd
			fmt.Println("SOAP")
			fmt.Println(cfg.Soap.SoapUrl)
			//bpmUrl = bpmUrlTest    		//bpmUrl = bpmUrlProd
			fmt.Println("BPM")
			fmt.Println(cfg.Bpm.BpmUrl)
			fmt.Println("")
			//Опрос устройств
			regionVcsSlice, errSurvey := polyUseCase.Survey(polyMap)
			if errSurvey != nil {
				l.Info(fmt.Errorf("app - Run - Download polyMap from DB: %w", errSurvey))
			}

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
					vcsInfo = append(vcsInfo, vcs.RoomName)
					vcsInfo = append(vcsInfo, vcs.IP)
					if vcs.PolyType == 1 {
						vcsInfo = append(vcsInfo, "Codec не отвечает на API-запросы")
					} else {
						vcsInfo = append(vcsInfo, "Visual недоступен по netdial")
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
				UpdateMapsToDBerr(polyConf.GlpiConnectStringITsupport, queries)
				fmt.Println("")
			}
			//
			//
			//Обновление мап раз в час. для контроля корректности ip-адресов
			if timeNow.Hour() != countHourFromDB {
				countHourFromDB = timeNow.Hour()

				//polyMap = make(map[string]PolyStruct{})
				//polyMap = map[string]PolyStruct{}
				//clear(polyMyMap)
				polyMap = DownloadMapFromDBvcsErr(polyConf.GlpiConnectStringITsupport)
				fmt.Println("")
			}

			//
			//
			//Перезагрузка
			if timeNow.Hour() == 7 && reboot == 0 {
				for _, v := range polyMap {
					if v.PolyType == 1 {
						fmt.Println(v.RoomName)
						apiSafeRestart2(v.IP, polyConf.PolyUsername, polyConf.PolyPassword)
					}
				}
				reboot = 1
				time.Sleep(2400 * time.Second) //40 minutes
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
