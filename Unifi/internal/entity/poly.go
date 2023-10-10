// Package entity defines main entities for business logic (services), data base mapping and
// HTTP response objects if suitable. Each logic group entities in own file.
package entity

type PolyStruct struct {
	IP        string `json:"ip"         example:"10.68.24.157"`
	Region    string `json:"region"     example:"Волгоград"`
	RoomName  string `json:"roomName"   example:"Ахтуба"`
	Login     string `json:"login"      example:"vasya.pupkin"`
	SrID      string `json:"srID"       example:"aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"`
	PolyType  int    `json:"polyType"   example:"1"`
	Comment   int    `json:"comment"    example:"1"`
	Exception int    `json:"exception"  example:"1"`
}
