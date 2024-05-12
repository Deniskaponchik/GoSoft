// Package postgres implements postgres connection.
package redis1

import (
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

const (
	_defaultMaxPoolSize  = 1
	_defaultConnAttempts = 5
	_defaultConnTimeout  = time.Second
)

type Redis struct {
	maxPoolSize  int
	connAttempts int
	connTimeout  time.Duration

	client *redis.Client

	//MySql
	//dataSource string
	//dataBase   string

	//PG
	//Builder squirrel.StatementBuilderType
	//Pool    *pgxpool.Pool
}

// New -.
func NewRedis(connectStr string, base string, opts ...Option) (*Redis, error) {
	rd := &Redis{
		//dataSource: connectStr + "/" + base,
		//dataBase:   base,

		maxPoolSize:  _defaultMaxPoolSize,
		connAttempts: _defaultConnAttempts,
		connTimeout:  _defaultConnTimeout,
	}

	// Custom options
	for _, opt := range opts {
		opt(rd)
	}

	log.Println(connectStr + "/" + base)

	//"redis://<user>:<pass>@localhost:6379/<db>"
	opt, err := redis.ParseURL(connectStr)
	if err != nil {
		//log.Fatal(err)
		log.Println(err)
		log.Println("Подключение к Redis НЕ установлено. Проверь доступность БД")
		return nil, err
	} else {
		rd.client = redis.NewClient(opt)
		return rd, nil
	}

}

// Close -.
func (r *Redis) Close() {
	//if p.Pool != nil {		p.Pool.Close()	}

}
