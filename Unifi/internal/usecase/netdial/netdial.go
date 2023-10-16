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

func (pnd *PolyNetDial) NetDialTmtErr(polyStruct entity.PolyStruct) (entity.PolyStruct, error) { //status string, err error) {
	//timeout := 1 * time.Second
	//url := ipString + ":" + "80"
	url := polyStruct.IP + ":" + "80"

	var err error
	myError := 1
	for myError != 0 {
		//	Dial("tcp", "golang.org:netdial")
		//	Dial("tcp", "192.0.2.1:netdial")
		//	Dial("tcp", "198.51.100.1:80")
		//conn, err := netdial.DialTimeout("tcp","mysyte:myport", timeout)
		//conn, err := netdial.DialTimeout("tcp", url, timeout)

		//_, errNetDial := net.DialTimeout("tcp", url, timeout)
		_, errNetDial := net.DialTimeout("tcp", url, pnd.timeout)
		if errNetDial == nil {
			//status = "Registered" //такой же статус возвращает Codec
			polyStruct.Status = "Registered" //такой же статус возвращает Codec
			myError = 0
			return polyStruct, nil
		} else {
			//log.Println("Site unreachable, error: ", err)
			fmt.Println("Visual не доступен по http")
			fmt.Println("Будет предпринята новая попытка отправки запроса через 30 сек.")
			time.Sleep(30 * time.Second)
			myError++
			err = errNetDial
		}
		//conn.Close()
		if myError == 4 {
			myError = 0
			fmt.Println("После 3 неудачных попыток идём дальше. Получить статус работы skype не удалось")
			//return "", err
			return polyStruct, err
		}
	}
	//return status, nil
	return polyStruct, nil
}
