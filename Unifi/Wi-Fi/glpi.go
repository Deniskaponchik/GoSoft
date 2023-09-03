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
	"time"
)

func UploadMapsToDBerr(datasource string, query string) {

	//datasource := "root:t2root@tcp(10.77.252.153:3306)/it_support_db"

	myError := 1
	for myError != 0 {
		if db, errSqlOpen := sql.Open("mysql", datasource); errSqlOpen == nil {
			errDBping := db.Ping()
			if errDBping == nil {
				defer db.Close() // defer the close till after the main function has finished

				//for myError != 0 { //зачем зацикливать выполнение запроса при корректном подключении к БД?
				_, errQuery := db.Exec(query)
				if errQuery == nil {
					//db.Close()
					if errDBclose := db.Close(); errDBclose != nil {
						fmt.Println("Закрытие подключения к БД завершилось не корректно")
					}
					myError = 0
					/*
						fmt.Println("Вывод мапы ВНУТРИ функции")
						for k, v := range m {
							fmt.Println("innerMap "+k, v.Name, v.Exception, v.SrID)
						}*/
				} else {
					//panic(errQuery.Error()) // proper error handling instead of panic in your app
					fmt.Println(errQuery.Error())
					fmt.Println("Запрос НЕ смог отработать. Проверь корректность всех данных в запросе")
					//fmt.Println("Будет предпринята новая попытка через 1 минут")
					//time.Sleep(60 * time.Second)
					myError = 0 //если такой таблицы нет в БД, то что она появится через 5 минут?
				}
				//} //db.Query
			} else {
				fmt.Println("db.Ping failed:", errDBping)
				fmt.Println("Подключение к БД НЕ установлено. Проверь доступность БД")
				fmt.Println("Будет предпринята новая попытка через 1 минут")
				time.Sleep(60 * time.Second)
				//myError = 1
				myError++
				if myError == 5 { //У меня всё равно будет повторная попытка выгрузки в БД через час. Не критично останавливаться на этом
					myError = 0
					//result = "denis.tirskikh"
				}
			}
		} else {
			//log.Print(errSqlOpen.Error())
			fmt.Println("Error creating DB:", errSqlOpen)
			fmt.Println("To verify, db is:", db)
			fmt.Println("Создание подключения к БД завершилось ошибкой. Часто возникает из-за не корректного драйвера")
			fmt.Println("Будет предпринята новая попытка через 1 минут")
			time.Sleep(60 * time.Second)
			//myError = 1
			myError++
			if myError == 5 {
				myError = 0
				//result = "denis.tirskikh"
			}
		}
	} //sql.Open
}

func UploadMapsToDBstring(dbName string, query string) {

	var datasource string
	datasource = "root:t2root@tcp(10.77.252.153:3306)/it_support_db"
	db, err := sql.Open("mysql", datasource)
	if err != nil {
		panic(err.Error())
	} // if there is an error opening the connection, handle it

	defer db.Close() // defer the close till after the main function has finished

	_, err = db.Exec(query)
	if err != nil {
		panic(err.Error())
	}
}

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

