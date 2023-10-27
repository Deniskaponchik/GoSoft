package main

import (
	//"bytes"
	"fmt"
	"github.com/unpoller/unifi"
	"io"
	"log"
	//"strconv"
	//"strings"
	"time"
)

type MachineMyStruct struct {
	Hostname  string
	Exception int
	SrID      string
	ApName    string
}
type Machine struct {
	site     string
	ApName   string
	Hostname string
	Count    int8
}

func main() {
	fmt.Println("")

	unifiController := 21 //10-Rostov Local; 11-Rostov ip; 20-Novosib Local; 21-Novosib ip
	var urlController string
	var bdController int8 //Да string, потому что значение пойдёт в replace для БД

	//ROSTOV
	if unifiController == 10 || unifiController == 11 {
		bdController = 1
		if unifiController == 10 {
			urlController = "https://localhost:8443/"
		} else {
			urlController = "https://:8443/"
		}

		//NOVOSIB
	} else if unifiController == 20 || unifiController == 21 {
		bdController = 2
		if unifiController == 20 {
			urlController = "https://localhost:8443/"
		} else {
			urlController = "https://:8443/"
		}

	}
	fmt.Println("Unifi controller")
	fmt.Println(urlController)
	fmt.Println(bdController)

	//machineMyMap := map[string]MachineMyStruct{}
	//machineMyMap := DownloadMapFromDBmachines(bdController)

	//fmt.Println("Вывод мапы СНАРУЖИ функции")
	/*
		for k, v := range siteApCutNameLogin {
			//fmt.Printf("key: %d, value: %t\n", k, v)
			fmt.Println("newMap "+k, v)
		}
		os.Exit(0)
	*/

	fmt.Println("")

	c := unifi.Config{
		//c := *unifi.Config{  //ORIGINAL
		User: "",
		Pass: "",
		//URL: "https://localhost:8443/"
		URL: urlController,
		// Log with log.Printf or make your own interface that accepts (msg, test_SOAP)
		ErrorLog: log.Printf,
		DebugLog: log.Printf,
	}

	log.SetOutput(io.Discard) //Отключить вывод лога

	timeNow := time.Now()
	fmt.Println(timeNow.Format("02 January, 15:04:05"))

	//uni, err := unifi.NewUnifi(c)
	uni, err := unifi.NewUnifi(&c)
	if err != nil {
		log.Fatalln("Error:", err)
	} else {
		fmt.Println("uni загрузился")
	}

	sites, err := uni.GetSites()
	if err != nil {
		log.Fatalln("Error:", err)
	} else {
		fmt.Println("sites загрузились")
	}
	/*
		devices, err := uni.GetDevices(sites) //devices = APs
		if err != nil {
			log.Fatalln("Error:", err)
		} else {
			fmt.Println("devices загрузились")
		}

		clients, err := uni.GetClients(sites) //client = Notebook or Mobile = machine
		if err != nil {
			log.Fatalln("Error:", err)
		} else {
			fmt.Println("clients загрузились")
		}
	*/

	count := 60 //минус 70 минут
	//count := 720 //минус 30 день
	//count := 3600
	//count := 36000 //+++
	//count := 86400
	//then := now.Add(time.Duration(-count) * time.Minute)
	//then := timeNow.Add(time.Duration(-count) * time.Minute)
	//then := timeNow.Add(time.Duration(-count) * time.)
	then := timeNow.Add(time.Duration(-count) * time.Hour)

	anomalies, err := uni.GetAnomalies(sites,
		//time.Date(2023, 7, 9, 0, 0, 0, 0, time.Local), //time.Now(),
		then,
	)
	if err != nil {
		log.Fatalln("Error:", err)
	} else {
		fmt.Println("anomalies загрузились")
	}
	fmt.Println("")

	for _, v := range anomalies {
		fmt.Println(v.Datetime)
	}

	//
	/*Для выгрузки в разрезе клиентов
	dateMac_site := map[string]string{}

	var siteName string
	var noutMac string
	var anomalyStr string
	var anomalyDatetime time.Time
	for _, anomaly := range anomalies {
		siteName = anomaly.SiteName
		noutMac = anomaly.DeviceMAC
		anomalyStr = anomaly.Anomaly
		anomalyDatetime = anomaly.Datetime
		fmt.Println(anomalyDatetime, siteName, noutMac, anomalyStr)

		anomalyDatetime.String()
		uniqKey := anomalyDatetime.Format("2006-01-02") + "_" + noutMac
		//dateMac_mac[uniqKey] = noutMac
		dateMac_site[uniqKey] = siteName
	}
	for k, v := range dateMac_site {
		kMac := strings.Split(k, "_")[1]
		for ke, va := range machineMyMap {
			if kMac == ke {
				va.Exception++
				//va.SrID = v
				va.SrID = v[:len(v)-11]
				machineMyMap[ke] = va
			}
		}
	}
	var login string
	var count string
	for _, v := range machineMyMap {
		if v.Exception != 0 {
			login = GetLoginPC(v.Hostname)
			count = strconv.Itoa(int(v.Exception))
			fmt.Println(v.SrID + ";" + v.ApName + ";" + v.Hostname + ";" + login + ";" + count)
		}
	}*/

	//
	/*Для выгрузки в разрезе точек
	dateName_site := map[string]string{}

	var siteName string
	var noutMac string
	//var anomalyStr string
	var anomalyDatetime time.Time
	for _, anomaly := range anomalies {
		anomalyDatetime = anomaly.Datetime
		siteName = anomaly.SiteName
		noutMac = anomaly.DeviceMAC
		//anomalyStr = anomaly.Anomaly
		//fmt.Println(anomalyDatetime, siteName, noutMac, anomalyStr)
		for k, v := range machineMyMap {
			if k == noutMac {
				anomalyDatetime.String()
				uniqKey := anomalyDatetime.Format("2006-01-02") + "+" + v.ApName
				//dateMac_mac[uniqKey] = noutMac
				dateName_site[uniqKey] = siteName //[:len(siteName)-11]
				break
			}
		}
	}
	for k, v := range dateName_site {
		kName := strings.Split(k, "+")[1]
		for ke, va := range machineMyMap {
			if kName == va.ApName {
				va.Exception++
				va.SrID = v[:len(v)-11]
				machineMyMap[ke] = va
				break
			}
		}
	}
	var count string
	for _, v := range machineMyMap {
		if v.Exception != 0 {
			count = strconv.Itoa(int(v.Exception))
			fmt.Println(v.SrID + ";" + v.ApName + ";" + count)
		}
	}*/

} //main func

func cointains(slice []string, compareString string) bool {
	for _, v := range slice {
		if v == compareString {
			return true
		}
	}
	return false
}
