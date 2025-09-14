package base

import (
	"custom-go/pkg/types"
	"github.com/labstack/echo/v4"
	"github.com/spf13/cast"
	"golang.org/x/exp/slices"
	"strconv"
	"strings"
	"time"
)

var (
	eventTopic      = new(topicTrie)
	eventAroundFunc func(*types.InternalClient, *EventAroundInput) (func(error), error)
)

type EventAroundInput struct {
	Properties       *ApplicationProperties
	DeviceName       string
	Data             []byte
	ExtractedValue   []byte
	CheckRepeatValue bool
}

func SetEventAroundFunc(aroundFunc func(*types.InternalClient, *EventAroundInput) (func(error), error)) {
	eventAroundFunc = aroundFunc
}

func (am *AmqpManager) isIncludeDevice(deviceName string) (bool, bool) {
	return len(am.includeDeviceNames) > 0, slices.Contains(am.includeDeviceNames, deviceName)
}

func (am *AmqpManager) isExcludeDevice(deviceName string) (bool, bool) {
	return len(am.excludeDeviceNames) > 0, slices.Contains(am.excludeDeviceNames, deviceName)
}

func (am *AmqpManager) executeHandler(appProperties *ApplicationProperties, data []byte) (err error) {
	handler, options, params := eventTopic.findHandler(appProperties.Topic)
	if handler == nil {
		return
	}
	deviceName := params["deviceName"]
	// 包含设备功能开启且不匹配当前设备
	enabled, matched := am.isIncludeDevice(deviceName)
	if enabled && !matched {
		return
	}
	// 排除设备功能开启且匹配当前设备
	enabled, matched = am.isExcludeDevice(deviceName)
	if enabled && matched {
		return
	}
	eventAroundInput := &EventAroundInput{
		Properties: appProperties,
		DeviceName: deviceName,
		Data:       data,
	}
	if options.ExtractValue != nil {
		eventAroundInput.ExtractedValue = options.ExtractValue(data)
		eventAroundInput.CheckRepeatValue = options.CheckRepeatValue
	}
	internalClient := types.NewEmptyInternalClient()
	if eventAroundFunc != nil {
		deferFunc, _err := eventAroundFunc(internalClient, eventAroundInput)
		if _err != nil {
			return _err
		}
		defer deferFunc(err)
	}
	handlerInput := &topicHandlerInput{
		Properties:     appProperties,
		Data:           data,
		ExtractedValue: eventAroundInput.ExtractedValue,
	}
	err = handler(am.logger, internalClient, handlerInput)
	return
}

type (
	topicHandler      func(echo.Logger, *types.InternalClient, *topicHandlerInput) error
	topicHandlerInput struct {
		Properties     *ApplicationProperties
		Data           []byte
		ExtractedValue []byte
	}
	topicTrie struct {
		staticItems map[string]*topicTrie
		paramItems  map[string]*topicTrie
		handler     topicHandler
		options     TopicTrieOptions
	}
	TopicTrieOptions struct {
		ExtractValue     func([]byte) []byte
		CheckRepeatValue bool
	}
)

func (t *topicTrie) addHandler(topic string, handler topicHandler, options ...TopicTrieOptions) {
	node := t
	for _, v := range strings.Split(topic, "/") {
		if node.staticItems == nil {
			node.staticItems = make(map[string]*topicTrie)
		}
		if node.paramItems == nil {
			node.paramItems = make(map[string]*topicTrie)
		}
		var items map[string]*topicTrie
		if strings.HasPrefix(v, "$") {
			items = node.paramItems
		} else {
			items = node.staticItems
		}
		nextNode, found := items[v]
		if found {
			node = nextNode
		} else {
			nextNode = new(topicTrie)
			items[v] = nextNode
			node = nextNode
		}
	}
	node.handler = handler
	if len(options) > 0 {
		node.options = options[0]
	}
}

func (t *topicTrie) findHandler(topic string) (handler topicHandler, options TopicTrieOptions, params map[string]string) {
	node, params := t, make(map[string]string)
	for _, v := range strings.Split(topic, "/") {
		nextNode, found := node.staticItems[v]
		if !found {
			for k, item := range node.paramItems {
				params[strings.Trim(k, "${}")] = v
				nextNode = item
			}
		}
		if node = nextNode; node == nil {
			break
		}
	}
	if node != nil {
		handler, options = node.handler, node.options
	}
	return
}

type MicrosecondTime struct{ time.Time }

func (e *MicrosecondTime) UnmarshalJSON(bytes []byte) (err error) {
	e.Time = castToTime(cast.ToInt64(string(bytes)))
	return
}

func castToTime(value int64) time.Time {
	if value == 0 {
		return time.Time{}
	}
	return time.Unix(0, value*int64(time.Millisecond))
}

func CastToStringValue(data []byte) string {
	value := string(data)
	if strings.HasPrefix(value, `"`) {
		value, _ = strconv.Unquote(value)
	}
	return value
}
