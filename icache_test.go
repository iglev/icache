package icache

import (
	"context"
	"fmt"
	"testing"
	"time"

	json "github.com/json-iterator/go"
)

var (
	disableTestObjLRU       = false
	disableTestCacheByteLRU = false
	disableTestBenchmark1   = false
	disableTestBenchmark2   = false
	disableTestBenchmark3   = false
)

type nodeObj struct {
	Num int
	Key string
	Vec []int
}

func getter(ctx context.Context, strKey string, dest SinkIf) error {
	t := ctx.Value("testing").(*testing.T)
	t.Logf("key=%s\n", strKey)
	switch strKey {
	case "byteKey":
		dest.SetBytes([]byte("byte val"))
		dest.SetTTL(2)
	case "byteSink_SetString":
		dest.SetString(string("byteSink set string"))
		dest.SetTTL(2)
	case "stringKey":
		dest.SetString("string val")
		dest.SetTTL(2)
	case "stringSink_SetBytes":
		dest.SetBytes([]byte("stringSink set bytes"))
		dest.SetTTL(2)
	case "jsonKey":
		obj := &nodeObj{
			Num: 10,
			Key: "key",
			Vec: []int{1, 2, 3},
		}
		b, _ := json.Marshal(obj)
		dest.SetString(string(b))
		dest.SetTTL(2)
	case "objKey":
		obj := &nodeObj{
			Num: 10,
			Key: "key",
			Vec: []int{1, 2, 3},
		}
		dest.SetObj(obj)
		dest.SetTTL(2)
	default:
		return fmt.Errorf("test error")
	}
	return nil
}

func TestObjLRU(t *testing.T) {
	t.Logf("TestObjLRU begin-----------------")
	defer t.Logf("TestObjLRU end-------------------")

	if disableTestObjLRU {
		return
	}

	objLRU := NewLRUObjCache(10)
	ic, err := NewICache(
		SetCache(objLRU),
		SetGetter(GetterIfFunc(getter)),
		SetRateLimit(10),
	)
	if err != nil {
		t.Logf("NewICache fail, err=%+v\n", err)
		return
	}

	ctx := context.WithValue(context.Background(), "testing", t)
	var stringValue string
	err = ic.Get(ctx, "stringKey", StringSink(&stringValue))
	t.Logf("stringKey, val=%s err=%+v\n", stringValue, err)

	var stringByteValue string
	err = ic.Get(ctx, "stringSink_SetBytes", StringSink(&stringByteValue))
	t.Logf("stringSink_SetBytes, val=%s err=%+v\n", stringByteValue, err)

	var byteValue []byte
	err = ic.Get(ctx, "byteKey", ByteSink(&byteValue))
	t.Logf("byteValue, val=%s err=%+v\n", byteValue, err)

	var byteStringValue []byte
	err = ic.Get(ctx, "byteSink_SetString", ByteSink(&byteStringValue))
	t.Logf("byteSink_SetString, val=%s err=%+v\n", byteStringValue, err)

	var objValue nodeObj
	err = ic.Get(ctx, "objKey", ObjSink(&objValue))
	t.Logf("objValue, val=%+v err=%+v\n", objValue, err)

	objValue.Key = "new Key"               // warning!!! change cache obj
	objValue.Vec = append(objValue.Vec, 4) // warning!!! change cache obj

	var objValue2 nodeObj
	err = ic.Get(ctx, "objKey", ObjSink(&objValue2))
	t.Logf("objValue2, val=%+v err=%+v\n", objValue2, err)

	t.Logf("for begin-----------")
	for i := 0; i < 10; i++ {
		var stringValue string
		err = ic.Get(ctx, "stringKey", StringSink(&stringValue))
		t.Logf("idx=%d stringKey, val=%s err=%+v\n", i, stringValue, err)
		time.Sleep(1 * time.Second)
	}
	t.Logf("for end-------------")

	/*
	   icache_test.go:66: TestObjLRU begin-----------------
	   icache_test.go:28: key=stringKey
	   icache_test.go:87: stringKey, val=string val err=<nil>
	   icache_test.go:28: key=stringSink_SetBytes
	   icache_test.go:91: stringSink_SetBytes, val=stringSink set bytes err=<nil>
	   icache_test.go:28: key=byteKey
	   icache_test.go:95: byteValue, val=byte val err=<nil>
	   icache_test.go:28: key=byteSink_SetString
	   icache_test.go:99: byteSink_SetString, val=byteSink set string err=<nil>
	   icache_test.go:28: key=objKey
	   icache_test.go:103: objValue, val={Num:10 Key:key Vec:[1 2 3]} err=<nil>
	   icache_test.go:110: objValue2, val={Num:10 Key:new Key Vec:[1 2 3 4]} err=<nil>
	   icache_test.go:112: for begin-----------
	   icache_test.go:116: idx=0 stringKey, val=string val err=<nil>
	   icache_test.go:116: idx=1 stringKey, val=string val err=<nil>
	   icache_test.go:116: idx=2 stringKey, val=string val err=<nil>
	   icache_test.go:28: key=stringKey
	   icache_test.go:116: idx=3 stringKey, val=string val err=<nil>
	   icache_test.go:116: idx=4 stringKey, val=string val err=<nil>
	   icache_test.go:116: idx=5 stringKey, val=string val err=<nil>
	   icache_test.go:28: key=stringKey
	   icache_test.go:116: idx=6 stringKey, val=string val err=<nil>
	   icache_test.go:116: idx=7 stringKey, val=string val err=<nil>
	   icache_test.go:116: idx=8 stringKey, val=string val err=<nil>
	   icache_test.go:28: key=stringKey
	   icache_test.go:116: idx=9 stringKey, val=string val err=<nil>
	   icache_test.go:119: for end-------------
	   icache_test.go:151: TestObjLRU end-------------------
	*/
}

