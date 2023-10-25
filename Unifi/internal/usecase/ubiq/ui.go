package ubiq

import (
	"fmt"
	"github.com/deniskaponchik/GoSoft/Unifi/internal/entity"
	"github.com/unpoller/unifi"
	"log"
	"strings"
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
	sitesException := map[string]bool{
		"5f2285f3a1a7693ae6139c00": true, //Novosib. Резерв/Склад
		"5f5b49d1a9f6167b55119c9b": true, //Ростов. Резерв/Склад
	}
	//devices, errGetDevices := uni.GetDevices(sites) //devices = APs
	devices, errGetDevices := ui.Uni.GetDevices(ui.Sites) //devices = APs
	if errGetDevices == nil {
		fmt.Println("devices загрузились")
		fmt.Println("")
		for _, ap := range devices.UAPs {
			siteID := ap.SiteID
			if !sitesException[siteID] { // НЕ Резерв/Склад

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
					kap.StateInt = ap.State.Int()
					//k.Exception = ap. //в идеале должен прилетать от контроллера, но в жизни вношу его в БД руками
					//mapAp[ap.Mac] = kap
				} else {
					mapAp[ap.Mac] = &entity.Ap{
						Mac:       ap.Mac,
						SiteName:  siteName,
						Name:      ap.Name,
						StateInt:  ap.State.Int(),
						SrID:      "",
						Exception: 0,
						Comment:   0,
					}
				}
			}
		}
	} else {
		fmt.Println("devices НЕ загрузились")
		return errGetDevices
	}
	return
}

func (ui *Ui) AddClients(mapAp map[string]*entity.Ap, mapClient map[string]*entity.Client) (err error) {
	//clients, errGetClients := uni.GetClients(sites) //client = Notebook or Mobile = machine
	clients, errGetClients := ui.Uni.GetClients(ui.Sites) //client = Notebook or Mobile = machine
	if errGetClients == nil {
		fmt.Println("clients загрузились")
		fmt.Println("")

		var apName string //		var clientMac string		var clientName string

		for _, client := range clients {
			//apName = client.ApName //НИЧЕГО не выводит и не содержит. Имя точки берётся ниже на основании сравнения мапой точек
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
				//пробегаемся по всей мапе точек и получаем имя соответствию мака
				for k, v := range mapAp { //apMyMap {
					if k == client.ApMac { //clientMac {
						apName = v.Name
						apException := v.Exception
						//пробегаемся по всей мапе клиентов и назначаем имя точки клиенту

						kcl, exis := mapClient[client.Mac]
						if exis {
							ke.Hostname = client.Hostname

						} else {
							mapClient[client.Mac] = &entity.Client{
								Mac:       client.Mac,
								Hostname:  client.Name,
								SrID:      "",
								Exception: clExInt + apException,
								ApName:    apName,
							}
						}

						_, exisNoutMyMap := machineMyMap[clientMac]
						if !exisNoutMyMap { //если записи клиента НЕТ
							machineMyMap[clientMac] = MachineMyStruct{
								clientName,
								clExInt + apException,
								"",
								apName,
							}
						} else { //если запись клиента создана, обновляем её
							for ke, va := range machineMyMap {
								if ke == client.Mac {
									va.Hostname = clientName
									va.ApName = apName
									va.Exception = clExInt + apException
									machineMyMap[ke] = va
									break //прекращаем цикл, когда найден клиент и имя точки присвоено ему
								}
							}
						}
						break //прекращаем цикл, когда найден мак точки
					}
				}
			} /* До будущих времён, когда буду обрабатывать Клиентов
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
	}
	return
}
