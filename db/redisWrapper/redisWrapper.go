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
}

//STRLEN [Strings Group] 获取key对应的val的长度
func (cli *Client) STRLEN(key string) (len int64, err error) {
	conn := cli.Pool.Get()
	defer conn.Close()
	len, err = redis.Int64(conn.Do("STRLEN", key))
	return
}

//INCR [Strings Group] 如果key的val可以被解析成int或者float，就会增加1
func (cli *Client) INCR(key string) (val int64, err error) {
	conn := cli.Pool.Get()
	defer conn.Close()
	val, err = redis.Int64(conn.Do("INCR", key))
	return
}

//DECR [Strings Group] 如果key的val可以被解析成int或者float，就会减去1
func (cli *Client) DECR(key string) (val int64, err error) {
	conn := cli.Pool.Get()
	defer conn.Close()
	val, err = redis.Int64(conn.Do("DECR", key))
	return
}

//INCRBY [Strings Group] key的val增加num
func (cli *Client) INCRBY(key string, num int64) (val int64, err error) {
	conn := cli.Pool.Get()
	defer conn.Close()
	val, err = redis.Int64(conn.Do("INCRBY", key, num))
	return
}

//DECRBY [Strings Group] key的val减少num
func (cli *Client) DECRBY(key string, num int64) (val int64, err error) {
	conn := cli.Pool.Get()
	defer conn.Close()
	val, err = redis.Int64(conn.Do("DECRBY", key, num))
	return
}

//INCRBYFLOAT [Strings Group] key的val增加num
func (cli *Client) INCRBYFLOAT(key string, num float64) (val float64, err error) {
	conn := cli.Pool.Get()
	defer conn.Close()
	val, err = redis.Float64(conn.Do("INCRBYFLOAT", key, num))
	return
}

//APPEND [Strings Group] append val return len of new val
func (cli *Client) APPEND(key, val string) (len int64, err error) {
	conn := cli.Pool.Get()
	defer conn.Close()
	return redis.Int64(conn.Do("APPEND", key, val))
}

//RPUSH [Lists Group] 从右边放入一个值
func (cli *Client) RPUSH(key string, val string) (len int64, err error) {
	conn := cli.Pool.Get()
	defer conn.Close()
	return redis.Int64(conn.Do("RPUSH", key, val))
}

//LPUSH [Lists Group] 从左边放入一个值
func (cli *Client) LPUSH(key string, val string) (len int64, err error) {
	conn := cli.Pool.Get()
	defer conn.Close()
	return redis.Int64(conn.Do("LPUSH", key, val))
}

//RPOP [Lists Group]从右边取出一个值
func (cli *Client) RPOP(key string) (val string, err error) {
	conn := cli.Pool.Get()
	defer conn.Close()
	return redis.String(conn.Do("RPOP", key))
}

//LPOP [Lists Group] 从左边取出一个值
func (cli *Client) LPOP(key string) (val string, err error) {
	conn := cli.Pool.Get()
	defer conn.Close()
	b, err := redis.Bytes(conn.Do("LPOP", key))
	return string(b), err
}

//LINDEX [Lists Group] 按索引取值
func (cli *Client) LINDEX(key string, index int64) (val string, err error) {
	conn := cli.Pool.Get()
	defer conn.Close()
	return redis.String(conn.Do("LINDEX", key, index))
}

//LRANGE [Lists Group] 读取start,end区间的值
func (cli *Client) LRANGE(key string, start, end int64) (vals []string, err error) {
	conn := cli.Pool.Get()
	defer conn.Close()
	return redis.Strings(conn.Do("LRANGE", key, start, end))
}

//LTRIM [Lists Group] 切除start,end之外的元素
func (cli *Client) LTRIM(key string, start, end int64) (err error) {
	conn := cli.Pool.Get()
	defer conn.Close()
	_, err = conn.Do("LTRIM", key, start, end)
	return err
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
