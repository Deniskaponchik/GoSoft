package ubiq

import (
	"fmt"
	"github.com/deniskaponchik/GoSoft/Unifi/internal/entity"
	"github.com/unpoller/unifi"
	"log"
	"strings"
	"time"
)

type Ui struct {
	//unf unifi.Unifi
	Conf  unifi.Config
	Uni   *unifi.Unifi
	Sites []*unifi.Site
}

func NewUi(u string, p string, url string) *Ui {
	fmt.Println(url)

	unfConf := unifi.Config{
		User:     u,
		Pass:     p,
		URL:      url,
		ErrorLog: log.Printf,
		DebugLog: log.Printf,
	}
	return &Ui{
		Conf: unfConf,
	}
}

func (ui *Ui) GetSites() (err error) { //unifi.Unifi, error){
	uni, errNewUnifi := unifi.NewUnifi(&ui.Conf) //&c)
	if errNewUnifi == nil {
		fmt.Println("uni загрузился")
		ui.Uni = uni
		sites, errGetSites := uni.GetSites()
		if errGetSites == nil {
			fmt.Println("sites загрузились")
			ui.Sites = sites
			return nil
		} else {
			fmt.Println("sites НЕ загрузились")
			return errGetSites
		}
	} else {
		fmt.Println("uni НЕ загрузился")
		return errNewUnifi
	}
	//return nil
}

func (ui *Ui) AddAps(mapAp map[string]*entity.Ap) error {

	//devices, errGetDevices := uni.GetDevices(sites) //devices = APs
	devices, errGetDevices := ui.Uni.GetDevices(ui.Sites) //devices = APs
	if errGetDevices == nil {
		fmt.Println("devices загрузились")
		fmt.Println("")
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
				//fmt.Println(kap.Mac + " kap есть в мапе. Обновление данных")
				//fmt.Println(ap.Mac + " ap есть в мапе. Обновление данных")
				//fmt.Println(ap.Name, ap.SiteName, ap.State.Int())
				kap.Name = ap.Name
				kap.SiteName = siteName
				kap.SiteID = siteID
				kap.StateInt = ap.State.Int()
				//k.Exception = ap. //исключение должно приходить от контроллера, но по факту вношу единички в БД
				//Подгрузка единичек исключений по точкам из БД реализована пока что только в самом начале скрипта
				//Периодического обновления из БД пока что нет
			} else {
				//fmt.Println(ap.Mac + "в мапе нет. Создание новой сущности")
				//fmt.Println(ap.Name, ap.SiteName, ap.State.Int())
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
		fmt.Println("devices НЕ загрузились")
		//return mapAp, errGetDevices
		return errGetDevices
	}
}

// При обработке каждого клиента к мапе точек НЕ ПОДКЛЮЧАЮСЬ для получения имени. 2 мапы заполняется
func (ui *Ui) Update2MapClientsWithoutApMap(macClient map[string]*entity.Client, hostnameClient map[string]*entity.Client, date string) error {
	//Загружает в мапу Клиентов: Hostname, exception, mac Ap

	//clients, errGetClients := uni.GetClients(sites) //client = Notebook or Mobile = machine
	clients, errGetClients := ui.Uni.GetClients(ui.Sites) //client = Notebook or Mobile = machine
	if errGetClients == nil {
		fmt.Println("clients загрузились")
		fmt.Println("")
		var clExInt int
		var clPointer *entity.Client //клиент создаётся при каждом взятии из массива
		var client1 *entity.Client   //клиент из мапы macClient
		//var client2 *entity.Client   //клиент из мапы hostnameClient
		var exisClient1 bool
		var exisClient2 bool

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

				//проверяем доступность бакета в мапе macClient
				client1, exisClient1 = macClient[client0.Mac]
				if exisClient1 {
					//если бакет есть, обновляем данные клиента (меняются и в мапе hostnameClient также)
					client1.Hostname = client0.Hostname
					client1.Exception = clExInt
					client1.ApMac = client0.ApMac
					client1.Modified = date

					//проверяем доступность в мапе hostnameClient
					_, exisClient2 = hostnameClient[client0.Name]
					if !exisClient2 {
						//Если бакета нет, значит, только сменилось сетевое имя - переназначаем ссылку на client1
						hostnameClient[client0.Name] = client1
					}

				} else {
					//если мак не бьётся, создаём нового клиента
					clPointer = &entity.Client{
						Mac:       client0.Mac,
						Hostname:  client0.Name,
						Exception: clExInt,
						ApMac:     client0.ApMac,
						SrID:      "",
						Modified:  date,
					}

					macClient[client0.Mac] = clPointer
					//если мак не бьётся, значит в hostname клиента не будет
					hostnameClient[client0.Name] = clPointer
				}

			} else {
				//Если клиент Guest
			}
		}
		//return mapClient, nil
		return nil
	} else {
		fmt.Println("clients НЕ загрузились")
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
		fmt.Println("clients загрузились")
		fmt.Println("")
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
		fmt.Println("clients НЕ загрузились")
		//return mapClient, errGetClients
		return errGetClients
	}
	//return
}

