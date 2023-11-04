package repo

//https://tutorialedge.net/golang/golang-mysql-tutorial/
//https://github.com/evrone/go-clean-template/blob/master/internal/usecase/repo/translation_postgres.go
import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/deniskaponchik/GoSoft/Unifi/internal/entity"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"strings"
	"time"
)

type UnifiRepo struct {
	dataSourceITsup string
	databaseITsup   string
	dataSourceGLPI  string
	databaseGLPI    string
	controller      int
}

// реализуем Инъекцию зависимостей DI. Используется в app
func NewUnifiRepo(i string, g string, c int) (*UnifiRepo, error) {
	pr := &UnifiRepo{
		dataSourceITsup: i,
		databaseITsup:   strings.Split(i, "/")[1],
		dataSourceGLPI:  g,
		databaseGLPI:    strings.Split(g, "/")[1],
		controller:      c,
	}

	if db, errSqlOpen := sql.Open("mysql", pr.dataSourceITsup); errSqlOpen == nil {
		errDBping := db.Ping()
		if errDBping == nil {
			return pr, nil
		} else {
			fmt.Println("db.Ping failed:", errDBping)
			fmt.Println("Подключение к БД НЕ установлено. Проверь доступность БД")
			return nil, errDBping
		}
	} else {
		fmt.Println("Error creating DB:", errSqlOpen)
		fmt.Println("To verify, db is:", db)
		fmt.Println("Создание подключения к БД завершилось ошибкой. Часто возникает из-за не корректного драйвера")
		return nil, errSqlOpen
	}
	//return pr, nil
}

// под логику, где у Клиентов есть массив Аномалий
func (ur *UnifiRepo) UpdateDbAnomaly(mac_Anomaly map[string]*entity.Anomaly) (err error) {
	//приходит мапа типа: мак адрес клиента _ аномалии клиента за 1 час

	bdCntrl := strconv.Itoa(int(ur.controller)) //bdController))
	var anomSliceString string
	var query string
	var b1 bytes.Buffer
	b1.WriteString("INSERT INTO  " + ur.databaseITsup + ".anomalies VALUES ")
	countB1 := 0

	for _, v := range mac_Anomaly { //mac_DateSiteAnom {

		if len(v.SliceAnomStr) > 1 {
			//если аномалий за час накопилось 2 и более, то такие заносим в БД
			countB1++
			anomSliceString = strings.Join(v.SliceAnomStr, ";")
			b1.WriteString("('" + v.DateHour + "','" + v.ClientMac + "','" + bdCntrl + "','" + v.SiteName + "','" + anomSliceString +
				"','" + v.ApMac + "','" + v.ApName + "','" + strconv.Itoa(int(v.Exception)) + "'),")
		}

	}
	query = b1.String()
	//Возможно, не самый эффективный метод обрезать строку с конца, но рабочий
	if last := len(query) - 1; last >= 0 && query[last] == ',' {
		query = query[:last]
	}
	fmt.Println(query)
	if countB1 != 0 {
		//UploadMapsToDBerr(wifiConf.GlpiConnectStringITsupport, query)
		err = ur.UploadMapsToDBerr(query)
	} else {
		fmt.Println("Передана пустая карта. Запрос не выполнен")
	}
	fmt.Println("")

	return nil
}

