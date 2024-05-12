// Package postgres implements postgres connection.
package mysql1

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql" //для установки драйвера mysql. Сам пакет как бы не используется явно в коде
	"log"
	"time"
	//_ "github.com/go-sql-driver/mysql" //для установки драйвера mysql. Сам пакет как бы не используется явно в коде
)

const (
	_defaultMaxPoolSize  = 1
	_defaultConnAttempts = 5
	_defaultConnTimeout  = time.Second
)

// MySql
type SqlMy struct {
	maxPoolSize  int
	connAttempts int
	connTimeout  time.Duration

	dataSource string
	dataBase   string

	//Builder squirrel.StatementBuilderType
	//Pool    *pgxpool.Pool
}

// New -.
func NewSqlMy(connectStr string, base string, opts ...Option) (*SqlMy, error) {
	sm := &SqlMy{
		dataSource: connectStr + "/" + base,
		dataBase:   base,

		maxPoolSize:  _defaultMaxPoolSize,
		connAttempts: _defaultConnAttempts,
		connTimeout:  _defaultConnTimeout,
	}

	// Custom options
	for _, opt := range opts {
		opt(sm)
	}

	log.Println(connectStr + "/" + base)

	if db, errSqlOpen := sql.Open("mysql", sm.dataSource); errSqlOpen == nil {
		errDBping := db.Ping()
		if errDBping == nil {
			return sm, nil
		} else {
			log.Println("db.Ping failed:", errDBping)
			log.Println("Подключение к БД НЕ установлено. Проверь доступность БД")
			return nil, errDBping
		}
	} else {
		log.Println("Error creating DB:", errSqlOpen)
		log.Println("To verify, db is:", db)
		log.Println("Создание подключения к БД завершилось ошибкой. Часто возникает из-за не корректного драйвера")
		return nil, errSqlOpen
	}

	//return sm, nil
}

func (sm *SqlMy) dbExec(query string) (err error) {
	myError := 1
	for myError != 0 {
		if db, errSqlOpen := sql.Open("mysql", sm.dataSource); errSqlOpen == nil {
			errDBping := db.Ping()
			if errDBping == nil {
				defer db.Close() // defer the close till after the main function has finished
				//for myError != 0 { //зачем зацикливать выполнение запроса при корректном подключении к БД?
				_, errQuery := db.Exec(query)
				if errQuery == nil {
					myError = 0
					return nil
				} else {
					//panic(errQuery.Error()) // proper error handling instead of panic in your app
					log.Println(errQuery.Error())
					log.Println("Запрос НЕ смог отработать. Проверь корректность всех данных в запросе")
					//log.Println("Будет предпринята новая попытка через 1 минут")
					//time.Sleep(60 * time.Second)
					myError = 0 //если такой таблицы нет в БД, то что она появится через 5 минут?
					err = errQuery
				}
			} else {
				log.Println("db.Ping failed:", errDBping)
				log.Println("Подключение к БД НЕ установлено. Проверь доступность БД")
				log.Println("Будет предпринята новая попытка через 1 минут")
				time.Sleep(60 * time.Second)
				myError++
				err = errDBping
			}
		} else {
			//log.Print(errSqlOpen.Error())
			log.Println("Error creating DB:", errSqlOpen)
			log.Println("To verify, db is:", db)
			log.Println("Создание подключения к БД завершилось ошибкой. Часто возникает из-за не корректного драйвера")
			log.Println("Будет предпринята новая попытка через 1 минут")
			time.Sleep(60 * time.Second)
			myError++
			err = errSqlOpen
		}
		//if myError == 5 {
		if myError == sm.connAttempts {
			myError = 0
			return err
		}
	} //sql.Open
	return nil
}

// Close -.
func (p *SqlMy) Close() {
	//if p.Pool != nil {		p.Pool.Close()	}

}
