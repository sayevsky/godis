package client

import (
	"github.com/sayevsky/godis/server"
	"testing"
)

func TestGetEmpty(t *testing.T) {
	s := server.NewServer()
	s.Start(true)
	client, _ := NewClient("localhost:6380")

	res := client.Get("a")
	if res.Err.Error() != "NE" {
		t.Errorf("fail to get a key", res)
	}
	s.Stop()
}

func TestSetGet(t *testing.T) {
	s := server.NewServer()
	s.Start(true)
	client, _ := NewClient("localhost:6380")

	res := client.Set("a", "b")

	res = client.Get("a")

	if res.Result != "b"{
		t.Errorf("fail to get a key that was set", res)
	}

	s.Stop()
}

