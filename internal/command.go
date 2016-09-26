package internal

import "strconv"
import "bufio"
import "io"
import "log"
import "fmt"
import (
	"time"
	"bytes"
)

type Response struct {
	Result interface{}
	Err error
}

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
	Duration time.Duration
	Update bool
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

type OK byte

const  delim = '\n'

func ParseCommand(reader *bufio.Reader) (Commander, error) {
	// <command>\r\n<number of bytes>\r\n<key>...

	com, err := reader.ReadString(delim)
	if len(com) < 2 {
		return nil, fmt.Errorf("wrong format")
	}
	com = com[:len(com)-2]
	if err != nil {
		log.Println("Error reading request, wrong format? " + com, err)
		return nil, err
	}
	switch com {
	case "GET":
		// <command>\r\n<numberOfBytesOfValue>\r\n<key>\r\n
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
		command.Update = true
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

	value, err := readValue(reader, size)
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

func readValue(reader *bufio.Reader, size int) (value interface{}, err error) {
	// value could be a string (starts with '@'), a list ( starts with '*')
	// or dict (starts with '>')
	typ, err := reader.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("Can't read type of value", err)
	}
	switch typ {
	case byte('@'):
		return readDataGivenSize(reader, size - 1)
	case '*':
		// "*<size of array>\r\n<sizeOfBytesOfFirstElement>\r\n<First element>\r\n
		// ...<sizeOfBytesOfLastElement>\r\n<LastElement>\r\n"
		sizeOfArray, err := readIntByDelim(reader)
		if err != nil {
			return nil, fmt.Errorf("array is broken at sizeOfArray", err)
		}
		array := make([]string, sizeOfArray)
		for i := 0; i < sizeOfArray; i ++ {
			ithSize, err := readIntByDelim(reader)
			if err != nil {
				return nil, fmt.Errorf("array is broken at i=" +
					strconv.Itoa(i) + " ithSize=" + strconv.Itoa(ithSize) , err)
			}
			ithElement, err := readDataGivenSize(reader, ithSize)
			if err != nil {
				return nil, fmt.Errorf("array is broken at ith element " + ithElement, err)
			}
			array[i] = ithElement
		}
		return array, nil
	case '>':
	// "><size of dict>\r\n<sizeOfBytesOfFirstKey>\r\n
	// <First key>\r\n<sizeOfBytesOfFirstValue>\r\n<First value>\r\n...
	// <sizeOfBytesOfLastKey>\r\n<Last key>\r\n
	// <sizeOfBytesOfLastValue>\r\n<Last value>\r\n
		sizeOfDictionary, err := readIntByDelim(reader)
		if err != nil {
			return nil, fmt.Errorf("dict is broken at sizeOfDictionary", err)
		}
		dict := make(map[string]string, sizeOfDictionary)
		for i := 0; i < sizeOfDictionary; i++ {
			ithKeySize, err := readIntByDelim(reader)
			if err != nil {
				return nil, fmt.Errorf("dict is broken at ithKeySize=" + strconv.Itoa(ithKeySize), err)
			}
			ithKey, err := readDataGivenSize(reader, ithKeySize)
			if err != nil {
				return nil, fmt.Errorf("dict is broken at ithKey=" + ithKey, err)
			}
			ithValueSize, err := readIntByDelim(reader)
			if err != nil {
				return nil, fmt.Errorf("dict is broken at ithValueSize=" + strconv.Itoa(ithValueSize), err)
			}
			ithValue, err := readDataGivenSize(reader, ithValueSize)
			if err != nil {
				return nil, fmt.Errorf("dict is broken at ithValue=" + ithValue, err)
			}

			dict[ithKey] = ithValue
		}
		return dict, nil
	default:
		return nil, fmt.Errorf("Unknown value type " + string(typ))



	}

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

func (r Response) Serialize() ([]byte){
	// start with
	// error: -\r\n<numberOfBytes>\r\n<message>
	// success with result: +\r\n<result>
	// if result starts with @ it's a string +\r\n@<numberOfBytes>\r\n<mesage>
	// if result starts with * it's an array
	// 	+\r\n*<numberOfElements>\r\n<sizeOfFirstElement>\r\n<FirstElement>\r\n...<sizeOfLastElement>\r\n<LastElement>\r\n
	// if result starts with > it's an dict
	// +\r\n*<numberOfElements>\r\n<sizeOfFirstElement>\r\n<FirstElement>\r\n...<sizeOfLastElement>\r\n<LastElement>\r\n
	// if result starts with $ it's int $\r\n<number>\r\n

	var buffer bytes.Buffer

	if r.Err != nil {
		size := len(r.Err.Error())
		buffer.WriteString("-\r\n")
		buffer.WriteString(strconv.Itoa(size))
		buffer.WriteString("\r\n@")
		buffer.WriteString(r.Err.Error())
		buffer.WriteString("\r\n")
		return buffer.Bytes()
	}

	switch value := r.Result.(type) {
	case string:
		buffer.WriteString("+\r\n")
		buffer.WriteString(strconv.Itoa(len(value)))
		buffer.WriteString("\r\n@")
		buffer.WriteString(value)
		buffer.WriteString("\r\n")
	case []string:
		buffer.WriteString("+\r\n*")
		buffer.WriteString(strconv.Itoa(len(value)))
		buffer.WriteString("\r\n")
		for _, element := range value {
			buffer.WriteString(strconv.Itoa(len(element)))
			buffer.WriteString("\r\n")
			buffer.WriteString(element)
			buffer.WriteString("\r\n")
		}
	case map[string]string:
		buffer.WriteString("+\r\n>")
		buffer.WriteString(strconv.Itoa(len(value)))
		buffer.WriteString("\r\n")
		for k,v := range value {
			buffer.WriteString(strconv.Itoa(len(k)))
			buffer.WriteString("\r\n")
			buffer.WriteString(k)
			buffer.WriteString("\r\n")
			buffer.WriteString(strconv.Itoa(len(v)))
			buffer.WriteString("\r\n")
			buffer.WriteString(v)
			buffer.WriteString("\r\n")
		}
	case int:
		buffer.WriteString("+\r\n")
		buffer.WriteString(strconv.Itoa(value))
		buffer.WriteString("\r\n")
	case OK:
		buffer.WriteString("+\r\nOK\r\n")

	default:
		message := "Unxpected type"
		buffer.WriteString("-\r\n")
		buffer.WriteString(strconv.Itoa(len(message)))
		buffer.WriteString("\r\n@")
		buffer.WriteString(message)
		buffer.WriteString("\r\n")
	}

	return buffer.Bytes()
}