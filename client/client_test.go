package client

import (
	"github.com/sayevsky/godis/server"
	"reflect"
	"sort"
	"testing"
	"strconv"
	"log"
	"time"
)

func TestGetEmpty(t *testing.T) {
	s := server.NewServer()
	s.Start(true)
	client, _ := NewClient("localhost:6380")
	res := client.Get("a")
	if res.Err.Error() != "NE" {
		t.Error("fail to get a key", res)
	}
	s.Stop()
}

func TestUpdateGet(t *testing.T) {
	s := server.NewServer()
	s.Start(true)
	client, _ := NewClient("localhost:6380")
	res := client.Update("a", "b", 10 * time.Second)
	res = client.Get("a")
	if res.Err == nil {
		t.Error("managed to get a key though it should not exist", res)
	}
	s.Stop()
}

func TestSetUpdateGet(t *testing.T) {
	s := server.NewServer()
	s.Start(true)
	client, _ := NewClient("localhost:6380")
	res := client.Set("a>", "b", 10 * time.Second)
	res = client.Get("a>")
	if res.Err != nil || res.Result != "b" {
		t.Error("fail to get a key that was set", res)
	}
	res = client.Update("a>", "bb", 10 * time.Second)
	res = client.Get("a>")

	if res.Result != "bb" || res.Err != nil {
		t.Error("fail to get a key that was set then updated", res)
	}
	s.Stop()
}

func TestSetUpdateGetMap(t *testing.T) {
	s := server.NewServer()
	s.Start(true)
	client, _ := NewClient("localhost:6380")
	value := make(map[string]string)
	client.Set("key", value, 0)
	res := client.Get("key")
	if !reflect.DeepEqual(value, res.Result) || res.Err != nil {
		t.Error("fail to get a key that was set then updated", res)
	}
	value["a"] = "b"
	client.Update("key", value, 0)
	res = client.Get("key")
	if !reflect.DeepEqual(value, res.Result) || res.Err != nil {
		t.Error("fail to get a key that was set then updated", res)
	}
	s.Stop()
}

func TestSetUpdateGetArray(t *testing.T) {
	s := server.NewServer()
	s.Start(true)
	client, _ := NewClient("localhost:6380")
	value := make([]string, 1)
	client.Set("key", value, 0)
	res := client.Get("key")
	if !reflect.DeepEqual(value, res.Result) || res.Err != nil {
		t.Error("fail to get a key that was set then updated", res)
	}
	value[0] = "b"
	client.Update("key", value, 0)
	res = client.Get("key")
	if !reflect.DeepEqual(value, res.Result) || res.Err != nil {
		t.Error("fail to get a key that was set then updated", res)
	}
	s.Stop()
}

func TestDelete(t *testing.T) {
	s := server.NewServer()
	s.Start(true)
	client, _ := NewClient("localhost:6380")
	res := client.Delete("key")
	if res.Err == nil {
		t.Error("can't delete unxisted key with success")
	}
	client.Set("a", "b", 0)
	client.Delete("a")
	res = client.Get("a")
	if res.Err == nil {
		t.Error("can't get deleted key with success")
	}
	s.Stop()
}

func TestKeys(t *testing.T) {
	s := server.NewServer()
	s.Start(true)
	client, _ := NewClient("localhost:6380")
	client.Set("a", ">b", 0)
	client.Set("abc", ">abc", 0)
	client.Set("ab", ">bb", 0)

	res := client.Keys("ab+")
	expect := []string{"ab", "abc"}
	sort.Strings(res.Result.([]string))
	if res.Err != nil || !reflect.DeepEqual(expect, res.Result) {
		t.Error("keys was not filtered by regex")
	}

	s.Stop()
}

func TestCount(t *testing.T) {
	s := server.NewServer()
	s.Start(true)
	client, _ := NewClient("localhost:6380")
	client.Set("a", ">b", 0)
	client.Set("abc", ">abc", 0)
	client.Set("ab", ">bb", 0)

	res := client.Count()
	if res.Err != nil || res.Result != 3 {
		t.Error("error with count", res)
	}

	s.Stop()
}

func TestSetUpdateGGetInArray(t *testing.T) {
	s := server.NewServer()
	s.Start(true)
	client, _ := NewClient("localhost:6380")
	value := make([]string, 1)
	value[0] = "value1"
	client.Set("key", value, 0)
	res := client.GetIth("key", 0)
	if !reflect.DeepEqual(value[0], res.Result) || res.Err != nil {
		t.Error("fail to get a key that was set", res)
	}

	res = client.GetIth("key", 1)
	if res.Err == nil || res.Err.Error() != "OOR" {
		t.Error("OutOfRange was not fired", res)
	}

	value2 := "value"
	client.Set("key", value2, 0)
	res = client.GetIth("key", 0)
	if res.Err == nil || res.Err.Error() != "WT" {
		t.Error("WrongType was not fired", res)
	}

	s.Stop()
}

func TestSetUpdateGGetInMap(t *testing.T) {
	s := server.NewServer()
	s.Start(true)
	client, _ := NewClient("localhost:6380")
	value := make(map[string]string, 1)
	value1 := "value1"
	value["key1"] = value1
	client.Set("key", value, 0)
	res := client.GetKeyInValue("key", "key1")
	if !reflect.DeepEqual(value["key1"], res.Result) || res.Err != nil {
		t.Error("fail to get a key that was set", res)
	}
	res = client.GetKeyInValue("key", "key2")
	if res.Err == nil || res.Err.Error() != "NE" {
		t.Error("'Not exist' was not fired", res)
	}

	value2 := "value"
	client.Set("key", value2, 0)
	res = client.GetKeyInValue("key", "key1")
	if res.Err == nil || res.Err.Error() != "WT" {
		t.Error("WrongType was not fired", res)
	}

	s.Stop()
}

func BenchmarkBasic(b *testing.B) {
	log.Println("start bench")

	for n := 0; n < b.N; n++ {
		log.Println("start bench" , n)
		s:= server.NewServer()
		s.Start(true)
		client, _ := NewClient("localhost:6380")
		setValues(&client, 100000)
		s.Stop()
		log.Println("stop bench" , n)
	}
	log.Println("stop bench")
}

func setValues(client *Client, size int) {
	for i := 0; i < size; i++ {
		value := strconv.Itoa(i*2)
		client.Set(strconv.Itoa(i), value, 1 * time.Second)
		response := client.Get(strconv.Itoa(i))
		if response.Result != value {
			log.Print("fails in set get")
		}
	}

}
