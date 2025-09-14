package config

import (
	"custom-go/iot/base"
	"encoding/json"
	"fmt"
	"github.com/buger/jsonparser"
	"golang.org/x/exp/slices"
)

type properties Properties

func (t *Properties) MarshalJSON() ([]byte, error) {
	bytes, err := json.Marshal(properties(*t))
	if err != nil {
		return nil, err
	}

	newBytes := []byte("{}")
	err = jsonparser.ObjectEach(bytes, func(key []byte, value []byte, _ jsonparser.ValueType, _ int) error {
		newBytes, err = jsonparser.Set(newBytes, value, fmt.Sprintf("config:%s", string(key)))
		return err
	})
	return newBytes, err
}

func SetDeviceDesiredProperty(deviceName string, data *Properties) (err error) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return
	}

	err = base.SetDeviceDesiredProperty(deviceName, string(dataBytes))
	return
}

func GetSleepTimePeriod(period ...*TimePeriod) (sleepTimePeriod []string) {
	var timePeriod, excludePeriod []string
	for _, item := range period {
		if item.End <= item.Start {
			continue
		}
		timePeriod = append(timePeriod, item.Start, item.End)
		excludePeriod = append(excludePeriod, item.getTimePeriod())
	}
	if len(timePeriod) == 0 {
		return []string{}
	}

	slices.Sort(timePeriod)
	if timePeriod[0] != dailyStartTime {
		timePeriod = append([]string{dailyStartTime}, timePeriod...)
	}
	if timePeriod[len(timePeriod)-1] != dailyEndTime {
		timePeriod = append(timePeriod, dailyEndTime)
	}
	for i := 0; i < len(timePeriod)-1; i++ {
		start, end := timePeriod[i], timePeriod[i+1]
		if end <= start {
			continue
		}
		itemPeriod := formatSleepTimePeriod(start, end)
		if !slices.Contains(excludePeriod, itemPeriod) {
			sleepTimePeriod = append(sleepTimePeriod, itemPeriod)
		}
	}
	return
}

const (
	dailyStartTime = "00:00"
	dailyEndTime   = "24:00"
)

type TimePeriod struct {
	Start string
	End   string
}

func NewTimePeriod(start, end string) *TimePeriod {
	return &TimePeriod{Start: start, End: end}
}

func (t *TimePeriod) getTimePeriod() string {
	return formatSleepTimePeriod(t.Start, t.End)
}

func formatSleepTimePeriod(start, end string) string {
	return fmt.Sprintf("%s-%s", start, end)
}
