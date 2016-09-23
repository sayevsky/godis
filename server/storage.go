package server

type Result struct {
	Value interface{}
	TTL int
}

func (Result) Serialize() ([]byte){
	return []byte("+OK")
}

func ProcessCommands(dbCannel chan interface{}) {
	storage := make(map[string] *Result)

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
			//storage
		}
	}
}