// При обработке кажого клиента идёт подключение к мапе точек для получения имени. НЕ ИСПОЛЬЗУЕТСЯ
func (ui *Ui) AddClientsDeprecated(mapAp map[string]*entity.Ap, mapClient map[string]*entity.Client) (err error) {
	//clients, errGetClients := uni.GetClients(sites) //client = Notebook or Mobile = machine
	clients, errGetClients := ui.Uni.GetClients(ui.Sites) //client = Notebook or Mobile = machine
	if errGetClients == nil {
		fmt.Println("clients загрузились")
		fmt.Println("")
		var apName string //		var clientMac string		var clientName string
		var apException int

		for _, client := range clients {
			//client.ApName //!!! НИЧЕГО не выводит и не содержит!!! Имя точки берётся ниже на основании сравнения мапой точек
			//clientMac = client.Mac 	clientName = client.Name  		//clientIP = client.IP
			//siteName = client.SiteName
			//SiteName нужен только на этапе создания заявок по клиентам. Поэтому при обработке каждого клиента его не получаю.

			if !client.IsGuest.Val {
				var clExInt int
				if client.Noted.Val {
					clientExceptionStr := strings.Split(client.Note, " ")[0]
					if clientExceptionStr == "Exception" {
						clExInt = 1
					} else {
						clExInt = 0
					}
				}

				k, exisApMac := mapAp[client.ApMac]
				if exisApMac {
					apName = k.Name
					apException = k.Exception

					//пробегаемся по всей мапе клиентов и добавляем имя точки клиенту
					kcl, exis := mapClient[client.Mac]
					if exis {
						kcl.Hostname = client.Hostname
						kcl.ApName = apName
						kcl.Exception = clExInt + apException
					} else {
						mapClient[client.Mac] = &entity.Client{
							Mac:       client.Mac,
							Hostname:  client.Name,
							ApName:    apName,
							Exception: clExInt + apException,
							SrID:      "",
						}
					}
				} else {
					fmt.Println("В мапе точек не удалось найти соответствие с маком точки, взятым у клиента")
				}

				/*код предыдущего поколения
				//пробегаемся по всей мапе точек и получаем имя соответствию мака
				for k, v := range mapAp { //apMyMap {
					if k == client.ApMac { //clientMac {
						apName = v.Name
						apException := v.Exception

						//пробегаемся по всей мапе клиентов и назначаем имя точки клиенту
						kcl, exis := mapClient[client.Mac]
						if exis {
							kcl.Hostname = client.Hostname
							kcl.ApName = apName
							kcl.Exception = clExInt + apException
						} else {
							mapClient[client.Mac] = &entity.Client{
								Mac:       client.Mac,
								Hostname:  client.Name,
								SrID:      "",
								Exception: clExInt + apException,
								ApName:    apName,
							}
						}
						break //прекращаем цикл, когда найден мак точки
					}
				}
				*/
			} /* До будущих времён, когда буду обрабатывать Гостевых Клиентов
			else {
				//Если клиент Guest
				splitIP := strings.Split(clientIP, ".")[0]
				if splitIP == "169" {
					forGuestClientTicket := ForGuestClientTicket{
						clientMac,
						clientName,
						clientIP,
					}

					//Заносим в мапу для заявки
					_, exisRegion := region_guestClients[region]
					if exisRegion {
						for k, v := range region_guestClients {
							if k == region {
								v = append(v, forGuestClientTicket)
								region_guestClients[k] = v
								break
							}
						}
					} else {
						forGuestClientTicketSlice := []ForGuestClientTicket{
							forGuestClientTicket,
						}
						region_guestClients[region] = forGuestClientTicketSlice
					}
				}
			}*/
		}
		return nil
	} else {
		fmt.Println("clients НЕ загрузились")
		return errGetClients
	}
	//return
}

