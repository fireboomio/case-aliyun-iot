package config

import "custom-go/iot/base"

var (
    SetService = base.NewRRpcRequest[SetServiceInput, SetServiceOutput]("set", "set", "async", "thing.service.property.set", "属性设置", true)
    GetService = base.NewRRpcRequest[GetServiceInput, GetServiceOutput]("get", "get", "async", "thing.service.property.get", "属性获取", true)
)