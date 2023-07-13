package main

import (
	"fmt"
	"github.com/unpoller/unifi"
	"log"
	"time"
)

func main() {
	//c := *unifi.Config{
	c := unifi.Config{
		User: "unifi",
		Pass: "FORCEpower23",
		//URL:  "https://10.78.221.142:8443/", //ROSTOV
		URL: "https://10.8.176.8:8443/", //NOVOSIB
		// Log with log.Printf or make your own interface that accepts (msg, test)
		ErrorLog: log.Printf,
		DebugLog: log.Printf,
	}

	clientMacName := map[string]string{}  // clientMAC  -> clientName
	apMacName := map[string]string{}      // apMac      -> apName
	namesClientAps := map[string]string{} // clientName -> apName

	// DO WHILE раз в час
	//clientNameAnomaly := map[string]string{}

	//
	//
	//uni, err := unifi.NewUnifi(c)
	uni, err := unifi.NewUnifi(&c) //в аргументах функций обычно всегда используется &
	if err != nil {
		log.Fatalln("Error:", err)
	}
	sites, err := uni.GetSites()
	if err != nil {
		log.Fatalln("Error:", err)
	}
	log.Println(len(sites), "Unifi Sites Found: ", sites)

	//
	//
	devices, err := uni.GetDevices(sites) //devices = APs
	if err != nil {
		log.Fatalln("Error:", err)
	}
	/* ORIGINAL
	log.Println(len(devices.UAPs), "Unifi Wireless APs Found:")
	for i, uap := range devices.UAPs {
		log.Println(i+1, uap.Name, uap.IP, uap.Mac)
	}*/
	// Добавляем маки и имена точек в map
	for _, uap := range devices.UAPs {
		_, existence := apMacName[uap.Mac] //проверяем, есть ли мак в мапе
		if !existence {
			apMacName[uap.Mac] = uap.Name
		}
	}
	//Вывести AP мапу на экран
	for k, v := range apMacName {
		//fmt.Printf("key: %d, value: %t\n", k, v)
		fmt.Println(k, v)
	}

	//
	//
	clients, err := uni.GetClients(sites)
	if err != nil {
		log.Fatalln("Error:", err)
	}
	//Выбрать из клиентов только КОРП
	clientsCorp := make([]*unifi.Client, 0)
	//
	for _, clientCorp := range clients {
		if !clientCorp.IsGuest.Val {
			clientsCorp = append(clientsCorp, clientCorp)
		}
	}
	//
	for _, clientCorp := range clientsCorp {
		/* Проверка на существование мака в мапе не обязательна
		_, existence := clientMacName[clientCorp.Mac] //проверяем, есть ли мак в мапе
		if !existence { //если нет, добавляем в мапу
			clientMacName[clientCorp.Mac] = clientCorp.Hostname
		}*/
		clientMacName[clientCorp.Mac] = clientCorp.Hostname //Добавить КОРП клиентов в map
	}
	/* ORIGINAL
	log.Println(len(clients), "Clients connected:")
	for i, client := range clients {
		log.Println(i+1, client.SiteName, client.IsGuest.Val, client.Mac, client.Hostname, client.IP, client.LastSeen, client.Anomalies) //i+1
	}*/
	for _, clientCorp := range clientsCorp {
		siteName := clientCorp.SiteName[:len(clientCorp.SiteName)-11]
		apHostName := apMacName[clientCorp.ApMac]
		fmt.Println(siteName, apHostName, clientCorp.Hostname, clientCorp.Mac, clientCorp.IP, clientCorp.LastSeen)

		clientMacName[clientCorp.Mac] = clientCorp.Hostname //Добавить КОРП клиентов в map
		namesClientAps[clientCorp.Name] = clientCorp.ApName //Соответсвие имён клиентов и точек
	}
	//Вывести CLIENT мапу на экран
	for k, v := range clientMacName {
		//fmt.Printf("key: %d, value: %t\n", k, v)
		fmt.Println(k, v)
	}
	//Вывести соответсвие имён клиентов и имён точек на экран
	for k, v := range namesClientAps {
		//fmt.Printf("key: %d, value: %t\n", k, v)
		fmt.Println(k, v)
	}

	//
	//
	anomalies, err := uni.GetAnomalies(sites,
		time.Date(2023, 07, 10, 7, 0, 0, 0, time.Local), time.Now())
	if err != nil {
		log.Fatalln("Error:", err)
	}
	/* ORIGINAL
	log.Println(len(anomalies), "Anomalies:")
	for i, anomaly := range anomalies {
		log.Println(i+1, anomaly.Datetime, anomaly.DeviceMAC, anomaly.Anomaly) //i+1
	}*/
	for _, anomaly := range anomalies {
		_, existence := clientMacName[anomaly.DeviceMAC] //проверяем, есть ли мак в мапе corp clients
		if existence {                                   //если есть, выводим на экран с именем ПК, взятым из мапы
			siteName := anomaly.SiteName[:len(anomaly.SiteName)-11]
			clientHostName := clientMacName[anomaly.DeviceMAC]
			apHostName := namesClientAps[clientHostName]
			fmt.Println(siteName, clientHostName, apHostName, anomaly.Datetime, anomaly.Anomaly)
			/*
				reflect.ValueOf(clientMacName[anomaly.DeviceMAC]) := BpmTicket{
					anomaly.SiteName[:len(anomaly.SiteName)-11],
					clientMacName[anomaly.DeviceMAC],
					anomaly.DeviceMAC,
					append(anomalies, anomaly.Anomaly)
				}*/
		}
	}
}

/*
func GetClientsCorpWithAnomalies(anoms []*Anomaly) ([]*ClientCorp) {
	return
}*/

type BpmTicket struct {
	site          string
	apName        string
	mac           string
	corpAnomalies []string
}
