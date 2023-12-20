package repo

//https://tutorialedge.net/golang/golang-mysql-tutorial/
//https://github.com/evrone/go-clean-template/blob/master/internal/usecase/repo/translation_postgres.go
import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql" //для установки драйвера mysql. Сам пакет как бы не используется явно в коде
	"log"
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
func NewUnifiRepo(connectStr string, base string) (*UnifiRepo, error) {
	//log.Println(connectStr + "/" + db)
	log.Println(connectStr + "/" + base)

	pr := &UnifiRepo{
		dataSourceITsup: connectStr + "/" + base, //i,
		databaseITsup:   base,                    //strings.Split(i, "/")[1],
		dataSourceGLPI:  connectStr + "/glpi_db", // g,
		databaseGLPI:    "glpi_db",               //strings.Split(g, "/")[1],
		controller:      0,                       //c,
	}

	if db, errSqlOpen := sql.Open("mysql", pr.dataSourceITsup); errSqlOpen == nil {
		errDBping := db.Ping()
		if errDBping == nil {
			return pr, nil
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
	//return pr, nil
}

func (ur *UnifiRepo) ChangeCntrlNumber(newCntrlNumber int) {
	ur.controller = newCntrlNumber
}

func (ur *UnifiRepo) dbExec(query string) (err error) { //UploadMapsToDBerr

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
		if myError == 5 {
			myError = 0
			return err
		}
	} //sql.Open
	return nil
}
