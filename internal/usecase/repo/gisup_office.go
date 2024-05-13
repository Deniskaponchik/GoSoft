package repo

import (
	"database/sql"
	"github.com/deniskaponchik/GoSoft/internal/entity"
	"log"
	"time"
)

func (gr *Repo) InsertOffice(newOffice *entity.Office) (err error) {
	//site_apcut_login
	query := "INSERT INTO  " + gr.repoMySql.DataBase + ".office VALUES " +
		"('" + newOffice.Site_ApCutName + "','" + newOffice.UserLogin + "','" + newOffice.TimeZoneStr + "','0');"

	err = gr.repoMySql.DbExec(query)
	if err == nil {
		return nil
	} else {
		return err
	}
}

func (gr *Repo) UpdateOfficeLogin(sapcn string, newLogin string) (err error) {
	//site_apcut_login
	query := "UPDATE " + gr.repoMySql.DataBase + ".office SET login = '" + newLogin + "' WHERE site_apcut = '" + sapcn + "';"
	err = gr.repoMySql.DbExec(query)
	if err == nil {
		return nil
	} else {
		return err
	}
}

func (gr *Repo) UpdateOfficeException(sapcn string, newException string) (err error) {
	//site_apcut_login
	query := "UPDATE " + gr.repoMySql.DataBase + ".office SET exception = '" + newException + "' WHERE site_apcut = '" + sapcn + "';"
	err = gr.repoMySql.DbExec(query)
	if err == nil {
		return nil
	} else {
		return err
	}
}

func (gr *Repo) DownloadMapOffice() (map[string]*entity.Office, error) {

	officeMap := make(map[string]*entity.Office) //https://yourbasic.org/golang/gotcha-assignment-entry-nil-map/
	var err error
	myError := 1
	for myError != 0 {
		if db, errSqlOpen := sql.Open("mysql", gr.repoMySql.DataSource); errSqlOpen == nil {
			errDBping := db.Ping()
			if errDBping == nil {
				defer db.Close() // defer the close till after the main function has finished
				//queryAfter := "SELECT * FROM " + ur.database + ".poly"
				//queryAfter := "SELECT * FROM it_support_db.ap WHERE controller = " + strconv.Itoa(int(bdController))
				queryAfter := "SELECT * FROM " + gr.repoMySql.DataBase + ".office"
				log.Println(queryAfter)

				for myError != 0 { //зацикливание выполнения запроса
					results, errQuery := db.Query(queryAfter)
					if errQuery == nil {
						//var tag TagPoly
						var tag entity.Office

						for results.Next() {
							errScan := results.Scan(&tag.Site_ApCutName, &tag.UserLogin, &tag.TimeZone) //, &tag.Exception)
							//errScan := results.Scan(&tag.Mac, &tag.Name, &tag.Controller, &tag.SrID)
							if errScan == nil {
								//log.Println(tag.Mac, tag.Name, tag.Controller, tag.Exception, tag.SrID)

								//тоже рабочее
								//var ap entity.Ap
								//ap = tag
								//apMap[tag.Mac] = &ap //&tag

								officeMap[tag.Site_ApCutName] = &entity.Office{
									Site_ApCutName: tag.Site_ApCutName,
									UserLogin:      tag.UserLogin,
									TimeZone:       tag.TimeZone,
									//Exception:      tag.Exception,
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
	return officeMap, err
}
