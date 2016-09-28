package internal

import "strconv"
import "bufio"
import "io"
import "log"
import "fmt"
import "time"

type Commander interface {
	GetBaseCommand() BaseCommand
	//Serialize() []byte
}

type BaseCommand struct {
	IsAsync           bool
	ChannelWithResult chan Response
}

type Count struct {
	Base BaseCommand
}

func (c *Count) GetBaseCommand() BaseCommand {
	return c.Base
}

// command will initiate active eviction
type Evict struct{}

func (c *Evict) GetBaseCommand() BaseCommand {
	return BaseCommand{}
}

type Quit struct{}

func (c *Quit) GetBaseCommand() BaseCommand {
	return BaseCommand{}
}

const Delim = '\n'

func ParseCommand(reader *bufio.Reader) (Commander, error) {
	// <command>\r\n<number of bytes>\r\n<key>...
	com, err := reader.ReadString(Delim)
	if len(com) < 2 {
		return nil, fmt.Errorf("wrong format")
	}
	com = com[:len(com)-2]
	if err != nil {
		log.Println("Error reading request, wrong format? "+com, err)
		return nil, err
	}
	switch com {
	case "GET":
		return DeserializeGet(reader)

	case "SET":
		return DeserializeSetUpd(reader, false)
	case "UPD":
		return DeserializeSetUpd(reader, true)

	case "DEL":
		return DeserializeDelete(reader)
	case "KEYS":
		return DeserializeKeys(reader)
	case "COUNT":
		//<command>\r\n
		return &Count{BaseCommand{false, make(chan Response)}}, nil
	}

	return nil, fmt.Errorf("Unknown incoming command.")
}

func readIntByDelim(reader *bufio.Reader) (size int, err error) {
	bytesNumber, err := ReadByDelim(reader)
	size, err = strconv.Atoi(string(bytesNumber))
	if err != nil {
		log.Println("Error to parse bytesNumber " + string(bytesNumber), err)
		return
	}
	return
}

func readDurationByDelim(reader *bufio.Reader) (duration time.Duration, err error) {
	str, err := ReadByDelim(reader)
	duration, err = time.ParseDuration(str)
	if err != nil {
		log.Println("Error to parse duration "+str, err)
		return
	}
	return
}

func ReadByDelim(reader *bufio.Reader) (data string, err error) {
	data, err = reader.ReadString(Delim)
	if err != nil {
		log.Println("Error reading request, wrong format? ", err)
		return
	}
	if len(data) == 2 {
		return data, fmt.Errorf("Empty payload")
	}
	data = data[:len(data)-2]

	return
}

func readDataGivenSize(reader *bufio.Reader, size int) (value string, err error) {
	bytes := make([]byte, size)
	n, err := io.ReadFull(reader, bytes)
	if n != size || err != nil {
		log.Println("Can't read data with expected size "+strconv.Itoa(n), err)
		return
	}
	value = string(bytes)
	// verify ending \r\n
	verify := make([]byte, 2)
	_, err = io.ReadFull(reader, verify)
	if err != nil || string(verify) != "\r\n" {
		log.Println("requested command malformed")
		return
	}

	return
}

func Btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}