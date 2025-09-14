package device

import (
	"custom-go/iot/base"
	"custom-go/pkg/plugins"
	"custom-go/pkg/types"
	"errors"
	"github.com/alibabacloud-go/tea/tea"
	"strings"
)

type (
	deviceSecretInput struct {
		DeviceName string `json:"deviceName"`
		AshcanCode string `json:"ashcanCode,omitempty"`
	}
	deviceSecretOutput struct {
		Secret string `json:"secret"`
	}
	deviceSecretBody = *types.OperationBody[deviceSecretInput, deviceSecretOutput]
)

func init() {
	plugins.RegisterFunction[deviceSecretInput, deviceSecretOutput](deviceSecret, types.OperationType_MUTATION)
}

func deviceSecret(_ *types.HookRequest, body deviceSecretBody) (resp deviceSecretBody, err error) {
	var (
		secret      *string
		businessErr *base.BusinessError
	)
	defer func() {
		if secret != nil {
			body.ResetResponse(deviceSecretOutput{Secret: tea.StringValue(secret)})
		}
	}()
	deviceName := body.Input.DeviceName
	queryData, err := base.QueryDeviceDetail(deviceName)
	if errors.As(err, &businessErr) && businessErr.Code == "iot.device.NotExistedDevice" {
		err = nil
	}
	if err != nil {
		return
	}
	if queryData != nil && queryData.DeviceSecret != nil {
		secret = queryData.DeviceSecret
		return
	}

	registerData, err := base.RegisterDevice(deviceName, strings.ReplaceAll(body.Input.AshcanCode, "-", "_"))
	if err != nil {
		return
	}

	secret = registerData.DeviceSecret
	return body, nil
}
