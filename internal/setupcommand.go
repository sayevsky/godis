package internal

import (
	"bufio"
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
	// SET\r\n<numberOfBytes>\r\n<key>\r\n<numberOfBytes>\r\n<value>\r\n<TTL in duration format>\r\n<async>\r\n
	size, err := readIntByDelim(reader)
	if err != nil {
		return
	}

	key, err := readDataGivenSize(reader, size)
	if err != nil {
		return
	}
	size, err = readIntByDelim(reader)
	if err != nil {
		return
	}

	value, err := readValue(reader, size)
	if err != nil {
		return
	}

	ttl, err := readDurationByDelim(reader)
	if err != nil {
		return
	}

	var async bool

	size, err = readIntByDelim(reader)

	if err != nil {
		return nil, err
	}
	if size == 1 {
		async = true
	}

	return &SetUpd{key, value, ttl, update, BaseCommand{async, make(chan Response)}}, nil
}
