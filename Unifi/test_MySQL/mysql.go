package main

import (
	"bytes"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strings"
)

func main() {
	/*
		newMap := DownloadMapFromDB("glpi_db", "name", "contact", "glpi_db.glpi_computers", "date_mod")
		for k, v := range newMap {
			//fmt.Printf("key: %d, value: %t\n", k, v)
			fmt.Println("newMap "+k, v)
		}*/
	newMap := map[string]string{}
	newMap["KMS-TECH-07"] = "dsfsdfsdsd"
	newMap["KMS-TECH-08"] = "dsfsdfsdsd"
	newMap["KMS-TECH-09"] = "dsfsdfsdsd"
	UploadMapsToDB(newMap, "wifi_db", "wifi_db.ap_name_srid", "TRUNCATE")
}
func UploadMapsToDB(uploadMap map[string]string, dbName string, tableName string, delType string) {
	datasource := ""
	if dbName == "glpi_db" {
		datasource = "root:t2root@tcp(10.77.252.153:3306)/glpi_db"
	} else {
		datasource = "root:t2root@tcp(10.77.252.153:3306)/wifi_db"
	}
	//db, err := sql.Open("mysql", "root:password1@tcp(127.0.0.1:3306)/test")
	db, err := sql.Open("mysql", datasource)
	if err != nil {
		panic(err.Error())
	} // if there is an error opening the connection, handle it
	// executing
	defer db.Close() // defer the close till after the main function has finished

	//!!! We should always use db.Exec whenever we want to do an insert or update or delete !!!
	//TRUNCATE удаляет полностью таблицу и создаёт точно такую же новую. Это быстрее для большого кол-ва данных
	//DELETE удаляет построчно. Можно сделать rolle back
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

	//!!! We should always use db.Exec whenever we want to do an insert or update or delete !!!
	//РАБОЧИЙ вариант сейчас
	var b bytes.Buffer
	//insertQuery := "INSERT INTO " + tableName + " VALUES "
	b.WriteString("INSERT INTO " + tableName + " VALUES ")

	//var megaQuery strings.Builder
	lenMap := len(uploadMap)
	count := 0
	for k, v := range uploadMap {
		count++
		if count != lenMap {
			//megaQuery.WriteString("('" + k + "`,`" + v + "`),")
			//insertQuery += "('" + k + "','" + v + "'),"
			b.WriteString("('" + k + "','" + v + "'),")
		} else {
			b.WriteString("('" + k + "','" + v + "')") //в конце НЕ ставим запятую
		}
	}
	//trim the last ,
	//megaQuery = megaQuery[0 : len(megaQuery)-1]  //String.Builder нельзя обрезать с конца
	//insertQuery = insertQuery[0 : len(insertQuery)-1]
	//b = b[0 : len(b)-1]
	//fmt.Println(megaQuery.String())
	fmt.Println(b.String())
	//Вроде как, Exec нужно выполнять, а не Query
	/*exec, err := db.Exec(megaQuery.String())
	if err != nil {
		return
	}*/
	//db.Exec("INSERT INTO wifi_db.ap_name_srid (apname, srid) VALUES ('ws01', 'wso1'),('ws02','ws02')")
	//db.Exec("INSERT INTO wifi_db.ap_name_srid VALUES ('ws03', 'ws03'),('ws04','ws04')")
	_, err = db.Exec(b.String())
	if err != nil {
		panic(err.Error())
	}

	/*https://stackoverflow.com/questions/21108084/how-to-insert-multiple-data-at-once
	//sqlStr := "INSERT INTO test(n1, n2, n3) VALUES "
	sqlStr := "INSERT INTO wifi_db.ap_name_srid (apname, srid) VALUES "
	vals := []interface{}{}

	for _, row := range data {
		sqlStr += "(?, ?, ?),"
		vals = append(vals, row["v1"], row["v2"], row["v3"])
	}
	for k, v := range uploadMap {
		sqlStr += "('" + k + "`,`" + v + "`),"
		//vals = append(vals, row["v1"], row["v2"], row["v3"])
	}
	//trim the last ,
	sqlStr = sqlStr[0 : len(sqlStr)-1]
	//prepare the statement
	stmt, _ := db.Prepare(sqlStr)
	//format all vals at once
	res, _ := stmt.Exec(vals...)
	//stmt.Exec(sqlStr)
	*/
	/*
		//https://ruby-china.org/topics/38906
		valueStrings := []string{}
		valueArgs := []interface{}{}
		for _, w := range wallets {
			valueStrings = append(valueStrings, "(?, ?, ?, ?)")
			valueArgs = append(valueArgs, w.Address)
			valueArgs = append(valueArgs, asset.Symbol)
			valueArgs = append(valueArgs, asset.Identify)
			valueArgs = append(valueArgs, asset.Decimal)
		}
		for k, v := range uploadMap {
			valueStrings = append(valueStrings, "(?, ?)")
			valueArgs = append(valueArgs, k)
			valueArgs = append(valueArgs, v)
		}
		//smt := `INSERT INTO balances(address, symbol, identify, decimal) VALUES %s ON CONFLICT (address, symbol) DO UPDATE SET address = excluded.address`
		smt := `INSERT INTO wifi_db.ap_name_srid (apname, srid) VALUES %s`
		smt = fmt.Sprintf(smt, strings.Join(valueStrings, ","))
		fmt.Println("smttt:", smt)
		//tx := repo.db.Begin()
		tx := db.Begin()
		err := tx.Exec(smt, valueArgs...).Error
		if err != nil {
			tx.Rollback()
			return err
		}
		return tx.Commit().Error
	*/

	/*
		//https://stackoverflow.com/questions/51130658/how-to-insert-multiple-rows-into-postgres-sql-in-a-go
		samples := // the slice of samples you want to insert
		query := `insert into samples (<the list of columns>) values `
		values := []interface{}{}
		for i, s := range samples {
			values = append(values, s.<field1>, s.<field2>, < ... >)

			numFields := 15 // the number of fields you are inserting
			n := i * numFields

			query += `(`
			for j := 0; j < numFields; j++ {
				query += `$`+strconv.Itoa(n+j+1) + `,`
			}
			query = query[:len(query)-1] + `),`
		}
		query = query[:len(query)-1] // remove the trailing comma
		db.Exec(query, values...)
	*/

	/*Prepared Statements. НЕ сработало
	//http://go-database-sql.org/prepared.html
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()
	//stmt, err := tx.Prepare("INSERT INTO foo VALUES (?)")
	stmt, err := tx.Prepare("INSERT INTO ? VALUES (?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	//defer stmt.Close() // danger!

	for k, v := range uploadMap {
		//_, err = stmt.Exec(tableName, k, v)
		_, err = stmt.Exec("INSERT INTO ? VALUES (?, ?)", dbName, k, v)
		if err != nil {
				log.Fatal(err)
			}
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
	stmt.Close() //runs here!
	//We need to close the statement explicitly when we don’t need the prepared statement anymore.
	//Else, we’ll fail to free up allocated resources both on the client as well as server!
	*/

	/*insert, err := db.Exec(megaQuery.String())
	if err != nil {
		panic(err.Error())
	} // if there is an error inserting, handle it
	//defer insert.Close() // be careful deferring Queries if you are using transactions
	*/

	/*РАБОТАЕТ, но медленно
	for k, v := range uploadMap {
		queryBefore = "INSERT INTO tableName VALUES ('key', 'value')"
		replacer = strings.NewReplacer("key", k, "value", v, "tableName", tableName)
		queryAfter = replacer.Replace(queryBefore)
		fmt.Println(queryAfter)
		insert, err := db.Query(queryAfter)
		if err != nil {
			panic(err.Error())
		} // if there is an error inserting, handle it
		defer insert.Close() // be careful deferring Queries if you are using transactions
	}*/
	/*Блок для одиночного запуска
	//insert, err := db.Query("INSERT INTO test VALUES ( 'TEST', 'TEST' )")
	insert, err := db.Query("INSERT INTO wifi_db.ap_name_srid VALUES ( 'TEST', 'TEST' );")
	//insert, err := db.Query(megaQuery.String())
	if err != nil {
		panic(err.Error())
	} // if there is an error inserting, handle it
	defer insert.Close() // be careful deferring Queries if you are using transactions
	*/
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

	/*https://www.google.com/url?sa=t&rct=j&q=&esrc=s&source=web&cd=&cad=rja&uact=8&ved=2ahUKEwitgNfW4JqAAxXrGxAIHWKVAcAQFnoECAMQAQ&url=https%3A%2F%2Faloksinhanov.medium.com%2Fquery-vs-exec-vs-prepare-in-golang-e7c49212c36c&usg=AOvVaw3Ah11daev2ZHoF47XuAmQD&opi=89978449
	We should always use db.Query whenever we want to do a select and we should never ignore the returned rows of Query but iterate over it (else we’ll	leak the db connection!)
	*/

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
	results.Close()
	return m
}

func GetUserLogin(dbName string, keySelect string, tableName string, whereKey string, orderBY string) string {
	type Target struct {
		UserLogin string `json:"user_name"`
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

	var target Target
	queryBefore := "SELECT keySelect FROM dbName where whereKey = ? ORDER BY orderBY DESC"
	replacer := strings.NewReplacer("keySelect", keySelect, "tableName", tableName, "orderBY", orderBY)
	queryAfter := replacer.Replace(queryBefore)
	fmt.Println(queryAfter)

	//err = db.QueryRow("SELECT contact FROM glpi_db.glpi_computers where name = ? ORDER BY date_mod DESC", keySelect).Scan(&target.UserLogin)
	err = db.QueryRow(queryAfter, keySelect).Scan(&target.UserLogin)
	// после запятой указываем значение, которое будет подставляться заместо вопроса + ОБЯЗАТЕЛЬНО в Scan использовать &

	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
		return "denis.tirskikh"
	} else {
		return target.UserLogin
	}
	//log.Println(pc.ID)//log.Println(pc.UserName)
	//return pc.UserName
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
