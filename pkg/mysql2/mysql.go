// Package postgres implements postgres connection.
package mysql2

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"

	//"github.com/uptrace/bun/driver/sqliteshim"
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
	database     string

	bun *bun.DB

	//PG
	//Builder squirrel.StatementBuilderType
	//Pool    *pgxpool.Pool
}

// New -.
func NewSqlMy(connectStr string, base string, opts ...Option) (*MySql2, error) {
	mysql2 := &MySql2{

		maxPoolSize:  _defaultMaxPoolSize,
		connAttempts: _defaultConnAttempts,
		connTimeout:  _defaultConnTimeout,
	}

	// Custom options
	for _, opt := range opts {
		opt(mysql2)
	}

	dataSource := connectStr + "/" + base
	//log.Println(connectStr + "/" + base)
	log.Println(dataSource)

	//sqldb, err := sql.Open("mysql", "root:pass@/test")
	sqldb, err := sql.Open("mysql", dataSource)
	if err != nil {
		//panic(err)
		log.Println("Error creating Bun MySql DB:", err)
		return nil, err
	} else {
		db := bun.NewDB(sqldb, mysqldialect.New())
		mysql2.bun = db
		return mysql2, nil
	}
}

// Close -.
func (b *MySql2) Close() {
	//if p.Pool != nil {		p.Pool.Close()	}

}
