// Package entity defines main entities for business logic (services), data base mapping and
// HTTP response objects if suitable. Each logic group entities in own file.
package entity

type Zabbix struct {
	Description string `example:"Зафиксированы сбои в работе системы"`
	Reason      string `example:"Устройство недоступно"`
	MacroRegion string `example:"Центр"`
	Monitoring  string `example:"https://zabbix.com"`
	AlarmType   string `example:"Устройство недоступно"` //Новый или закрытие старой аварии
	Comment     string `example:"любой текст"`

	ID     string `example:""` //Может быть есть ID у аларма
	Status string `example:"Решено"`
	Number string `example:""`
	Url    string `example:"https://zabbix.tele2.ru"`
	Server string `example:"t2rn-mon"`

	//CreateWithoutUserLogin bool //использовать ли UserLogin. Если не использовать, то будут создаваться из под denis.tirskikh

	//SliceAps []*Ap
	//Client   *Client
	VCS *Vcs
}
