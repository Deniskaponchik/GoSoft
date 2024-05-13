package mysql2

import "time"

// Option -.
type Option func(*MySql2)

func Base(db string) Option {
	return func(c *MySql2) {
		c.database = db
	}
}

// MaxPoolSize -.
func MaxPoolSize(size int) Option {
	return func(c *MySql2) {
		c.maxPoolSize = size
	}
}

// ConnAttempts -.
func ConnAttempts(attempts int) Option {
	return func(c *MySql2) {
		c.connAttempts = attempts
	}
}

// ConnTimeout -.
func ConnTimeout(timeout time.Duration) Option {
	return func(c *MySql2) {
		c.connTimeout = timeout
	}
}
