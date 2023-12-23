package repo

import (
	"bytes"
	"database/sql"
	"github.com/deniskaponchik/GoSoft/internal/entity"
	"log"
	"strconv"
	"strings"
	"time"
)

// под логику, где у Клиентов есть массив Аномалий
func (ur *UnifiRepo) UpdateDbAnomaly(mac_Anomaly map[string]*entity.Anomaly) (err error) {
	//приходит мапа типа: мак адрес клиента _ аномалии клиента за 1 час

	//bdCntrl := strconv.Itoa(int(ur.controller)) //bdController))
	var anomSliceString string
	var query string
	var b1 bytes.Buffer
	b1.WriteString("INSERT INTO  " + ur.databaseITsup + ".anomaly VALUES ")
	countB1 := 0

	for _, v := range mac_Anomaly { //mac_DateSiteAnom {

		if len(v.SliceAnomStr) > 1 {
			//если аномалий за час накопилось 2 и более, то такие заносим в БД
			countB1++
			anomSliceString = strings.Join(v.SliceAnomStr, ";")
			b1.WriteString("('" + v.DateHour + "','" + v.ClientMac + "','" + strconv.Itoa(int(v.Controller)) +
				"','" + v.SiteName + "','" + anomSliceString + "','" + v.ApMac + "','" + v.ApName + "','" + strconv.Itoa(int(v.Exception)) + "'),")
		}

	}
	query = b1.String()
	//Возможно, не самый эффективный метод обрезать строку с конца, но рабочий
	if last := len(query) - 1; last >= 0 && query[last] == ',' {
		query = query[:last]
	}
	log.Println(query)
	if countB1 != 0 {
		//UploadMapsToDBerr(wifiConf.GlpiConnectStringITsupport, query)
		//err = ur.UploadMapsToDBerr(query)
		err = ur.dbExec(query)
		if err == nil {
			return nil
		} else {
			return err
		}
	} else {
		log.Println("Передана пустая карта. Запрос не выполнен")
	}
	log.Println("")

	return nil
}

// Создаёт query  и передаёт в функцию UploadMapsToDBerr
func (ur *UnifiRepo) UpdateDbClient(mac_Client map[string]*entity.Client) (err error) {
	//bdCntrl := strconv.Itoa(int(ur.controller)) //bdController))
	var lenMap int
	var count int
	var exception string
	var b1 bytes.Buffer
	var query string
	lenMap = len(mac_Client)
	count = 0
	apName := ""
	var modified string

	//b1.WriteString("REPLACE INTO " + "it_support_db.ap" + " VALUES ")
	b1.WriteString("REPLACE INTO " + ur.databaseITsup + ".client" + " VALUES ")

	for k, v := range mac_Client {
		exception = strconv.Itoa(int(v.Exception))
		if v.Modified == "" {
			modified = "2001-01-01"
		} else {
			modified = v.Modified
		}
		count++
		if count != lenMap {
			// mac, hostname, controller, exception, srid, ap_name(empty), ap_mac, modified
			b1.WriteString("('" + k + "','" + v.Hostname + "','" + strconv.Itoa(int(v.Controller)) +
				"','" + exception + "','" + v.SrID + "','" + apName + "','" + v.ApMac + "','" + modified + "'),")
		} else {
			b1.WriteString("('" + k + "','" + v.Hostname + "','" + strconv.Itoa(int(v.Controller)) +
				"','" + exception + "','" + v.SrID + "','" + apName + "','" + v.ApMac + "','" + modified + "')")
			//в конце НЕ ставим запятую
		}
	}
	query = b1.String()
	log.Println(query)
	if count != 0 {
		//UploadMapsToDBstring("it_support_db", query)
		//UploadMapsToDBerr(wifiConf.GlpiConnectStringITsupport, query)
		//ur.UploadMapsToDBerr(query)
		err = ur.dbExec(query)
		if err == nil {
			return nil
		} else {
			return err
		}
	} else {
		log.Println("Передана пустая карта. Запрос не выполнен")
	}
	log.Println("")
	return nil
}

