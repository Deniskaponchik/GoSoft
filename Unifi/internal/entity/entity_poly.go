// Package entity defines main entities for business logic (services), data base mapping and
// HTTP response objects if suitable. Each logic group entities in own file.
package entity

type PolyStruct struct {
	Mac       string `json:"mac"        example:"a0-b1-c2-d3-e4-f5"`
	IP        string `json:"ip"         example:"10.68.24.157"`
	Region    string `json:"region"     example:"Волгоград"`
	RoomName  string `json:"room_name"  example:"Ахтуба"`
	Login     string `json:"login"      example:"vasya.pupkin"`
	SrID      string `json:"srid"       example:"aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"`
	PolyType  int    `json:"type"       example:"1"`
	Comment   int    `json:"comment"    example:"1"`
	Exception int    `json:"exception"  example:"1"`

	Status string `example:"Registered"`
}

/*
type PolyTicket struct {
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
}
*/
