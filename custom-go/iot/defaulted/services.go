package defaulted

import "custom-go/iot/base"

var (
    SetService = base.NewRRpcRequest[SetServiceInput, SetServiceOutput]("set", "set", "async", "thing.service.property.set", "属性设置", true)
    GetService = base.NewRRpcRequest[GetServiceInput, GetServiceOutput]("get", "get", "async", "thing.service.property.get", "属性获取", true)
    OtaService = base.NewRRpcRequest[OtaServiceInput, OtaServiceOutput]("ota", "ota升级", "sync", "thing.service.ota", "ota升级", false)
    RebootService = base.NewRRpcRequest[RebootServiceInput, RebootServiceOutput]("Reboot", "重启", "sync", "thing.service.Reboot", "特殊情况下，重启设备", false)
    RefreshPropertiesService = base.NewRRpcRequest[RefreshPropertiesServiceInput, RefreshPropertiesServiceOutput]("refreshProperties", "主动请求设备上报属性", "sync", "thing.service.refreshProperties", "", false)
    BindBluetoothService = base.NewRRpcRequest[BindBluetoothServiceInput, BindBluetoothServiceOutput]("bindBluetooth", "蓝牙绑定/解绑", "sync", "thing.service.bindBluetooth", "绑定蓝牙设备mac为空意味着解绑", false)
)