package config

import (
	"custom-go/iot/base"
	"custom-go/pkg/types"
	"errors"
	"sync"
	"time"
)

var eventArrivedInfos sync.Map

type eventArrivedInfo struct {
	time.Time
	*sync.Mutex
	value string
}

func init() {
	base.SetEventAroundFunc(func(client *types.InternalClient, input *base.EventAroundInput) (deferFunc func(error), err error) {
		var previousInfo *eventArrivedInfo
		arrivedInfo := &eventArrivedInfo{
			Time:  input.Properties.GenerateTime,
			Mutex: &sync.Mutex{},
		}
		if input.CheckRepeatValue {
			arrivedInfo.value = string(input.ExtractedValue)
		}
		value, loaded := eventArrivedInfos.Swap(input.Properties.Topic, arrivedInfo)
		if loaded {
			previousInfo = value.(*eventArrivedInfo)
			arrivedInfo.Mutex = previousInfo.Mutex
			if input.CheckRepeatValue && arrivedInfo.value == previousInfo.value {
				input.Properties.Repeated = true
			}
		}

		eventOutdated := loaded && input.Properties.GenerateTime.Before(previousInfo.Time)
		if eventOutdated {
			err = errors.New("event arrived time expired")
			return
		}

		arrivedInfo.Mutex.Lock()
		deferFunc = func(_err error) {
			arrivedInfo.Mutex.Unlock()
		}
		return
	})
}