func DownloadMapFromDBanomaliesErr(datasource string, bdController int8, beforeDays string) map[string]DateSiteAnom {
	type TagAnomaly struct {
		Mac        string `json:"mac"`
		Datetime   string `json:"date_hour"`
		Controller int    `json:"controller"`
		Sitename   string `json:"sitename"`
		Anomalies  string `json:"anomalies"`
	}
	m := make(map[string]DateSiteAnom)
	//datasource := "root:t2root@tcp(10.77.252.153:3306)/it_support_db"
	var anomSlice []string
	var dayMac string

	myError := 1
	for myError != 0 {
		if db, errSqlOpen := sql.Open("mysql", datasource); errSqlOpen == nil {
			errDBping := db.Ping()
			if errDBping == nil {
				defer db.Close() // defer the close till after the main function has finished
				//queryAfter := "SELECT * FROM it_support_db.anomalies WHERE controller = " + strconv.Itoa(int(bdController))
				queryAfter := "SELECT * FROM it_support_db.anomalies WHERE date_hour >= '" + beforeDays + "' AND controller = " + strconv.Itoa(int(bdController))
				fmt.Println(queryAfter)
				for myError != 0 { //зацикливание выполнения запроса
					results, errQuery := db.Query(queryAfter)
					if errQuery == nil {
						var tag TagAnomaly
						for results.Next() {
							//errScan := results.Scan(&tag.Mac, &tag.Name, &tag.Controller, &tag.Exception, &tag.SrID)
							errScan := results.Scan(&tag.Datetime, &tag.Mac, &tag.Controller, &tag.Sitename, &tag.Anomalies)
							if errScan == nil {
								anomSlice = strings.Split(tag.Anomalies, ";")
								if len(anomSlice) > 2 {
									dayMac = strings.Split(tag.Datetime, " ")[0] + tag.Mac
									m[dayMac] = DateSiteAnom{
										tag.Mac,
										tag.Datetime,
										tag.Sitename,
										anomSlice,
									}
								}
							} else {
								//panic(errScan.Error()) // proper error handling instead of panic in your app
								fmt.Println(errScan.Error())
								fmt.Println("Сканирование строки и занесение в переменные структуры завершилось ошибкой")
								fmt.Println("Проверь, что не изменилась структура таблицы и кол-во полей")
								myError = 0
								//break
							}
						}
						if errRowsNext := results.Err(); errRowsNext != nil {
							fmt.Println("Цикл прохода по результирующим рядам завершился не корректно")
							//если есть ошибка прохода по строкам, отправляем на перезапрос
							myError = 0
						}
						if myError != 1 {
							//results.Close()
							if errRowsClose := results.Close(); errRowsClose != nil {
								fmt.Println("Закрытие процесса прохода по результирующим полям завершилось не корректно")
							}
							//db.Close()
							if errDBclose := db.Close(); errDBclose != nil {
								fmt.Println("Закрытие подключения к БД завершилось не корректно")
							}
							myError = 0
							/*
								fmt.Println("Вывод мапы ВНУТРИ функции")
								for k, v := range m {
									fmt.Println("innerMap "+k, v.Name, v.Exception, v.SrID)
								}*/
						} else {
							//fmt.Println("Будет предпринята новая попытка запроса через 1 минут")
							//time.Sleep(60 * time.Second)
							myError = 0
						}
					} else {
						//panic(errQuery.Error()) // proper error handling instead of panic in your app
						fmt.Println(errQuery.Error())
						fmt.Println("Запрос НЕ смог отработать. Проверь корректность всех данных в запросе")
						//fmt.Println("Будет предпринята новая попытка через 1 минут")
						//time.Sleep(60 * time.Second)
						myError = 0 //если такой таблицы нет в БД, то что она появится через 5 минут?
					}
				} //db.Query
			} else {
				fmt.Println("db.Ping failed:", errDBping)
				fmt.Println("Подключение к БД НЕ установлено. Проверь доступность БД")
				fmt.Println("Будет предпринята новая попытка через 1 минут")
				time.Sleep(60 * time.Second)
				//myError = 1
				myError++
				if myError == 300 { //Если ночью сервер перезагрузился + нет доступа к БД = в ЦОДЕ коллапс. Могу подождать 5 часов
					myError = 0
				}
			}
		} else {
			//log.Print(errSqlOpen.Error())
			fmt.Println("Error creating DB:", errSqlOpen)
			fmt.Println("To verify, db is:", db)
			fmt.Println("Создание подключения к БД завершилось ошибкой. Часто возникает из-за не корректного драйвера")
			fmt.Println("Будет предпринята новая попытка через 1 минут")
			time.Sleep(60 * time.Second)
			//myError = 1
			myError++
			if myError == 300 { //Если ночью сервер перезагрузился + нет доступа к БД = в ЦОДЕ коллапс. Могу подождать 5 часов
				myError = 0
			}
		}
	} //sql.Open
	return m
}

