package api_web

import (
	"github.com/deniskaponchik/GoSoft/internal/entity"
	"github.com/unpoller/unifi"
	"log"
	"strings"
	"time"
)

type Ui struct {
	//unf unifi.Unifi
	Conf       unifi.Config
	Uni        *unifi.Unifi
	Sites      []*unifi.Site
	Controller int //1-Rostov, 2-Novosib
}

func NewUi(u string, p string, url string, cntrlInt int) *Ui {
	log.Println(url)

	unfConf := unifi.Config{
		User:     u,
		Pass:     p,
		URL:      url,
		ErrorLog: log.Printf,
		DebugLog: log.Printf,
	}
	return &Ui{
		Conf:       unfConf,
		Controller: cntrlInt,
	}
}

func (ui *Ui) GetSites() (err error) { //unifi.Unifi, error){
	uni, errNewUnifi := unifi.NewUnifi(&ui.Conf) //&c)
	if errNewUnifi == nil {
		log.Println("uni загрузился")
		ui.Uni = uni
		sites, errGetSites := uni.GetSites()
		if errGetSites == nil {
			log.Println("sites загрузились")
			ui.Sites = sites
			return nil
		} else {
			log.Println("sites НЕ загрузились")
			return errGetSites
		}
	} else {
		log.Println("uni НЕ загрузился")
		return errNewUnifi
	}
	//return nil
}

// Не создаёт заявки по взятым точкам. Создание заявок - логика usecase
func (ui *Ui) AddAps2Maps(macAp map[string]*entity.Ap, hostnameAp map[string]*entity.Ap) error {

	//devices, errGetDevices := uni.GetDevices(sites) //devices = APs
	devices, errGetDevices := ui.Uni.GetDevices(ui.Sites) //devices = APs
	if errGetDevices == nil {
		log.Println("devices загрузились")
		log.Println("")

		var apPointer *entity.Ap //точка создаётся при каждом взятии из массива
		var ap1 *entity.Ap       //точка из мапы macAp

		var exisAp1 bool
		var exisAp2 bool

		var siteName string
		var apNameUpperCase string

		for _, ap0 := range devices.UAPs {

			if ap0.Name != "" { //не знаю, откуда могут взяться точки в Бд без имени, но логику с двумя мапами может сломать такая точка6

				siteID := ap0.SiteID
				//if !sitesException[siteID] { // НЕ Резерв/Склад
				if siteID == "5e74aaa6a1a76964e770815c" {
					siteName = "Урал" //именно с дефолтными сайтами так почему-то
				} else if siteID == "5e758bdca9f6163bb0c3c962" {
					siteName = "Волга" //именно с дефолтными сайтами так почему-то
				} else {
					siteName = ap0.SiteName[:len(ap0.SiteName)-11]
				}

				apNameUpperCase = strings.ToUpper(ap0.Name)

				ap1, exisAp1 = macAp[ap0.Mac]
				if exisAp1 {
					ap1.Name = apNameUpperCase //ap0.Name
					ap1.SiteName = siteName
					ap1.SiteID = siteID
					ap1.StateInt = ap0.State.Int()
					//k.Exception = ap. //исключение должно приходить от контроллера, но по факту вношу единички в БД
					//Подгрузка единичек исключений по точкам из БД реализована пока что только в самом начале скрипта
					//Периодического обновления из БД пока что нет
					ap1.Controller = ui.Controller

					//проверяем доступность в мапе hostnameAp
					_, exisAp2 = hostnameAp[apNameUpperCase] //ap0.Name]
					if !exisAp2 {
						//Если бакета нет, значит, только сменилось сетевое имя - переназначаем ссылку на client1
						hostnameAp[apNameUpperCase] = ap1
					}

				} else {
					//если мак не бьётся, создаём новую точку
					apPointer = &entity.Ap{
						Mac:          ap0.Mac,
						SiteName:     siteName,
						SiteID:       siteID,
						Name:         apNameUpperCase, //ap0.Name,
						StateInt:     ap0.State.Int(),
						SrID:         "",
						Exception:    0, //исключение для аномалий клиентов
						CommentCount: 0,
						Controller:   ui.Controller,
					}

					macAp[ap0.Mac] = apPointer
					//если мак не бьётся, значит в hostname точки не будет
					hostnameAp[apNameUpperCase] = apPointer
				}
				//} //НЕ Резерв/Склад

			} //if ap0.Name != ""

		} //range devices.UAPs

		return nil //return mapAp, nil
	} else {
		log.Println("devices НЕ загрузились")
		//return mapAp, errGetDevices
		return errGetDevices
	}
}

