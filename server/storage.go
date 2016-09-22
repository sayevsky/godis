package server

type Value struct {
	Value interface{}
	TTL int
}

func ProcessCommands(dbCannel chan interface{}) {
	storage := make(map[string] *Value)

	for {
		command := <-dbCannel
		switch command := command.(type) {
		case *SetUpd:
			// sends +OK if value was updated successfully
			if command.update && storage[command.Key] == nil {
				break
			}
			value := &Value{command.Value, command.TTL}
			storage[command.Key] = value
			command.result <- "+OK\r\n"
		case *Get:
			value := storage[command.Key]

			if value == nil {
				command.result <- "-NE\r\n"
			} else {
				command.result <- value.Value
			}
		}
	}
}
