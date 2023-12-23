package unpoller

import (
	"github.com/unpoller/unifi"
	"log"
)

type Unpoller struct {
	//username string
	//password string
	uconf *unifi.Config
}

func NewUnpoller(u string, p string, url string) *Unpoller {
	//up := &Unpoller{}
	uc := unifi.Config{
		User: u,   //wifiConf.UnifiUsername,
		Pass: p,   //wifiConf.UnifiPassword,
		URL:  url, //urlController,
		// Log with log.Printf or make your own interface that accepts (msg, test_SOAP)
		ErrorLog: log.Printf,
		DebugLog: log.Printf,
	}
	return &Unpoller{
		uconf: &uc,
	}
}
