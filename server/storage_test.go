package server

import (
	"github.com/sayevsky/godis/internal"
	"testing"
	"time"
)

func TestPutStringThenGet(t *testing.T) {
	dbChannel := make(chan interface{})
	go ProcessCommands(dbChannel)

	setCommand := &internal.SetUpd{"golang", "awesome", 0, false, internal.BaseCommand{false, make(chan internal.Response)}}
	dbChannel <- setCommand
	<-setCommand.GetBaseCommand().ChannelWithResult
	getCommand := &internal.Get{"golang", internal.BaseCommand{false, make(chan internal.Response)}}
	dbChannel <- getCommand
	result := (<-getCommand.GetBaseCommand().ChannelWithResult)

	if s, _ := result.Result.(string); s != setCommand.Value {
		t.Errorf("can't get after set ")
	}
}

func TestActiveEviction(t *testing.T) {
	dbChannel := make(chan interface{})
	go sendEvictMessages(dbChannel, make(chan bool))
	go ProcessCommands(dbChannel)

	setCommand := &internal.SetUpd{"golang", "awesome", 1 * time.Nanosecond, false,
		internal.BaseCommand{false, make(chan internal.Response)}}
	dbChannel <- setCommand
	<-setCommand.GetBaseCommand().ChannelWithResult
	countCommand := &internal.Count{internal.BaseCommand{false, make(chan internal.Response)}}
	dbChannel <- countCommand
	result := (<-countCommand.GetBaseCommand().ChannelWithResult)
	if s, _ := result.Result.(int); s == 0 {
		t.Errorf("not evicted")
	}
	time.Sleep(300 * time.Millisecond)
	dbChannel <- countCommand
	result = (<-countCommand.GetBaseCommand().ChannelWithResult)
	if s, _ := result.Result.(int); s != 0 {
		t.Errorf("not evicted")
	}
}

func TestPassiveEviction(t *testing.T) {

	dbChannel := make(chan interface{})
	go ProcessCommands(dbChannel)

	setCommand := &internal.SetUpd{"golang", "awesome", 1 * time.Nanosecond, false,
		internal.BaseCommand{false, make(chan internal.Response)}}
	dbChannel <- setCommand
	<-setCommand.GetBaseCommand().ChannelWithResult

	getCommand := &internal.Get{"golang", internal.BaseCommand{false, make(chan internal.Response)}}
	dbChannel <- getCommand
	result := (<-getCommand.GetBaseCommand().ChannelWithResult)

	if result.Result != nil {
		t.Errorf("can't get after set ")
	}
}
