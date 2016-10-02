package internal

import (
	"bufio"
	"bytes"
	"strconv"
)

type Del struct {
	Key  string
	Base BaseCommand
}

func (c *Del) GetBaseCommand() BaseCommand {
	return c.Base
}

func DeserializeDelete(reader *bufio.Reader) (command *Del, err error) {
	//<command>\r\n<numberOdBytesOfKey>\r\n<key><\r\n<async>\r\n
	size, err := readIntByDelim(reader)
	if err != nil {
		return nil, err
	}
	key, err := readDataGivenSize(reader, size)
	if err != nil {
		return nil, err
	}
	var async bool
	size, err = readIntByDelim(reader)
	if err != nil {
		return nil, err
	}
	if size == 1 {
		async = true
	}
	return &Del{key, BaseCommand{async, make(chan Response)}}, nil
}

func (c Del) Serialize() ([]byte, error) {
	// <command>\r\n<numberOfBytesOfKey>\r\n<key>\r\n<async>\r\n
	var buffer bytes.Buffer
	buffer.WriteString("DEL\r\n")
	buffer.WriteString(strconv.Itoa(len(c.Key)))
	buffer.WriteString("\r\n" + c.Key + "\r\n")
	buffer.WriteString(strconv.Itoa(Btoi(c.GetBaseCommand().IsAsync)))
	buffer.WriteString("\r\n")
	return buffer.Bytes(), nil
}
