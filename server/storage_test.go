package server

import (
	"testing"
	"time"
)

func TestPutThenGet(t *testing.T) {
	dbChannel := make(chan interface{})
	go ProcessCommands(dbChannel, true)

	setCommand := &SetUpd{"golang", "awesome", 0, false, BaseCommand{false, make(chan Response)}}
	dbChannel <- setCommand
	<- setCommand.GetBaseCommand().ChannelWithResult
	getCommand := &Get{"golang", BaseCommand{false,  make(chan Response)}}
	dbChannel <- getCommand
	result := (<- getCommand.GetBaseCommand().ChannelWithResult)

	if s, _ := result.Result.(string); s != setCommand.Value {
		t.Errorf("can't get after set ")
	}
}

func TestActiveEviction(t *testing.T) {
	dbChannel := make(chan interface{})
	go ProcessCommands(dbChannel, true)

	setCommand := &SetUpd{"golang", "awesome", 1 * time.Nanosecond, false, BaseCommand{false, make(chan Response)}}
	dbChannel <- setCommand
	<- setCommand.GetBaseCommand().ChannelWithResult
	countCommand := &Count{BaseCommand{false,  make(chan Response)}}
	dbChannel <- countCommand
	result := (<- countCommand.GetBaseCommand().ChannelWithResult)
	if s, _ := result.Result.(int); s == 0 {
		t.Errorf("not evicted")
	}
	time.Sleep(300 * time.Millisecond)
	dbChannel <- countCommand
	result = (<- countCommand.GetBaseCommand().ChannelWithResult)
	if s, _ := result.Result.(int); s != 0 {
		t.Errorf("not evicted")
	}
}

func TestPassiveEviction(t *testing.T) {
	dbChannel := make(chan interface{})
	go ProcessCommands(dbChannel, false)

	setCommand := &SetUpd{"golang", "awesome", 1 * time.Nanosecond, false, BaseCommand{false, make(chan Response)}}
	dbChannel <- setCommand
	<- setCommand.GetBaseCommand().ChannelWithResult

	getCommand := &Get{"golang", BaseCommand{false,  make(chan Response)}}
	dbChannel <- getCommand
	result := (<- getCommand.GetBaseCommand().ChannelWithResult)

	if result.Result != nil {
		t.Errorf("can't get after set ")
	}
}
