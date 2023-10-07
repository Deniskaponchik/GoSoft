package main

import (
	"fmt"
	"github.com/unpoller/unifi"
	"io"
	"log"
	"strings"
)

func main() {
	//c := *unifi.Config{
	c := unifi.Config{
		User: "unifi",
		Pass: "FORCEpower23",
		//URL:  "https://10.78.221.142:8443/", //ROSTOV
		URL: "https://10.8.176.8:8443/", //NOVOSIB
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
	log.Println(len(sites), "Unifi Sites Found: ", sites)

	devices, err := uni.GetDevices(sites)
	if err != nil {
		log.Fatalln("Error:", err)
	}
	log.Println(len(devices.UAPs), "Unifi Wireless APs Found:")
	for i, uap := range devices.UAPs {
		log.Println(i+1, uap.Name, uap.IP, uap.)
	}

	clients, err := uni.GetClients(sites)
	if err != nil {
		log.Fatalln("Error:", err)
	}
	log.Println(len(clients), "Clients connected:")
	for _, client := range clients {
		client.

		splitIP := strings.Split(client.IP, ".")[0]
		//if splitIP == "169" {
		if splitIP != "10" && splitIP != "192" {
			fmt.Println(client.IP, client.Hostname, client.Mac, client.SiteName)
		}
		//log.Println(i+1, client.ID, client.Hostname, client.IP, client.Name)
	}

}
