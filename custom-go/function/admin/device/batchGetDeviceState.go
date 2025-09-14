package device

import (
	"custom-go/iot/base"
	"custom-go/pkg/plugins"
	"custom-go/pkg/types"
	iot "github.com/alibabacloud-go/iot-20180120/v6/client"
	"github.com/alibabacloud-go/tea/tea"
)

type (
	batchGetDeviceStateInput struct {
		DeviceNames []string `json:"deviceNames"`
	}
	batchGetDeviceStateOutput = map[string]*iot.BatchGetDeviceStateResponseBodyDeviceStatusListDeviceStatus
	batchGetDeviceStateBody   = *types.OperationBody[batchGetDeviceStateInput, batchGetDeviceStateOutput]
)

func init() {
	plugins.RegisterFunction[batchGetDeviceStateInput, batchGetDeviceStateOutput](batchGetDeviceState, types.OperationType_MUTATION)
}

func batchGetDeviceState(hook *types.HookRequest, body batchGetDeviceStateBody) (resp batchGetDeviceStateBody, err error) {
	data, err := base.BatchGetDeviceState(body.Input.DeviceNames...)
	if err != nil {
		return
	}

	onlineDevices := make([]string, 0, len(data))
	for _, item := range data {
		if tea.StringValue(item.Status) == string(base.DeviceStatusOnline) {
			onlineDevices = append(onlineDevices, tea.StringValue(item.DeviceName))
		}
	}
	body.ResetResponse(data)
	resp = body
	return
}
