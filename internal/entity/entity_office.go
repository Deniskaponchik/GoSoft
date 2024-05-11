// Package entity defines main entities for business logic (services), database mapping and
// HTTP response objects if suitable. Each logic group entities in own file.
package entity

type Office struct {
	Site_ApCutName string `json:"site_apcut" example:"Иркутск_IRK"`
	UserLogin      string `json:"user_login" example:"vasya.pupkin"`
	TimeZone       int    `json:"time_zone"  example:"+5"`
	TimeZoneStr    string //для приёма из формы на странице adminka.html
	Exception      int    `json:"exception" example:"1"` //1-true, 0-false
}
