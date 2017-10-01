package redisWrapper_test

import (
	"testing"
	"time"

	"github.com/yunyuncc/tools/db/redisWrapper"
)

var (
	testAddr = "120.77.64.153:6379"
)

func TestInitDefaultClient(t *testing.T) {
	redisWrapper.InitDefaultClient(testAddr, 200)
}
func TestPING(t *testing.T) {
	cli := redisWrapper.GetDefaultClient()
	res, err := cli.PING()
	if res != "PONG" {
		t.Fatalf("PING fail:%v\n", err)
	}
}

func TestNewClient(t *testing.T) {
	cli := redisWrapper.NewClient(testAddr, 200)
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

func TestSTRLEN(t *testing.T) {
	cli := redisWrapper.GetDefaultClient()
	_ = cli.SET("testStrLen", "val")
	len, _ := cli.STRLEN("testStrLen")
	if len != 3 {
		t.Fatalf("expire 3 get %v\n", len)
	}
	err := cli.DEL("testStrLen")
	if err != nil {
		t.Fatal("del error:", err)
	}
}

func TestINCRDECR(t *testing.T) {
	cli := redisWrapper.GetDefaultClient()
	err := cli.DEL("testCount")
	if err != nil {
		t.Fatal("del fail:", err)
	}
	val, _ := cli.INCR("testCount")
	if val != 1 {
		t.Fatalf("val should equal 1 got:%v\n", val)
	}
	val, _ = cli.INCR("testCount")
	if val != 2 {
		t.Fatalf("val should equal 2 got:%v\n", val)
	}
	val, _ = cli.DECR("testCount")
	if val != 1 {
		t.Fatalf("2 - 1 should be 1\n")
	}
	val, _ = cli.INCRBY("testCount", 10)
	if val != 11 {
		t.Fatalf("1 + 10 should be 11\n")
	}
	val, _ = cli.DECRBY("testCount", 5)
	if val != 6 {
		t.Fatalf("11 - 5 should be 6")
	}
	fval, err := cli.INCRBYFLOAT("testCount", 5.5)
	if fval != 11.5 {
		t.Fatalf("6 + 5.5 should be 11.5 got %v,  err :%v", fval, err)
	}
	err = cli.DEL("testCount")
	if err != nil {
		t.Fatal("del fail:", err)
	}
}
