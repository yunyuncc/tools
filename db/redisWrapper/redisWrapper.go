package redisWrapper

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
)

type Client struct {
	Pool *redis.Pool
}

var (
	defaultPool   *redis.Pool
	DefaultClient *Client
)

//NewPool 初始化redis连接池
func newPool(serverAddr string, maxIdle int) *redis.Pool {
	pool := redis.NewPool(func() (redis.Conn, error) {
		return redis.Dial("tcp", serverAddr)
	}, maxIdle)
	return pool
}

//GetDefaultClient 返回默认的客户端，如果此前没有调用InitDefaultClient会panic
func GetDefaultClient() *Client {
	if DefaultClient == nil {
		panic(fmt.Errorf("DefaultClient has not init"))
	}
	return DefaultClient
}

//InitDefaultClient 在使用默认客户端之前先初始化
func InitDefaultClient(serverAddr string, maxIdle int) {
	initDefaultPool(serverAddr, maxIdle)
	DefaultClient = &Client{
		Pool: defaultPool,
	}
}

//NewClient 新建一个redis客户端
func NewClient(serverAddr string, maxIdle int) *Client {
	cli := &Client{
		Pool: newPool(serverAddr, maxIdle),
	}
	return cli
}

//initDefaultPool 初始化默认连接池
func initDefaultPool(serverAddr string, maxIdle int) {
	defaultPool = newPool(serverAddr, maxIdle)
}

//SET [Strings Group] exp 是可选参数，为过期时间second为单位，过期后取出的值为nil
func (cli *Client) SET(key string, val string, exp ...int) error {
	conn := cli.Pool.Get()
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

//GET [Strings Group] 如果是空值val 为"nil"
func (cli *Client) GET(key string) (val string, err error) {
	conn := cli.Pool.Get()
	defer conn.Close()
	val, err = redis.String(conn.Do("GET", key))
	return
	/*
		if err != nil {
			return "", err
		}
		if valInf == nil {
			return "", nil
		}
		return redis.String(valInf, err)
	*/
}

//STRLEN [Strings Group] 获取key对应的val的长度
func (cli *Client) STRLEN(key string) (len int64, err error) {
	conn := cli.Pool.Get()
	defer conn.Close()
	len, err = redis.Int64(conn.Do("STRLEN", key))
	return
}

//DEL 删除之后再去GET会取出nil
func (cli *Client) DEL(key string) (err error) {
	conn := cli.Pool.Get()
	defer conn.Close()
	_, err = conn.Do("DEL", key)
	return err
}

//PING 检查连接是否正常
func (cli *Client) PING() (string, error) {
	conn := cli.Pool.Get()
	defer conn.Close()
	return redis.String(conn.Do("PING"))
}
