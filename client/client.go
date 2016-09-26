package client

import (
	"bufio"
	"net"
	"github.com/sayevsky/godis/internal"
	"fmt"
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
	request, _ := internal.Get{key, nil}.Serialize()
	c.conn.Write(request)
	response := bufio.NewReader(c.conn)
	status, err := response.ReadString(internal.Delim)
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

	return response, err

}
