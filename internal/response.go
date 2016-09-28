package internal

import (
	"bufio"
	"bytes"
	"fmt"
)

type Response struct {
	Result interface{}
	Err    error
}


func DeserializeResponse(reader *bufio.Reader) (*Response) {
	// start with
	// error: -\r\n@<numberOfBytes>\r\n<message>
	// success with result: +\r\n<result>
	// result serialized as 'value'
	var response Response
	status, err := ReadByDelim(reader)
	if err != nil {
		response.Err = err
		return &response
	}
	value, err := ReadValue(reader)
	if status == "-" {
		response.Err = fmt.Errorf(value.(string))
	} else {
		response.Result = value
	}

	return &response

}

func (r Response) Serialize() []byte {
	// start with
	// error: -\r\n@<numberOfBytes>\r\n<message>
	// success with result: +\r\n<result>
	// result serialized as 'value'

	var b bytes.Buffer
	buffer := &b

	if r.Err != nil {
		failStatus(buffer)
		buf, _ := SerializeValue(r.Err)
		buf.WriteTo(buffer)
		return buffer.Bytes()
	}
	buf, err := SerializeValue(r.Result)

	if err != nil {
		failStatus(buffer)
		buf, _ = SerializeValue(err)
		buf.WriteTo(buffer)
		return buffer.Bytes()

	}
	okStatus(buffer)
	buf.WriteTo(buffer)
	return buffer.Bytes()
}

func okStatus(buffer *bytes.Buffer) {
	buffer.WriteString("+\r\n")
}

func failStatus(buffer *bytes.Buffer) {
	buffer.WriteString("-\r\n")
}
