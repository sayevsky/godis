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

const delim = '\n'

func ParseCommand(reader *bufio.Reader) (Commander, error) {
	// <command>\r\n<number of bytes>\r\n<key>...
	com, err := reader.ReadString(delim)
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
		DeserializeGet(reader)

	case "SET":
		return DeserializeSetUpd(reader, false)
	case "UPD":
		return DeserializeSetUpd(reader, true)

	case "DEL":
		DeserializeDelete(reader)
	case "KEYS":
		DeserializeKeys(reader)
	case "COUNT":
		//<command>\r\n
		return &Count{BaseCommand{false, make(chan Response)}}, nil
	}

	return nil, fmt.Errorf("Unknown incoming command.")
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
		return readDataGivenSize(reader, size-1)
	case '*':
		// "*<size of array>\r\n<sizeOfBytesOfFirstElement>\r\n<First element>\r\n
		// ...<sizeOfBytesOfLastElement>\r\n<LastElement>\r\n"
		sizeOfArray, err := readIntByDelim(reader)
		if err != nil {
			return nil, fmt.Errorf("array is broken at sizeOfArray", err)
		}
		array := make([]string, sizeOfArray)
		for i := 0; i < sizeOfArray; i++ {
			ithSize, err := readIntByDelim(reader)
			if err != nil {
				return nil, fmt.Errorf("array is broken at i="+
					strconv.Itoa(i)+" ithSize="+strconv.Itoa(ithSize), err)
			}
			ithElement, err := readDataGivenSize(reader, ithSize)
			if err != nil {
				return nil, fmt.Errorf("array is broken at ith element "+ithElement, err)
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
				return nil, fmt.Errorf("dict is broken at ithKeySize="+strconv.Itoa(ithKeySize), err)
			}
			ithKey, err := readDataGivenSize(reader, ithKeySize)
			if err != nil {
				return nil, fmt.Errorf("dict is broken at ithKey="+ithKey, err)
			}
			ithValueSize, err := readIntByDelim(reader)
			if err != nil {
				return nil, fmt.Errorf("dict is broken at ithValueSize="+strconv.Itoa(ithValueSize), err)
			}
			ithValue, err := readDataGivenSize(reader, ithValueSize)
			if err != nil {
				return nil, fmt.Errorf("dict is broken at ithValue="+ithValue, err)
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
		log.Println("Error to parse bytesNumber "+string(bytesNumber), err)
		return
	}
	return
}

func readDurationByDelim(reader *bufio.Reader) (duration time.Duration, err error) {
	str, err := readByDelim(reader)
	duration, err = time.ParseDuration(str)
	if err != nil {
		log.Println("Error to parse duration "+str, err)
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
