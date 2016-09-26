package internal

import "bufio"

type Keys struct {
	Pattern string
	Base    BaseCommand
}

func (c *Keys) GetBaseCommand() BaseCommand {
	return c.Base
}

func DeserializeKeys(reader *bufio.Reader) (command *Keys, err error) {
	//<command>\r\n<numberOdBytesOfPattern>\r\n<pattern><\r\n
	size, err := readIntByDelim(reader)
	if err != nil {
		return nil, err
	}
	pattern, err := readDataGivenSize(reader, size)
	if err != nil {
		return nil, err
	}
	return &Keys{pattern, BaseCommand{false, make(chan Response)}}, nil
}
