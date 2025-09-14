package defaulted

import "custom-go/iot/base"

var (
    PostEvent = base.NewEventTslHandler[PostEventOutput]("post", "post", "info", "thing.event.property.post", "属性上报", true)
    ErrorReportEvent = base.NewEventTslHandler[ErrorReportEventOutput]("ErrorReport", "故障上报", "error", "thing.event.ErrorReport.post", "故障和故障恢复上报", false)
    LowBatteyEventEvent = base.NewEventTslHandler[LowBatteyEventEventOutput]("LowBatteyEvent", "电量低告警", "alert", "thing.event.LowBatteyEvent.post", "低电量告警电压低于阀值时触发，高于阀值时解除", false)
)

func (e *ErrorReport_type) UnmarshalJSON(data []byte) (err error) {
    *e = ErrorReport_type(base.CastToStringValue(data))
    return
}
