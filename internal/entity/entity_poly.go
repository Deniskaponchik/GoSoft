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
	TimeZone  int    `json:"time_zone"  example:"+5"`

	Status string `example:"Registered"`
}
