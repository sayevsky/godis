package server

import (
	"regexp"
	"time"
	"fmt"
	"strconv"
	"bytes"
)

type OK byte

type Response struct {
	Result interface{}
	err error
}

// TTL added as value to storage
type WrappedValue struct {
	Value interface{}
	TTL time.Time
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

	if r.err != nil {
		size := len(r.err.Error())
		buffer.WriteString("-\r\n")
		buffer.WriteString(strconv.Itoa(size))
		buffer.WriteString("\r\n@")
		buffer.WriteString(r.err.Error())
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
func (w WrappedValue) IsZero() (bool){
	return w.Value == nil && w.TTL.IsZero()
}


// if duraion is 0 then return zero-time
// otherwise duration + current time
func durationToTTL(duration time.Duration) time.Time {
	var ttl time.Time
	if duration == 0 {
		return ttl
	}
	ttl = time.Now().Add(duration)
	return ttl
}

func expiredKey(key string, storage map[string] *WrappedValue) bool {
	if storage[key] == nil {
		return false
	}
	value := storage[key]
	if(time.Now().After(value.TTL) && !value.TTL.IsZero()){
		return true
	} else{
		return false
	}
}

func sendEvictMessages(dbCannel chan interface{}) bool {
	evict := &Evict{}
	for {
		time.Sleep(200 * time.Millisecond)
		dbCannel <- evict
	}
}

func ProcessCommands(dbCannel chan interface{}, withActiveEviction bool) {

	if withActiveEviction {
		go sendEvictMessages(dbCannel)
	}
	storage := make(map[string] *WrappedValue)

	for {
		command := <-dbCannel
		switch command := command.(type) {
		case *SetUpd:
			if command.update && storage[command.Key] == nil {
				command.Base.ChannelWithResult <- Response {nil, fmt.Errorf("Fail to update. Key doesn't exist.")}
				break
			}
			ttl := durationToTTL(command.duration)
			wrappedValue := &WrappedValue{command.Value, ttl}
			storage[command.Key] = wrappedValue

			command.Base.ChannelWithResult <- Response{OK(0), nil}
		case *Get:
			value := storage[command.Key]
			//passive eviction
			if (expiredKey(command.Key, storage)) {
				delete(storage, command.Key)
			}
			value = storage[command.Key]
			var res interface{}
			if value != nil {
				res = value.Value
			}
			command.Base.ChannelWithResult <- Response{res, nil}
		case *Del:
			old := storage[command.Key]
			var res interface{}
			if old != nil {
				delete(storage, command.Key)
				res = old.Value
			}
			command.Base.ChannelWithResult <- Response{res, nil}

		case *Keys:
			keys := make([]string, 0)
			pattern := command.Pattern
			re, err := regexp.Compile(pattern)
			if (err != nil) {
				command.Base.ChannelWithResult <- Response{nil, err}
				break
			}
			i := 0
			for k := range storage {
				matched := re.Match([]byte(k))
				if matched {
					keys[i] = k
					i++
				}
			}
			command.Base.ChannelWithResult <- Response{keys, nil}
		case *Evict:
			// 20 (at most) randomly selected candidates to evict
			// see also (go sendEvictMessages(dbCannel))
			amountToSelect := len(storage)
			if amountToSelect > 20 {
				amountToSelect = 20
			}
			// range of map do not return uniform distribution
			// but probably it has enough 'ramdomness'
			// for real-life problems this should be reimplemented
			i := 0
			for key := range storage {
				if i > amountToSelect {
					break
				}
				if expiredKey(key, storage) {
					delete(storage, key)
				}
			}
			// here we can return number of evicted keys for statistics
		case *Count:
			size := len(storage)
			command.Base.ChannelWithResult <- Response{size, nil}


		}
	}
}
