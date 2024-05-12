// Package postgres implements postgres connection.
package mysql1

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/driver/sqliteshim"
	"log"
	"time"
)

const (
	_defaultMaxPoolSize  = 1
	_defaultConnAttempts = 5
	_defaultConnTimeout  = time.Second
)

// Bun
type MySql2 struct {
	maxPoolSize  int
	connAttempts int
	connTimeout  time.Duration

	bun *bun.

	//Builder squirrel.StatementBuilderType
	//Pool    *pgxpool.Pool
}

// New -.
func NewSqlMy(connectStr string, base string, opts ...Option) (*SqlMy, error) {
	bun := &MySql2{
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

	sqldb, err := sql.Open("mysql", "root:pass@/test")
	if err != nil {
		panic(err)
	}

	db := bun.NewDB(sqldb, mysqldialect.New())

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



// Close -.
func (p *SqlMy) Close() {
	//if p.Pool != nil {		p.Pool.Close()	}

}
