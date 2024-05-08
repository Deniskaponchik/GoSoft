// Package entity defines main entities for business logic (services), data base mapping and
// HTTP response objects if suitable. Each logic group entities in own file.
package entity

type Vcs struct {
	Mac         string `json:"mac"        example:"a0-b1-c2-d3-e4-f5"`
	IP          string `json:"ip"         example:"10.68.24.157"`
	RegionRus   string `json:"region"     example:"Волгоград"`
	RoomNameRus string `json:"room_name"  example:"Ахтуба"`
	Login       string `json:"login"      example:"vasya.pupkin"`
	SrID        string `json:"srid"       example:"aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"`
	TypeInt     int    `json:"type_int"   example:"1"`
	TypeStr     string `json:"type_str"   example:"AudioCodes"`
	Comment     int    `json:"comment"    example:"1"`  //был оставлен комментарий или нет
	Exception   int    `json:"exception"  example:"1"`  //добавлено ли в исключение на создание тикетов
	TimeZone    int    `json:"time_zone"  example:"+5"` //from Moscow
	UserLogin   string `example:"vasya.pupkin"`

	StatusTrueConf string `example:"Registered"` //мой выдуманный статус работы TrueConf Room
	StatusSkype    string `example:"Registered"`
}
