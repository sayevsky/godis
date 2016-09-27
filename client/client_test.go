package client

import (
	"github.com/sayevsky/godis/server"
	"log"
	"testing"
)

func TestGet(t *testing.T) {
	server.NewServer().Start(true)
	client, _ := NewClient("localhost:6380")

	res, err := client.Get("a")
	if err.Error() != "NE" {
		t.Errorf("fail to get a key", res, err)
	}
	log.Println(res)

}
