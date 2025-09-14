package defaulted

import (
	"errors"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/spf13/cast"
)

func Reboot(deviceName string) (err error) {
	rebootServiceOutput, err := RebootService.RRpc(deviceName, maxMessageId, &RebootServiceInput{})
	if err != nil {
		return
	}
	if !cast.ToBool(tea.StringValue(rebootServiceOutput.Result)) {
		err = errors.New("重启失败")
		return
	}
	return
}