func (ui *Ui) AddAps(mapAp map[string]*entity.Ap) error {

	//devices, errGetDevices := uni.GetDevices(sites) //devices = APs
	devices, errGetDevices := ui.Uni.GetDevices(ui.Sites) //devices = APs
	if errGetDevices == nil {
		log.Println("devices загрузились")
		log.Println("")
		for _, ap := range devices.UAPs {
			siteID := ap.SiteID
			//if !sitesException[siteID] { // НЕ Резерв/Склад
			//apSiteName := ap.SiteName
			var siteName string
			if siteID == "5e74aaa6a1a76964e770815c" {
				siteName = "Урал" //именно с дефолтными сайтами так почему-то
			} else if siteID == "5e758bdca9f6163bb0c3c962" {
				siteName = "Волга" //именно с дефолтными сайтами так почему-то
			} else {
				siteName = ap.SiteName[:len(ap.SiteName)-11]
			}

			kap, exis := mapAp[ap.Mac]
			if exis {
				//log.Println(kap.Mac + " kap есть в мапе. Обновление данных")
				//log.Println(ap.Mac + " ap есть в мапе. Обновление данных")
				//log.Println(ap.Name, ap.SiteName, ap.State.Int())
				kap.Name = ap.Name
				kap.SiteName = siteName
				kap.SiteID = siteID
				kap.StateInt = ap.State.Int()
				//k.Exception = ap. //исключение должно приходить от контроллера, но по факту вношу единички в БД
				//Подгрузка единичек исключений по точкам из БД реализована пока что только в самом начале скрипта
				//Периодического обновления из БД пока что нет
			} else {
				//log.Println(ap.Mac + "в мапе нет. Создание новой сущности")
				//log.Println(ap.Name, ap.SiteName, ap.State.Int())
				mapAp[ap.Mac] = &entity.Ap{
					Mac:          ap.Mac,
					SiteName:     siteName,
					SiteID:       siteID,
					Name:         ap.Name,
					StateInt:     ap.State.Int(),
					SrID:         "",
					Exception:    0, //исключение для аномалий клиентов
					CommentCount: 0,
				}
			}
			//} //НЕ Резерв/Склад
		}
		//return mapAp, nil
		return nil
	} else {
		log.Println("devices НЕ загрузились")
		//return mapAp, errGetDevices
		return errGetDevices
	}
}

// При обработке каждого клиента к мапе точек НЕ ПОДКЛЮЧАЮСЬ для получения имени. 2 мапы заполняется
func (ui *Ui) UpdateClients2MapWithoutApMap(macClient map[string]*entity.Client, hostnameClient map[string]*entity.Client, date string) error {
	//Загружает в мапу Клиентов: Hostname, exception, mac Ap

	//clients, errGetClients := uni.GetClients(sites) //client = Notebook or Mobile = machine
	clients, errGetClients := ui.Uni.GetClients(ui.Sites) //client = Notebook or Mobile = machine
	if errGetClients == nil {
		log.Println("clients загрузились")
		log.Println("")

		var clExInt int
		var clPointer *entity.Client //клиент создаётся при каждом взятии из массива
		var client1 *entity.Client   //клиент из мапы macClient

		var exisClient1 bool
		//var exisClient2 bool

		var clientNameUpperCase string

		for _, client0 := range clients {
			//client.ApName //!!! НИЧЕГО не выводит и не содержит!!! Имя точки берётся на основании сравнения мака мапой точек
			//clientMac = client.Mac 	clientName = client.Name  		//clientIP = client.IP		//siteName = client.SiteName

			if !client0.IsGuest.Val {

				if client0.Name != "" { //не знаю, откуда взялись клиенты в Бд без имени, но логику с двумя мапами может сломать такой клиент

					if client0.Noted.Val {
						clientExceptionStr := strings.Split(client0.Note, " ")[0]
						if clientExceptionStr == "Exception" {
							clExInt = 1
						} else {
							clExInt = 0
						}
					}

					clientNameUpperCase = strings.ToUpper(client0.Name)

					//проверяем доступность бакета в мапе macClient
					client1, exisClient1 = macClient[client0.Mac]
					if exisClient1 {
						//если бакет есть, обновляем данные клиента (меняются и в мапе hostnameClient также)
						client1.Hostname = clientNameUpperCase //client0.Hostname
						client1.Exception = clExInt
						client1.ApMac = client0.ApMac
						client1.Modified = date
						client1.Controller = ui.Controller

						/*Был случай с NBKG-BELYAEVA, когда заменили ноутбук, а имя оставил прежнее. Меняю этот блок
						_, exisClient2 = hostnameClient[clientNameUpperCase] //client0.Name]
						if !exisClient2 {
							//Если бакета нет, значит, только сменилось сетевое имя - переназначаем ссылку на client1
							hostnameClient[clientNameUpperCase] = client1
						}*/
						hostnameClient[clientNameUpperCase] = client1

					} else {
						//если мак не бьётся, создаём нового клиента
						clPointer = &entity.Client{
							Mac:        client0.Mac,
							Hostname:   clientNameUpperCase, //client0.Name,
							Exception:  clExInt,
							ApMac:      client0.ApMac,
							SrID:       "",
							Modified:   date,
							Controller: ui.Controller,
						}

						macClient[client0.Mac] = clPointer
						//если мак не бьётся, значит в hostname клиента не будет
						hostnameClient[clientNameUpperCase] = clPointer
					}

				} //if client0.Name != ""

			} else {
				//Если клиент Guest
				//check that ip starts with "192."
				//if not
			}
		}
		//return mapClient, nil
		return nil
	} else {
		log.Println("clients НЕ загрузились")
		//return mapClient, errGetClients
		return errGetClients
	}
	//return
}

