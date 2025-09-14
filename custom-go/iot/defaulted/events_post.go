package defaulted

import (
	"custom-go/customize"
	"custom-go/iot/base"
	"custom-go/pkg/types"
	"github.com/labstack/echo/v4"
)

func init() {
	PostEvent.SetHandler(func(logger echo.Logger, client *types.InternalClient, properties *base.ApplicationProperties, output *base.EventTslOutput[PostEventOutput]) error {
		go customize.PushEventData(PostEvent, properties, output)
		return nil
	})
}