func TestCacheByteLRU(t *testing.T) {
	t.Logf("TestCacheByteLRU begin-----------------")
	defer t.Logf("TestCacheByteLRU end-------------------")
	if disableTestCacheByteLRU {
		return
	}

	byteLRU := NewLRUByteCache(10)
	ic, err := NewICache(
		SetCache(byteLRU),
		SetGetter(GetterIfFunc(getter)),
	)
	if err != nil {
		t.Logf("NewICache fail, err=%+v\n", err)
		return
	}

	ctx := context.WithValue(context.Background(), "testing", t)
	var stringValue string
	err = ic.Get(ctx, "stringKey", StringSink(&stringValue))
	t.Logf("stringKey, val=%s err=%+v\n", stringValue, err)

	var stringByteValue string
	err = ic.Get(ctx, "stringSink_SetBytes", StringSink(&stringByteValue))
	t.Logf("stringSink_SetBytes, val=%s err=%+v\n", stringByteValue, err)

	var byteValue []byte
	err = ic.Get(ctx, "byteKey", ByteSink(&byteValue))
	t.Logf("byteValue, val=%s err=%+v\n", byteValue, err)

	var byteStringValue []byte
	err = ic.Get(ctx, "byteSink_SetString", ByteSink(&byteStringValue))
	t.Logf("byteSink_SetString, val=%s err=%+v\n", byteStringValue, err)

	var objValue nodeObj
	err = ic.Get(ctx, "objKey", ObjSink(&objValue))
	t.Logf("objValue, val=%+v err=%+v\n", objValue, err)

	objValue.Key = "new Key"
	objValue.Vec = append(objValue.Vec, 4)

	var objValue2 nodeObj
	err = ic.Get(ctx, "objKey", ObjSink(&objValue2)) // 再次回源，因为LRUByteCache缓存器支持[]byte和string类型
	t.Logf("objValue2, val=%+v err=%+v\n", objValue2, err)

	t.Logf("for begin-----------")
	for i := 0; i < 10; i++ {
		var stringValue string
		err = ic.Get(ctx, "stringKey", StringSink(&stringValue))
		t.Logf("idx=%d stringKey, val=%s err=%+v\n", i, stringValue, err)
		time.Sleep(1 * time.Second)
	}
	t.Logf("for end-------------")
	/*
	   icache_test.go:154: TestCacheByteLRU begin-----------------
	   icache_test.go:28: key=stringKey
	   icache_test.go:173: stringKey, val=string val err=<nil>
	   icache_test.go:28: key=stringSink_SetBytes
	   icache_test.go:177: stringSink_SetBytes, val=stringSink set bytes err=<nil>
	   icache_test.go:28: key=byteKey
	   icache_test.go:181: byteValue, val=byte val err=<nil>
	   icache_test.go:28: key=byteSink_SetString
	   icache_test.go:185: byteSink_SetString, val=byteSink set string err=<nil>
	   icache_test.go:28: key=objKey
	   icache_test.go:189: objValue, val={Num:10 Key:key Vec:[1 2 3]} err=<nil>
	   icache_test.go:28: key=objKey
	   icache_test.go:196: objValue2, val={Num:10 Key:key Vec:[1 2 3]} err=<nil>
	   icache_test.go:198: for begin-----------
	   icache_test.go:202: idx=0 stringKey, val=string val err=<nil>
	   icache_test.go:202: idx=1 stringKey, val=string val err=<nil>
	   icache_test.go:202: idx=2 stringKey, val=string val err=<nil>
	   icache_test.go:28: key=stringKey
	   icache_test.go:202: idx=3 stringKey, val=string val err=<nil>
	   icache_test.go:202: idx=4 stringKey, val=string val err=<nil>
	   icache_test.go:202: idx=5 stringKey, val=string val err=<nil>
	   icache_test.go:28: key=stringKey
	   icache_test.go:202: idx=6 stringKey, val=string val err=<nil>
	   icache_test.go:202: idx=7 stringKey, val=string val err=<nil>
	   icache_test.go:202: idx=8 stringKey, val=string val err=<nil>
	   icache_test.go:28: key=stringKey
	   icache_test.go:202: idx=9 stringKey, val=string val err=<nil>
	   icache_test.go:205: for end-------------
	   icache_test.go:237: TestCacheByteLRU end-------------------
	*/
}

