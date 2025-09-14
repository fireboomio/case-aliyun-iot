package base

import (
	"custom-go/iot/hbs"
	"custom-go/pkg/types"
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	iot "github.com/alibabacloud-go/iot-20180120/v6/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/labstack/echo/v4"
	"github.com/spf13/cast"
	"net/http"
	"os"
)

const (
	cloudRegionIdEnv   = "ALIBABA_CLOUD_REGION_ID"
	cloudProductKeyEnv = "ALIBABA_CLOUD_PRODUCT_KEY"

	// 用于控制是否模板生成IOT代码
	enableIotCodeGenerate = "ENABLE_IOT_CODE_GENERATE"
)

var client *Client

type Client struct {
	logger        echo.Logger
	iotClient     *iot.Client
	iotInstanceId *string
	productKey    *string
}

func init() {
	types.AddRegisteredHook(func(logger echo.Logger) {
		if cast.ToBool(os.Getenv(enableIotCodeGenerate)) {
			hbs.GenerateTslTemplate("defaulted")
			hbs.GenerateTslTemplate("config")
		}
		// 初始化配置对象 &openapi.Config 。 Config 对象存放 AccessKeyId、AccessKeySecret 、Endpoint 等配置
		apiConfig := &openapi.Config{
			RegionId: tea.String(os.Getenv(cloudRegionIdEnv)),
			// 您的AccessKey ID
			AccessKeyId: tea.String(os.Getenv(cloudAccessKeyEnv)),
			// 您的AccessKey Secret
			AccessKeySecret: tea.String(os.Getenv(cloudAccessSecretEnv)),
		}
		// 访问的域名
		apiConfig.Endpoint = tea.String(fmt.Sprintf("iot.%s.aliyuncs.com", os.Getenv(cloudRegionIdEnv)))

		// 实例化一个客户端，从 &iot.Client 类生成对象 client
		iotClient, _ := iot.NewClient(apiConfig)
		client = &Client{logger: logger, iotClient: iotClient}
		client.iotInstanceId = tea.String(os.Getenv(cloudIotInstanceIdEnv))
		client.productKey = tea.String(os.Getenv(cloudProductKeyEnv))
		logger.Debug("iot client init success")
	})
}

type BusinessError struct {
	method  string
	Code    string
	Message string
}

func (e *BusinessError) Error() string {
	return fmt.Sprintf("%s return body not success, %s", e.method, e.Message)
}

func QueryDevicePropertyStatus(deviceName string) (data []*iot.QueryDevicePropertyStatusResponseBodyDataListPropertyStatusInfo, err error) {
	req := &iot.QueryDevicePropertyStatusRequest{
		DeviceName:    tea.String(deviceName),
		IotInstanceId: client.iotInstanceId,
		ProductKey:    client.productKey,
	}
	resp, err := client.iotClient.QueryDevicePropertyStatus(req)
	if err != nil {
		return
	}
	if tea.Int32Value(resp.StatusCode) != http.StatusOK || resp.Body == nil {
		err = fmt.Errorf("QueryDevicePropertyStatus return status code: %d", tea.Int32Value(resp.StatusCode))
		return
	}
	if !tea.BoolValue(resp.Body.Success) {
		err = &BusinessError{method: "QueryDevicePropertyStatus", Code: tea.StringValue(resp.Body.Code), Message: tea.StringValue(resp.Body.ErrorMessage)}
		return
	}
	statusData := resp.Body.Data
	if statusData == nil || statusData.List == nil {
		err = fmt.Errorf("QueryDevicePropertyStatus return body statusDataList is empty")
		return
	}

	data = statusData.List.PropertyStatusInfo
	return
}

func SetDeviceDesiredProperty(deviceName string, items string) (err error) {
	req := &iot.SetDeviceDesiredPropertyRequest{
		DeviceName:    tea.String(deviceName),
		IotInstanceId: client.iotInstanceId,
		ProductKey:    client.productKey,
		Items:         tea.String(items),
		Versions:      tea.String("{}"),
	}
	resp, err := client.iotClient.SetDeviceDesiredProperty(req)
	if err != nil {
		return
	}
	if tea.Int32Value(resp.StatusCode) != http.StatusOK || resp.Body == nil {
		err = fmt.Errorf("SetDeviceDesiredProperty return status code: %d", tea.Int32Value(resp.StatusCode))
		return
	}
	if !tea.BoolValue(resp.Body.Success) {
		err = &BusinessError{method: "SetDeviceDesiredProperty", Code: tea.StringValue(resp.Body.Code), Message: tea.StringValue(resp.Body.ErrorMessage)}
		return
	}

	return
}

