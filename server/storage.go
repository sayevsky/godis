package server

import (
	"regexp"
	"time"
	"log"
)

// TTL added as value to storage
type WrappedValue struct {
	Value interface{}
	TTL time.Time
}

func (WrappedValue) Serialize() ([]byte){
	return []byte("+OK")
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
	if(time.Now().After(value.TTL)){
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
	var zeroValue WrappedValue

	for {
		command := <-dbCannel
		switch command := command.(type) {
		case *SetUpd:
			// sends previous data in cache
			if command.update && storage[command.Key] == nil {
				command.Base.ChannelWithResult <- zeroValue
				break
			}
			ttl := durationToTTL(command.duration)
			value := &WrappedValue{command.Value, ttl}
			old := storage[command.Key]
			if old == nil {
				old = &zeroValue
			}
			storage[command.Key] = value
			command.Base.ChannelWithResult <- *old
		case *Get:
			value := storage[command.Key]
			//passive eviction
			if (expiredKey(command.Key, storage)) {
				delete(storage, command.Key)
			}
			value = storage[command.Key]
			if value == nil {
				value = &zeroValue
			}
			command.Base.ChannelWithResult <- *value
		case *Del:
			old := storage[command.Key]
			delete(storage, command.Key)
			command.Base.ChannelWithResult <- *old

		case *Keys:
			keys := make([]string, 0)
			pattern := command.Pattern
			re, err := regexp.Compile(pattern)
			if (err != nil) {
				command.Base.ChannelWithResult <- zeroValue
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
		case *Evict:
			// 20 (at most) randomly selected candidates to evict
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
		case *Count:
			size := len(storage)
			value := WrappedValue{size, time.Time{}}
			command.Base.ChannelWithResult <- value


		}
	}
}