var (
	benchmarkSize = 20000000
)

func TestBenchmark1(t *testing.T) {
	t.Logf("TestBenchmark1 begin-----------------")
	defer t.Logf("TestBenchmark1 end-------------------")

	if disableTestBenchmark1 {
		return
	}

	objLRU := NewLRUObjCache(10)
	ic, err := NewICache(
		SetCache(objLRU),
		SetGetter(GetterIfFunc(getter)),
	)
	if err != nil {
		t.Logf("NewICache fail, err=%+v\n", err)
		return
	}

	ctx := context.WithValue(context.Background(), "testing", t)
	startTime := time.Now()
	for i := 0; i < benchmarkSize; i++ {
		var stringValue string
		err = ic.Get(ctx, "stringKey", StringSink(&stringValue))
		if err != nil {
			t.Logf("err=%+v\n", err)
		}
		// t.Logf("idx=%d stringKey, val=%s err=%+v\n", i, stringValue, err)
	}
	t.Logf("TestBenchmark1 string stats=%+v cost=%+v\n", ic.GetStat(), time.Since(startTime))
}

func TestBenchmark2(t *testing.T) {
	t.Logf("TestBenchmark2 begin-----------------")
	defer t.Logf("TestBenchmark2 end-------------------")

	if disableTestBenchmark2 {
		return
	}

	objLRU := NewLRUObjCache(10)
	ic, err := NewICache(
		SetCache(objLRU),
		SetGetter(GetterIfFunc(getter)),
	)
	if err != nil {
		t.Logf("NewICache fail, err=%+v\n", err)
		return
	}

	ctx := context.WithValue(context.Background(), "testing", t)
	startTime := time.Now()
	for i := 0; i < benchmarkSize; i++ {
		// var objValue nodeObj
		// err = ic.Get(ctx, "objKey", ObjSink(&objValue))
		var stringValue string
		err = ic.Get(ctx, "jsonKey", StringSink(&stringValue))
		if err != nil {
			t.Logf("err=%+v\n", err)
		}
		var obj nodeObj
		json.Unmarshal([]byte(stringValue), obj)
		// t.Logf("idx=%d stringKey, val=%s err=%+v\n", i, stringValue, err)
	}
	t.Logf("TestBenchmark2 obj stats=%+v cost=%+v\n", ic.GetStat(), time.Since(startTime))
}

func TestBenchmark3(t *testing.T) {
	t.Logf("TestBenchmark3 begin-----------------")
	defer t.Logf("TestBenchmark3 end-------------------")

	if disableTestBenchmark3 {
		return
	}

	objLRU := NewLRUObjCache(10)
	ic, err := NewICache(
		SetCache(objLRU),
		SetGetter(GetterIfFunc(getter)),
	)
	if err != nil {
		t.Logf("NewICache fail, err=%+v\n", err)
		return
	}

	ctx := context.WithValue(context.Background(), "testing", t)
	startTime := time.Now()
	for i := 0; i < benchmarkSize; i++ {
		var objValue nodeObj
		err = ic.Get(ctx, "objKey", ObjSink(&objValue))
		if err != nil {
			t.Logf("err=%+v\n", err)
		}
		// t.Logf("idx=%d stringKey, val=%s err=%+v\n", i, stringValue, err)
	}
	t.Logf("TestBenchmark3 obj stats=%+v cost=%+v\n", ic.GetStat(), time.Since(startTime))
}
