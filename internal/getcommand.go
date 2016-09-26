package internal

import (
	"bufio"
	"bytes"
	"strconv"
)

type Get struct {
	Key  string
	Base BaseCommand
}

func (c *Get) GetBaseCommand() BaseCommand {
	return c.Base
}

func DeserializeGet(reader *bufio.Reader) (command *Get, err error) {
	// <command>\r\n<numberOfBytesOfValue>\r\n<key>\r\n
	size, err := readIntByDelim(reader)
	if err != nil {
		return
	}

	key, err := readDataGivenSize(reader, size)
	if err != nil {
		return
	}

	return &Get{key, BaseCommand{false, make(chan Response)}}, nil
}

func (c Get ) Serialize() ([]byte, error) {
	// <command>\r\n<numberOfBytesOfValue>\r\n<key>\r\n
	var buffer bytes.Buffer
	buffer.WriteString("GET\r\n")
	buffer.WriteString(strconv.Itoa(len(c.Key)))
	buffer.WriteString("\r\n" + c.Key + "\r\n")
	return buffer.Bytes(), nil

}