func DownloadMapFromDBmachinesErr(datasource string, bdController int8) map[string]MachineMyStruct {
	type TagMachine struct {
		Mac        string `json:"mac"`
		Name       string `json:"name"`
		Controller int    `json:"controller"`
		Exception  int    `json:"exception"`
		SrID       string `json:"srid"`
		ApName     string `json:"apname"`
	}
	m := make(map[string]MachineMyStruct)
	//datasource := "root:t2root@tcp(10.77.252.153:3306)/it_support_db"

	myError := 1
	for myError != 0 {
		if db, errSqlOpen := sql.Open("mysql", datasource); errSqlOpen == nil {
			errDBping := db.Ping()
			if errDBping == nil {
				defer db.Close() // defer the close till after the main function has finished
				queryAfter := "SELECT * FROM it_support_db.machine WHERE controller = " + strconv.Itoa(int(bdController))
				//queryAfter := "SELECT * FROM it_support_db.a WHERE controller = " + strconv.Itoa(int(bdController))
				fmt.Println(queryAfter)
				for myError != 0 { //зацикливание выполнения запроса
					results, errQuery := db.Query(queryAfter)
					if errQuery == nil {
						var tag TagMachine
						for results.Next() {
							//errScan := results.Scan(&tag.Mac, &tag.Name, &tag.Controller, &tag.Exception, &tag.SrID)
							errScan := results.Scan(&tag.Mac, &tag.Name, &tag.Controller, &tag.Exception, &tag.SrID, &tag.ApName)
							if errScan == nil {
								//fmt.Println(tag.KeyDB.String, tag.ValueDB.String)
								//fmt.Println(tag.Mac, tag.Name, tag.Controller, tag.Exception, tag.SrID)
								m[tag.Mac] = MachineMyStruct{
									tag.Name,
									tag.Exception,
									tag.SrID,
									tag.ApName,
								}
							} else {
								//panic(errScan.Error()) // proper error handling instead of panic in your app
								fmt.Println(errScan.Error())
								fmt.Println("Сканирование строки и занесение в переменные структуры завершилось ошибкой")
								fmt.Println("Проверь, что не изменилась структура таблицы и кол-во полей")
								myError = 0
								//break
							}
						}
						if errRowsNext := results.Err(); errRowsNext != nil {
							fmt.Println("Цикл прохода по результирующим рядам завершился не корректно")
							//если есть ошибка прохода по строкам, отправляем на перезапрос
							myError = 0
						}
						if myError != 1 {
							//results.Close()
							if errRowsClose := results.Close(); errRowsClose != nil {
								fmt.Println("Закрытие процесса прохода по результирующим полям завершилось не корректно")
							}
							//db.Close()
							if errDBclose := db.Close(); errDBclose != nil {
								fmt.Println("Закрытие подключения к БД завершилось не корректно")
							}
							myError = 0
							/*
								fmt.Println("Вывод мапы ВНУТРИ функции")
								for k, v := range m {
									fmt.Println("innerMap "+k, v.Name, v.Exception, v.SrID)
								}*/
						} else {
							//fmt.Println("Будет предпринята новая попытка запроса через 1 минут")
							//time.Sleep(60 * time.Second)
							myError = 0
						}
					} else {
						//panic(errQuery.Error()) // proper error handling instead of panic in your app
						fmt.Println(errQuery.Error())
						fmt.Println("Запрос НЕ смог отработать. Проверь корректность всех данных в запросе")
						//fmt.Println("Будет предпринята новая попытка через 1 минут")
						//time.Sleep(60 * time.Second)
						myError = 0 //если такой таблицы нет в БД, то что она появится через 5 минут?
					}
				} //db.Query
			} else {
				fmt.Println("db.Ping failed:", errDBping)
				fmt.Println("Подключение к БД НЕ установлено. Проверь доступность БД")
				fmt.Println("Будет предпринята новая попытка через 1 минут")
				time.Sleep(60 * time.Second)
				//myError = 1
				myError++
				if myError == 300 { //Если ночью сервер перезагрузился + нет доступа к БД = в ЦОДЕ коллапс. Могу подождать 5 часов
					myError = 0
				}
			}
		} else {
			//log.Print(errSqlOpen.Error())
			fmt.Println("Error creating DB:", errSqlOpen)
			fmt.Println("To verify, db is:", db)
			fmt.Println("Создание подключения к БД завершилось ошибкой. Часто возникает из-за не корректного драйвера")
			fmt.Println("Будет предпринята новая попытка через 1 минут")
			time.Sleep(60 * time.Second)
			//myError = 1
			myError++
			if myError == 300 { //Если ночью сервер перезагрузился + нет доступа к БД = в ЦОДЕ коллапс. Могу подождать 5 часов
				myError = 0
			}
		}
	} //sql.Open
	return m
}

