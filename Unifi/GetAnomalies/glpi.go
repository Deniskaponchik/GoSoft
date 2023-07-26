package main

//https://tutorialedge.net/golang/golang-mysql-tutorial/
import (
	"bytes"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strconv"
	"strings"
)

func UploadMapsToDBreplace(uploadMap map[string]string, dbName string, tableName string, valueDB string, bdController int8) {

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

	/*заместо DELETE делаем UPDATE
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
	}*/

	bdCntrl := strconv.Itoa(int(bdController))
	//Если передаём параметр valueDB, значит хотим обнулить это поле. Актуально для таблиц с номерами заявок
	if valueDB != "" {
		//обнуляем ВСЕ значения ключей
		updateQuery := "UPDATE " + tableName + " SET " + valueDB + " = NULL WHERE controller = " + bdCntrl
		fmt.Println(updateQuery)
		_, err = db.Exec(updateQuery)
		if err != nil {
			panic(err.Error())
		}
	}

	var b bytes.Buffer
	b.WriteString("REPLACE INTO " + tableName + " VALUES ")
	lenMap := len(uploadMap)
	count := 0
	for k, v := range uploadMap {
		count++
		// ('k','v','bdCntrl'),
		if count != lenMap {
			b.WriteString("('" + k + "','" + v + "','" + bdCntrl + "'),")
		} else {
			b.WriteString("('" + k + "','" + v + "','" + bdCntrl + "')") //в конце НЕ ставим запятую
		}
	}
	fmt.Println(b.String())
	if count != 0 {
		_, err = db.Exec(b.String())
		if err != nil {
			panic(err.Error())
		}
	} else {
		fmt.Println("Передана пустая карта. Запрос не выполнен")
	}
	fmt.Println("")
}

func UploadsMapsToDBdelete(uploadMap map[string]string, dbName string, tableName string, delType string) {

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

func DownloadMapFromDB(dbName string, keyDB string, valueDB string, tableName string, bdController int8, orderBY string) map[string]string {
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

	var queryAfter string
	if bdController != 0 {
		//bdController = strconv.Itoa(int(bdController))
		queryBefore := "SELECT keyDB, valueDB FROM tableName WHERE controller = bdController ORDER BY orderBY DESC"
		replacer := strings.NewReplacer("keyDB", keyDB, "valueDB", valueDB, "tableName", tableName, "bdController", strconv.Itoa(int(bdController)), "orderBY", orderBY)
		queryAfter = replacer.Replace(queryBefore)
	} else {
		//без WHERE и bdController
		queryBefore := "SELECT keyDB, valueDB FROM tableName ORDER BY orderBY DESC"
		replacer := strings.NewReplacer("keyDB", keyDB, "valueDB", valueDB, "tableName", tableName, "orderBY", orderBY)
		queryAfter = replacer.Replace(queryBefore)
	}
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

func GetLoginAP(siteApCutName string) string {
	//Единичный запрос для получения логина вк БД уже не использую. Только массовая загрузка/выгрузка в мапы
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

func GetLoginPC(pcName string) string {
	//Единичный запрос для получения логина вк БД уже не использую. Только массовая загрузка/выгрузка в мапы
	type PC struct {
		//ID   int    `json:"id"`
		UserName string `json:"user_name"`
		//Date_Mod string `json:"date_mod"`
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
		//panic(err.Error()) // proper error handling instead of panic in your app
		fmt.Println("В БД нет доступного соответствия имени ПК и логина")
		return "denis.tirskikh"
	} else {
		return pc.UserName
	}
	//log.Println(pc.ID)//log.Println(pc.UserName)
	//return pc.UserName
}
