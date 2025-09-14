package defaulted

import (
	"github.com/alibabacloud-go/tea/tea"
)

const maxMessageId = "90000000"

func BindBluetooth(deviceName, mac string) (err error) {
	bindBluetoothServiceInput := &BindBluetoothServiceInput{Mac: tea.String(mac)}
	_, err = BindBluetoothService.RRpc(deviceName, maxMessageId, bindBluetoothServiceInput)
	return
}
