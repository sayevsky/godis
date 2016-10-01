package internal

import (
	"bufio"
	"bytes"
	"strconv"
	"time"
)

type SetUpd struct {
	Key      string
	Value    interface{}
	Duration time.Duration
	Update   bool
	Base     BaseCommand
}

func (c *SetUpd) GetBaseCommand() BaseCommand {
	return c.Base
}

func DeserializeSetUpd(reader *bufio.Reader, update bool) (command *SetUpd, err error) {
	// SET\r\n<numberOfBytes>\r\n<key>\r\n<value>\r\n<TTL in duration format>\r\n<async>\r\n
	size, err := readIntByDelim(reader)
	if err != nil {
		return
	}

	key, err := readDataGivenSize(reader, size)
	if err != nil {
		return
	}
	value, err := ReadValue(reader)
	if err != nil {
		return
	}

	ttl, err := readDurationByDelim(reader)
	if err != nil {
		return
	}

	var async bool

	isasync, err := readIntByDelim(reader)

	if err != nil {
		return nil, err
	}
	if isasync == 1 {
		async = true
	}

	return &SetUpd{key, value, ttl, update, BaseCommand{async, make(chan Response)}}, nil
}

func (c SetUpd) Serialize() ([]byte, error) {
	// <SET/UPD>\r\n<numberOfBytes>\r\n<key>\r\n<value>\r\n<TTL in duration format>\r\n<async>\r\n
	var buffer bytes.Buffer
	if c.Update {
		buffer.WriteString("UPD\r\n")
	} else {
		buffer.WriteString("SET\r\n")
	}
	buffer.WriteString(strconv.Itoa(len(c.Key)))
	buffer.WriteString("\r\n" + c.Key + "\r\n")
	buf, err := SerializeValue(c.Value)
	if err != nil {
		return nil, err
	}
	buf.WriteTo(&buffer)
	buffer.WriteString(c.Duration.String())
	buffer.WriteString("\r\n")
	buffer.WriteString(strconv.Itoa(Btoi(c.GetBaseCommand().IsAsync)))
	buffer.WriteString("\r\n")

	return buffer.Bytes(), nil
}
