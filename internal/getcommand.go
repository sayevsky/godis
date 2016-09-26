package internal

import "bufio"

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
