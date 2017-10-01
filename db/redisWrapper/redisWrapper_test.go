package redisWrapper_test

import (
	"fmt"
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
func TestAPPEND(t *testing.T) {
	cli := redisWrapper.GetDefaultClient()
	err := cli.DEL("testKey")
	if err != nil {
		t.Fatal("del testKey fail:", err)
	}
	_ = cli.SET("testKey", "val1 ")
	len, _ := cli.APPEND("testKey", "val2")
	if len != 9 {
		t.Fatal("len should be 9")
	}
	val, _ := cli.GET("testKey")
	if val != "val1 val2" {
		t.Fatal("val should be val1 val2")
	}
	err = cli.DEL("testKey")
	if err != nil {
		t.Fatal("del testKey fail:", err)
	}
}

func TestLIST(t *testing.T) {
	cli := redisWrapper.GetDefaultClient()
	err := cli.DEL("testList")
	if err != nil {
		t.Fatal("del testKey fail:", err)
	}
	for i := int64(1); i <= 20; i++ {
		len, err := cli.RPUSH("testList", fmt.Sprintf("val%v", i))
		if len != i {
			t.Fatalf("len should be %v got:%v error:%v\n", i, len, err)
		}
	}
	for i := int64(20); i >= 1; i-- {
		val, _ := cli.RPOP("testList")
		if val != fmt.Sprintf("val%v", i) {
			t.Logf("should be val%v got %v", i, val)
		}
	}

	for i := int64(1); i <= 20; i++ {
		len, err := cli.LPUSH("testList", fmt.Sprintf("val%v", i))
		if len != i {
			t.Fatalf("len should be %v got:%v error:%v\n", i, len, err)
		}
	}
	for i := int64(20); i >= 1; i-- {
		val, _ := cli.LPOP("testList")
		if val != fmt.Sprintf("val%v", i) {
			t.Logf("should be val%v got %v", i, val)
		}
	}

	for i := int64(1); i <= 20; i++ {
		len, err := cli.RPUSH("testList", fmt.Sprintf("val%v", i))
		if len != i {
			t.Fatalf("len should be %v got:%v error:%v\n", i, len, err)
		}
	}
	for i := int64(1); i <= 20; i++ {
		val, _ := cli.LPOP("testList")
		if val != fmt.Sprintf("val%v", i) {
			t.Logf("should be val%v got %v", i, val)
		}
	}
	cli.RPUSH("testList", "1")
	cli.RPUSH("testList", "2")
	cli.RPUSH("testList", "3")
	cli.RPUSH("testList", "4")
	val, err := cli.LINDEX("testList", 3)
	if val != "4" {
		t.Fatalf("should be 4 got %s, err :%v", val, err)
	}
	vals, err := cli.LRANGE("testList", 0, 30)
	for i, v := range vals {
		if v != fmt.Sprintf("%v", i+1) {
			t.Fatalf("should be %v got %s", i, v)
		}
	}
	err = cli.LTRIM("testList", 0, 3)
	if err != nil {
		t.Fatal("LTRIM error:%v\n", err)
	}
	_, err = cli.LINDEX("testList", 4)
	if err == nil {
		t.Fatalf("after trim index 4 should be nil\n")
	}
	t.Log("trim success get index 4 get err:", err)
	val, _ = cli.LINDEX("testList", 3)
	if val != "4" {
		t.Fatalf("index 3 val should be 4 got %v\n", val)
	}
	err = cli.DEL("testList")
	if err != nil {
		t.Fatal("Del testList error:", err)
	}

}