func GetDeviceStatus(deviceName string) (data *iot.GetDeviceStatusResponseBodyData, err error) {
	req := &iot.GetDeviceStatusRequest{
		DeviceName:    tea.String(deviceName),
		IotInstanceId: client.iotInstanceId,
		ProductKey:    client.productKey,
	}
	resp, err := client.iotClient.GetDeviceStatus(req)
	if err != nil {
		return
	}
	if tea.Int32Value(resp.StatusCode) != http.StatusOK || resp.Body == nil {
		err = fmt.Errorf("GetDeviceStatus return status code: %d", tea.Int32Value(resp.StatusCode))
		return
	}
	if !tea.BoolValue(resp.Body.Success) {
		err = &BusinessError{method: "GetDeviceStatus", Code: tea.StringValue(resp.Body.Code), Message: tea.StringValue(resp.Body.ErrorMessage)}
		return
	}

	data = resp.Body.Data
	return
}

func BatchGetDeviceState(deviceName ...string) (data map[string]*iot.BatchGetDeviceStateResponseBodyDeviceStatusListDeviceStatus, err error) {
	req := &iot.BatchGetDeviceStateRequest{
		DeviceName:    tea.StringSlice(deviceName),
		IotInstanceId: client.iotInstanceId,
		ProductKey:    client.productKey,
	}
	resp, err := client.iotClient.BatchGetDeviceState(req)
	if err != nil {
		return
	}
	if tea.Int32Value(resp.StatusCode) != http.StatusOK || resp.Body == nil {
		err = fmt.Errorf("BatchGetDeviceState return status code: %d", tea.Int32Value(resp.StatusCode))
		return
	}
	if !tea.BoolValue(resp.Body.Success) {
		err = &BusinessError{method: "BatchGetDeviceState", Code: tea.StringValue(resp.Body.Code), Message: tea.StringValue(resp.Body.ErrorMessage)}
		return
	}
	statusList := resp.Body.DeviceStatusList
	if statusList == nil || len(statusList.DeviceStatus) == 0 {
		err = fmt.Errorf("BatchGetDeviceState return body deviceStatusList is empty")
		return
	}

	data = make(map[string]*iot.BatchGetDeviceStateResponseBodyDeviceStatusListDeviceStatus)
	for _, item := range statusList.DeviceStatus {
		data[tea.StringValue(item.DeviceName)] = item
	}
	return
}

func QueryDeviceDetail(deviceName string) (data *iot.QueryDeviceDetailResponseBodyData, err error) {
	req := &iot.QueryDeviceDetailRequest{
		DeviceName:    tea.String(deviceName),
		IotInstanceId: client.iotInstanceId,
		ProductKey:    client.productKey,
	}
	resp, err := client.iotClient.QueryDeviceDetail(req)
	if err != nil {
		return
	}
	if tea.Int32Value(resp.StatusCode) != http.StatusOK || resp.Body == nil {
		err = fmt.Errorf("QueryDeviceDetail return status code: %d", tea.Int32Value(resp.StatusCode))
		return
	}
	if !tea.BoolValue(resp.Body.Success) {
		err = &BusinessError{method: "QueryDeviceDetail", Code: tea.StringValue(resp.Body.Code), Message: tea.StringValue(resp.Body.ErrorMessage)}
		return
	}

	data = resp.Body.Data
	return
}

func RegisterDevice(deviceName, nickName string) (data *iot.RegisterDeviceResponseBodyData, err error) {
	if nickName == "" {
		nickName = "auto_create_by_api"
	}
	req := &iot.RegisterDeviceRequest{
		DeviceName:    tea.String(deviceName),
		Nickname:      tea.String(nickName),
		IotInstanceId: client.iotInstanceId,
		ProductKey:    client.productKey,
	}
	resp, err := client.iotClient.RegisterDevice(req)
	if err != nil {
		return
	}
	if tea.Int32Value(resp.StatusCode) != http.StatusOK || resp.Body == nil {
		err = fmt.Errorf("RegisterDevice return status code: %d", tea.Int32Value(resp.StatusCode))
		return
	}
	if !tea.BoolValue(resp.Body.Success) {
		err = &BusinessError{method: "RegisterDevice", Code: tea.StringValue(resp.Body.Code), Message: tea.StringValue(resp.Body.ErrorMessage)}
		return
	}

	data = resp.Body.Data
	return
}

func UpdateDeviceNickname(deviceName, nickName string) (err error) {
	req := &iot.BatchUpdateDeviceNicknameRequest{
		DeviceNicknameInfo: []*iot.BatchUpdateDeviceNicknameRequestDeviceNicknameInfo{{
			DeviceName: tea.String(deviceName),
			Nickname:   tea.String(nickName),
			ProductKey: client.productKey,
		}},
		IotInstanceId: client.iotInstanceId,
	}
	resp, err := client.iotClient.BatchUpdateDeviceNickname(req)
	if err != nil {
		return
	}
	if tea.Int32Value(resp.StatusCode) != http.StatusOK || resp.Body == nil {
		err = fmt.Errorf("UpdateDeviceNickname return status code: %d", tea.Int32Value(resp.StatusCode))
		return
	}
	if !tea.BoolValue(resp.Body.Success) {
		err = &BusinessError{method: "UpdateDeviceNickname", Code: tea.StringValue(resp.Body.Code), Message: tea.StringValue(resp.Body.ErrorMessage)}
		return
	}

	return
}
