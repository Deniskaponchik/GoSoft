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

func (ui *Ui) AddAps(mapAp map[string]*entity.Ap) (err error) {
	/*
		sitesException := map[string]bool{
			"5f2285f3a1a7693ae6139c00": true, //Novosib. Резерв/Склад
			"5f5b49d1a9f6167b55119c9b": true, //Ростов. Резерв/Склад
		}*/
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
				kap.Name = ap.Name
				kap.SiteName = siteName
				kap.SiteID = siteID
				kap.StateInt = ap.State.Int()
				//k.Exception = ap. //в идеале должен прилетать от контроллера, но в жизни вношу его в БД руками
			} else {
				mapAp[ap.Mac] = &entity.Ap{
					Mac:          ap.Mac,
					SiteName:     siteName,
					SiteID:       siteID,
					Name:         ap.Name,
					StateInt:     ap.State.Int(),
					SrID:         "",
					Exception:    0,
					CommentCount: 0,
				}
			}
			//} //НЕ Резерв/Склад
		}
		return nil
	} else {
		fmt.Println("devices НЕ загрузились")
		return errGetDevices
	}
	//return
}

func (ui *Ui) AddClients(mapAp map[string]*entity.Ap, mapClient map[string]*entity.Client) (err error) {
	//clients, errGetClients := uni.GetClients(sites) //client = Notebook or Mobile = machine
	clients, errGetClients := ui.Uni.GetClients(ui.Sites) //client = Notebook or Mobile = machine
	if errGetClients == nil {
		fmt.Println("clients загрузились")
		fmt.Println("")
		var apName string //		var clientMac string		var clientName string
		var apException int

		for _, client := range clients {
			//client.ApName //!!! НИЧЕГО не выводит и не содержит!!! Имя точки берётся ниже на основании сравнения мапой точек
			//clientMac = client.Mac 	clientName = client.Name  		//clientIP = client.IP		//siteName = client.SiteName

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

func (ui *Ui) GetHourAnomalies(mapClient map[string]*entity.Client) (mac_Anomaly map[string]*entity.Anomaly, err error) {

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

			k1, exisMac1 := mapClient[noutMac] //НЕОБХОДИМО подключение к мапе Клиентов каждый час, чтобы получать актуальный мак точки на этот период времени
			if exisMac1 {
				k2, exisMac2 := mac_Anomaly[noutMac]
				if !exisMac2 {

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
						ApMac:                k1.Mac, //в новой логике мапа Клиентов не будет запрашивать сразу имя точки. только мак
						Exception:            k1.Exception,
					}
				} else {
					//update slice of anomaly string
					k2.TimeStr_sliceAnomStr[dateTime] = append(k2.TimeStr_sliceAnomStr[dateTime], v.Anomaly)
					//k2.SiteName = siteNameCut
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

/* Логика, когда при каждой аномалии идёт подключение к мапе клиентов для получения имени точки
func (ui *Ui) GetHourAnomalies(mapClient map[string]*entity.Client) (maс_Anomaly map[string]*entity.Anomaly, err error) {
	var siteNameCut string
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
		//v.anomaly == всего 1 простая аномалия, Пример: USER_POOR_STREAM_EFF

		for _, v := range anomalies {
			noutMac = v.DeviceMAC
			k1, exisMac1 := mapClient[noutMac]
			if exisMac1 {
				siteNameCut = v.SiteName[:len(v.SiteName)-11]
				if strings.Contains(siteNameCut, "Волг") {
					siteNameCut = "Волга"
				} else if strings.Contains(siteNameCut, "Ура") {
					siteNameCut = "Урал"
				}

				k2, exisMac2 := mapAnomaly[noutMac]
				if !exisMac2 {
					mapAnomaly[noutMac] = &entity.Anomaly{
						ClientMac:    noutMac,
						SiteName:     siteNameCut, //v.SiteName,
						AnomalySlice: []string{v.Anomaly},
						ApName:       k1.ApName,
					}
				} else {
					k2.AnomalySlice = append(k2.AnomalySlice, v.Anomaly)
					k2.SiteName = siteNameCut
				}
			}
		}
		return mapAnomaly, nil
	} else {
		fmt.Println("anomalies НЕ загрузились")
		return nil, errGetAnomalies
	}
	//return
}
*/
