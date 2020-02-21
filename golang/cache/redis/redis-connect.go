package redis

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

type Connect struct {
	Pool *redis.Pool
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

func InitRedisConnect(address, db, passwd string, timeout, maxActive, maxIdle, idleTimeout int) *Connect {
	return &Connect{
		Pool: &redis.Pool{
			Dial:        dialFunc(timeout, address, db, passwd),
			MaxIdle:     maxIdle,
			MaxActive:   maxActive,
			IdleTimeout: time.Duration(idleTimeout) * time.Second,
		},
	}
}
