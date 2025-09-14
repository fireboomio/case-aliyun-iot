package defaulted

import (
	"errors"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/spf13/cast"
)

func RefreshProperties(deviceName string) (err error) {
	refreshPropertiesServiceInput := &RefreshPropertiesServiceInput{MessageID: tea.String(maxMessageId)}
	refreshPropertiesServiceOutput, err := RefreshPropertiesService.RRpc(deviceName, maxMessageId, refreshPropertiesServiceInput)
	if err != nil {
		return
	}
	if !cast.ToBool(tea.StringValue(refreshPropertiesServiceOutput.Result)) {
		err = errors.New("刷新属性失败")
		return
	}
	return
}
