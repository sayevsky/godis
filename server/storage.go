package server

import (
	"fmt"
	"github.com/sayevsky/godis/internal"
	"regexp"
	"time"
)

// TTL added as value to storage
type WrappedValue struct {
	Value interface{}
	TTL   time.Time
}

func (w WrappedValue) IsZero() bool {
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

func expiredKey(key string, storage map[string]*WrappedValue) bool {
	if storage[key] == nil {
		return false
	}
	value := storage[key]
	if time.Now().After(value.TTL) && !value.TTL.IsZero() {
		return true
	} else {
		return false
	}
}

func sendEvictMessages(dbCannel chan interface{}, poisonPill chan bool) {
	evict := &internal.Evict{}
	for {
		select {
		case <-poisonPill:
			return
		default:
			time.Sleep(200 * time.Millisecond)
			dbCannel <- evict
		}

	}
}

func get(key string, storage map[string]*WrappedValue) (interface{}, error) {
	value := storage[key]
	//passive eviction
	if expiredKey(key, storage) {
		delete(storage, key)
	}
	value = storage[key]
	var res interface{}
	var err error
	if value != nil {
		res = value.Value
	} else {
		// not exist
		err = fmt.Errorf("NE")
	}
	return res, err
}

func ProcessCommands(dbChannel chan interface{}) {

	storage := make(map[string]*WrappedValue)

	for {
		command := <-dbChannel
		switch command := command.(type) {
		case *internal.SetUpd:
			if command.Update && storage[command.Key] == nil {
				command.Base.ChannelWithResult <- internal.Response{nil, fmt.Errorf("NE")}
				break
			}
			ttl := durationToTTL(command.Duration)
			wrappedValue := &WrappedValue{command.Value, ttl}
			storage[command.Key] = wrappedValue

			command.Base.ChannelWithResult <- internal.Response{"OK", nil}
		case *internal.Get:
			res, err := get(command.Key, storage)
			command.Base.ChannelWithResult <- internal.Response{res, err}
		case *internal.Del:
			old := storage[command.Key]
			var res interface{}
			var err error
			if old != nil {
				delete(storage, command.Key)
				res = old.Value
			} else {
				err = fmt.Errorf("NE")
			}
			command.Base.ChannelWithResult <- internal.Response{res, err}

		case *internal.Keys:

			pattern := command.Pattern
			re, err := regexp.Compile(pattern)
			if err != nil {
				command.Base.ChannelWithResult <- internal.Response{nil, err}
				break
			}
			keys := make([]string, 0)
			i := 0
			for k := range storage {
				matched := re.Match([]byte(k))
				if matched {
					keys = append(keys, k)
					i++
				}
			}
			command.Base.ChannelWithResult <- internal.Response{keys, nil}
		case *internal.Evict:
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
		case *internal.Count:
			size := len(storage)
			command.Base.ChannelWithResult <- internal.Response{size, nil}
		case *internal.GGetI:
			value, err := get(command.Key, storage)
			if err != nil {
				command.Base.ChannelWithResult <- internal.Response{value, err}
				break
			}
			array, ok := value.([]string)
			if !ok {
				// wrong type
				command.Base.ChannelWithResult <- internal.Response{value, fmt.Errorf("WT")}
				break
			}
			if len(array) <= command.Index {
				// out of range
				command.Base.ChannelWithResult <- internal.Response{array, fmt.Errorf("OOR")}
				break
			}
			command.Base.ChannelWithResult <- internal.Response{array[command.Index], nil}
		case *internal.GGetK:
			value, err := get(command.Key, storage)
			if err != nil {
				command.Base.ChannelWithResult <- internal.Response{value, err}
				break
			}
			tomap, ok := value.(map[string]string)
			if !ok {
				// wrong type
				command.Base.ChannelWithResult <- internal.Response{value, fmt.Errorf("WT")}
				break
			}
			res, ok := tomap[command.KeyInValue]
			if !ok {
				// wrong type
				command.Base.ChannelWithResult <- internal.Response{value, fmt.Errorf("NE")}
				break
			}

			command.Base.ChannelWithResult <- internal.Response{res, nil}

		case *internal.Quit:
			return
		}
	}
}
