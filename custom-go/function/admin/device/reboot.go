package device

import (
	"custom-go/iot/defaulted"
	"custom-go/pkg/plugins"
	"custom-go/pkg/types"
)

type (
	rebootInput struct {
		DeviceName string `json:"deviceName"`
	}
	rebootOutput struct {
		Success bool `json:"success"`
	}
	rebootBody = *types.OperationBody[rebootInput, rebootOutput]
)

func init() {
	plugins.RegisterFunction[rebootInput, rebootOutput](reboot, types.OperationType_MUTATION)
}

func reboot(_ *types.HookRequest, body rebootBody) (resp rebootBody, err error) {
	if err = defaulted.Reboot(body.Input.DeviceName); err != nil {
		return
	}

	body.ResetResponse(rebootOutput{Success: true})
	resp = body
	return
}
