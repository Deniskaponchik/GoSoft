// Package entity defines main entities for business logic (services), data base mapping and
// HTTP response objects if suitable. Each logic group entities in own file.
package entity

import "time"

type Ap struct {
	Mac    string `json:"mac"        example:"a0-b1-c2-d3-e4-f5"`
	IP     string `json:"ip"         example:"10.68.24.157"`
	Region string `json:"region"     example:"Волгоград"`
	Name   string `json:"name"       example:"XXX-OPENSPACE"`

	UserLogin string `json:"login"      example:"vasya.pupkin"`
	SrID      string `json:"srid"       example:"aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"`
	Exception int    `json:"exception"  example:"1"`
	//Type  int    `json:"type"       example:"1"`
	//Comment   int    `json:"comment"    example:"1"`
	//Status    string `example:"Available"`
}

type Client struct {
	Mac        string `json:"mac"        example:"a0:b1:c2:d3:e4:f5"`
	Hostname   string `json:"hostname"   example:"XXXX-PUPKIN"`
	SrID       string `json:"srid"       example:"aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"`
	Controller int    `json:"controller" example:"1"`
	Exception  int    `json:"exception"  example:"1"`
	ApName     string `example:"XXX-OPENSPACE"`

	//Аномалии клиента за всё время
	Anomalies []Anomaly

	UserLogin  string `example:"vasya.pupkin"`
	Region     string `example:"Москва"`
	Monitoring string `example:"https://zabbix.com"`
	//Status      string `example:"Доступен"`
	//Comment     string `example:"любой текст"`
	//Description string `example:"Зафиксированы сбои в работе системы"`
}

// Аномалии клиента, накопившиеся за 1 час
type Anomaly struct {
	ClientMac    string    `json:"mac"        example:"a0:b1:c2:d3:e4:f5"`
	SiteName     string    `json:"sitename"   example:"Москва"`
	AnomalySlice []string  `json:"anomalies"  example:"USER_HIGH_TCP_LATENCY;USER_LOW_PHY_RATE;USER_SLEEPY_CLIENT;USER_HIGH_TCP_PACKET_LOSS;USER_HIGH_WIFI_RETRIES;USER_SIGNAL_STRENGTH_FAILURES;USER_DNS_TIMEOUT;USER_HIGH_WIFI_LATENCY;USER_POOR_STREAM_EFF;USER_HIGH_DNS_LATENCY"`
	Controller   int       `json:"controller" example:"1"`
	Exception    int       `json:"exception"  example:"1"`
	ApName       string    `json:"apname"     example:"XXX-OPENSPACE"`
	DateHour     time.Time `json:"date_hour"  example:"2023-09-01 12:00:00"`
}