func (ur *UnifiRepo) DownloadMacMapsClientApWithAnomaly(macClient map[string]*entity.Client, macAp map[string]*entity.Ap, beforeDays string, timeNow time.Time) (err error) {

	var anomalyRow *entity.Anomaly
	//beforeDays = "2023-09-01 12:00:00"
	todaydayInt := timeNow.Day() //нужен для массива аномалий
	//var date string //2023-09-01  нужна была для мапы аномалий за сутки
	//var date_Anomaly map[string]*entity.Anomaly     //date == 2023-09-01
	//mac_Client = make(map[string]*entity.Client)

	myError := 1
	for myError != 0 {
		if db, errSqlOpen := sql.Open("mysql", ur.dataSourceITsup); errSqlOpen == nil {
			errDBping := db.Ping()
			if errDBping == nil {
				defer db.Close() // defer the close till after the main function has finished
				//queryAfter := "SELECT * FROM it_support_db.anomalies WHERE controller = " + strconv.Itoa(int(bdController))
				//queryAfter := "SELECT * FROM " + ur.databaseITsup + ".anomaly WHERE date_hour >= '" + beforeDays + "' AND controller = " + strconv.Itoa(int(ur.controller)) + " AND exception = 0 order by date_hour DESC"
				queryAfter := "SELECT * FROM " + ur.databaseITsup + ".anomaly WHERE date_hour >= '" + beforeDays + "' AND exception = 0 order by date_hour"
				log.Println(queryAfter)

				for myError != 0 { //зацикливание выполнения запроса
					results, errQuery := db.Query(queryAfter)
					if errQuery == nil {

						var tag entity.Anomaly
						//var tag Tag

						for results.Next() {
							errScan := results.Scan(&tag.DateHour, &tag.ClientMac, &tag.Controller, &tag.SiteName, &tag.AnomStr,
								&tag.ApMac, &tag.ApName, &tag.Exception)

							if errScan == nil {
								//anomSlice = strings.Split(tag.Anomalies, ";")
								tag.SliceAnomStr = strings.Split(tag.AnomStr, ";")
								//date = strings.Split(tag.DateHour, " ")[0]

								//if len(anomSlice) > 2 { //в БД уже записи с двумя и более аномалиями.
								//комментирую на будущее, если захочу пропускать с тремя и более аномалиями

								anomalyRow = &entity.Anomaly{
									DateHour:     tag.DateHour,
									ClientMac:    tag.ClientMac,
									Controller:   tag.Controller, //не использую, если что, в дальнейшем
									SiteName:     tag.SiteName,
									SliceAnomStr: tag.SliceAnomStr,
									ApMac:        tag.ApMac,
									ApName:       tag.ApName,
									Exception:    tag.Exception, //по условию SELECT exception = 0
								}

								client, exisMacClient := macClient[tag.ClientMac]
								if !exisMacClient {
									/*мака в мапе НЕ МОЖЕТ НЕ БЫТЬ. Он создаётся до этой функции в загрузке из БД или в обновлении UI каждые 12 минут
									macClient[tag.ClientMac] = &entity.Client{
										Mac:            tag.ClientMac,
										SliceAnomalies: []*entity.Anomaly{anomalyRow},
										//Date_Anomaly: date_Anomaly,
									}*/

								} else {
									//мак клиента в мапе есть
									if client.Date30count != todaydayInt {
										//если массив сегодня ещё не обновлялся
										client.Date30count = todaydayInt
										client.SliceAnomalies = []*entity.Anomaly{anomalyRow}
									} else {
										//если масив сегодня обновился
										client.SliceAnomalies = append(client.SliceAnomalies, anomalyRow)
									}
								}

								ap, exisMacAp := macAp[tag.ApMac]
								if !exisMacAp {
									/*мака в мапе НЕ МОЖЕТ НЕ БЫТЬ. Он создаётся до этой функции в загрузке из БД или в обновлении UI каждые 12 минут
									macAp[tag.ApMac] = &entity.Ap{
										Mac:            tag.ApMac,
										SliceAnomalies: []*entity.Anomaly{anomalyRow},
										//Date_Anomaly: date_Anomaly,
									}*/
								} else {
									//мак клиента в мапе есть
									if ap.Date30count != todaydayInt {
										//если массив сегодня ещё не обновлялся
										ap.Date30count = todaydayInt
										ap.SliceAnomalies = []*entity.Anomaly{anomalyRow}
									} else {
										//если масив сегодня обновился
										ap.SliceAnomalies = append(ap.SliceAnomalies, anomalyRow)
									}
								}

								//}
							} else {
								//panic(errScan.Error()) // proper error handling instead of panic in your app
								log.Println(errScan.Error())
								log.Println("Сканирование строки и занесение в переменные структуры завершилось ошибкой")
								log.Println("Проверь, что не изменилась структура таблицы и кол-во полей")
								myError = 0
								//break
							}
						}
						if errRowsNext := results.Err(); errRowsNext != nil {
							log.Println("Цикл прохода по результирующим рядам завершился не корректно")
							//если есть ошибка прохода по строкам, отправляем на перезапрос
							myError = 0
						}
						if myError != 1 {
							//results.Close()
							if errRowsClose := results.Close(); errRowsClose != nil {
								log.Println("Закрытие процесса прохода по результирующим полям завершилось не корректно")
							}
							//db.Close()
							if errDBclose := db.Close(); errDBclose != nil {
								log.Println("Закрытие подключения к БД завершилось не корректно")
							}
							myError = 0

						} else {
							//log.Println("Будет предпринята новая попытка запроса через 1 минут")
							//time.Sleep(60 * time.Second)
							myError = 0
						}
					} else {
						//panic(errQuery.Error()) // proper error handling instead of panic in your app
						log.Println(errQuery.Error())
						log.Println("Запрос НЕ смог отработать. Проверь корректность всех данных в запросе")
						//log.Println("Будет предпринята новая попытка через 1 минут")
						//time.Sleep(60 * time.Second)
						myError = 0 //если такой таблицы нет в БД, то что она появится через 5 минут?
					}
				} //db.Query
			} else {
				log.Println("db.Ping failed:", errDBping)
				log.Println("Подключение к БД НЕ установлено. Проверь доступность БД")
				log.Println("Будет предпринята новая попытка через 1 минут")
				time.Sleep(60 * time.Second)
				//myError = 1
				myError++
			}
		} else {
			//log.Print(errSqlOpen.Error())
			log.Println("Error creating DB:", errSqlOpen)
			log.Println("To verify, db is:", db)
			log.Println("Создание подключения к БД завершилось ошибкой. Часто возникает из-за не корректного драйвера")
			log.Println("Будет предпринята новая попытка через 1 минут")
			time.Sleep(60 * time.Second)
			//myError = 1
			myError++
		}
		if myError == 300 { //Если ночью сервер перезагрузился + нет доступа к БД = в ЦОДЕ коллапс. Могу подождать 5 часов
			myError = 0
		}
	} //sql.Open
	return nil
}

