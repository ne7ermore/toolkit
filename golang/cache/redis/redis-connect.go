package redis

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

type redisConfig struct {
	address, db, password                    string
	wait                                     bool
	timeout, maxActive, maxIdle, idleTimeout int
}

func dialFunc(timeout int, address, db, passwd string) func() (redis.Conn, error) {
	return func() (redis.Conn, error) {
		var (
			cto = time.Duration(timeout) * time.Second
			rto = time.Duration(timeout) * time.Second
			wto = time.Duration(timeout) * time.Second
		)

		c, err := redis.DialTimeout("tcp", address, cto, rto, wto)
		if err != nil {
			return nil, err
		}
		if len(passwd) > 0 {
			if _, err := c.Do("AUTH", passwd); err != nil {
				c.Close()
				return nil, err
			}
		}
		if len(db) > 0 {
			if _, err = c.Do("SELECT", db); err != nil {
				c.Close()
				return nil, err
			}
		}
		return c, err
	}
}

func InitRedisConnect(address, db, passwd string,
	wait bool,
	timeout, maxActive, maxIdle, idleTimeout int) *redis.Pool {
	return &redis.Pool{
		Dial:        dialFunc(timeout, address, db, passwd),
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: time.Duration(idleTimeout) * time.Second}
}