//Логика, когда в аномалии была мапа со временем и аномалиями
/*func (ur *UnifiRepo) UpdateDbAnomaly(mac_Anomaly map[string]*entity.Anomaly) (err error) {
	//приходит мапа типа: dayMac_Anomaly

	bdCntrl := strconv.Itoa(int(ur.controller)) //bdController))
	var anomSliceString string
	var query string
	var b1 bytes.Buffer
	b1.WriteString("INSERT INTO  " + ur.databaseITsup + ".anomalies VALUES ")
	//lenMap := len(mapAnomaly) //mac_DateSiteAnom)
	var lenSlice int
	countB1 := 0

	for _, v := range mac_Anomaly { //mac_DateSiteAnom {
		for dateHour, va := range v.TimeStr_sliceAnomStr {
			//dateHour = 2023-09-01 12:00:00
			//пробегаемся по каждой записи мапы. но она всего одна, так сгруппировано всё за один час
			lenSlice = len(va) //v.AnomalySlice) //AnomSlice)
			if lenSlice > 1 {
				//если аномалий за час накопилось 2 и более, то такие заносим в БД
				countB1++
				//anomSliceString = strings.Join(v.AnomalySlice, ";")
				anomSliceString = strings.Join(va, ";")
				//b1.WriteString("('" + v.DateTime + "','" + k + "','" + bdCntrl + "','" + siteNameCut + "','" + anomSliceString + "'),")
				b1.WriteString("('" + dateHour + "','" + v.ClientMac + "','" + bdCntrl + "','" + v.SiteName + "','" + anomSliceString +
					"','" + v.ApMac + "','" + v.ApName + "','" + strconv.Itoa(int(v.Exception)) + "'),")
			}
		}
	}
	query = b1.String()
	//Возможно, не самый эффективный метод обрезать строку с конца, но рабочий
	if last := len(query) - 1; last >= 0 && query[last] == ',' {
		query = query[:last]
	}
	fmt.Println(query)
	if countB1 != 0 {
		//UploadMapsToDBerr(wifiConf.GlpiConnectStringITsupport, query)
		err = ur.UploadMapsToDBerr(query)
	} else {
		fmt.Println("Передана пустая карта. Запрос не выполнен")
	}
	fmt.Println("")

	return nil
}*/

func (ur *UnifiRepo) UpdateDbAp(mapAp map[string]*entity.Ap) (err error) {

	bdCntrl := strconv.Itoa(int(ur.controller)) //bdController))
	var lenMap int
	var count int
	var exception string
	var b1 bytes.Buffer
	var query string

	//b1.WriteString("REPLACE INTO " + "it_support_db.ap" + " VALUES ")
	b1.WriteString("REPLACE INTO " + ur.databaseITsup + ".ap" + " VALUES ")
	//lenMap := len(uploadMap)
	lenMap = len(mapAp)
	count = 0
	//for k, v := range uploadMap {
	for k, v := range mapAp {
		exception = strconv.Itoa(int(v.Exception))
		count++
		if count != lenMap {
			// mac, name, controller, exception, srid
			b1.WriteString("('" + k + "','" + v.Name + "','" + bdCntrl + "','" + exception + "','" + v.SrID + "'),")
		} else {
			b1.WriteString("('" + k + "','" + v.Name + "','" + bdCntrl + "','" + exception + "','" + v.SrID + "')")
			//в конце НЕ ставим запятую
		}
	}
	query = b1.String()
	fmt.Println(query)
	if count != 0 {
		//UploadMapsToDBstring("it_support_db", query)
		//UploadMapsToDBerr(wifiConf.GlpiConnectStringITsupport, query)
		ur.UploadMapsToDBerr(query)
	} else {
		fmt.Println("Передана пустая карта. Запрос не выполнен")
	}
	fmt.Println("")

	return nil
}

func (ur *UnifiRepo) UploadMapsToDBerr(query string) (err error) {

	myError := 1
	for myError != 0 {
		if db, errSqlOpen := sql.Open("mysql", ur.dataSourceITsup); errSqlOpen == nil {
			errDBping := db.Ping()
			if errDBping == nil {
				defer db.Close() // defer the close till after the main function has finished
				//for myError != 0 { //зачем зацикливать выполнение запроса при корректном подключении к БД?
				_, errQuery := db.Exec(query)
				if errQuery == nil {
					myError = 0
					return nil
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
					err = errQuery
				}
			} else {
				fmt.Println("db.Ping failed:", errDBping)
				fmt.Println("Подключение к БД НЕ установлено. Проверь доступность БД")
				fmt.Println("Будет предпринята новая попытка через 1 минут")
				time.Sleep(60 * time.Second)
				myError++
				err = errDBping
			}
		} else {
			//log.Print(errSqlOpen.Error())
			fmt.Println("Error creating DB:", errSqlOpen)
			fmt.Println("To verify, db is:", db)
			fmt.Println("Создание подключения к БД завершилось ошибкой. Часто возникает из-за не корректного драйвера")
			fmt.Println("Будет предпринята новая попытка через 1 минут")
			time.Sleep(60 * time.Second)
			myError++
			err = errSqlOpen
		}
		if myError == 5 {
			myError = 0
			return err
		}
	} //sql.Open
	return nil
}