func (ur *UnifiRepo) DownloadClientsWithAnomalySlice(mac_Client map[string]*entity.Client, beforeDays string, timeNow time.Time) (err error) {

	var anomalyRow *entity.Anomaly
	//beforeDays = "2023-09-01 12:00:00"
	todaydayInt := timeNow.Day() //нужен для массива аномалий
	//var date string //2023-09-01  нужна была для мапы аномалий за сутки
	//var date_Anomaly map[string]*entity.Anomaly     //date == 2023-09-01
	//mac_Client = make(map[string]*entity.Client)

	myError := 1
	for myError != 0 {
		if db, errSqlOpen := sql.Open("mysql", ur.dataSourceITsup); errSqlOpen == nil {
			errDBping := db.Ping()
			if errDBping == nil {
				defer db.Close() // defer the close till after the main function has finished
				//queryAfter := "SELECT * FROM it_support_db.anomalies WHERE controller = " + strconv.Itoa(int(bdController))
				queryAfter := "SELECT * FROM " + ur.databaseITsup + ".anomaly WHERE date_hour >= '" + beforeDays + "' AND controller = " +
					strconv.Itoa(int(ur.controller)) + " AND exception = 0 order by date_hour DESC"
				log.Println(queryAfter)

				for myError != 0 { //зацикливание выполнения запроса
					results, errQuery := db.Query(queryAfter)
					if errQuery == nil {

						var tag entity.Anomaly
						//var tag Tag

						for results.Next() {
							errScan := results.Scan(&tag.DateHour, &tag.ClientMac, &tag.Controller, &tag.SiteName, &tag.AnomStr,
								&tag.ApMac, &tag.ApName, &tag.Exception)

							if errScan == nil {
								//anomSlice = strings.Split(tag.Anomalies, ";")
								tag.SliceAnomStr = strings.Split(tag.AnomStr, ";")
								//date = strings.Split(tag.DateHour, " ")[0]

								//if len(anomSlice) > 2 { //в БД уже записи с двумя и более аномалиями.
								//комментирую на будущее, если захочу пропускать с тремя и более аномалиями

								anomalyRow = &entity.Anomaly{
									DateHour:     tag.DateHour,
									ClientMac:    tag.ClientMac,
									Controller:   tag.Controller, //не использую, если что, в дальнейшем
									SiteName:     tag.SiteName,
									SliceAnomStr: tag.SliceAnomStr,
									ApMac:        tag.ApMac,
									ApName:       tag.ApName,
									Exception:    tag.Exception, //по условию SELECT exception = 0
								}

								client, exisMacClient := mac_Client[tag.ClientMac]
								if !exisMacClient {
									mac_Client[tag.ClientMac] = &entity.Client{
										Mac:            tag.ClientMac,
										SliceAnomalies: []*entity.Anomaly{anomalyRow},
										//Date_Anomaly: date_Anomaly,
									}
								} else {
									//мак клиента в мапе есть
									if client.Date30count != todaydayInt {
										//если массив сегодня ещё не обновлялся
										client.Date30count = todaydayInt
										client.SliceAnomalies = []*entity.Anomaly{anomalyRow}
									} else {
										//если масив сегодня обновился
										client.SliceAnomalies = append(client.SliceAnomalies, anomalyRow)
									}
								}
								//}
							} else {
								//panic(errScan.Error()) // proper error handling instead of panic in your app
								log.Println(errScan.Error())
								log.Println("Сканирование строки и занесение в переменные структуры завершилось ошибкой")
								log.Println("Проверь, что не изменилась структура таблицы и кол-во полей")
								myError = 0
								//break
							}
						}
						if errRowsNext := results.Err(); errRowsNext != nil {
							log.Println("Цикл прохода по результирующим рядам завершился не корректно")
							//если есть ошибка прохода по строкам, отправляем на перезапрос
							myError = 0
						}
						if myError != 1 {
							//results.Close()
							if errRowsClose := results.Close(); errRowsClose != nil {
								log.Println("Закрытие процесса прохода по результирующим полям завершилось не корректно")
							}
							//db.Close()
							if errDBclose := db.Close(); errDBclose != nil {
								log.Println("Закрытие подключения к БД завершилось не корректно")
							}
							myError = 0

						} else {
							//log.Println("Будет предпринята новая попытка запроса через 1 минут")
							//time.Sleep(60 * time.Second)
							myError = 0
						}
					} else {
						//panic(errQuery.Error()) // proper error handling instead of panic in your app
						log.Println(errQuery.Error())
						log.Println("Запрос НЕ смог отработать. Проверь корректность всех данных в запросе")
						//log.Println("Будет предпринята новая попытка через 1 минут")
						//time.Sleep(60 * time.Second)
						myError = 0 //если такой таблицы нет в БД, то что она появится через 5 минут?
					}
				} //db.Query
			} else {
				log.Println("db.Ping failed:", errDBping)
				log.Println("Подключение к БД НЕ установлено. Проверь доступность БД")
				log.Println("Будет предпринята новая попытка через 1 минут")
				time.Sleep(60 * time.Second)
				//myError = 1
				myError++
			}
		} else {
			//log.Print(errSqlOpen.Error())
			log.Println("Error creating DB:", errSqlOpen)
			log.Println("To verify, db is:", db)
			log.Println("Создание подключения к БД завершилось ошибкой. Часто возникает из-за не корректного драйвера")
			log.Println("Будет предпринята новая попытка через 1 минут")
			time.Sleep(60 * time.Second)
			//myError = 1
			myError++
		}
		if myError == 300 { //Если ночью сервер перезагрузился + нет доступа к БД = в ЦОДЕ коллапс. Могу подождать 5 часов
			myError = 0
		}
	} //sql.Open
	return nil
}

