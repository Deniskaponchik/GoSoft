package repo

//https://tutorialedge.net/golang/golang-mysql-tutorial/
//https://github.com/evrone/go-clean-template/blob/master/internal/usecase/repo/translation_postgres.go
import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql" //для установки драйвера mysql. Сам пакет как бы не используется явно в коде
	"log"
)

type GisupRepo struct {
	dataSourceITsup string
	databaseITsup   string
	dataSourceGLPI  string
	databaseGLPI    string
	//controller      int
}

// реализуем Инъекцию зависимостей DI. Используется в app
func NewGisupRepo(connectStr string, base string) (*GisupRepo, error) {
	//log.Println(connectStr + "/" + db)
	log.Println(connectStr + "/" + base)

	gr := &GisupRepo{
		dataSourceITsup: connectStr + "/" + base, //i,
		databaseITsup:   base,                    //strings.Split(i, "/")[1],
		dataSourceGLPI:  connectStr + "/glpi_db", // g,
		databaseGLPI:    "glpi_db",               //strings.Split(g, "/")[1],
		//controller:      0,                       //c,
	}

	if db, errSqlOpen := sql.Open("mysql", gr.dataSourceITsup); errSqlOpen == nil {
		errDBping := db.Ping()
		if errDBping == nil {
			return gr, nil
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
