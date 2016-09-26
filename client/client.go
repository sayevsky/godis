package client

import (
	"bufio"
	"bytes"
	"net"
	"strconv"
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
	var buffer bytes.Buffer
	buffer.WriteString("GET\r\n")
	buffer.WriteString(strconv.Itoa(len(key)))
	buffer.WriteString("\r\n" + key + "\r\n")
	c.conn.Write(buffer.Bytes())
	result := bufio.NewReader(c.conn)
	resType, err := result.ReadByte()
	return string(resType), err

}