func DownloadMapFromDBmachines(bdController int8) map[string]MachineMyStruct {
	type TagMachine struct {
		Mac        string `json:"mac"`
		Name       string `json:"name"`
		Controller int    `json:"controller"`
		Exception  int    `json:"exception"`
		SrID       string `json:"srid"`
		ApName     string `json:"apname"`
	}

	//var ap ApMyStruct
	//var machine MachineMyStruct
	m := make(map[string]MachineMyStruct)

	db, err := sql.Open("mysql", "root:t2root@tcp(10.77.252.153:3306)/it_support_db")
	//db, err := sql.Open("mysql", datasource)
	if err != nil {
		log.Print(err.Error())
	}
	defer db.Close() // defer the close till after the main function has finished

	queryAfter := "SELECT * FROM it_support_db.machine WHERE controller = " + strconv.Itoa(int(bdController))
	fmt.Println(queryAfter)

	results, err := db.Query(queryAfter)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	for results.Next() {
		var tag TagMachine
		//err = results.Scan(&tag.ID, &tag.Name)
		err = results.Scan(&tag.Mac, &tag.Name, &tag.Controller, &tag.Exception, &tag.SrID, &tag.ApName)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		//fmt.Println(tag.KeyDB.String, tag.ValueDB.String)
		//fmt.Println(tag.Mac, tag.Name, tag.Controller, tag.Exception, tag.SrID)
		m[tag.Mac] = MachineMyStruct{
			tag.Name,
			tag.Exception,
			tag.SrID,
			tag.ApName,
		}
	}
	results.Close()
	/*
		fmt.Println("Вывод мапы ВНУТРИ функции")
		for k, v := range m {
			fmt.Println("innerMap "+k, v.Hostname, v.Exception, v.SrID, v.ApName)
		}*/
	return m
}

