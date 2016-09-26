package internal

import (
	"bytes"
	"strconv"
)

type Response struct {
	Result interface{}
	Err    error
}

func (r Response) Serialize() []byte {
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
		buffer.WriteString("+\r\n")
		buffer.WriteString(strconv.Itoa(value))
		buffer.WriteString("\r\n")

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
