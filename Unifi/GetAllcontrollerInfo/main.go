package main

import (
	"fmt"
	"github.com/unpoller/unifi"
	"io"
	"log"
)

func main() {
	//c := *unifi.Config{
	c := unifi.Config{
		User: "unifi",
		Pass: "FORCEpower23",
		URL:  "https://10.78.221.142:8443/", //ROSTOV
		//URL: "https://10.8.176.8:8443/", //NOVOSIB
		// Log with log.Printf or make your own interface that accepts (msg, test_SOAP)
		ErrorLog: log.Printf,
		DebugLog: log.Printf,
	}

	log.SetOutput(io.Discard) //Отключить вывод лога

	//uni, err := unifi.NewUnifi(c)
	uni, err := unifi.NewUnifi(&c)

	if err != nil {
		log.Fatalln("Error:", err)
	}

	sites, err := uni.GetSites()
	if err != nil {
		log.Fatalln("Error:", err)
	}
	for _, site := range sites {
		fmt.Println(site.SiteName)
		fmt.Println(site.ID)
		fmt.Println(site.Name)
		fmt.Println(site.Desc)
		fmt.Println("")
	}

	clients, err := uni.GetClients(sites)
	if err != nil {
		log.Fatalln("Error:", err)
	}
	devices, err := uni.GetDevices(sites)
	if err != nil {
		log.Fatalln("Error:", err)
	}

	log.Println(len(sites), "Unifi Sites Found: ", sites)

	log.Println(len(clients), "Clients connected:")
	for i, client := range clients {
		log.Println(i+1, client.ID, client.Hostname, client.IP, client.Name)
	}

	log.Println(len(devices.USWs), "Unifi Switches Found")

	log.Println(len(devices.USGs), "Unifi Gateways Found")

	log.Println(len(devices.UAPs), "Unifi Wireless APs Found:")
	for i, uap := range devices.UAPs {
		log.Println(i+1, uap.Name, uap.IP)
	}
}
