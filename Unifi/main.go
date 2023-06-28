package main

import (
	"log"

	"github.com/unpoller/unifi"
)

func main() {
	c := *unifi.Config{
		User: "admin",
		Pass: "superSecret1234",
		URL:  "https://127.0.0.1:8443/",
		// Log with log.Printf or make your own interface that accepts (msg, fmt)
		ErrorLog: log.Printf,
		DebugLog: log.Printf,
	}
	uni, err := unifi.NewUnifi(c)
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
	devices, err := uni.GetDevices(sites)
	if err != nil {
		log.Fatalln("Error:", err)
	}

	log.Println(len(sites), "Unifi Sites Found: ", sites)
	log.Println(len(clients), "Clients connected:")
	for i, client := range clients {
		log.Println(i+1, client.ID, client.Hostname, client.IP, client.Name, client.LastSeen)
	}

	log.Println(len(devices.USWs), "Unifi Switches Found")
	log.Println(len(devices.USGs), "Unifi Gateways Found")

	log.Println(len(devices.UAPs), "Unifi Wireless APs Found:")
	for i, uap := range devices.UAPs {
		log.Println(i+1, uap.Name, uap.IP)
	}
}
