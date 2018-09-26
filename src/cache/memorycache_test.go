package cache

import (
	"testing"
	"time"

	"define"
	"errcode"
)

func TestMemoryCache(t *testing.T) {

	// init memory cache client
	memoryClient, err := NewMemoryCache()
	if err != nil {
		t.Fatal("expected err:nil, actual:", err)
	}
	defer memoryClient.Clear()

	{
		// test data
		cdn := define.CDNProvider("test01")
		data := define.AuthRes{}

		// insert new data
		err = memoryClient.Upsert(cdn, data, time.Second)
		if err != nil {
			t.Fatal("expected err:nil, actual:", err)
		}

		// get timeout data
		time.Sleep(time.Second * time.Duration(2))
		_, err = memoryClient.Get(cdn)
		if err != errcode.NoSuchAuthInfoCacheErr {
			t.Fatal(`expected err:{"code":404000,"error":"no such auth info cache"}, actual:`, err)
		}

		// delete timeout data
		err = memoryClient.Delete(cdn)
		if err != nil {
			t.Fatal("expected err:nil, actual:", err)
		}
	}

	{
		// test data
		cdn := define.CDNProvider("test02")
		data := define.AuthRes{Auth: []define.AuthInfo{
			define.AuthInfo{Name: "key", Value: "value", Location: define.LocationTypeHeader},
		}}

		// insert new data
		err = memoryClient.Upsert(cdn, data, time.Minute)
		if err != nil {
			t.Fatal("expected err:nil, actual:", err)
		}

		// get valid data
		_, err := memoryClient.Get(cdn)
		if err != nil {
			t.Fatal("expected err:nil, actual:", err)
		}

		// delete invalid data
		err = memoryClient.Delete(cdn)
		if err != nil {
			t.Fatal("expected err:nil, actual:", err)
		}
	}
}