func DownloadMapFromDBapsErr(datasource string, bdController int8) map[string]ApMyStruct {
	type TagAp struct {
		Mac        string `json:"mac"`
		Name       string `json:"name"`
		Controller int    `json:"controller"`
		Exception  int    `json:"exception"`
		SrID       string `json:"srid"`
	}
	m := make(map[string]ApMyStruct)
	//datasource := "root:t2root@tcp(10.77.252.153:3306)/it_support_db"

	myError := 1
	for myError != 0 {
		if db, errSqlOpen := sql.Open("mysql", datasource); errSqlOpen == nil {
			errDBping := db.Ping()
			if errDBping == nil {
				defer db.Close() // defer the close till after the main function has finished
				queryAfter := "SELECT * FROM it_support_db.ap WHERE controller = " + strconv.Itoa(int(bdController))
				//queryAfter := "SELECT * FROM it_support_db.a WHERE controller = " + strconv.Itoa(int(bdController))
				fmt.Println(queryAfter)
				for myError != 0 { //зацикливание выполнения запроса
					results, errQuery := db.Query(queryAfter)
					if errQuery == nil {
						var tag TagAp
						for results.Next() {
							errScan := results.Scan(&tag.Mac, &tag.Name, &tag.Controller, &tag.Exception, &tag.SrID)
							//errScan := results.Scan(&tag.Mac, &tag.Name, &tag.Controller, &tag.Exception)
							if errScan == nil {
								//fmt.Println(tag.KeyDB.String, tag.ValueDB.String)
								//fmt.Println(tag.Mac, tag.Name, tag.Controller, tag.Exception, tag.SrID)
								m[tag.Mac] = ApMyStruct{
									tag.Name,
									tag.Exception,
									tag.SrID,
									0,
								}
							} else {
								//panic(errScan.Error()) // proper error handling instead of panic in your app
								fmt.Println(errScan.Error())
								fmt.Println("Сканирование строки и занесение в переменные структуры завершилось ошибкой")
								fmt.Println("Проверь, что не изменилась структура таблицы и кол-во полей")
								myError = 0
							}
						}
						if errRowsNext := results.Err(); errRowsNext != nil {
							fmt.Println("Цикл прохода по результирующим рядам завершился не корректно")
							//если есть ошибка прохода по строкам, отправляем на перезапрос
							myError = 0
						}
						if myError != 1 {
							//results.Close()
							if errRowsClose := results.Close(); errRowsClose != nil {
								fmt.Println("Закрытие процесса прохода по результирующим полям завершилось не корректно")
							}
							//db.Close()
							if errDBclose := db.Close(); errDBclose != nil {
								fmt.Println("Закрытие подключения к БД завершилось не корректно")
							}
							myError = 0
							/*
								fmt.Println("Вывод мапы ВНУТРИ функции")
								for k, v := range m {
									fmt.Println("innerMap "+k, v.Name, v.Exception, v.SrID)
								}*/
						} else {
							//fmt.Println("Будет предпринята новая попытка запроса через 1 минут")
							//time.Sleep(60 * time.Second)
							myError = 0
						}
					} else {
						//panic(errQuery.Error()) // proper error handling instead of panic in your app
						fmt.Println(errQuery.Error())
						fmt.Println("Запрос НЕ смог отработать. Проверь корректность всех данных в запросе")
						//fmt.Println("Будет предпринята новая попытка через 1 минут")
						//time.Sleep(60 * time.Second)
						myError = 0 //если такой таблицы нет в БД, то что она появится через 5 минут?
					}
				} //db.Query
			} else {
				fmt.Println("db.Ping failed:", errDBping)
				fmt.Println("Подключение к БД НЕ установлено. Проверь доступность БД")
				fmt.Println("Будет предпринята новая попытка через 1 минут")
				time.Sleep(60 * time.Second)
				//myError = 1
				myError++
				if myError == 300 { //Если ночью сервер перезагрузился + нет доступа к БД = в ЦОДЕ коллапс. Могу подождать 5 часов
					myError = 0
				}
			}
		} else {
			//log.Print(errSqlOpen.Error())
			fmt.Println("Error creating DB:", errSqlOpen)
			fmt.Println("To verify, db is:", db)
			fmt.Println("Создание подключения к БД завершилось ошибкой. Часто возникает из-за не корректного драйвера")
			fmt.Println("Будет предпринята новая попытка через 1 минут")
			time.Sleep(60 * time.Second)
			//myError = 1
			myError++
			if myError == 300 { //Если ночью сервер перезагрузился + нет доступа к БД = в ЦОДЕ коллапс. Могу подождать 5 часов
				myError = 0
			}
		}
	} //sql.Open
	return m
}

