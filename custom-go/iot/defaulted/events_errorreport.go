package defaulted

import (
	"custom-go/iot/base"
	"custom-go/pkg/types"
	"github.com/labstack/echo/v4"
)

func init() {
	ErrorReportEvent.SetHandler(func(logger echo.Logger, client *types.InternalClient, properties *base.ApplicationProperties, output *base.EventTslOutput[ErrorReportEventOutput]) (err error) {
		// todo: 错误上报
		return
	})
}
