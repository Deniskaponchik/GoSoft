package main

import (
	"github.com/unpoller/unifi"
	"log"
	"time"
)

func main() {
	//c := *unifi.Config{
	c := unifi.Config{
		User: "unifi",
		Pass: "FORCEpower23",
		URL:  "https://10.78.221.142:8443/",
		// Log with log.Printf or make your own interface that accepts (msg, test)
		ErrorLog: log.Printf,
		DebugLog: log.Printf,
	}
	//uni, err := unifi.NewUnifi(c)
	uni, err := unifi.NewUnifi(&c)
	if err != nil {
		log.Fatalln("Error:", err)
	}
	sites, err := uni.GetSites()
	if err != nil {
		log.Fatalln("Error:", err)
	}
	clients, err := uni.GetClients(sites)
	if err != nil {
		log.Fatalln("Error:", err)
	}
	//Выбрать из клиентов только КОРП
	clientsCorp := make([]unifi.Client, 0)
	for _, client := range clients {
		if !client.IsGuest.Val {
			clientsCorp = append(clientsCorp, client)
		}
	}

	//Добавить КОРП клиентов в таблицу БД

	/*
		devices, err := uni.GetDevices(sites) //devices = APs
		if err != nil {
			log.Fatalln("Error:", err)
		}
	*/

	//anomalies, err := uni.GetAnomalies(sites, time.Now(), time.Date(2023, 07, 10, 16, 0, 0, 0, time.Local))
	anomalies, err := uni.GetAnomalies(sites, time.Date(2023, 07, 10, 16, 0, 0, 0, time.Local), time.Now())
	if err != nil {
		log.Fatalln("Error:", err)
	}
	//anomaliesExt :=

	log.Println(len(sites), "Unifi Sites Found: ", sites)

	log.Println(len(clients), "Clients connected:")
	for i, client := range clients {
		log.Println(i+1, client.SiteName, client.IsGuest.Val, client.Mac, client.Hostname, client.IP, client.LastSeen, client.Anomalies) //i+1
	}
	log.Println(len(anomalies), "Anomalies:")
	for i, anomaly := range anomalies {
		//for
		log.Println(i+1, anomaly.Datetime, anomaly.DeviceMAC, anomaly.Anomaly) //i+1
	}

	/*
		log.Println(len(devices.UAPs), "Unifi Wireless APs Found:")
		for i, uap := range devices.UAPs {
			log.Println(i+1, uap.Name, uap.IP)
		}
	*/


	func getClientsCorp (clients) ([]unifi.Client clientsCorp) {
		data := make([]unifi.Client, 0)
		for _, client := range clients { // value only
			if !client.IsGuest.Val {
				data = append(data, client)
			}
		}
		return data
	}
	//func getAnomaliesExtend (clients, anomalies) ([]*AnomalyExt) {}


}