func DownloadMapFromDBaps(bdController int8) map[string]ApMyStruct {
	type TagAp struct {
		Mac        string `json:"mac"`
		Name       string `json:"name"`
		Controller int    `json:"controller"`
		Exception  int    `json:"exception"`
		SrID       string `json:"srid"`
	}

	//var ap ApMyStruct
	//var machine MachineMyStruct
	m := make(map[string]ApMyStruct)

	db, err := sql.Open("mysql", "root:t2root@tcp(10.77.252.153:3306)/it_support_db")
	//db, err := sql.Open("mysql", datasource)
	if err != nil {
		log.Print(err.Error())
	}

	defer db.Close() // defer the close till after the main function has finished

	queryAfter := "SELECT * FROM it_support_db.ap WHERE controller = " + strconv.Itoa(int(bdController))
	fmt.Println(queryAfter)

	results, err := db.Query(queryAfter)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	for results.Next() {
		var tag TagAp
		//err = results.Scan(&tag.ID, &tag.Name)
		err = results.Scan(&tag.Mac, &tag.Name, &tag.Controller, &tag.Exception, &tag.SrID)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		//fmt.Println(tag.KeyDB.String, tag.ValueDB.String)
		//fmt.Println(tag.Mac, tag.Name, tag.Controller, tag.Exception, tag.SrID)
		m[tag.Mac] = ApMyStruct{
			tag.Name,
			tag.Exception,
			tag.SrID,
			0,
		}
	}
	results.Close()
	/*
		fmt.Println("Вывод мапы ВНУТРИ функции")
		for k, v := range m {
			fmt.Println("innerMap "+k, v.Name, v.Exception, v.SrID)
		}*/
	return m
}

