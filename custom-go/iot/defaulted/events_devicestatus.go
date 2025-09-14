package defaulted

import (
	"custom-go/iot/base"
	"custom-go/iot/config"
	"custom-go/pkg/types"
	iot "github.com/alibabacloud-go/iot-20180120/v6/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/labstack/echo/v4"
	"strings"
)

func init() {
	base.SetDeviceStatusHandler(func(logger echo.Logger, client *types.InternalClient, properties *base.ApplicationProperties, output *base.EventDeviceStatusOutput) (err error) {
		// todo: 查询业务数据

		// 修复离线上报异常的问题
		if output.Status == base.DeviceStatusOffline {
			var statusResp *iot.GetDeviceStatusResponseBodyData
			if statusResp, err = base.GetDeviceStatus(output.DeviceName); err != nil {
				return
			}
			if strings.EqualFold(tea.StringValue(statusResp.Status), string(base.DeviceStatusOnline)) {
				output.Status = base.DeviceStatusOnline
			}
		}
		configProperties := &config.Properties{
			// todo: 根据业务需要设置设备属性
		}
		err = config.SetDeviceDesiredProperty(output.DeviceName, configProperties)
		return
	})
}
