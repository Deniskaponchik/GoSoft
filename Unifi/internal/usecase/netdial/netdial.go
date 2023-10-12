package netdial

import (
	"fmt"
	"github.com/deniskaponchik/GoSoft/Unifi/internal/entity"
	"net"
	"time"
)

type PolyNetDial struct {
	timeout time.Duration
	//url     string
}

func New() *PolyNetDial {
	return &PolyNetDial{
		timeout: 1 * time.Second,
	}
}

func (pnd *PolyNetDial) NetDialTmtErr(polyStruct entity.PolyStruct) (status string, err error) {
	//timeout := 1 * time.Second
	//url := ipString + ":" + "80"
	url := polyStruct.IP + ":" + "80"

	myError := 1
	for myError != 0 {
		//	Dial("tcp", "golang.org:netdial")
		//	Dial("tcp", "192.0.2.1:netdial")
		//	Dial("tcp", "198.51.100.1:80")
		//conn, err := netdial.DialTimeout("tcp","mysyte:myport", timeout)
		//conn, err := netdial.DialTimeout("tcp", url, timeout)

		//_, errNetDial := net.DialTimeout("tcp", url, timeout)
		_, errNetDial := net.DialTimeout("tcp", url, pnd.timeout)
		if errNetDial != nil {
			//log.Println("Site unreachable, error: ", err)
			fmt.Println("Visual не доступен по netdial")
			fmt.Println("Будет предпринята новая попытка отправки запроса через 1 минут")
			time.Sleep(30 * time.Second)
			myError++
			err = errNetDial
		} else {
			//fmt.Println("Visual доступен")
			status = "ok"
			myError = 0
		}
		//conn.Close()
		if myError == 4 {
			myError = 0
			fmt.Println("После 3 неудачных попыток идём дальше. Получить статус работы skype не удалось")
			//status = ""
			return "", err
		}
	}

	return status, nil
}

func netDialTmt(ipString string) (status string) {
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
