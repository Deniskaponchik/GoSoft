// Package entity defines main entities for business logic (services), data base mapping and
// HTTP response objects if suitable. Each logic group entities in own file.
package entity

type Ticket struct {
	UserLogin    string `example:"vasya.pupkin"`
	Description  string `example:"Зафиксированы сбои в работе системы"`
	Reason       string `example:"Устройство недоступно"`
	Region       string `example:"Москва"`
	Monitoring   string `example:"https://zabbix.com"`
	IncidentType string `example:"Устройство недоступно"`
	Comment      string `example:"любой текст"`

	ID        string `example:"aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"`
	Status    string `example:"Решено"`
	Number    string `example:"SR12345678"`
	Url       string `example:"https://bpm.com/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"`
	BpmServer string `example:"https://bpm.com/"`

	//CreateWithoutUserLogin bool //использовать ли UserLogin. Если не использовать, то будут создаваться из под denis.tirskikh

	SliceAps []*Ap
	Client   *Client
}
