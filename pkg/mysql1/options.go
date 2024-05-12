package mysql1

import "time"

// Option -.
type Option func(*SqlMy)

func Base(db string) Option {
	return func(c *SqlMy) {
		c.dataBase = db
	}
}

// MaxPoolSize -.
func MaxPoolSize(size int) Option {
	return func(c *SqlMy) {
		c.maxPoolSize = size
	}
}

// ConnAttempts -.
func ConnAttempts(attempts int) Option {
	return func(c *SqlMy) {
		c.connAttempts = attempts
	}
}

// ConnTimeout -.
func ConnTimeout(timeout time.Duration) Option {
	return func(c *SqlMy) {
		c.connTimeout = timeout
	}
}
