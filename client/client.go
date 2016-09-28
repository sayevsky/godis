package client

import (
	"bufio"
	"github.com/sayevsky/godis/internal"
	"net"
	"log"
)

func NewClient(addr string) (Client, error) {
	conn, err := net.Dial("tcp", addr)
	return Client{conn}, err
}

type Client struct {
	conn net.Conn
}

func (c Client) Get(key string) (*internal.Response) {
	// <command>\r\n<numberOfBytesOfValue>\r\n<key>\r\n
	request, _ := internal.Get{key, internal.BaseCommand{false, nil}}.Serialize()
	c.conn.Write(request)
	response := bufio.NewReader(c.conn)
	return internal.DeserializeResponse(response)
}

func (c Client) Set(key string, value string) (*internal.Response) {
	// <command>\r\n<numberOfBytesOfValue>\r\n<key>\r\n
	request, _ := internal.SetUpd{key, value, 0, false, internal.BaseCommand{false, nil}}.Serialize()
	log.Println(string(request))
	c.conn.Write(request)
	response := bufio.NewReader(c.conn)
	return internal.DeserializeResponse(response)
}
