package internal

import (
	"bufio"
	"bytes"
	"strconv"
)

// get by index anf get by key in value

type GGetI struct {
	Key   string
	Index int
	Base  BaseCommand
}

func (c *GGetI) GetBaseCommand() BaseCommand {
	return c.Base
}

func DeserializeGGetI(reader *bufio.Reader) (command *GGetI, err error) {
	// <GGETI>\r\n<index>\r\n<numberOfBytesOfKey>\r\n<key>\r\n
	index, err := readIntByDelim(reader)
	if err != nil {
		return
	}
	get, err := DeserializeGet(reader)
	if err != nil {
		return
	}

	return &GGetI{get.Key, index, get.GetBaseCommand()}, nil
}

func (c GGetI) Serialize() ([]byte, error) {
	// <GGETI>\r\n<index>\r\n<numberOfBytesOfKey>\r\n<key>\r\n
	var buffer bytes.Buffer
	buffer.WriteString("GGETI\r\n")
	buffer.WriteString(strconv.Itoa(c.Index))
	buffer.WriteString("\r\n")
	buffer.WriteString(strconv.Itoa(len(c.Key)))
	buffer.WriteString("\r\n" + c.Key + "\r\n")
	return buffer.Bytes(), nil
}

type GGetK struct {
	Key        string
	KeyInValue string
	Base       BaseCommand
}

func (c *GGetK) GetBaseCommand() BaseCommand {
	return c.Base
}

func DeserializeGGetK(reader *bufio.Reader) (command *GGetK, err error) {
	// <GGETK>\r\n<numberOfBytesOfKeyInValue>\r\n<keyInValue>\r\n<numberOfBytesOfKey>\r\n<key>\r\n
	size, err := readIntByDelim(reader)
	if err != nil {
		return
	}
	keyInValue, nil := readDataGivenSize(reader, size)
	get, err := DeserializeGet(reader)
	if err != nil {
		return
	}

	return &GGetK{get.Key, keyInValue, get.GetBaseCommand()}, nil
}

func (c GGetK) Serialize() ([]byte, error) {
	// <GGETK>\r\n<numberOfBytesOfKeyInValue>\r\n<keyInValue>\r\n<numberOfBytesOfKey>\r\n<key>\r\n
	var buffer bytes.Buffer
	buffer.WriteString("GGETK\r\n")
	buffer.WriteString(strconv.Itoa(len(c.KeyInValue)))
	buffer.WriteString("\r\n")
	buffer.WriteString(c.KeyInValue)
	buffer.WriteString("\r\n")
	buffer.WriteString(strconv.Itoa(len(c.Key)))
	buffer.WriteString("\r\n" + c.Key + "\r\n")
	return buffer.Bytes(), nil
}
