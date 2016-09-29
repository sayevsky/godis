package internal

import (
	"bufio"
	"bytes"
)

type Count struct {
	Base BaseCommand
}

func (c *Count) GetBaseCommand() BaseCommand {
	return c.Base
}

func DeserializeCount(reader *bufio.Reader) (command *Count, err error) {
	//<command>\r\n
	return &Count{BaseCommand{false, make(chan Response)}}, nil
}

func (c Count) Serialize() ([]byte, error) {
	//<command><\r\n
	var buffer bytes.Buffer
	buffer.WriteString("COUNT\r\n")
	return buffer.Bytes(), nil
}
