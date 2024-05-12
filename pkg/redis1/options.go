package redis1

import "time"

// Option -.
type Option func(*Redis)

func Base(db string) Option {
	return func(c *Redis) {
		c.dataBase = db
	}
}

// MaxPoolSize -.
func MaxPoolSize(size int) Option {
	return func(c *Redis) {
		c.maxPoolSize = size
	}
}

// ConnAttempts -.
func ConnAttempts(attempts int) Option {
	return func(c *Redis) {
		c.connAttempts = attempts
	}
}

// ConnTimeout -.
func ConnTimeout(timeout time.Duration) Option {
	return func(c *Redis) {
		c.connTimeout = timeout
	}
}