func DownloadMapFromDBerr(datasource string) map[string]string {
	type Tag struct {
		KeyDB   sql.NullString `json:"keyDB""`
		ValueDB sql.NullString `json:"valueDB"`
	}
	m := make(map[string]string)
	//datasource := "root:t2root@tcp(10.77.252.153:3306)/it_support_db"

	myError := 1
	for myError != 0 {
		if db, errSqlOpen := sql.Open("mysql", datasource); errSqlOpen == nil {
			errDBping := db.Ping()
			if errDBping == nil {
				defer db.Close() // defer the close till after the main function has finished
				queryAfter := "SELECT * FROM it_support_db.site_apcut_login"
				//queryAfter := "SELECT * FROM it_support_db.a WHERE controller = " + strconv.Itoa(int(bdController))
				fmt.Println(queryAfter)
				for myError != 0 { //зацикливание выполнения запроса
					results, errQuery := db.Query(queryAfter)
					if errQuery == nil {
						var tag Tag
						for results.Next() {
							//errScan := results.Scan(&tag.Mac, &tag.Name, &tag.Controller, &tag.Exception, &tag.SrID)
							errScan := results.Scan(&tag.KeyDB, &tag.ValueDB)
							if errScan == nil {
								//fmt.Println(tag.KeyDB.String, tag.ValueDB.String)
								//fmt.Println(tag.Mac, tag.Name, tag.Controller, tag.Exception, tag.SrID)
								m[tag.KeyDB.String] = tag.ValueDB.String
							} else {
								//panic(errScan.Error()) // proper error handling instead of panic in your app
								fmt.Println(errScan.Error())
								fmt.Println("Сканирование строки и занесение в переменные структуры завершилось ошибкой")
								fmt.Println("Проверь, что не изменилась структура таблицы и кол-во полей")
								myError = 0 //если изменилась структура полей табл, то они изменятся за 5 минут? думаю, нет
								//break
							}
						}
						if errRowsNext := results.Err(); errRowsNext != nil {
							fmt.Println("Цикл прохода по результирующим рядам завершился не корректно")
							//если есть ошибка прохода по строкам, отправляем на перезапрос. отключено
							myError = 0
						}
						if myError != 1 {
							//results.Close()
							if errRowsClose := results.Close(); errRowsClose != nil {
								fmt.Println("Закрытие процесса прохода по результирующим полям завершилось не корректно")
							}
							//db.Close()
							if errDBclose := db.Close(); errDBclose != nil {
								fmt.Println("Закрытие подключения к БД завершилось не корректно")
							}
							myError = 0
							/*
								fmt.Println("Вывод мапы ВНУТРИ функции")
								for k, v := range m {
									fmt.Println("innerMap "+k, v.Name, v.Exception, v.SrID)
								}*/
						} else {
							//fmt.Println("Будет предпринята новая попытка запроса через 1 минут")
							//time.Sleep(60 * time.Second)
							myError = 0
						}
					} else {
						//panic(errQuery.Error()) // proper error handling instead of panic in your app
						fmt.Println(errQuery.Error())
						fmt.Println("Запрос НЕ смог отработать. Проверь корректность всех данных в запросе")
						//fmt.Println("Будет предпринята новая попытка через 1 минут")
						//time.Sleep(60 * time.Second)
						myError = 0 //если такой таблицы нет в БД, то что она появится через 5 минут?
					}
				} //db.Query
			} else {
				fmt.Println("db.Ping failed:", errDBping)
				fmt.Println("Подключение к БД НЕ установлено. Проверь доступность БД")
				fmt.Println("Будет предпринята новая попытка через 1 минут")
				time.Sleep(60 * time.Second)
				//myError = 1
				myError++
				if myError == 5 {
					myError = 0
					//result = "denis.tirskikh"
				}
			}
		} else {
			//log.Print(errSqlOpen.Error())
			fmt.Println("Error creating DB:", errSqlOpen)
			fmt.Println("To verify, db is:", db)
			fmt.Println("Создание подключения к БД завершилось ошибкой. Часто возникает из-за не корректного драйвера")
			fmt.Println("Будет предпринята новая попытка через 1 минут")
			time.Sleep(60 * time.Second)
			//myError = 1
			myError++
			if myError == 5 {
				myError = 0
				//result = "denis.tirskikh"
			}
		}
	} //sql.Open
	return m
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

func GetLoginPCerr(datasource string, pcName string) string {
	type PC struct {
		UserName string `json:"user_name"`
	}
	var pc PC
	var result string
	//datasource := "root:t2root@tcp(10.77.252.153:3306)/glpi_db"
	myError := 1

	for myError != 0 {
		if db, errSqlOpen := sql.Open("mysql", datasource); errSqlOpen == nil {
			errDBping := db.Ping()
			if errDBping == nil {
				defer db.Close() // defer the close till after the main function has finished
				//queryAfter := "SELECT * FROM it_support_db.a WHERE controller = " + strconv.Itoa(int(bdController))
				queryAfter := "SELECT contact FROM glpi_db.glpi_computers where name = ? ORDER BY date_mod DESC"

				errQuery := db.QueryRow(queryAfter, pcName).Scan(&pc.UserName)
				if errQuery != nil {
					fmt.Println(errQuery.Error())
					//fmt.Println("В БД нет доступного соответствия имени ПК и логина")
					//return "denis.tirskikh"
					result = "denis.tirskikh"
				} else {
					//Если изменилась имя или структура таблицы, то нет смысла зацикливать на 5 минут SELECT
					result = pc.UserName
				}
				myError = 0
				//db.Close()
			} else {
				fmt.Println("db.Ping failed:", errDBping)
				fmt.Println("Подключение к БД НЕ установлено. Проверь доступность БД")
				fmt.Println("Будет предпринята новая попытка через 1 минут")
				time.Sleep(60 * time.Second)
				//myError = 1

				myError++
				if myError == 5 {
					myError = 0
					result = "denis.tirskikh"
				}
			}
		} else {
			//По факту подключения к БД НЕ происходит на этом этапе
			//https://stackoverflow.com/questions/32345124/why-does-sql-open-return-nil-as-error-when-it-should-not
			fmt.Println("Error creating DB:", errSqlOpen)
			fmt.Println("To verify, db is:", db)
			fmt.Println("Создание подключения к БД завершилось ошибкой. Часто возникает из-за не корректного драйвера")
			fmt.Println("Будет предпринята новая попытка через 1 минут")
			time.Sleep(60 * time.Second)
			//myError = 1

			myError++
			if myError == 5 {
				myError = 0
				result = "denis.tirskikh"
			}
		}
	} //sql.Open
	return result
}

func GetLoginPC(pcName string) string {
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
