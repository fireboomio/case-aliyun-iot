package device

import (
	"custom-go/iot/defaulted"
	"custom-go/pkg/plugins"
	"custom-go/pkg/types"
)

type (
	queryDevicePropertyInput struct {
		DeviceName string `json:"deviceName"`
	}
	queryDevicePropertyOutput = defaulted.Properties
	queryDevicePropertyBody   = *types.OperationBody[queryDevicePropertyInput, queryDevicePropertyOutput]
)

func init() {
	plugins.RegisterFunction[queryDevicePropertyInput, queryDevicePropertyOutput](queryDeviceProperty, types.OperationType_MUTATION)
}

func queryDeviceProperty(hook *types.HookRequest, body queryDevicePropertyBody) (resp queryDevicePropertyBody, err error) {
	data, err := defaulted.QueryDeviceProperty(body.Input.DeviceName)
	if err != nil {
		return
	}

	body.ResetResponse(*data)
	resp = body
	return
}
