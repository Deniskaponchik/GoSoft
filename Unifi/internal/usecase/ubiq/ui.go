package ubiq

import (
	"fmt"
	"github.com/deniskaponchik/GoSoft/Unifi/internal/entity"
	"github.com/unpoller/unifi"
	"log"
)

type Ui struct{
	//unf unifi.Unifi
	unfConf unifi.Config
}

func NewUi(u string, p string, url string) *Ui{
	unfConf := unifi.Config{
		User: u,
		Pass: p,
		URL: url,
		ErrorLog: log.Printf,
		DebugLog: log.Printf,
	}
	return &Ui{
		unfConf: unfConf,
	}
}

func (ui *Ui) GetUni()(unifi.Unifi, error){
	uni, errNewUnifi := unifi.NewUnifi(&ui.unfConf) //&c)
	if errNewUnifi == nil {
		fmt.Println("uni загрузился")
	} else{

	}
	return *uni, nil
}

func (ui *Ui) GetApsArr() ([]entity.Ap, error){
	ui.unf.
	return
}