// НЕ ИСПОЛЬЗУЕТСЯ. При обработке каждого клиента к мапе точек НЕ ПОДКЛЮЧАЮСЬ для получения имени
func (ui *Ui) UpdateClientsWithoutApMap(mapClient map[string]*entity.Client, date string) error {
	//Загружает в мапу Клиентов: Hostname, exception, mac Ap

	//clients, errGetClients := uni.GetClients(sites) //client = Notebook or Mobile = machine
	clients, errGetClients := ui.Uni.GetClients(ui.Sites) //client = Notebook or Mobile = machine
	if errGetClients == nil {
		log.Println("clients загрузились")
		log.Println("")
		var clExInt int

		for _, client0 := range clients {
			//client.ApName //!!! НИЧЕГО не выводит и не содержит!!! Имя точки берётся на основании сравнения мака мапой точек
			//clientMac = client.Mac 	clientName = client.Name  		//clientIP = client.IP		//siteName = client.SiteName

			if !client0.IsGuest.Val {

				if client0.Noted.Val {
					clientExceptionStr := strings.Split(client0.Note, " ")[0]
					if clientExceptionStr == "Exception" {
						clExInt = 1
					} else {
						clExInt = 0
					}
				}

				//пробегаемся по всей мапе клиентов и обновляем Hostname, exception, mac Ap, modified
				client1, exisClient1 := mapClient[client0.Mac]
				if exisClient1 {
					client1.Hostname = client0.Hostname
					client1.Exception = clExInt
					client1.ApMac = client0.ApMac
					client1.Modified = date
				} else {
					mapClient[client0.Mac] = &entity.Client{
						Mac:       client0.Mac,
						Hostname:  client0.Name,
						Exception: clExInt,
						ApMac:     client0.ApMac,
						SrID:      "",
						Modified:  date,
					}
				}
			} else {
				//Если клиент Guest
			}
		}
		//return mapClient, nil
		return nil
	} else {
		log.Println("clients НЕ загрузились")
		//return mapClient, errGetClients
		return errGetClients
	}
	//return
}

