package device

import (
	"custom-go/iot/base"
	"custom-go/pkg/plugins"
	"custom-go/pkg/types"
	iot "github.com/alibabacloud-go/iot-20180120/v6/client"
)

type (
	getDeviceStatusInput struct {
		DeviceName string `json:"deviceName"`
	}
	getDeviceStatusOutput = iot.GetDeviceStatusResponseBodyData
	getDeviceStatusBody   = *types.OperationBody[getDeviceStatusInput, getDeviceStatusOutput]
)

func init() {
	plugins.RegisterFunction[getDeviceStatusInput, getDeviceStatusOutput](getDeviceStatus, types.OperationType_MUTATION)
}

func getDeviceStatus(hook *types.HookRequest, body getDeviceStatusBody) (resp getDeviceStatusBody, err error) {
	data, err := base.GetDeviceStatus(body.Input.DeviceName)
	if err != nil {
		return
	}

	body.ResetResponse(*data)
	resp = body
	return
}
