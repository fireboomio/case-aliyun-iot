package defaulted

import (
	"custom-go/iot/base"
	"encoding/json"
	iot "github.com/alibabacloud-go/iot-20180120/v6/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/buger/jsonparser"
	"strconv"
)

func QueryDeviceProperty(deviceName string) (data *Properties, err error) {
	statusListData, err := base.QueryDevicePropertyStatus(deviceName)
	if err != nil {
		return
	}

	data, err = translateDeviceProperties(statusListData)
	return
}

func translateDeviceProperties(data []*iot.QueryDevicePropertyStatusResponseBodyDataListPropertyStatusInfo) (properties *Properties, err error) {
	dataBytes := []byte("{}")
	for _, item := range data {
		if item.Value == nil {
			continue
		}

		var setValue []byte
		itemValue := tea.StringValue(item.Value)
		switch tea.StringValue(item.DataType) {
		case "enum", "text":
			setValue = []byte(strconv.Quote(itemValue))
		default:
			setValue = []byte(itemValue)
		}
		dataBytes, _ = jsonparser.Set(dataBytes, setValue, tea.StringValue(item.Identifier))
	}

	err = json.Unmarshal(dataBytes, &properties)
	return
}
