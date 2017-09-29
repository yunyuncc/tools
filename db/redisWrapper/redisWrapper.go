package redisWrapper

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
)

//DefaultPool 默认连接池
var DefaultPool *redis.Pool

//NewPool 初始化redis连接池
func NewPool(serverAddr string, maxIdle int) *redis.Pool {
	pool := redis.NewPool(func() (redis.Conn, error) {
		return redis.Dial("tcp", serverAddr)
	}, maxIdle)
	return pool
}

//InitDefaultPool 初始化默认连接池
func InitDefaultPool(serverAddr string, maxIdle int) {
	DefaultPool = NewPool(serverAddr, maxIdle)
}

//PING 检查连接是否正常
func PING() (string, error) {
	if DefaultPool == nil {
		return "", fmt.Errorf("default pool has not init")
	}
	conn := DefaultPool.Get()
	defer conn.Close()
	return redis.String(conn.Do("PING"))
}

//SET exp 是可选参数，为过期时间second为单位，过期后取出的值为nil
func SET(key string, val string, exp ...int) error {
	if DefaultPool == nil {
		return fmt.Errorf("default pool has not init")
	}
	conn := DefaultPool.Get()
	defer conn.Close()
	if len(exp) == 0 {
		_, err := conn.Do("SET", key, val)
		if err != nil {
			return err
		}
	} else if len(exp) == 1 {
		_, err := conn.Do("SET", key, val, "ex", exp[0])
		if err != nil {
			return err
		}
	}
	return nil
}

//GET 如果是空值val 为"nil"
func GET(key string) (val string, err error) {
	if DefaultPool == nil {
		return "", fmt.Errorf("default pool has not init")
	}
	conn := DefaultPool.Get()
	defer conn.Close()
	valInf, err := conn.Do("GET", key)
	if err != nil {
		return "", err
	}
	if valInf == nil {
		return "nil", nil
	}
	return redis.String(valInf, err)
}

//DEL 删除之后再去GET会取出nil
func DEL(key string) (err error) {
	if DefaultPool == nil {
		return fmt.Errorf("default pool has not init")
	}
	conn := DefaultPool.Get()
	defer conn.Close()
	_, err = conn.Do("DEL", key)
	return err
}