func (ur *UnifiRepo) DownloadMacClientsWithAnomalies(mac_Client map[string]*entity.Client, beforeDays string, timeNow time.Time) (err error) {

	//beforeDays = "2023-09-01 12:00:00"
	//var anomSlice []string
	var date string //2023-09-01
	//var date_Anomaly map[string]*entity.Anomaly     //date == 2023-09-01
	//mac_Client = make(map[string]*entity.Client)

	myError := 1
	for myError != 0 {
		if db, errSqlOpen := sql.Open("mysql", ur.dataSourceITsup); errSqlOpen == nil {
			errDBping := db.Ping()
			if errDBping == nil {
				defer db.Close() // defer the close till after the main function has finished
				//queryAfter := "SELECT * FROM it_support_db.anomalies WHERE controller = " + strconv.Itoa(int(bdController))
				queryAfter := "SELECT * FROM " + ur.databaseITsup + ".anomaly WHERE date_hour >= '" + beforeDays +
					"' AND controller = " + strconv.Itoa(int(ur.controller)) + " AND exception = 0"
				log.Println(queryAfter)

				for myError != 0 { //зацикливание выполнения запроса
					results, errQuery := db.Query(queryAfter)
					if errQuery == nil {

						var tag entity.Anomaly
						//var tag Tag

						for results.Next() {
							errScan := results.Scan(&tag.DateHour, &tag.ClientMac, &tag.Controller, &tag.SiteName, &tag.AnomStr,
								&tag.ApMac, &tag.ApName, &tag.Exception)

							if errScan == nil {
								entityAnomaly := &entity.Anomaly{
									DateHour:     tag.DateHour,
									ClientMac:    tag.ClientMac,
									Controller:   tag.Controller, //не использую, если что, в дальнейшем
									SiteName:     tag.SiteName,
									SliceAnomStr: tag.SliceAnomStr,
									ApMac:        tag.ApMac,
									ApName:       tag.ApName,
									Exception:    tag.Exception, //по условию SELECT exception = 0
								}

								//anomSlice = strings.Split(tag.Anomalies, ";")
								tag.SliceAnomStr = strings.Split(tag.AnomStr, ";")
								date = strings.Split(tag.DateHour, " ")[0]

								//if len(anomSlice) > 2 { //в БД уже записи с двумя и более аномалиями.
								//комментирую на будущее, если захочу пропускать с тремя и более аномалиями

								client, exisMacClient := mac_Client[tag.ClientMac]
								if !exisMacClient {
									//log.Println(tag.ClientMac + " не был создан в мапе Клиентов. СОЗДАНИЕ клиента")

									date_Anomaly := make(map[string]*entity.Anomaly) //date == 2023-09-01
									date_Anomaly[date] = entityAnomaly

									mac_Client[tag.ClientMac] = &entity.Client{
										Mac: tag.ClientMac,
										//Hostname не хватает только его. Можно добавить вручную точечно из старой базы потом
										Date_Anomaly:            date_Anomaly,
										DateTicketCreateAttempt: timeNow.Day(),
									}

								} else {
									//если мак клиента уже был мапе. Самый распространённый случай.

									if client.DateTicketCreateAttempt != timeNow.Day() {
										//если заходов на создание мапы аномалий сегодня ещё не было
										//log.Println(tag.ClientMac + "заходов на создание мапы сегодня ещё не было. СОЗДАНИЕ мапы")

										client.DateTicketCreateAttempt = timeNow.Day()

										date_Anomaly := make(map[string]*entity.Anomaly) //date == 2023-09-01
										date_Anomaly[date] = entityAnomaly

										client.Date_Anomaly = date_Anomaly

									} else {
										//если новая мапа аномалий сегодня уже была создана
										//log.Println(tag.ClientMac + "мапа аномалий сегодня уже создана. ОБНОВЛЕНИЕ мапы")
										//создаём новый бакет
										//client.Date_Anomaly[date] = entityAnomaly
										client.Date_Anomaly[date] = &entity.Anomaly{
											DateHour:     tag.DateHour,
											ClientMac:    tag.ClientMac,
											Controller:   tag.Controller, //не использую, если что, в дальнейшем
											SiteName:     tag.SiteName,
											SliceAnomStr: tag.SliceAnomStr,
											ApMac:        tag.ApMac,
											ApName:       tag.ApName,
											Exception:    tag.Exception, //по условию SELECT exception = 0
										}
									}
								}

								//} //if len(anomSlice) > 2 {
							} else {
								//panic(errScan.Error()) // proper error handling instead of panic in your app
								log.Println(errScan.Error())
								log.Println("Сканирование строки и занесение в переменные структуры завершилось ошибкой")
								log.Println("Проверь, что не изменилась структура таблицы и кол-во полей")
								myError = 0
								//break
							}
						}
						if errRowsNext := results.Err(); errRowsNext != nil {
							log.Println("Цикл прохода по результирующим рядам завершился не корректно")
							//если есть ошибка прохода по строкам, отправляем на перезапрос
							myError = 0
						}
						if myError != 1 {
							//results.Close()
							if errRowsClose := results.Close(); errRowsClose != nil {
								log.Println("Закрытие процесса прохода по результирующим полям завершилось не корректно")
							}
							//db.Close()
							if errDBclose := db.Close(); errDBclose != nil {
								log.Println("Закрытие подключения к БД завершилось не корректно")
							}
							myError = 0

						} else {
							//log.Println("Будет предпринята новая попытка запроса через 1 минут")
							//time.Sleep(60 * time.Second)
							myError = 0
						}
					} else {
						//panic(errQuery.Error()) // proper error handling instead of panic in your app
						log.Println(errQuery.Error())
						log.Println("Запрос НЕ смог отработать. Проверь корректность всех данных в запросе")
						//log.Println("Будет предпринята новая попытка через 1 минут")
						//time.Sleep(60 * time.Second)
						myError = 0 //если такой таблицы нет в БД, то что она появится через 5 минут?
					}
				} //db.Query
			} else {
				log.Println("db.Ping failed:", errDBping)
				log.Println("Подключение к БД НЕ установлено. Проверь доступность БД")
				log.Println("Будет предпринята новая попытка через 1 минут")
				time.Sleep(60 * time.Second)
				//myError = 1
				myError++
			}
		} else {
			//log.Print(errSqlOpen.Error())
			log.Println("Error creating DB:", errSqlOpen)
			log.Println("To verify, db is:", db)
			log.Println("Создание подключения к БД завершилось ошибкой. Часто возникает из-за не корректного драйвера")
			log.Println("Будет предпринята новая попытка через 1 минут")
			time.Sleep(60 * time.Second)
			//myError = 1
			myError++
		}
		if myError == 300 { //Если ночью сервер перезагрузился + нет доступа к БД = в ЦОДЕ коллапс. Могу подождать 5 часов
			myError = 0
		}
	} //sql.Open
	//return mac_Client, nil
	return nil
}

