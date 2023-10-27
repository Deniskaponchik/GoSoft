package main

import "fmt"

func main545344645() {
	wifiConf := NewWiFiConfig()
	//fmt.Println(wifiConf.OneDrive)
	fmt.Println(wifiConf.UnifiUsername)
	fmt.Println(wifiConf.UnifiPassword)

	/*
		wifiConfExt := NewWiFiConfigExt()
		fmt.Println(wifiConf.WiFi.UnifiUsername)
		fmt.Println(wifiConf.WiFi.UnifiPassword)
		fmt.Println(wifiConf.DebugMode)
		fmt.Println(wifiConf.MaxUsers)

		// Print out each role
		for _, role := range wifiConf.UserRoles {
			fmt.Println(role)
		}
	*/
}
