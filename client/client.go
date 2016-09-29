package client

import (
	"bufio"
	"github.com/sayevsky/godis/internal"
	"net"
)

func NewClient(addr string) (Client, error) {
	conn, err := net.Dial("tcp", addr)
	return Client{conn}, err
}

type Client struct {
	conn net.Conn
}

func (c Client) Get(key string) (*internal.Response) {
	request, _ := internal.Get{key, internal.BaseCommand{false, nil}}.Serialize()
	c.conn.Write(request)
	response := bufio.NewReader(c.conn)
	return internal.DeserializeResponse(response)
}

func (c Client) Set(key string, value interface{}) (*internal.Response) {
	request, _ := internal.SetUpd{key, value, 0, false, internal.BaseCommand{false, nil}}.Serialize()
	c.conn.Write(request)
	response := bufio.NewReader(c.conn)
	return internal.DeserializeResponse(response)
}

func (c Client) Update(key string, value interface{}) (*internal.Response) {
	request, _ := internal.SetUpd{key, value, 0, true, internal.BaseCommand{false, nil}}.Serialize()
	c.conn.Write(request)
	response := bufio.NewReader(c.conn)
	return internal.DeserializeResponse(response)
}

func (c Client) Delete(key string) (*internal.Response) {
	request, _ := internal.Del{key, internal.BaseCommand{false, nil}}.Serialize()
	c.conn.Write(request)
	response := bufio.NewReader(c.conn)
	return internal.DeserializeResponse(response)
}

func (c Client) Keys(pattern string) (*internal.Response) {
	request, _ := internal.Keys{pattern, internal.BaseCommand{false, nil}}.Serialize()
	c.conn.Write(request)
	response := bufio.NewReader(c.conn)
	return internal.DeserializeResponse(response)
}

func (c Client) Count() (*internal.Response) {
	request, _ := internal.Count{internal.BaseCommand{false, nil}}.Serialize()
	c.conn.Write(request)
	response := bufio.NewReader(c.conn)
	return internal.DeserializeResponse(response)
}
