package base

import (
	"custom-go/iot/hbs"
	"encoding/base64"
	"encoding/json"
	"fmt"
	iot "github.com/alibabacloud-go/iot-20180120/v6/client"
	"github.com/alibabacloud-go/tea/tea"
	"net/http"
)

type (
	RRpcInput[I any] struct {
		Method  string `json:"method"`
		Id      string `json:"id"`
		Params  *I     `json:"params"`
		Version string `json:"version"`
	}
	RRpcOutput[O any] struct {
		Params *O `json:"params"`
	}
	RRpcRequest[I, O any] struct {
		service *hbs.TslSpecService
	}
)

func NewRRpcRequest[I, O any](identifier, name, callType, method, desc string, required bool) *RRpcRequest[I, O] {
	return &RRpcRequest[I, O]{&hbs.TslSpecService{
		Identifier: identifier,
		Name:       name,
		Required:   required,
		CallType:   callType,
		Desc:       desc,
		Method:     method,
	}}
}

func (r *RRpcRequest[I, O]) RRpc(deviceName, id string, params *I) (data *O, err error) {
	input := &RRpcInput[I]{Id: id, Method: r.service.Method, Params: params, Version: "1.0.0"}
	return RRpc[I, O](deviceName, input)
}

func RRpc[I, O any](deviceName string, input *RRpcInput[I]) (data *O, err error) {
	inputBytes, err := json.Marshal(input)
	if err != nil {
		return
	}

	req := &iot.RRpcRequest{
		DeviceName:        tea.String(deviceName),
		IotInstanceId:     client.iotInstanceId,
		ProductKey:        client.productKey,
		RequestBase64Byte: tea.String(base64.StdEncoding.EncodeToString(inputBytes)),
		Timeout:           tea.Int32(5000),
	}
	resp, err := client.iotClient.RRpc(req)
	if err != nil {
		return
	}
	if tea.Int32Value(resp.StatusCode) != http.StatusOK || resp.Body == nil {
		err = fmt.Errorf("RRpc return status code: %d", tea.Int32Value(resp.StatusCode))
		return
	}
	if !tea.BoolValue(resp.Body.Success) {
		err = &BusinessError{method: "RRpc", Code: tea.StringValue(resp.Body.Code), Message: tea.StringValue(resp.Body.ErrorMessage)}
		return
	}
	payloadBytes, err := base64.StdEncoding.DecodeString(tea.StringValue(resp.Body.PayloadBase64Byte))
	if err != nil {
		return
	}

	var output RRpcOutput[O]
	if err = json.Unmarshal(payloadBytes, &output); err != nil {
		return
	}

	data = output.Params
	return
}

func FormatMessageId(index int64) string {
	return fmt.Sprintf("%08d", index)
}
