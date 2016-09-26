package internal

import "bufio"

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
