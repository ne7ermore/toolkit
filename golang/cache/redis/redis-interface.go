package redis

import "github.com/gomodule/redigo/redis"

type Cache interface {
	Encode() string
	Decode(reply string) error
	RedisConn() redis.Conn
	Exp() int64
}

func SetCache(c Cache, key string) error {
	conn := c.RedisConn()
	if conn.Err() != nil {
		return conn.Err()
	}
	defer conn.Close()

	_, err := redis.String(conn.Do("setex", key, c.Exp(), c.Encode()))
	return err
}

func GetCache(c Cache, key string) error {
	conn := c.RedisConn()
	if conn.Err() != nil {
		return conn.Err()
	}
	defer conn.Close()

	reply, err := redis.String(conn.Do("get", key))
	if err != nil || reply == "" {
		return err
	}

	return c.Decode(reply)
}
