package redisWrapper_test

import (
	"testing"
	"time"

	"github.com/yunyuncc/tools/db/redisWrapper"
)

func TestInitDefaultClient(t *testing.T) {
	redisWrapper.InitDefaultClient("127.0.0.1:6379", 200)
}
func TestPING(t *testing.T) {
	cli := redisWrapper.GetDefaultClient()
	res, err := cli.PING()
	if res != "PONG" {
		t.Fatalf("PING fail:%v\n", err)
	}
}

func TestNewClient(t *testing.T) {
	cli := redisWrapper.NewClient("127.0.0.1:6379", 200)
	res, err := cli.PING()
	if res != "PONG" {
		t.Fatalf("new client PING fail:%v\n", err)
	}
}

func TestSETGETDEL(t *testing.T) {
	cli := redisWrapper.GetDefaultClient()
	err := cli.SET("testKey", "testVal")
	if err != nil {
		t.Fatalf("SET失败\n")
	}
	val, _ := cli.GET("testKey")
	if val != "testVal" {
		t.Fatalf("应取出testVal 实际取出:%s\n", val)
	}
	err = cli.SET("testKey2", "testVal2", 1)
	if err != nil {
		t.Fatalf("SET失败 error:%v\n", err)
	}
	val, _ = cli.GET("testKey2")
	if val != "testVal2" {
		t.Fatalf("应取出testVal2 实际取出:%s\n", val)
	}
	time.Sleep(2 * time.Second)
	val, err = cli.GET("testKey2")
	if err != nil {
		t.Logf("testKey2成功失效 val[%s] err[%v]", val, err)
	} else {
		t.Fatalf("testKey2失效失败")
	}

	err = cli.DEL("testKey")
	if err != nil {
		t.Fatalf("DEL失败%v\n", err)
	}
	val, _ = cli.GET("testKey")
	if val != "" {
		t.Fatalf("删除后应该取出nil\n")
	}
}
