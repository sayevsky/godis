package server

import (
	"testing"
)

func TestPutThenGet(t *testing.T) {
	dbChannel := make(chan interface{})
	go ProcessCommands(dbChannel)

	setCommand := &SetUpd{"golang", "awesome", 0, false, BaseCommand{false, make(chan WrappedValue)}}
	dbChannel <- setCommand
	<- setCommand.GetBaseCommand().ChannelWithResult
	getCommand := &Get{"golang", BaseCommand{false,  make(chan WrappedValue)}}
	dbChannel <- getCommand
	result := (<- getCommand.GetBaseCommand().ChannelWithResult)

	if s, _ := result.Value.(string); s != setCommand.Value {
		t.Errorf("can't get after set ")
	}
}
