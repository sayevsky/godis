package server

import "strconv"
import "bufio"
import "io"
import "log"
import "fmt"
import "strings"

type Get struct {
	Key string
}
type Del struct {
	Key string
}
type SetUpd struct {
	Key string
	Value string
	TTL int
	update bool
}

const  delim = '\n'

func ParseCommand(reader *bufio.Reader) (interface{}, error) {
	// <command>\r\n<number of bytes>\r\n<key>...

	com, err := reader.ReadString(delim)
	com = com[:len(com)-2]
	if err != nil {
		log.Println("Error reading request, wrong format? " + com, err)
		return nil, err
	}
	switch strings.ToUpper(com) {
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

 		return &Get{key}, nil

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
	case"DEL":
		//<command>\r\n<numberOdBytesOfKey>\r\n<key><\r\n
		size, err := readIntByDelim(reader)
		if err != nil {
			return nil, err
		}
		key, err := readDataGivenSize(reader, size)
		if(err != nil){
			return nil, err
		}
		return &Del{key}, nil

	}

	return nil, fmt.Errorf("Can't parse incoming command")
}

func parseSetUpd(reader *bufio.Reader) (setupd *SetUpd, err error) {
	// SET\r\n<numberOfBytes>\r\n<key>\r\n<numberOfBytes>\r\n<value>\r\n<TTL>\r\n
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

	ttl, err := readIntByDelim(reader)
	if(err != nil){
		return
	}

	return &SetUpd{key, value, ttl, false}, nil
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

func readByDelim(reader *bufio.Reader) (data string, err error) {
	data, err = reader.ReadString(delim)
	if err != nil {
		log.Println("Error reading request, wrong format? ", err)
		return
	}
	if len(data) == 2 {
		return
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