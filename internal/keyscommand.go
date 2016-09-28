package internal

import (
	"bufio"
	"bytes"
	"strconv"
)

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

func (c Keys) Serialize() ([]byte, error) {
	//<command>\r\n<numberOdBytesOfPattern>\r\n<pattern><\r\n
	var buffer bytes.Buffer
	buffer.WriteString("KEYS\r\n")
	buffer.WriteString(strconv.Itoa(len(c.Pattern)))
	buffer.WriteString("\r\n" + c.Pattern + "\r\n")
	return buffer.Bytes(), nil
}