// Изменения согласно логике, когда у Клиента несколько аномалий лежат в массиве
func (ui *Ui) GetHourAnomaliesAddSlice(mac_Client map[string]*entity.Client, mac_Ap map[string]*entity.Ap) (mac_Anomaly map[string]*entity.Anomaly, err error) {
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
		fmt.Println("anomalies загрузились")
		fmt.Println("")
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

		for macClient, anomalyStruct := range mac_Anomaly {
			client, exisClient := mac_Client[macClient]
			if exisClient {
				//мак клиента есть в мапе
				client.SliceAnomalies = append(client.SliceAnomalies, anomalyStruct)
			} else {
				//мака клиента из ряда аномалии нет в мапе Клиентов
				mac_Client[macClient] = &entity.Client{
					Mac:            macClient,
					SliceAnomalies: []*entity.Anomaly{anomalyStruct},
				}
			}
		}

		return mac_Anomaly, nil
	} else {
		fmt.Println("anomalies НЕ загрузились")
		return nil, errGetAnomalies
	}
	//return
}

// Изменения согласно логике, когда у Клиента несколько аномалий лежат в массиве
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
		fmt.Println("anomalies загрузились")
		fmt.Println("")
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
		fmt.Println("anomalies НЕ загрузились")
		return nil, errGetAnomalies
	}
	//return
}

// при получении первой аномалии по клиенту ничего не заполняю и никуда не подключаюсь. Всё делается при втором заходе
/*func (ui *Ui) GetHourAnomalies(mac_Client map[string]*entity.Client, mac_Ap map[string]*entity.Ap) (mac_Anomaly map[string]*entity.Anomaly, err error) {
	count := 1 //минус 1 час
	then := time.Now().Add(time.Duration(-count) * time.Hour)

	//anomalies, errGetAnomalies := uni.GetAnomalies(sites,
	anomalies, errGetAnomalies := ui.Uni.GetAnomalies(ui.Sites,
		//time.Date(2023, 07, 11, 7, 0, 0, 0, time.Local), time.Now()
		//time.Date(2023, 07, 01, 0, 0, 0, 0, time.Local), //time.Now(),
		then,
	)
	if errGetAnomalies == nil {
		fmt.Println("anomalies загрузились")
		fmt.Println("")
		var noutMac string
		var siteNameCut string
		var hourAnomalySlice map[string][]string
		var dateTime string

		for _, v := range anomalies {
			//v.Anomaly == всего 1 простая аномалия, Пример: USER_POOR_STREAM_EFF
			dateTime = v.Datetime.Format("2006-01-02 15:04:05")
			noutMac = v.DeviceMAC

			kAnom, exisMacAnom := mac_Anomaly[noutMac]
			if !exisMacAnom {

				hourAnomalySlice[dateTime] = []string{v.Anomaly}

				mac_Anomaly[noutMac] = &entity.Anomaly{
					ClientMac:            v.DeviceMAC,
					TimeStr_sliceAnomStr: hourAnomalySlice,
					SiteName:             v.SiteName, //изменим на втором заходе
					ApMac:                "",         //Бакет с одной записью аномалии в массиве  TimeStr_sliceAnomStr мне не интересен
					ApName:               "",         //Поэтому буду заполнять эти поля, когда появится вторая запись за час по тому же маку
					Exception:            0,          //
				}

			} else {
				//если размер массива == 1
				if len(kAnom.TimeStr_sliceAnomStr[dateTime]) == 1 {
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

							kAnom.SiteName = siteNameCut
							kAnom.ApMac = kClient.ApMac
							kAnom.ApName = kAp.Name
							kAnom.Exception = kClient.Exception + kAp.Exception
							kAnom.TimeStr_sliceAnomStr[dateTime] = append(kAnom.TimeStr_sliceAnomStr[dateTime], v.Anomaly)
						}
					}
				} else {
					kAnom.TimeStr_sliceAnomStr[dateTime] = append(kAnom.TimeStr_sliceAnomStr[dateTime], v.Anomaly)
				}
			}
		}
		return mac_Anomaly, nil
	} else {
		fmt.Println("anomalies НЕ загрузились")
		return nil, errGetAnomalies
	}
	//return
}*/

