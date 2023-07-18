package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func UploadsMapsToDB() {

}

func DownloadMapsFromDB() {

}

func GetLogin(pcName string) string {
	type PC struct {
		//ID   int    `json:"id"`
		UserName string `json:"user_name"`
		Date_Mod string `json:"date_mod"`
	}

	db, err := sql.Open("mysql", "root:t2root@tcp(10.77.252.153:3306)/glpi_db")
	if err != nil {
		log.Print(err.Error())
	}
	defer db.Close()

	var pc PC
	err = db.QueryRow("SELECT contact FROM glpi_db.glpi_computers where name = ? ORDER BY date_mod DESC", pcName).Scan(&pc.UserName)
	// после запятой указываем значение, которое будет подставляться заметсо вопроса + ОБЯЗАТЕЛЬНО в Scan использовать &

	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
		return "denis.tirskikh"
	} else {
		return pc.UserName
	}
	//log.Println(pc.ID)//log.Println(pc.UserName)
	//return pc.UserName
}
