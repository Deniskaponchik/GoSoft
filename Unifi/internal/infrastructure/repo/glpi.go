package repo

//https://tutorialedge.net/golang/golang-mysql-tutorial/
import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

func UpdateMapsToDBerr(datasource string, queries []string) {
	//datasource := ""

	myError := 1
	for myError != 0 {
		if db, errSqlOpen := sql.Open("mysql", datasource); errSqlOpen == nil {
			errDBping := db.Ping()
			if errDBping == nil {
				defer db.Close() // defer the close till after the main function has finished

				for _, query := range queries {
					fmt.Println(query)
					_, errQuery := db.Exec(query)
					if errQuery == nil {
						myError = 0
					} else {
						fmt.Println(errQuery.Error())
						fmt.Println("Запрос НЕ смог отработать. Проверь корректность всех данных в запросе")
						myError = 0
					}
				}
				if errDBclose := db.Close(); errDBclose != nil {
					fmt.Println("Закрытие подключения к БД завершилось не корректно")
				}

			} else {
				fmt.Println("db.Ping failed:", errDBping)
				fmt.Println("Подключение к БД НЕ установлено. Проверь доступность БД")
				fmt.Println("Будет предпринята новая попытка через 1 минут")
				time.Sleep(60 * time.Second)
				//myError = 1
				myError++
				if myError == 5 { //У меня всё равно будет повторная попытка выгрузки в БД через час. Не критично останавливаться на этом
					myError = 0
				}
			}
		} else {
			fmt.Println("Error creating DB:", errSqlOpen)
			fmt.Println("To verify, db is:", db)
			fmt.Println("Создание подключения к БД завершилось ошибкой. Часто возникает из-за не корректного драйвера")
			fmt.Println("Будет предпринята новая попытка через 1 минут")
			time.Sleep(60 * time.Second)
			//myError = 1
			myError++
			if myError == 5 {
				myError = 0
			}
		}
	} //sql.Open
}

func UploadMapsToDBerr(query string) {

	datasource := ""

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

func DownloadMapFromDBvcsErr(datasource string) map[string]PolyStruct {
	type TagAp struct {
		Mac       string `json:"mac"`
		IP        string `json:"ip"`
		Region    string `json:"region"`
		RoomName  string `json:"room_name"`
		Login     string `json:"login"`
		SrID      string `json:"srid"`
		PolyType  int    `json:"type"`
		Comment   int    `json:"comment"`
		Exception int    `json:"exception"`
	}

	m := make(map[string]PolyStruct)
	//datasource := ""

	myError := 1
	for myError != 0 {
		if db, errSqlOpen := sql.Open("mysql", datasource); errSqlOpen == nil {
			errDBping := db.Ping()
			if errDBping == nil {
				defer db.Close()                                 // defer the close till after the main function has finished
				queryAfter := "SELECT * FROM it_support_db.poly" // WHERE type = " + strconv.Itoa(int(polyType))
				//queryAfter := "SELECT * FROM it_support_db.a WHERE controller = " + strconv.Itoa(int(bdController))
				fmt.Println(queryAfter)
				for myError != 0 { //зацикливание выполнения запроса
					results, errQuery := db.Query(queryAfter)
					if errQuery == nil {
						var tag TagAp
						for results.Next() {
							errScan := results.Scan(&tag.Mac, &tag.IP, &tag.Region, &tag.RoomName, &tag.Login, &tag.SrID, &tag.PolyType, &tag.Comment, &tag.Exception)
							if errScan == nil {
								//fmt.Println(tag.Mac, tag.Name, tag.Controller, tag.Exception, tag.SrID)
								m[tag.Mac] = PolyStruct{
									tag.IP,
									tag.Region,
									tag.RoomName,
									tag.Login,
									tag.SrID,
									tag.PolyType,
									tag.Comment,
									tag.Exception,
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
