package base

import (
	"custom-go/iot/hbs"
	"custom-go/pkg/types"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"strings"
)

type (
	EventTslHandler[O any] struct {
		Event *hbs.TslSpecEvent
	}
	EventTslOutput[O any] struct {
		IotId           string          `json:"iotId"`
		ProductKey      string          `json:"productKey"`
		DeviceName      string          `json:"deviceName"`
		GmtCreate       MicrosecondTime `json:"gmtCreate"`
		Value           *O              `json:"value"`
		CheckFailedData interface{}     `json:"checkFailedData"`
	}
)

func NewEventTslHandler[O any](identifier, name, eventType, method, desc string, required bool) *EventTslHandler[O] {
	eventHandler := &EventTslHandler[O]{Event: &hbs.TslSpecEvent{
		Identifier: identifier,
		Name:       name,
		Type:       eventType,
		Required:   required,
		Desc:       desc,
		Method:     method,
	}}
	return eventHandler
}

const eventTslTopicFormat = "/${productKey}/${deviceName}/%s"

func (e *EventTslHandler[O]) SetHandler(
	handler func(echo.Logger, *types.InternalClient, *ApplicationProperties, *EventTslOutput[O]) error,
	options ...TopicTrieOptions) {
	topic := fmt.Sprintf(eventTslTopicFormat, strings.ReplaceAll(e.Event.Method, ".", "/"))
	eventTopic.addHandler(topic, func(logger echo.Logger, client *types.InternalClient, input *topicHandlerInput) error {
		var output EventTslOutput[O]
		if err := json.Unmarshal(input.Data, &output); err != nil {
			return err
		}

		if output.Value == nil && input.ExtractedValue != nil {
			if err := json.Unmarshal(input.ExtractedValue, &output.Value); err != nil {
				return err
			}
		}
		return handler(logger, client, input.Properties, &output)
	}, options...)
}
