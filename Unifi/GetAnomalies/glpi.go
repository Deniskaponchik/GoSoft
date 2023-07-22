package main

//https://tutorialedge.net/golang/golang-mysql-tutorial/
import (
	"bytes"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strings"
)

func UploadsMapsToDB(uploadMap map[string]string, dbName string, tableName string, delType string) {

	var datasource string
	if dbName == "glpi_db" {
		datasource = "root:t2root@tcp(10.77.252.153:3306)/glpi_db"
	} else {
		datasource = "root:t2root@tcp(10.77.252.153:3306)/wifi_db"
	}
	db, err := sql.Open("mysql", datasource)
	if err != nil {
		panic(err.Error())
	} // if there is an error opening the connection, handle it
	defer db.Close() // defer the close till after the main function has finished

	var delQuery string
	if delType == "DELETE" {
		delQuery = "DELETE FROM " + tableName
	} else {
		delQuery = "TRUNCATE TABLE " + tableName
	}
	fmt.Println(delQuery)
	_, err = db.Exec(delQuery)
	if err != nil {
		panic(err.Error())
	}

	var b bytes.Buffer
	b.WriteString("INSERT INTO " + tableName + " VALUES ")
	lenMap := len(uploadMap)
	count := 0
	for k, v := range uploadMap {
		count++
		if count != lenMap {
			b.WriteString("('" + k + "','" + v + "'),")
		} else {
			b.WriteString("('" + k + "','" + v + "')") //в конце НЕ ставим запятую
		}
	}
	fmt.Println(b.String())
	fmt.Println("")
	_, err = db.Exec(b.String())
	if err != nil {
		panic(err.Error())
	}
}

func DownloadMapFromDB(dbName string, keyDB string, valueDB string, tableName string, orderBY string) map[string]string {
	m := make(map[string]string)

	type Tag struct {
		KeyDB   sql.NullString `json:"keyDB""`
		ValueDB sql.NullString `json:"valueDB"`
	}

	datasource := ""
	if dbName == "glpi_db" {
		datasource = "root:t2root@tcp(10.77.252.153:3306)/glpi_db"
	} else {
		datasource = "root:t2root@tcp(10.77.252.153:3306)/wifi_db"
	}

	//db, err := sql.Open("mysql", "root:t2root@tcp(10.77.252.153:3306)/glpi_db")
	db, err := sql.Open("mysql", datasource)
	if err != nil {
		log.Print(err.Error())
	}
	defer db.Close() // defer the close till after the main function has finished

	queryBefore := "SELECT keyDB, valueDB FROM tableName ORDER BY orderBY DESC"
	replacer := strings.NewReplacer("keyDB", keyDB, "valueDB", valueDB, "tableName", tableName, "orderBY", orderBY)
	queryAfter := replacer.Replace(queryBefore)
	fmt.Println(queryAfter)

	//("SELECT id, name FROM tags")
	//results, err := db.Query("SELECT id, contact FROM glpi_db.glpi_computers ORDER BY date_mod DESC")
	results, err := db.Query(queryAfter)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	//count := 0
	for results.Next() {
		//fmt.Println(count)
		var tag Tag
		//err = results.Scan(&tag.ID, &tag.Name)
		err = results.Scan(&tag.KeyDB, &tag.ValueDB)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		//log.Println(tag.Name)
		//fmt.Println(tag.KeyDB.String, tag.ValueDB.String)
		m[tag.KeyDB.String] = tag.ValueDB.String //добавляем строку в map
		//count++
	}
	results.Close()
	return m
}

func GetUserLogin(siteApCutName string) string {
	type User struct {
		//ID   int    `json:"id"`
		UserLogin string `json:"login"`
	}

	db, err := sql.Open("mysql", "root:t2root@tcp(10.77.252.153:3306)/wifi_db")
	if err != nil {
		log.Print(err.Error())
	}
	defer db.Close() // defer the close till after the main function has finished

	var user User
	err = db.QueryRow("SELECT login FROM wifi_db.site_apcut_login where site_apcut = ?", siteApCutName).Scan(&user.UserLogin)
	// после запятой указываем значение, которое будет подставляться заместо вопроса + ОБЯЗАТЕЛЬНО в Scan использовать &

	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
		return "denis.tirskikh"
	} else {
		return user.UserLogin
	}
	//log.Println(pc.ID)//log.Println(pc.UserName)
	//return pc.UserName
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
	defer db.Close() // defer the close till after the main function has finished

	var pc PC
	err = db.QueryRow("SELECT contact FROM glpi_db.glpi_computers where name = ? ORDER BY date_mod DESC", pcName).Scan(&pc.UserName)
	// после запятой указываем значение, которое будет подставляться заместо вопроса + ОБЯЗАТЕЛЬНО в Scan использовать &

	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
		return "denis.tirskikh"
	} else {
		return pc.UserName
	}
	//log.Println(pc.ID)//log.Println(pc.UserName)
	//return pc.UserName
}