// у Клиента аномалии лежат в массиве, а не мапе и добавляются каждый час последним элнементом
// Новых клиентов и Точек в данном методе не создаю, т.к. это всё будет делаться раз в сутки при загрузке за 30 дней
// Без apName аномалия не нужна. А apName можно получить только через мапу Клиентов. Создав на данном этапе Клиента, всё равно нужных данных не получу
func (ui *Ui) GetHourAnomaliesAddSlice(anomalyHourTime string, mac_Client map[string]*entity.Client, mac_Ap map[string]*entity.Ap) (
	mac_Anomaly map[string]*entity.Anomaly, err error) {
	count := 1 //минус 1 час
	then := time.Now().Add(time.Duration(-count) * time.Hour)

	mac_Anomaly = make(map[string]*entity.Anomaly) //panic: assignment to entry in nil map

	//anomalies, errGetAnomalies := uni.GetAnomalies(sites,
	anomalies, errGetAnomalies := ui.Uni.GetAnomalies(ui.Sites,
		//time.Date(2023, 07, 11, 7, 0, 0, 0, time.Local), time.Now()
		//time.Date(2023, 07, 01, 0, 0, 0, 0, time.Local), //time.Now(),
		then,
	)
	if errGetAnomalies == nil {
		log.Println("anomalies загрузились")
		log.Println("")
		var noutMac string
		var siteNameCut string
		//var hourAnomalySlice map[string][]string

		for _, v := range anomalies {
			//v.Anomaly == всего 1 простая аномалия, Пример: USER_POOR_STREAM_EFF
			noutMac = v.DeviceMAC

			kAnom, exisMacAnom := mac_Anomaly[noutMac]
			if !exisMacAnom {

				//В связи с ведением двух мап маков и hostname это усложняет процесс создания отдельных элементов в различных функциях
				mac_Anomaly[noutMac] = &entity.Anomaly{
					ClientMac:    v.DeviceMAC,
					SliceAnomStr: []string{v.Anomaly},
					SiteName:     v.SiteName, //изменим на втором заходе
					ApMac:        "",         //Бакет с одной записью аномалии в массиве  TimeStr_sliceAnomStr мне не интересен
					ApName:       "",         //Поэтому буду заполнять эти поля, когда появится вторая запись за час по тому же маку
					Exception:    0,          //
					//Controller: ui.Controller,
				}

			} else {
				//выполнится ТОЛЬКО на ВТОРОМ заходе. не на первом, не на третьем
				if len(kAnom.SliceAnomStr) == 1 {

					//Добавить для избавления от двойного USER_DNS_TIMEOUT
					//if !(kAnom.SliceAnomStr[0] == "USER_DNS_TIMEOUT" && v.Anomaly == "USER_DNS_TIMEOUT") {

					//подключаемся к мапе Клиентов
					kClient, exisMacClient := mac_Client[noutMac]
					if exisMacClient {
						//подключаемся к мапе Точек
						kAp, exisMacAp := mac_Ap[kClient.ApMac]
						if exisMacAp {

							siteNameCut = v.SiteName[:len(v.SiteName)-11]
							if strings.Contains(siteNameCut, "Волг") {
								siteNameCut = "Волга"
							} else if strings.Contains(siteNameCut, "Ура") {
								siteNameCut = "Урал"
							}

							//kAnom.DateHour = v.Datetime.Format("2006-01-02 15:04:05") //время опрашивающего сервера
							kAnom.DateHour = anomalyHourTime //время контроллера, откуда взята аномалия
							kAnom.SiteName = siteNameCut
							kAnom.ApMac = kClient.ApMac
							kAnom.ApName = kAp.Name
							kAnom.Exception = kClient.Exception + kAp.Exception
							//kAnom.TimeStr_sliceAnomStr[dateTime] = append(kAnom.TimeStr_sliceAnomStr[dateTime], v.Anomaly)
							kAnom.SliceAnomStr = append(kAnom.SliceAnomStr, v.Anomaly)
							kAnom.Controller = ui.Controller

						} else {
							//Если в мапе macAp нет мака? Создав точку здесь, всё равно не получу её имени и не создам её в мапе hostnameAp
						}
					} else {
						//Если в мапе macClient нет мака клиента из аномалии?
						// Без apName аномалия не нужна. А apName можно получить только через мапу Клиентов. Создав на данном этапе Клиента, всё равно нужных данных не получу
					}
					//} //"USER_DNS_TIMEOUT"

				} else {
					//kAnom.TimeStr_sliceAnomStr[dateTime] = append(kAnom.TimeStr_sliceAnomStr[dateTime], v.Anomaly)
					kAnom.SliceAnomStr = append(kAnom.SliceAnomStr, v.Anomaly)
				}
			}
		}

		for macClient, anomalyStruct := range mac_Anomaly {
			if len(anomalyStruct.SliceAnomStr) > 1 {
				client, exisClient := mac_Client[macClient]
				if exisClient {
					//мак клиента есть в мапе
					client.SliceAnomalies = append(client.SliceAnomalies, anomalyStruct)
				} else {
					/*мака клиента из ряда аномалии нет в мапе Клиентов
					//Не имеешь права создавать структуру клиента в мапе macClient, не продублировав эту же запись в мапу hostnameClient
					mac_Client[macClient] = &entity.Client{
						Mac:            macClient,
						SliceAnomalies: []*entity.Anomaly{anomalyStruct},
					}*/
				}
			}
		}

		return mac_Anomaly, nil
	} else {
		log.Println("anomalies НЕ загрузились")
		return nil, errGetAnomalies
	}
	//return
}

