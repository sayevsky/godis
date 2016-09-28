package internal

import (
	"strconv"
	"fmt"
	"bytes"
	"bufio"
)

func ReadValue(reader *bufio.Reader) (value interface{}, err error) {
	// value could be a string (starts with '@'), a list ( starts with '*')
	// or dict (starts with '>')
	typ, err := reader.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("Can't read type of value", err)
	}
	switch typ {
	case '@':
		size, err := readIntByDelim(reader)
		if err != nil {
			return nil, err
		}
		return readDataGivenSize(reader, size)
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

// if value starts with @ it's a string +\r\n@<numberOfBytes>\r\n<mesage>
// if value starts with * it's an array
// 	+\r\n*<numberOfElements>\r\n<sizeOfFirstElement>\r\n<FirstElement>\r\n...<sizeOfLastElement>\r\n<LastElement>\r\n
// if value starts with > it's an dict
// +\r\n><numberOfElements>\r\n<sizeOfFirstElement>\r\n<FirstElement>\r\n...<sizeOfLastElement>\r\n<LastElement>\r\n
// if value starts with $ it's int $\r\n<number>\r\n

func SerializeValue(value interface{}, ) (buffer *bytes.Buffer, err error) {
	var b bytes.Buffer
	buffer = &b
	//TODO: error handling if buffer are too large
	switch value := value.(type) {
	case string:
		buffer.WriteString("@")
		buffer.WriteString(strconv.Itoa(len(value)))
		buffer.WriteString("\r\n")
		buffer.WriteString(value)
		buffer.WriteString("\r\n")
	case []string:
		buffer.WriteString("*")
		buffer.WriteString(strconv.Itoa(len(value)))
		buffer.WriteString("\r\n")
		for _, element := range value {
			buffer.WriteString(strconv.Itoa(len(element)))
			buffer.WriteString("\r\n")
			buffer.WriteString(element)
			buffer.WriteString("\r\n")
		}
	case map[string]string:
		buffer.WriteString(strconv.Itoa(len(value)))
		buffer.WriteString("\r\n")
		for k, v := range value {
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
		buffer.WriteString(strconv.Itoa(value))
		buffer.WriteString("\r\n")
	case error:
		size := len(value.Error())
		buffer.WriteString("@")
		buffer.WriteString(strconv.Itoa(size))
		buffer.WriteString("\r\n")
		buffer.WriteString(value.Error())
		buffer.WriteString("\r\n")

	default:
		err = fmt.Errorf("Unxpected type")
	}
	return buffer, err
}
