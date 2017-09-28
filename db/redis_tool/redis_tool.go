package redis_tool

import (
	"fmt"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
)

//NewRedis 新建一个redis连接池
func NewRedis(password, host string, port, db, maxIdle, maxOpen int) *redis.Pool {
	connStr := fmt.Sprintf("%s:%s", host, strconv.Itoa(port))
	redis := &redis.Pool{
		MaxIdle:     maxIdle,
		MaxActive:   maxOpen,
		Wait:        true,
		IdleTimeout: 180 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", connStr)
			if err != nil {
				return nil, err
			}
			if password != "" {
				_, err = c.Do("AUTH", password)
				if err != nil {
					c.Close()
					return nil, err
				}
			}
			c.Do("SELECT", db)
			return c, nil
		},
	}
	return redis
}
