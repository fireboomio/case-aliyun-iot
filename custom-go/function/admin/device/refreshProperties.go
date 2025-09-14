package device

import (
	"custom-go/iot/defaulted"
	"custom-go/pkg/plugins"
	"custom-go/pkg/types"
)

type (
	refreshPropertiesInput struct {
		DeviceName string `json:"deviceName"`
	}
	refreshPropertiesOutput struct {
		Success bool `json:"success"`
	}
	refreshPropertiesBody = *types.OperationBody[refreshPropertiesInput, refreshPropertiesOutput]
)

func init() {
	plugins.RegisterFunction[refreshPropertiesInput, refreshPropertiesOutput](refreshProperties, types.OperationType_MUTATION)
}

func refreshProperties(_ *types.HookRequest, body refreshPropertiesBody) (resp refreshPropertiesBody, err error) {
	if err = defaulted.RefreshProperties(body.Input.DeviceName); err != nil {
		return
	}

	body.ResetResponse(refreshPropertiesOutput{Success: true})
	resp = body
	return
}
