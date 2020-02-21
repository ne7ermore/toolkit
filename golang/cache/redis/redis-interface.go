package redis

import "github.com/gomodule/redigo/redis"

type Cache interface {
	Encode() string
	Decode(reply string) error
	Key() string
	RedisConn() redis.Conn
	Exp() int64
}

func SetCache(c Cache) error {
	conn := c.RedisConn()
	if conn.Err() != nil {
		return conn.Err()
	}
	defer conn.Close()

	_, err := redis.String(conn.Do("setex", c.Key(), c.Exp(), c.Encode()))
	return err
}

func GetCache(c Cache) error {
	conn := c.RedisConn()
	if conn.Err() != nil {
		return conn.Err()
	}
	defer conn.Close()

	reply, err := redis.String(conn.Do("get", c.Key()))
	if err != nil || reply == "" {
		return err
	}

	return c.Decode(reply)
}
