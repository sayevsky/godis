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
			// sends +OK if value was updated successfully
			if command.update && storage[command.Key] == nil {
				break
			}
			value := &Result{command.Value, command.TTL}
			storage[command.Key] = value
			command.Base.ChannelWithResult <- Result{"", 0}
		case *Get:
			value := storage[command.Key]

			if value == nil {
				command.Base.ChannelWithResult <- Result{"", 0}
			} else {
				command.Base.ChannelWithResult <- Result{"", 0}
			}
		}
	}
}
