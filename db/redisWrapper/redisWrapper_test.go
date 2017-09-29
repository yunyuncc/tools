package redisWrapper_test

import (
	"testing"
	"time"

	"github.com/yunyuncc/tools/db/redisWrapper"
)

func TestNewPool(t *testing.T) {
	pool := redisWrapper.NewPool("127.0.0.1:6379", 200)
	if pool == nil {
		t.Fatal("NewPool失败")
	}
}

func TestInitDefaultPool(t *testing.T) {
	redisWrapper.InitDefaultPool("127.0.0.1:6379", 200)
	if redisWrapper.DefaultPool == nil {
		t.Fatalf("初始化连接池失败\n")
	}
}
func TestPING(t *testing.T) {
	res, err := redisWrapper.PING()
	if res != "PONG" {
		t.Fatalf("PING fail:%v\n", err)
	}
}
func TestSETGETDEL(t *testing.T) {
	err := redisWrapper.SET("testKey", "testVal")
	if err != nil {
		t.Fatalf("SET失败\n")
	}
	val, _ := redisWrapper.GET("testKey")
	if val != "testVal" {
		t.Fatalf("应取出testVal 实际取出:%s\n", val)
	}
	err = redisWrapper.SET("testKey2", "testVal2", 1)
	if err != nil {
		t.Fatalf("SET失败 error:%v\n", err)
	}
	val, _ = redisWrapper.GET("testKey2")
	if val != "testVal2" {
		t.Fatalf("应取出testVal2 实际取出:%s\n", val)
	}
	time.Sleep(2 * time.Second)
	val, err = redisWrapper.GET("testKey2")
	if err != nil {
		t.Fatalf("GET失败\n")
	}
	if val != "nil" {
		t.Fatalf("失效失败\n")
	}

	err = redisWrapper.DEL("testKey")
	if err != nil {
		t.Fatalf("DEL失败%v\n", err)
	}
	val, _ = redisWrapper.GET("testKey")
	if val != "nil" {
		t.Fatalf("删除后应该取出nil\n")
	}
}
