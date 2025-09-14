package config

import "custom-go/iot/base"

var (
    PostEvent = base.NewEventTslHandler[PostEventOutput]("post", "post", "info", "thing.event.property.post", "属性上报", true)
)

func (e *Post_PowerType) UnmarshalJSON(data []byte) (err error) {
    *e = Post_PowerType(base.CastToStringValue(data))
    return
}
func (e *Post_sleepMode) UnmarshalJSON(data []byte) (err error) {
    *e = Post_sleepMode(base.CastToStringValue(data))
    return
}
