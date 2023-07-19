package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strings"
)

func main() {

	newMap := DownloadMapFromDB("glpi_db", "name", "contact", "glpi_db.glpi_computers", "date_mod")
	for k, v := range newMap {
		//fmt.Printf("key: %d, value: %t\n", k, v)
		fmt.Println("newMap "+k, v)
	}

}

func DownloadMapFromDB(dbName string, keyDB string, valueDB string, tableName string, orderBY string) map[string]string {
	m := make(map[string]string)
	/*
		type Tag struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}*/
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
	defer db.Close()

	queryBefore := "SELECT keyDB, valueDB FROM nameDB ORDER BY orderBY DESC"
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
		fmt.Println(tag.KeyDB.String, tag.ValueDB.String)
		m[tag.KeyDB.String] = tag.ValueDB.String //добавляем строку в map
		//count++
	}
	return m
}

//https://tutorialedge.net/golang/golang-mysql-tutorial/
/* Many ROWS.
type PC struct {
	//ID   int    `json:"id"`
	UserName string `json:"user_name"`
	Date_Mod string `json:"date_mod"`
}
func main() {
	db, err := sql.Open("mysql", "root:t2root@tcp(10.77.252.153:3306)/glpi_db")
	if err != nil {
		log.Print(err.Error())
	}
	defer db.Close()

	//results, err := db.Query("SELECT id, contact login FROM glpi_db.glpi_computers where name = 'NBAB-FISHER'")
	results, err := db.Query("SELECT date_mod, contact FROM glpi_db.glpi_computers where name = 'NBAB-FISHER' ORDER BY date_mod DESC")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	for results.Next() {
		var pc PC
		// for each row, scan the result into our tag composite object
		err = results.Scan(&pc.Date_Mod, &pc.UserName)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		//log.Printf(tag.Date_Mod)
		log.Printf(pc.UserName)
		break
	}
}*/
/* Many ROWS. Полностью рабочее. !!! НЕ РЕДАКТИРОВАТЬ !!!
type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	//Login string `json:"name"`
}
func main() {
	db, err := sql.Open("mysql", "root:t2root@tcp(10.77.252.153:3306)/glpi_db")
	if err != nil {
		log.Print(err.Error())
	}
	defer db.Close()

	results, err := db.Query("SELECT id, contact login FROM glpi_db.glpi_computers where name = 'NBAB-FISHER'")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	for results.Next() {
		var tag Tag
		// for each row, scan the result into our tag composite object
		err = results.Scan(&tag.ID, &tag.Name)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		// and then print out the tag's Name attribute
		log.Printf(tag.Name)
		break
	}
}*/
/*One row. Рабочее
type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	//Login string `json:"name"`
}
func main() {
	db, err := sql.Open("mysql", "root:t2root@tcp(10.77.252.153:3306)/glpi_db")
	if err != nil {
		log.Print(err.Error())
	}
	defer db.Close()

	pcname := "NBAB-FISHER"
	fmt.Println(pcname)
	var tag Tag
	err = db.QueryRow("SELECT id, contact FROM glpi_db.glpi_computers where name = ? ORDER BY date_mod DESC", pcname).Scan(&tag.ID, &tag.Name)
	// после запятой указываем значение, которое будет подставляться заметсо вопроса
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	log.Println(tag.ID)
	log.Println(tag.Name)
}*/
/*One row. Рабочее. не менять
type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	//Login string `json:"name"`
}
func main() {
	db, err := sql.Open("mysql", "root:t2root@tcp(10.77.252.153:3306)/glpi_db")
	if err != nil {
		log.Print(err.Error())
	}
	defer db.Close()

	var tag Tag
	err = db.QueryRow("SELECT id, contact login FROM glpi_db.glpi_computers where id = ?", 2).Scan(&tag.ID, &tag.Name)
	// после запятой указываем значение, которое будет подставляться заметсо вопроса
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	log.Println(tag.ID)
	log.Println(tag.Name)
}*/
