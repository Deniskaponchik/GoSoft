// Package entity defines main entities for business logic (services), database mapping and
// HTTP response objects if suitable. Each logic group entities in own file.
package entity

type Office struct {
	ID                int    `json:"ID" example:"1"`
	UsrHelpDeskRegion string `json:"UsrHelpDeskRegion" example:"Москва ЦФ"`
	Site_ApCutName    string `json:"unifi" example:"Иркутск_IRK"`  //"site_apcut"
	UserLogin         string `json:"login" example:"vasya.pupkin"` //"user_login"
	TimeZone          int    `json:"timezone"  example:"+5"`       //"time_zone"
	TimeZoneStr       string //для приёма из формы на странице adminka.html
	//Exception    		int    `json:"exception" example:"1"` //1-true, 0-false
}
