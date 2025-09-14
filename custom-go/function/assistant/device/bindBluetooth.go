package device

import (
	"custom-go/iot/defaulted"
	"custom-go/pkg/plugins"
	"custom-go/pkg/types"
)

type (
	bindBluetoothInput struct {
		DeviceName string `json:"deviceName"`
		Mac        string `json:"mac"`
	}
	bindBluetoothOutput struct {
		Success bool `json:"success"`
	}
	bindBluetoothBody = *types.OperationBody[bindBluetoothInput, bindBluetoothOutput]
)

func init() {
	plugins.RegisterFunction[bindBluetoothInput, bindBluetoothOutput](bindBluetooth, types.OperationType_MUTATION)
}

func bindBluetooth(hook *types.HookRequest, body bindBluetoothBody) (resp bindBluetoothBody, err error) {
	if err = defaulted.BindBluetooth(body.Input.DeviceName, body.Input.Mac); err != nil {
		return
	}

	body.ResetResponse(bindBluetoothOutput{Success: true})
	resp = body
	return
}
