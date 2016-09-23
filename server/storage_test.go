package server

import (
	"testing"
	"log"
)

func TestPutThenGet(t *testing.T) {
	dbChannel := make(chan interface{})
	go ProcessCommands(dbChannel)

	setCommand := &SetUpd{"golang", "awesome", 10, false, BaseCommand{false, make(chan Resulter)}}
	dbChannel <- setCommand
	<- setCommand.GetBaseCommand().ChannelWithResult
	log.Println("set commnad finished")
	getCommand := &Get{"golang", BaseCommand{false,  make(chan Resulter)}}
	dbChannel <- getCommand
	result, _ := (<- getCommand.GetBaseCommand().ChannelWithResult).(*Result)

	if s, _ := result.Value.(string); s != setCommand.Value {
		t.Errorf("can't get after set ")
	}
}