func (ur *UnifiRepo) Download2MapFromDBclient() (map[string]*entity.Client, map[string]*entity.Client, error) {

	macClient := make(map[string]*entity.Client) //https://yourbasic.org/golang/gotcha-assignment-entry-nil-map/
	hostnameClient := make(map[string]*entity.Client)
	var clPointer *entity.Client //клиент создаётся при каждом взятии из массива
	var upperCaseHostName string
	var err error

	myError := 1
	for myError != 0 {
		if db, errSqlOpen := sql.Open("mysql", ur.dataSourceITsup); errSqlOpen == nil {
			errDBping := db.Ping()
			if errDBping == nil {
				defer db.Close() // defer the close till after the main function has finished
				//queryAfter := "SELECT * FROM it_support_db.machine WHERE controller = " + strconv.Itoa(int(bdController))
				//queryAfter := "SELECT * FROM " + ur.databaseITsup + ".client WHERE controller = " + strconv.Itoa(int(ur.controller))
				queryAfter := "SELECT * FROM " + ur.databaseITsup + ".client" // WHERE controller = " + strconv.Itoa(int(ur.controller))
				log.Println(queryAfter)

				for myError != 0 { //зацикливание выполнения запроса
					results, errQuery := db.Query(queryAfter)
					if errQuery == nil {
						//var tag TagPoly
						var tag entity.Client

						for results.Next() {
							errScan := results.Scan(&tag.Mac, &tag.Hostname, &tag.Controller, &tag.Exception, &tag.SrID,
								&tag.ApName, &tag.ApMac, &tag.Modified)
							if errScan == nil {
								//log.Println(tag.Mac, tag.Name, tag.Controller, tag.Exception, tag.SrID)
								//machineMap[tag.Mac] = &tag

								upperCaseHostName = strings.ToUpper(tag.Hostname)

								clPointer = &entity.Client{
									Mac:        tag.Mac,
									Hostname:   upperCaseHostName,
									Controller: tag.Controller,
									Exception:  tag.Exception,
									SrID:       tag.SrID,
									ApName:     tag.ApName,
									ApMac:      tag.ApMac,
									Modified:   tag.Modified,
								}

								macClient[tag.Mac] = clPointer
								hostnameClient[upperCaseHostName] = clPointer

							} else {
								log.Println(errScan.Error())
								log.Println("Сканирование СТРОКИ и занесение в переменные структуры завершилось ошибкой")
								log.Println("Проверь, что не изменилась структура таблицы и кол-во полей")
								myError = 0
							}
						}
						if errRowsNext := results.Err(); errRowsNext != nil {
							log.Println("Цикл прохода по результирующим рядам завершился не корректно")
							//если есть ошибка прохода по строкам, отправляем на перезапрос
							myError = 0
						}
						if myError != 1 {
							//results.Close()
							if errRowsClose := results.Close(); errRowsClose != nil {
								log.Println("Закрытие процесса прохода по результирующим полям завершилось не корректно")
							}
							//db.Close()
							if errDBclose := db.Close(); errDBclose != nil {
								log.Println("Закрытие подключения к БД завершилось не корректно")
							}
							myError = 0
							/*
								log.Println("Вывод мапы ВНУТРИ функции")
								for k, v := range m {
									log.Println("innerMap "+k, v.Name, v.Exception, v.SrID)
								}*/
						} else {
							//log.Println("Будет предпринята новая попытка запроса через 1 минут")
							//time.Sleep(60 * time.Second)
							myError = 0
						}
					} else {
						log.Println(errQuery.Error())
						log.Println("Запрос НЕ смог отработать. Проверь корректность всех данных в запросе")
						myError = 0 //если такой таблицы нет в БД, то что она появится через 5 минут?
						err = errQuery
					}
				} //db.Query
			} else {
				log.Println("db.Ping failed:", errDBping)
				log.Println("Подключение к БД НЕ установлено. Проверь доступность БД")
				log.Println("Будет предпринята новая попытка через 1 минут")
				time.Sleep(60 * time.Second)
				//myError = 1
				myError++
				err = errDBping
				//Если ночью сервер перезагрузился + нет доступа к БД = в ЦОДЕ коллапс. Могу подождать 5 часов
				//if myError == 300 { 	myError = 0				}
			}
		} else {
			log.Println("Error creating DB:", errSqlOpen)
			log.Println("To verify, db is:", db)
			log.Println("Создание подключения к БД завершилось ошибкой. Часто возникает из-за не корректного драйвера")
			log.Println("Будет предпринята новая попытка через 1 минут")
			time.Sleep(60 * time.Second)
			//myError = 1
			myError++
			err = errSqlOpen
		}
		if myError == 5 { //&& marker == 1 {
			//Если ночью нет доступа к БД = в ЦОДЕ коллапс. Могу подождать 5 часов при условии, что это ежечасовая актуализация ip-адресов
			myError = 0
			return nil, nil, err //errors.New("подключение к бд не удалось")
		}
	} //sql.Open
	return macClient, hostnameClient, err
}

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
				queryAfter := "SELECT * FROM " + ur.databaseITsup + ".client WHERE controller = " + strconv.Itoa(int(ur.controller))
				log.Println(queryAfter)

				for myError != 0 { //зацикливание выполнения запроса
					results, errQuery := db.Query(queryAfter)
					if errQuery == nil {
						//var tag TagPoly
						var tag entity.Client

						for results.Next() {
							errScan := results.Scan(&tag.Mac, &tag.Hostname, &tag.Controller, &tag.Exception, &tag.SrID,
								&tag.ApName, &tag.ApMac, &tag.Modified)
							if errScan == nil {
								//log.Println(tag.Mac, tag.Name, tag.Controller, tag.Exception, tag.SrID)
								//machineMap[tag.Mac] = &tag
								machineMap[tag.Mac] = &entity.Client{
									Mac:        tag.Mac,
									Hostname:   tag.Hostname,
									Controller: tag.Controller,
									Exception:  tag.Exception,
									SrID:       tag.SrID,
									ApName:     tag.ApName,
									ApMac:      tag.ApMac,
									Modified:   tag.Modified,
								}

							} else {
								log.Println(errScan.Error())
								log.Println("Сканирование СТРОКИ и занесение в переменные структуры завершилось ошибкой")
								log.Println("Проверь, что не изменилась структура таблицы и кол-во полей")
								myError = 0
							}
						}
						if errRowsNext := results.Err(); errRowsNext != nil {
							log.Println("Цикл прохода по результирующим рядам завершился не корректно")
							//если есть ошибка прохода по строкам, отправляем на перезапрос
							myError = 0
						}
						if myError != 1 {
							//results.Close()
							if errRowsClose := results.Close(); errRowsClose != nil {
								log.Println("Закрытие процесса прохода по результирующим полям завершилось не корректно")
							}
							//db.Close()
							if errDBclose := db.Close(); errDBclose != nil {
								log.Println("Закрытие подключения к БД завершилось не корректно")
							}
							myError = 0
							/*
								log.Println("Вывод мапы ВНУТРИ функции")
								for k, v := range m {
									log.Println("innerMap "+k, v.Name, v.Exception, v.SrID)
								}*/
						} else {
							//log.Println("Будет предпринята новая попытка запроса через 1 минут")
							//time.Sleep(60 * time.Second)
							myError = 0
						}
					} else {
						log.Println(errQuery.Error())
						log.Println("Запрос НЕ смог отработать. Проверь корректность всех данных в запросе")
						myError = 0 //если такой таблицы нет в БД, то что она появится через 5 минут?
						err = errQuery
					}
				} //db.Query
			} else {
				log.Println("db.Ping failed:", errDBping)
				log.Println("Подключение к БД НЕ установлено. Проверь доступность БД")
				log.Println("Будет предпринята новая попытка через 1 минут")
				time.Sleep(60 * time.Second)
				//myError = 1
				myError++
				err = errDBping
				//Если ночью сервер перезагрузился + нет доступа к БД = в ЦОДЕ коллапс. Могу подождать 5 часов
				//if myError == 300 { 	myError = 0				}
			}
		} else {
			log.Println("Error creating DB:", errSqlOpen)
			log.Println("To verify, db is:", db)
			log.Println("Создание подключения к БД завершилось ошибкой. Часто возникает из-за не корректного драйвера")
			log.Println("Будет предпринята новая попытка через 1 минут")
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

// из БД GLPI
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
					return nil
				} else {
					log.Println(errQuery.Error())
					//log.Println("В БД нет доступного соответствия имени ПК и логина")
					//client.UserLogin = "denis.tirskikh"
					return errQuery
				}
				myError = 0
				//db.Close()
			} else {
				log.Println("db.Ping failed:", errDBping)
				log.Println("Подключение к БД НЕ установлено. Проверь доступность БД")
				log.Println("Будет предпринята новая попытка через 1 минут")
				time.Sleep(60 * time.Second)
				myError++
				err = errDBping
			}
		} else {
			//По факту подключения к БД НЕ происходит на этом этапе
			//https://stackoverflow.com/questions/32345124/why-does-sql-open-return-nil-as-error-when-it-should-not
			log.Println("Error creating DB:", errSqlOpen)
			log.Println("To verify, db is:", db)
			log.Println("Создание подключения к БД завершилось ошибкой. Часто возникает из-за не корректного драйвера")
			log.Println("Будет предпринята новая попытка через 1 минут")
			time.Sleep(60 * time.Second)
			//myError = 1
			myError++
			err = errSqlOpen
		}
		if myError == 5 {
			myError = 0
			//client.UserLogin = "denis.tirskikh"
			return err
		}
	} //sql.Open
	return nil //result
}
