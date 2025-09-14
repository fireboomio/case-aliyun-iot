package base

import (
	"custom-go/pkg/types"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"time"
)

type EventDeviceStatusOutput struct {
	Status      DeviceStatus `json:"status"`
	IotId       string       `json:"iotId"`
	ProductKey  string       `json:"productKey"`
	DeviceName  string       `json:"deviceName"`
	UtcTime     time.Time    `json:"utcTime"`
	UtcLastTime time.Time    `json:"utcLastTime"`
	ClientIp    string       `json:"clientIp"`
}

type DeviceStatus string

const (
	DeviceStatusOnline  DeviceStatus = "online"
	DeviceStatusOffline DeviceStatus = "offline"
)

const eventDeviceStatusTopic = "/as/mqtt/status/${productKey}/${deviceName}"

func SetDeviceStatusHandler(handler func(echo.Logger, *types.InternalClient, *ApplicationProperties, *EventDeviceStatusOutput) error) {
	eventTopic.addHandler(eventDeviceStatusTopic, func(logger echo.Logger, client *types.InternalClient, input *topicHandlerInput) error {
		var output EventDeviceStatusOutput
		if err := json.Unmarshal(input.Data, &output); err != nil {
			return err
		}

		return handler(logger, client, input.Properties, &output)
	})
}
