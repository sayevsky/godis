package server

import "strconv"
import "bufio"
import "io"
import "log"
import "fmt"
import "time"



type Commander interface {
	GetBaseCommand() BaseCommand
}

type BaseCommand struct {
	IsAsync bool
	ChannelWithResult chan Response
}

type Get struct {
	Key string
	Base BaseCommand
}

func (c *Get) GetBaseCommand() (BaseCommand){
	return c.Base
}

type Del struct {
	Key string
	Base BaseCommand
}

func (c *Del) GetBaseCommand() (BaseCommand){
	return c.Base
}

type SetUpd struct {
	Key string
	Value interface{}
	duration time.Duration
	update bool
	Base BaseCommand
}

func (c *SetUpd) GetBaseCommand() (BaseCommand){
	return c.Base
}

type Keys struct {
	Pattern string
	Base BaseCommand
}

type Count struct {
	Base BaseCommand
}

func (c *Count) GetBaseCommand() (BaseCommand){
	return c.Base
}

// command will initiate active eviction
type Evict struct {}

func (c *Evict) GetBaseCommand() (BaseCommand){
	return BaseCommand{}
}


func (c *Keys) GetBaseCommand() (BaseCommand){
	return c.Base
}

const  delim = '\n'

func ParseCommand(reader *bufio.Reader) (Commander, error) {
	// <command>\r\n<number of bytes>\r\n<key>...

	com, err := reader.ReadString(delim)
	if len(com) < 2{
		return nil, fmt.Errorf("wrong format")
	}
	com = com[:len(com)-2]
	if err != nil {
		log.Println("Error reading request, wrong format? " + com, err)
		return nil, err
	}
	switch com {
	case "GET":
		// <command>\r\n<numberOfBytesOfValue>\r\n<value>\r\n
		size, err := readIntByDelim(reader)
		if(err != nil){
			return nil, err
		}

		key, err := readDataGivenSize(reader, size)
		if(err != nil){
			return nil, err
		}

 		return &Get{key, BaseCommand{false,  make(chan Response)}}, nil

	case "SET":
		command, err := parseSetUpd(reader)

		return command, err
	case "UPD":
		command, err := parseSetUpd(reader)
		if err != nil {
			return nil, err
		}
		command.update = true
		return command, err
	case "DEL":
		//<command>\r\n<numberOdBytesOfKey>\r\n<key><\r\n<async>\r\n
		size, err := readIntByDelim(reader)
		if err != nil {
			return nil, err
		}
		key, err := readDataGivenSize(reader, size)
		if(err != nil){
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
	case "KEYS":
		//<command>\r\n<numberOdBytesOfPattern>\r\n<pattern><\r\n
		size, err := readIntByDelim(reader)
		if err != nil {
			return nil, err
		}
		pattern, err := readDataGivenSize(reader, size)
		if(err != nil){
			return nil, err
		}
		return &Keys{pattern, BaseCommand{false, make(chan Response)}}, nil
	case "COUNT":
		//<command>\r\n
	return &Count{BaseCommand{false, make(chan Response)}}, nil
	}

	return nil, fmt.Errorf("Unknown incoming command.")
}

func parseSetUpd(reader *bufio.Reader) (setupd *SetUpd, err error) {
	// SET\r\n<numberOfBytes>\r\n<key>\r\n<numberOfBytes>\r\n<value>\r\n<TTL in duration format>\r\n<async>\r\n
	size, err := readIntByDelim(reader)
	if(err != nil){
		return
	}

	key, err := readDataGivenSize(reader, size)
	if(err != nil){
		return
	}
	size, err = readIntByDelim(reader)
	if(err != nil){
		return
	}

	value, err := readDataGivenSize(reader, size)
	if(err != nil){
		return
	}

	ttl, err := readDurationByDelim(reader)
	if(err != nil){
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

	return &SetUpd{key, value, ttl, false, BaseCommand{async, make(chan Response)}}, nil
}

func readIntByDelim(reader *bufio.Reader) (size int, err error) {
	bytesNumber, err := readByDelim(reader)
	size, err = strconv.Atoi(string(bytesNumber))
	if err != nil {
		log.Println("Error to parse bytesNumber " + string(bytesNumber), err)
		return
	}
	return
}

func readDurationByDelim(reader *bufio.Reader) (duration time.Duration, err error) {
	str, err := readByDelim(reader)
	duration, err = time.ParseDuration(str)
	if err != nil {
		log.Println("Error to parse duration " + str, err)
		return
	}
	return
}

func readByDelim(reader *bufio.Reader) (data string, err error) {
	data, err = reader.ReadString(delim)
	if err != nil {
		log.Println("Error reading request, wrong format? ", err)
		return
	}
	if len(data) == 2 {
		return data, fmt.Errorf("Empty payload")
	}
	data = data[:len(data) - 2]

	return
}

func readDataGivenSize(reader *bufio.Reader, size int) (value string, err error) {
	bytes := make([]byte, size)
	n, err := io.ReadFull(reader, bytes)
	if n != size || err != nil{
		log.Println("Can't read data with expected size " + strconv.Itoa(n), err)
		return
	}
	value = string(bytes)
	// verify ending \r\n
	verify := make ([]byte, 2)
	_, err = io.ReadFull(reader, verify)
	if(err != nil || string(verify) != "\r\n") {
		log.Println("requested command malformed")
		return
	}

	return
}