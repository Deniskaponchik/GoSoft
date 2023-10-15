package repo

//https://tutorialedge.net/golang/golang-mysql-tutorial/
//https://github.com/evrone/go-clean-template/blob/master/internal/usecase/repo/translation_postgres.go
import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	//"github.com/go-sql-driver/mysql"
	"github.com/deniskaponchik/GoSoft/Unifi/internal/entity"
	"time"
)

type PolyRepo struct {
	dataSource string
}

// реализуем Инъекцию зависимостей DI. Используется в app
func New(d string) (*PolyRepo, error) {
	pr := &PolyRepo{
		dataSource: d,
	}

	if db, errSqlOpen := sql.Open("mysql", pr.dataSource); errSqlOpen == nil {
		errDBping := db.Ping()
		if errDBping == nil {

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
	return pr, nil
}

func (pr *PolyRepo) UpdateMapsToDBerr(polyMap map[string]entity.PolyStruct) (err error) {

	//взято из app
	var queries []string
	for k, v := range polyMap {
		queries = append(queries, "UPDATE it_support_db.poly SET srid = '"+v.SrID+"', comment = "+strconv.Itoa(int(v.Comment))+" WHERE mac = '"+k+"';")
	}

	myError := 1
	for myError != 0 {
		if db, errSqlOpen := sql.Open("mysql", pr.dataSource); errSqlOpen == nil {
			errDBping := db.Ping()
			if errDBping == nil {
				defer db.Close() // defer the close till after the main function has finished

				for _, query := range queries {
					fmt.Println(query)
					_, errQuery := db.Exec(query)
					if errQuery != nil {
						fmt.Println(errQuery.Error())
						fmt.Println("Запрос НЕ смог отработать. Проверь корректность всех данных в запросе")
						//myError = 0
						//err = errQuery //ошибка не критическая. запросов много и не все будут не корректными
					}
				}
				if errDBclose := db.Close(); errDBclose != nil {
					fmt.Println("Закрытие подключения к БД завершилось не корректно")
					//return errDBclose //ошибка не критическая
				}
				return nil
			} else {
				fmt.Println("db.Ping failed:", errDBping)
				fmt.Println("Подключение к БД НЕ установлено. Проверь доступность БД")
				fmt.Println("Будет предпринята новая попытка через 1 минут")
				time.Sleep(60 * time.Second)
				myError++
			}
		} else {
			fmt.Println("Error creating DB:", errSqlOpen)
			fmt.Println("To verify, db is:", db)
			fmt.Println("Создание подключения к БД завершилось ошибкой. Часто возникает из-за не корректного драйвера")
			fmt.Println("Будет предпринята новая попытка через 1 минут")
			time.Sleep(60 * time.Second)
			myError++
		}
		if myError == 5 {
			myError = 0
			fmt.Println("После 5 неудачных попыток идём дальше. Подключение к БД не удалось")
			return errors.New("подключение к бд не удалось")
		}
	} //sql.Open
	return nil
}

func (pr *PolyRepo) DownloadMapFromDBvcsErr(marker int) (polyMap map[string]entity.PolyStruct, err error) {
	//Функция вызывается НЕ только в начале скрипта, но каждый час для корректности ip-адресов

	/*
		type TagPoly struct {
			Mac       string `json:"mac"`
			IP        string `json:"ip"`
			Region    string `json:"region"`
			RoomName  string `json:"room_name"`
			Login     string `json:"login"`
			SrID      string `json:"srid"`
			PolyType  int    `json:"type"`
			Comment   int    `json:"comment"`
			Exception int    `json:"exception"`
		} */
	//m := make(map[string]entity.PolyStruct)

	myError := 1
	for myError != 0 {
		if db, errSqlOpen := sql.Open("mysql", pr.dataSource); errSqlOpen == nil {
			errDBping := db.Ping()
			if errDBping == nil {
				defer db.Close()                                 // defer the close till after the main function has finished
				queryAfter := "SELECT * FROM it_support_db.poly" // WHERE type = " + strconv.Itoa(int(polyType))
				//queryAfter := "SELECT * FROM it_support_db.a WHERE controller = " + strconv.Itoa(int(bdController))
				fmt.Println(queryAfter)

				for myError != 0 { //зацикливание выполнения запроса
					results, errQuery := db.Query(queryAfter)
					if errQuery == nil {
						//var tag TagPoly
						var tag entity.PolyStruct

						for results.Next() {
							errScan := results.Scan(&tag.Mac, &tag.IP, &tag.Region, &tag.RoomName, &tag.Login, &tag.SrID, &tag.PolyType, &tag.Comment, &tag.Exception)
							if errScan == nil {
								//fmt.Println(tag.Mac, tag.Name, tag.Controller, tag.Exception, tag.SrID)
								polyMap[tag.Mac] = entity.PolyStruct{
									tag.Mac,
									tag.IP,
									tag.Region,
									tag.RoomName,
									tag.Login,
									tag.SrID,
									tag.PolyType,
									tag.Comment,
									tag.Exception,
									"",
								}
							} else {
								//panic(errScan.Error()) // proper error handling instead of panic in your app
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
		if myError == 5 && marker == 1 {
			//Если ночью нет доступа к БД = в ЦОДЕ коллапс. Могу подождать 5 часов при условии, что это ежечасовая актуализация ip-адресов
			myError = 0
			return nil, errors.New("подключение к бд не удалось")
		}
	} //sql.Open
	return polyMap, err
}