//DownloadApWithAnomalies

func (ur *UnifiRepo) DownloadClientsWithAnomalies(beforeDays string) (mac_Client map[string]*entity.Client, err error) {

	//beforeDays = "2023-09-01 12:00:00"
	//var anomSlice []string
	var date string                             //2023-09-01
	var date_Anomaly map[string]*entity.Anomaly //date == 2023-09-01
	//var dateStr_sliceAnomStr map[string][]string

	myError := 1
	for myError != 0 {
		if db, errSqlOpen := sql.Open("mysql", ur.dataSourceITsup); errSqlOpen == nil {
			errDBping := db.Ping()
			if errDBping == nil {
				defer db.Close() // defer the close till after the main function has finished
				//queryAfter := "SELECT * FROM it_support_db.anomalies WHERE controller = " + strconv.Itoa(int(bdController))
				queryAfter := "SELECT * FROM " + ur.databaseITsup + ".anomaly WHERE date_hour >= `" + beforeDays + "` AND controller = " +
					strconv.Itoa(int(ur.controller)) + " AND exception = 0"
				fmt.Println(queryAfter)

				for myError != 0 { //зацикливание выполнения запроса
					results, errQuery := db.Query(queryAfter)
					if errQuery == nil {

						var tag *entity.Anomaly
						//var tag Tag

						for results.Next() {
							errScan := results.Scan(&tag.DateHour, &tag.ClientMac, &tag.Controller, &tag.SiteName, &tag.AnomStr,
								&tag.ApMac, &tag.ApName, &tag.Exception)

							if errScan == nil {
								//anomSlice = strings.Split(tag.Anomalies, ";")
								tag.SliceAnomStr = strings.Split(tag.AnomStr, ";")

								//if len(anomSlice) > 2 { //в БД уже записи с двумя и более аномалиями.
								//комментирую на будущее, если захочу пропускать с тремя и более аномалиями

								date = strings.Split(tag.DateHour, " ")[0]
								date_Anomaly[date] = tag //date == 2023-09-01

								client, exisMacClient := mac_Client[tag.ClientMac]
								if !exisMacClient {
									mac_Client[tag.ClientMac] = &entity.Client{
										Mac:          tag.ClientMac,
										Date_Anomaly: date_Anomaly,
									}
								} else {
									client.Date_Anomaly[date] = tag
								}
								//} //if len(anomSlice) > 2 {
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
		}
		if myError == 300 { //Если ночью сервер перезагрузился + нет доступа к БД = в ЦОДЕ коллапс. Могу подождать 5 часов
			myError = 0
		}
	} //sql.Open
	return mac_Client, nil
}

/*
func (ur *UnifiRepo) DownloadMapFromDBanomalies(beforeDays string) (mac_Anomaly map[string]*entity.Anomaly, err error) {

	type Tag struct {
		DateHour   string `json:"date_hour"`
		MacClient  string `json:"mac_client"`
		Controller int    `json:"controller"`
		SiteName   string `json:"sitename"`
		Anomalies  string `json:"anomalies"`
		ApMac      string `json:"ap_mac"`
		ApName     string `json:"ap_name"`
		Exception  int    `json:"exception"`
	}

	//m := make(map[string]DateSiteAnom)
	//dayMac_anomaly := make(map[string]*entity.Anomaly) //dayMac = 2023-09-01_a0:b1:c2:d3:e4:f5
	//mac_Anomaly = make(map[string]*entity.Anomaly)

	//beforeDays = "2023-09-01 12:00:00"
	var anomSlice []string
	var date string //2023-09-01
	//var day string
	var dateStr_sliceAnomStr map[string][]string

	myError := 1
	for myError != 0 {
		if db, errSqlOpen := sql.Open("mysql", ur.dataSourceITsup); errSqlOpen == nil {
			errDBping := db.Ping()
			if errDBping == nil {
				defer db.Close() // defer the close till after the main function has finished
				//queryAfter := "SELECT * FROM it_support_db.anomalies WHERE controller = " + strconv.Itoa(int(bdController))
				queryAfter := "SELECT * FROM " + ur.databaseITsup + ".anomaly WHERE date_hour >= `" + beforeDays + "` AND controller = " +
					strconv.Itoa(int(ur.controller)) + " AND exception = 0"
				fmt.Println(queryAfter)

				for myError != 0 { //зацикливание выполнения запроса
					results, errQuery := db.Query(queryAfter)
					if errQuery == nil {

						//var tag TagAnomaly
						var tag Tag

						for results.Next() {
							//errScan := results.Scan(&tag.DateHour, &tag.ClientMac, &tag.Controller, &tag.SiteName, &tag.AnomalySlice)
							//errScan := results.Scan(&tag.DateHour, &tag.ClientMac, &tag.Controller, &tag.SiteName, anomStr)
							errScan := results.Scan(&tag.DateHour, &tag.MacClient, &tag.Controller, &tag.SiteName, &tag.Anomalies,
								&tag.ApMac, &tag.ApName, &tag.Exception)

							if errScan == nil {
								if tag.Exception == 0 {
									//anomSlice = strings.Split(tag.AnomalySlice, ";")
									anomSlice = strings.Split(tag.Anomalies, ";")
									if len(anomSlice) > 2 {

										date = strings.Split(tag.DateHour, " ")[0]
										anomaly, exis1 := mac_Anomaly[tag.MacClient]
										if !exis1 {
											dateStr_sliceAnomStr[date] = anomSlice

											mac_Anomaly[tag.ApMac] = &entity.Anomaly{
												ClientMac:            tag.MacClient,
												SiteName:             tag.SiteName,
												Controller:           tag.Controller,
												Exception:            tag.Exception, //параметр всегда ноль. Проверка на не ноль заложена при загрузке аномалий раз в час
												ApMac:                tag.ApMac,
												ApName:               tag.ApName,
												TimeStr_sliceAnomStr: dateStr_sliceAnomStr,
											}
										} else {
											_, exis2 := anomaly.TimeStr_sliceAnomStr[date]
											if !exis2 {
												//добавляем в мапу новую запись дня
												anomaly.TimeStr_sliceAnomStr[date] = anomSlice
											} else {
												//добавляем в мапе бакет дня новые аномалии
												for _, v := range anomSlice {
													anomaly.TimeStr_sliceAnomStr[date] = append(anomaly.TimeStr_sliceAnomStr[date], v)
												}

											}
										}
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
		}
		if myError == 300 { //Если ночью сервер перезагрузился + нет доступа к БД = в ЦОДЕ коллапс. Могу подождать 5 часов
			myError = 0
		}
	} //sql.Open
	return mac_Anomaly, nil
}
*/

/*Старая логика
func (ur *UnifiRepo) DownloadMapFromDBanomaliesErr(beforeDays string) (map[string]*entity.Anomaly, error) {

	//m := make(map[string]DateSiteAnom)
	dayMac_anomaly := make(map[string]*entity.Anomaly) //dayMac = 2023-09-01_a0:b1:c2:d3:e4:f5
	//beforeDays = ""
	var anomSlice []string
	var anomStr string
	var dayMac string
	//var dayHourStr string

	myError := 1
	for myError != 0 {
		if db, errSqlOpen := sql.Open("mysql", ur.dataSourceITsup); errSqlOpen == nil {
			errDBping := db.Ping()
			if errDBping == nil {
				defer db.Close() // defer the close till after the main function has finished
				//queryAfter := "SELECT * FROM it_support_db.anomalies WHERE controller = " + strconv.Itoa(int(bdController))
				queryAfter := "SELECT * FROM " + ur.databaseITsup + ".anomalies WHERE date_hour >= '" + beforeDays + "' AND controller = " + strconv.Itoa(int(ur.controller))
				fmt.Println(queryAfter)
				for myError != 0 { //зацикливание выполнения запроса
					results, errQuery := db.Query(queryAfter)
					if errQuery == nil {
						//var tag TagAnomaly
						var tag entity.Anomaly
						for results.Next() {
							//errScan := results.Scan(&tag.DateHour, &tag.ClientMac, &tag.Controller, &tag.SiteName, &tag.AnomalySlice)
							errScan := results.Scan(&tag.DateHour, &tag.ClientMac, &tag.Controller, &tag.SiteName, anomStr)
							if errScan == nil {
								//anomSlice = strings.Split(tag.AnomalySlice, ";")
								anomSlice = strings.Split(anomStr, ";")
								if len(anomSlice) > 2 {
									//Если за час более двух аномалий, то заносим
									dayMac = strings.Split(tag.DateHour, " ")[0] + tag.ClientMac
									//dayMac = strings.Split(dayHourStr, " ")[0] + tag.ClientMac
									//tag.DateHour, _ = time.Parse("2006-01-02 15:04:05", dayHourStr)
									tag.AnomalySlice = anomSlice
									dayMac_anomaly[dayMac] = &tag
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
		}
		if myError == 300 { //Если ночью сервер перезагрузился + нет доступа к БД = в ЦОДЕ коллапс. Могу подождать 5 часов
			myError = 0
		}
	} //sql.Open
	return dayMac_anomaly, nil
}
*/

func (ur *UnifiRepo) DownloadMapFromDBmachinesErr() (map[string]*entity.Client, error) {

	machineMap := make(map[string]*entity.Client) //https://yourbasic.org/golang/gotcha-assignment-entry-nil-map/
	var err error
	myError := 1
	for myError != 0 {
		if db, errSqlOpen := sql.Open("mysql", ur.dataSourceITsup); errSqlOpen == nil {
			errDBping := db.Ping()
			if errDBping == nil {
				defer db.Close() // defer the close till after the main function has finished
				//queryAfter := "SELECT * FROM it_support_db.machine WHERE controller = " + strconv.Itoa(int(bdController))
				queryAfter := "SELECT * FROM " + ur.databaseITsup + ".machine WHERE controller = " + strconv.Itoa(int(ur.controller))
				fmt.Println(queryAfter)

				for myError != 0 { //зацикливание выполнения запроса
					results, errQuery := db.Query(queryAfter)
					if errQuery == nil {
						//var tag TagPoly
						var tag *entity.Client

						for results.Next() {
							errScan := results.Scan(&tag.Mac, &tag.Hostname, &tag.Controller, &tag.Exception, &tag.SrID, &tag.ApName, &tag.ApMac, &tag.Modified)
							if errScan == nil {
								//fmt.Println(tag.Mac, tag.Name, tag.Controller, tag.Exception, tag.SrID)
								//m[tag.Mac] = MachineMyStruct{
								machineMap[tag.Mac] = tag
								/*
									machineMap[tag.Mac] = entity.Client{
										tag.Mac,
										tag.Hostname,
										tag.SrID,
										tag.Controller,
										tag.Exception,
										tag.ApName,
										nil,
										"",
									}*/
							} else {
								fmt.Println(errScan.Error())
								fmt.Println("Сканирование СТРОКИ и занесение в переменные структуры завершилось ошибкой")
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
						fmt.Println(errQuery.Error())
						fmt.Println("Запрос НЕ смог отработать. Проверь корректность всех данных в запросе")
						myError = 0 //если такой таблицы нет в БД, то что она появится через 5 минут?
						err = errQuery
					}
				} //db.Query
			} else {
				fmt.Println("db.Ping failed:", errDBping)
				fmt.Println("Подключение к БД НЕ установлено. Проверь доступность БД")
				fmt.Println("Будет предпринята новая попытка через 1 минут")
				time.Sleep(60 * time.Second)
				//myError = 1
				myError++
				err = errDBping
				//Если ночью сервер перезагрузился + нет доступа к БД = в ЦОДЕ коллапс. Могу подождать 5 часов
				//if myError == 300 { 	myError = 0				}
			}
		} else {
			fmt.Println("Error creating DB:", errSqlOpen)
			fmt.Println("To verify, db is:", db)
			fmt.Println("Создание подключения к БД завершилось ошибкой. Часто возникает из-за не корректного драйвера")
			fmt.Println("Будет предпринята новая попытка через 1 минут")
			time.Sleep(60 * time.Second)
			//myError = 1
			myError++
			err = errSqlOpen
		}
		if myError == 5 { //&& marker == 1 {
			//Если ночью нет доступа к БД = в ЦОДЕ коллапс. Могу подождать 5 часов при условии, что это ежечасовая актуализация ip-адресов
			myError = 0
			return nil, err //errors.New("подключение к бд не удалось")
		}
	} //sql.Open
	return machineMap, err
}

func (ur *UnifiRepo) DownloadMapFromDBapsErr() (map[string]*entity.Ap, error) {

	apMap := make(map[string]*entity.Ap) //https://yourbasic.org/golang/gotcha-assignment-entry-nil-map/
	var err error
	myError := 1
	for myError != 0 {
		if db, errSqlOpen := sql.Open("mysql", ur.dataSourceITsup); errSqlOpen == nil {
			errDBping := db.Ping()
			if errDBping == nil {
				defer db.Close() // defer the close till after the main function has finished
				//queryAfter := "SELECT * FROM " + ur.database + ".poly"
				//queryAfter := "SELECT * FROM it_support_db.ap WHERE controller = " + strconv.Itoa(int(bdController))
				queryAfter := "SELECT * FROM " + ur.databaseITsup + ".ap WHERE controller = " + strconv.Itoa(int(ur.controller))
				fmt.Println(queryAfter)

				for myError != 0 { //зацикливание выполнения запроса
					results, errQuery := db.Query(queryAfter)
					if errQuery == nil {
						//var tag TagPoly
						var tag entity.Ap

						for results.Next() {
							errScan := results.Scan(&tag.Mac, &tag.Name, &tag.Controller, &tag.SrID, &tag.Exception)
							if errScan == nil {
								//fmt.Println(tag.Mac, tag.Name, tag.Controller, tag.Exception, tag.SrID)
								//m[tag.Mac] = ApMyStruct{
								apMap[tag.Mac] = &tag
								/*
									apMap[tag.Mac] = entity.Ap{
										tag.Mac,
										"",
										tag.Name,
										"",
										tag.SrID,
										tag.Exception,
										tag.Controller,
										0,
									}*/
							} else {
								fmt.Println(errScan.Error())
								fmt.Println("Сканирование СТРОКИ и занесение в переменные структуры завершилось ошибкой")
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
						err = errQuery
					}
				} //db.Query
			} else {
				fmt.Println("db.Ping failed:", errDBping)
				fmt.Println("Подключение к БД НЕ установлено. Проверь доступность БД")
				fmt.Println("Будет предпринята новая попытка через 1 минут")
				time.Sleep(60 * time.Second)
				//myError = 1
				myError++
				err = errDBping
				//Если ночью сервер перезагрузился + нет доступа к БД = в ЦОДЕ коллапс. Могу подождать 5 часов
				//if myError == 300 { 	myError = 0				}
			}
		} else {
			fmt.Println("Error creating DB:", errSqlOpen)
			fmt.Println("To verify, db is:", db)
			fmt.Println("Создание подключения к БД завершилось ошибкой. Часто возникает из-за не корректного драйвера")
			fmt.Println("Будет предпринята новая попытка через 1 минут")
			time.Sleep(60 * time.Second)
			//myError = 1
			myError++
			err = errSqlOpen
		}
		if myError == 5 { //&& marker == 1 {
			//Если ночью нет доступа к БД = в ЦОДЕ коллапс. Могу подождать 5 часов при условии, что это ежечасовая актуализация ip-адресов
			myError = 0
			return nil, err //errors.New("подключение к бд не удалось")
		}
	} //sql.Open
	return apMap, err
}

func (ur *UnifiRepo) DownloadMapFromDBerr() (siteApcut_login map[string]string, err error) {
	type Tag struct {
		KeyDB   sql.NullString `json:"keyDB""`
		ValueDB sql.NullString `json:"valueDB"`
	}
	//siteApcut_login := make(map[string]string)

	myError := 1
	for myError != 0 {
		if db, errSqlOpen := sql.Open("mysql", ur.dataSourceITsup); errSqlOpen == nil {
			errDBping := db.Ping()
			if errDBping == nil {
				defer db.Close() // defer the close till after the main function has finished
				queryAfter := "SELECT * FROM " + ur.databaseITsup + ".site_apcut_login"
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
								siteApcut_login[tag.KeyDB.String] = tag.ValueDB.String
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
				myError++
				err = errDBping
			}
		} else {
			//log.Print(errSqlOpen.Error())
			fmt.Println("Error creating DB:", errSqlOpen)
			fmt.Println("To verify, db is:", db)
			fmt.Println("Создание подключения к БД завершилось ошибкой. Часто возникает из-за не корректного драйвера")
			fmt.Println("Будет предпринята новая попытка через 1 минут")
			time.Sleep(60 * time.Second)
			myError++
			err = errSqlOpen
		}
		if myError == 5 {
			myError = 0
			return nil, err
		}
	} //sql.Open
	return siteApcut_login, nil
}

func (ur *UnifiRepo) GetLoginPCerr(client *entity.Client) (err error) { //(entity.Client, error) {

	//type PC struct {		UserName string `json:"user_name"`	}
	//var err error

	myError := 1
	for myError != 0 {
		if db, errSqlOpen := sql.Open("mysql", ur.dataSourceGLPI); errSqlOpen == nil {
			errDBping := db.Ping()
			if errDBping == nil {
				defer db.Close() // defer the close till after the main function has finished
				//queryAfter := "SELECT * FROM it_support_db.a WHERE controller = " + strconv.Itoa(int(bdController))
				//queryAfter := "SELECT contact FROM glpi_db.glpi_computers where name = ? ORDER BY date_mod DESC"
				queryAfter := "SELECT contact FROM " + ur.databaseGLPI + ".glpi_computers where name = ? ORDER BY date_mod DESC"

				//errQuery := db.QueryRow(queryAfter, pcName).Scan(&pc.UserName)
				errQuery := db.QueryRow(queryAfter, client.Hostname).Scan(&client.UserLogin)
				if errQuery == nil {
					//Если изменилась имя или структура таблицы, то нет смысла зацикливать на 5 минут SELECT
					//result = pc.UserName
					return nil
				} else {
					fmt.Println(errQuery.Error())
					//fmt.Println("В БД нет доступного соответствия имени ПК и логина")
					//result = "denis.tirskikh"
					client.UserLogin = "denis.tirskikh"
					return err
				}
				myError = 0
				//db.Close()
			} else {
				fmt.Println("db.Ping failed:", errDBping)
				fmt.Println("Подключение к БД НЕ установлено. Проверь доступность БД")
				fmt.Println("Будет предпринята новая попытка через 1 минут")
				time.Sleep(60 * time.Second)
				myError++
				err = errDBping
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
			err = errSqlOpen
		}
		if myError == 5 {
			myError = 0
			//result = "denis.tirskikh"
			client.UserLogin = "denis.tirskikh"
			return err
		}
	} //sql.Open
	return nil //result
}
