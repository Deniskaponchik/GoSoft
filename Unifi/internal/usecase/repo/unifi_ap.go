package repo

import (
	"bytes"
	"database/sql"
	"github.com/deniskaponchik/GoSoft/Unifi/internal/entity"
	"log"
	"strconv"
	"strings"
	"time"
)

// Создаёт query  и передаёт в функцию UploadMapsToDBerr
func (ur *UnifiRepo) UpdateDbAp(mapAp map[string]*entity.Ap) (err error) {
	//bdCntrl := strconv.Itoa(int(ur.controller)) //bdController))
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
			b1.WriteString("('" + k + "','" + v.Name + "','" + strconv.Itoa(int(v.Controller)) + "','" + exception + "','" + v.SrID + "'),")
			// mac, name, controller, srid
			//b1.WriteString("('" + k + "','" + v.Name + "','" + bdCntrl + "','" + v.SrID + "'),")
		} else {
			b1.WriteString("('" + k + "','" + v.Name + "','" + strconv.Itoa(int(v.Controller)) + "','" + exception + "','" + v.SrID + "')")
			//b1.WriteString("('" + k + "','" + v.Name + "','" + bdCntrl + "','" + v.SrID + "')")
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

func (ur *UnifiRepo) Download2MapFromDBaps() (map[string]*entity.Ap, map[string]*entity.Ap, error) {

	macAp := make(map[string]*entity.Ap)      //https://yourbasic.org/golang/gotcha-assignment-entry-nil-map/
	hostnameAp := make(map[string]*entity.Ap) //https://yourbasic.org/golang/gotcha-assignment-entry-nil-map/
	var apPointer *entity.Ap                  //клиент создаётся при каждом взятии из массива
	//var ticketPointer *entity.Ticket  //если буду вкладывать сущность Тикета, а не srID
	var upperCaseHostName string
	var err error

	myError := 1
	for myError != 0 {
		if db, errSqlOpen := sql.Open("mysql", ur.dataSourceITsup); errSqlOpen == nil {
			errDBping := db.Ping()
			if errDBping == nil {
				defer db.Close() // defer the close till after the main function has finished
				//queryAfter := "SELECT * FROM " + ur.databaseITsup + ".ap WHERE controller = " + strconv.Itoa(int(ur.controller))
				queryAfter := "SELECT * FROM " + ur.databaseITsup + ".ap" // WHERE controller = " + strconv.Itoa(int(ur.controller))
				log.Println(queryAfter)

				for myError != 0 { //зацикливание выполнения запроса
					results, errQuery := db.Query(queryAfter)
					if errQuery == nil {
						//var tag TagPoly
						var tag entity.Ap

						for results.Next() {
							errScan := results.Scan(&tag.Mac, &tag.Name, &tag.Controller, &tag.Exception, &tag.SrID)
							//errScan := results.Scan(&tag.Mac, &tag.Name, &tag.Controller, &tag.SrID)
							if errScan == nil {
								//log.Println(tag.Mac, tag.Name, tag.Controller, tag.Exception, tag.SrID)

								/*если буду вкладывать сущность Тикета, а не srID
								if tag.SrID != "" {
									ticketPointer = &entity.Ticket{
										ID: tag.SrID,
									}
								}*/

								upperCaseHostName = strings.ToUpper(tag.Name)

								apPointer = &entity.Ap{
									Mac:        tag.Mac,
									Name:       upperCaseHostName,
									Controller: tag.Controller,
									Exception:  tag.Exception,
									SrID:       tag.SrID,
								}

								macAp[tag.Mac] = apPointer
								hostnameAp[upperCaseHostName] = apPointer

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
	return macAp, hostnameAp, err
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
				log.Println(queryAfter)

				for myError != 0 { //зацикливание выполнения запроса
					results, errQuery := db.Query(queryAfter)
					if errQuery == nil {
						//var tag TagPoly
						var tag entity.Ap

						for results.Next() {
							errScan := results.Scan(&tag.Mac, &tag.Name, &tag.Controller, &tag.Exception, &tag.SrID)
							//errScan := results.Scan(&tag.Mac, &tag.Name, &tag.Controller, &tag.SrID)
							if errScan == nil {
								//log.Println(tag.Mac, tag.Name, tag.Controller, tag.Exception, tag.SrID)

								//тоже рабочее
								//var ap entity.Ap
								//ap = tag
								//apMap[tag.Mac] = &ap //&tag

								apMap[tag.Mac] = &entity.Ap{
									Mac:        tag.Mac,
									Name:       tag.Name,
									Controller: tag.Controller,
									Exception:  tag.Exception,
									SrID:       tag.SrID,
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
	return apMap, err
}