// Изменил порядок подключения к мапам: mac_Anomaly -> mac_Client -> mac_Ap
/*func (ui *Ui) GetHourAnomalies(mac_Client map[string]*entity.Client, mac_Ap map[string]*entity.Ap) (mac_Anomaly map[string]*entity.Anomaly, err error) {
	count := 1 //минус 1 час
	then := time.Now().Add(time.Duration(-count) * time.Hour)

	//anomalies, errGetAnomalies := uni.GetAnomalies(sites,
	anomalies, errGetAnomalies := ui.Uni.GetAnomalies(ui.Sites,
		//time.Date(2023, 07, 11, 7, 0, 0, 0, time.Local), time.Now()
		//time.Date(2023, 07, 01, 0, 0, 0, 0, time.Local), //time.Now(),
		then,
	)
	if errGetAnomalies == nil {
		fmt.Println("anomalies загрузились")
		fmt.Println("")
		var noutMac string
		var siteNameCut string
		var hourAnomalySlice map[string][]string
		var dateTime string

		for _, v := range anomalies {
			//v.Anomaly == всего 1 простая аномалия, Пример: USER_POOR_STREAM_EFF
			dateTime = v.Datetime.Format("2006-01-02 15:04:05")
			noutMac = v.DeviceMAC

			kAnom, exisMacAnom := mac_Anomaly[noutMac]
			if !exisMacAnom {
				kClient, exisMacClient := mac_Client[noutMac]
				if exisMacClient {

					siteNameCut = v.SiteName[:len(v.SiteName)-11]
					if strings.Contains(siteNameCut, "Волг") {
						siteNameCut = "Волга"
					} else if strings.Contains(siteNameCut, "Ура") {
						siteNameCut = "Урал"
					}

					hourAnomalySlice[dateTime] = []string{v.Anomaly}

					mac_Anomaly[noutMac] = &entity.Anomaly{
						ClientMac:            v.DeviceMAC,
						SiteName:             siteNameCut, //v.SiteName,
						TimeStr_sliceAnomStr: hourAnomalySlice,
						ApMac:                kClient.Mac,       //в новой логике мапа Клиентов не будет запрашивать сразу имя точки. только мак
						ApName:               "",                //на первом текущем заходе неизвестно. Будем подключаться к мапе точек на втором заходе
						Exception:            kClient.Exception, //пока заносим исключения клиента. на втором заходе добавим искл точки
					}

				} else {
					fmt.Println("мак в мапе Клиентов не найден. Дальнейшее создание записи аномалии невозможно")
				}
			} else {
				//если размер массива == 1
				//подключаемся к мапе Клиентов
				//подключаемся к мапе Точек
				//Увеличиваем размер массива
				//в обратном случае
				//Увеличиваем размер массива

				//если запись в мапе Аномалий уже создана
				kAnom.TimeStr_sliceAnomStr[dateTime] = append(kAnom.TimeStr_sliceAnomStr[dateTime], v.Anomaly)

				if len(kAnom.TimeStr_sliceAnomStr[dateTime]) == 2 {

					//подключаемся к мапе точек для получения Exception, а заодно и имени точки
					kAp, exisMac3 := mac_Ap[kAnom.ApMac]
					if exisMac3 {
						kAnom.ApName = kAp.Name
						kAnom.Exception = kAnom.Exception + kAp.Exception
					} else {
						kAnom.ApName = "undefined"
						//exceptionSumm = k1.Exception
					}
				}
			}
		}
		return mac_Anomaly, nil
	} else {
		fmt.Println("anomalies НЕ загрузились")
		return nil, errGetAnomalies
	}
	//return
}*/