// Подключаться к мапе Точек нужно как можно раньше здесь и сейчас, пока актуальная привязка Клиент-Точка
func (ui *Ui) GetHourAnomalies(mac_Client map[string]*entity.Client, mac_Ap map[string]*entity.Ap) (mac_Anomaly map[string]*entity.Anomaly, err error) {
	count := 1 //минус 1 час
	then := time.Now().Add(time.Duration(-count) * time.Hour)

	//anomalies, errGetAnomalies := uni.GetAnomalies(sites,
	anomalies, errGetAnomalies := ui.Uni.GetAnomalies(ui.Sites,
		//time.Date(2023, 07, 11, 7, 0, 0, 0, time.Local), time.Now()
		//time.Date(2023, 07, 01, 0, 0, 0, 0, time.Local), //time.Now(),
		then,
	)

	mac_Anomaly = make(map[string]*entity.Anomaly) //panic: assignment to entry in nil map

	if errGetAnomalies == nil {
		log.Println("anomalies загрузились")
		log.Println("")
		var noutMac string
		var siteNameCut string
		//var hourAnomalySlice map[string][]string

		for _, v := range anomalies {
			//v.Anomaly == всего 1 простая аномалия, Пример: USER_POOR_STREAM_EFF
			noutMac = v.DeviceMAC

			kAnom, exisMacAnom := mac_Anomaly[noutMac]
			if !exisMacAnom {

				//hourAnomalySlice[dateTime] = []string{v.Anomaly}

				mac_Anomaly[noutMac] = &entity.Anomaly{
					ClientMac:    v.DeviceMAC,
					SliceAnomStr: []string{v.Anomaly},
					SiteName:     v.SiteName, //изменим на втором заходе
					ApMac:        "",         //Бакет с одной записью аномалии в массиве  TimeStr_sliceAnomStr мне не интересен
					ApName:       "",         //Поэтому буду заполнять эти поля, когда появится вторая запись за час по тому же маку
					Exception:    0,          //
				}

			} else {
				//если размер массива == 1
				if len(kAnom.SliceAnomStr) == 1 {
					//подключаемся к мапе Клиентов
					kClient, exisMacClient := mac_Client[noutMac]
					if exisMacClient {
						//подключаемся к мапе Точек
						kAp, exisMacAp := mac_Ap[kClient.ApMac]
						if exisMacAp {

							siteNameCut = v.SiteName[:len(v.SiteName)-11]
							if strings.Contains(siteNameCut, "Волг") {
								siteNameCut = "Волга"
							} else if strings.Contains(siteNameCut, "Ура") {
								siteNameCut = "Урал"
							}

							kAnom.DateHour = v.Datetime.Format("2006-01-02 15:04:05")
							kAnom.SiteName = siteNameCut
							kAnom.ApMac = kClient.ApMac
							kAnom.ApName = kAp.Name
							kAnom.Exception = kClient.Exception + kAp.Exception
							//kAnom.TimeStr_sliceAnomStr[dateTime] = append(kAnom.TimeStr_sliceAnomStr[dateTime], v.Anomaly)
							kAnom.SliceAnomStr = append(kAnom.SliceAnomStr, v.Anomaly)
						}
					}
				} else {
					//kAnom.TimeStr_sliceAnomStr[dateTime] = append(kAnom.TimeStr_sliceAnomStr[dateTime], v.Anomaly)
					kAnom.SliceAnomStr = append(kAnom.SliceAnomStr, v.Anomaly)
				}
			}
		}
		return mac_Anomaly, nil
	} else {
		log.Println("anomalies НЕ загрузились")
		return nil, errGetAnomalies
	}
	//return
}
