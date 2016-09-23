package server

import (
	"regexp"
)

type Result struct {
	Value interface{}
	TTL int
}

func (Result) Serialize() ([]byte){
	return []byte("+OK")
}

func ProcessCommands(dbCannel chan interface{}) {
	storage := make(map[string] Resulter)

	for {
		command := <-dbCannel
		switch command := command.(type) {
		case *SetUpd:
			// sends previous data in cache
			if command.update && storage[command.Key] == nil {
				command.Base.ChannelWithResult <- nil
				break
			}
			value := &Result{command.Value, command.TTL}
			old := storage[command.Key]
			storage[command.Key] = value
			command.Base.ChannelWithResult <- old
		case *Get:
			command.Base.ChannelWithResult <- storage[command.Key]
		case *Del:
			old := storage[command.Key]
			delete(storage, command.Key)
			command.Base.ChannelWithResult <- old

		case *Keys:
			keys := make([]string, 0)
			pattern := command.Pattern
			re, err := regexp.Compile(pattern)
			if (err != nil) {
				command.Base.ChannelWithResult <- Result{err, 0}
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
		}
	}
}
