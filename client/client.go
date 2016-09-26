package client

import (
	"bufio"
	"bytes"
	"net"
	"strconv"
	"github.com/sayevsky/godis/internal"
)

func NewClient(addr string) (Client, error) {
	conn, err := net.Dial("tcp", addr)
	return Client{conn}, err
}

type Client struct {
	conn net.Conn
}

func (c Client) Get(key string) (interface{}, error) {
	// <command>\r\n<numberOfBytesOfValue>\r\n<key>\r\n
	datum, _ := internal.Get{key}.Serialize()
	c.conn.Write(datum)
	result := bufio.NewReader(c.conn)
	readValue(result)
	return string(resType), err

}
