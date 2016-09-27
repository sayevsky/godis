package client

import (
	"bufio"
	"fmt"
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

func (c Client) Get(key string) (interface{}, error) {
	// <command>\r\n<numberOfBytesOfValue>\r\n<key>\r\n
	request, _ := internal.Get{key, internal.BaseCommand{false, nil}}.Serialize()
	c.conn.Write(request)
	response := bufio.NewReader(c.conn)
	status, err := internal.ReadByDelim(response)
	if err != nil {
		return nil, err
	}

	result, readErr := internal.ReadValue(response)
	if readErr != nil {
		return result, readErr
	}
	if status == "-" {
		return nil, fmt.Errorf(result.(string))
	}

	return result, err

}
