package main

//https://tutorialedge.net/golang/golang-mysql-tutorial/
import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strings"
)

func UploadsMapsToDB(downMap map[string]string) {

	// Open up our database connection.
	// I've set up a database on my local machine using phpmyadmin.
	// The database is called testDb
	db, err := sql.Open("mysql", "root:password1@tcp(127.0.0.1:3306)/test")

	// if there is an error opening the connection, handle it
	if err != nil {
		panic(err.Error())
	}

	// defer the close till after the main function has finished
	// executing
	defer db.Close()

	// perform a db.Query insert
	insert, err := db.Query("INSERT INTO test VALUES ( 2, 'TEST' )")

	// if there is an error inserting, handle it
	if err != nil {
		panic(err.Error())
	}
	// be careful deferring Queries if you are using transactions
	defer insert.Close()
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
	defer db.Close()

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
		fmt.Println(tag.KeyDB.String, tag.ValueDB.String)
		m[tag.KeyDB.String] = tag.ValueDB.String //добавляем строку в map
		//count++
	}
	return m
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
