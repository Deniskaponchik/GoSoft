package main

import (
	"fmt"
	"net"
	"time"
)

func netDialTmtErr(ipString string) (status string) {
	timeout := 1 * time.Second
	url := ipString + ":" + "80"

	myError := 1
	for myError != 0 {

		//	Dial("tcp", "golang.org:netdial")
		//	Dial("tcp", "192.0.2.1:netdial")
		//	Dial("tcp", "198.51.100.1:80")
		//conn, err := netdial.DialTimeout("tcp","mysyte:myport", timeout)
		//conn, err := netdial.DialTimeout("tcp", url, timeout)
		_, err := net.DialTimeout("tcp", url, timeout)
		if err != nil {
			//log.Println("Site unreachable, error: ", err)
			fmt.Println("Visual не доступен по netdial")
			fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
			time.Sleep(30 * time.Second)
			myError++
		} else {
			//fmt.Println("Visual доступен")
			status = "ok"
			myError = 0
		}
		//conn.Close()
		if myError == 4 {
			myError = 0
			fmt.Println("После 3 неудачных попыток идём дальше. Получить статус работы skype не удалось")
			status = ""
			//statuses = append(statuses, 0)
			//statuses = append(statuses, 0)
		}
	}

	return
}

func netDialTmt(ipString string) (status string) {
	timeout := 1 * time.Second

	url := ipString + ":" + "80"

	//	Dial("tcp", "golang.org:netdial")
	//	Dial("tcp", "192.0.2.1:netdial")
	//	Dial("tcp", "198.51.100.1:80")
	//conn, err := netdial.DialTimeout("tcp","mysyte:myport", timeout)
	//conn, err := netdial.DialTimeout("tcp", url, timeout)
	_, err := net.DialTimeout("tcp", url, timeout)
	if err != nil {
		//log.Println("Site unreachable, error: ", err)
		fmt.Println("Visual не доступен по netdial")
		status = ""
	} else {
		//fmt.Println("Visual доступен")
		status = "ok"
	}
	//conn.Close()
	return
}
