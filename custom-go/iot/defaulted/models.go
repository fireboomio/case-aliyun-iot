package defaulted



type ErrorReportEventOutput struct {
    /* 故障设备 */
    Type *ErrorReport_type `json:"type,omitempty"`
    /* 是故障还是故障恢复 */
    Error *ErrorReport_error `json:"error,omitempty"`
}

type LowBatteyEventEventOutput struct {
    /* 电压【{"min":"0","max":"100000","step":"1"}】 */
    BatteryLevel *int64 `json:"BatteryLevel,omitempty"`
    /* trigger */
    Trigger *LowBatteyEvent_trigger `json:"trigger,omitempty"`
}

type RebootServiceInput struct {
}

type RebootServiceOutput struct {
    /* result【{"length":"10240"}】 */
    Result *string `json:"result,omitempty"`
}

type RunningState struct {
    /* 蓝牙连接状态 */
    Bluetooth *Properties_bluetooth `json:"bluetooth,omitempty"`
}

type BindBluetoothServiceInput struct {
    /* mac地址【{"length":"100"}】 */
    Mac *string `json:"mac,omitempty"`
}

type BindBluetoothServiceOutput struct {
}

type GetServiceInput struct {
}

type GetServiceOutput struct {
    /* 电池电压【{"min":"0","max":"80000","step":"1","unit":"mV","unitName":"毫伏"}】 */
    BatteryVoltage *int64 `json:"BatteryVoltage,omitempty"`
    /* 设备信号强度【{"min":"0","max":"100","step":"1"}】 */
    NetState *int64 `json:"NetState,omitempty"`
    /* 设备固件版本【{"length":"64"}】 */
    FirmwareVersion *string `json:"FirmwareVersion,omitempty"`
    /* 设备状态 */
    RunningState *RunningState `json:"RunningState,omitempty"`
}

type OtaServiceInput struct {
    /* aa【{"min":"-2147483648","max":"2147483647","step":"1"}】 */
    Aa *int64 `json:"aa,omitempty"`
    /* bc【{"min":"0","max":"4294967295","step":"1"}】 */
    Bc *float64 `json:"bc,omitempty"`
    /* bd【{"length":"10240"}】 */
    Bd *string `json:"bd,omitempty"`
    /* bx【{"min":"-2147483648","max":"2147483647","step":"1"}】 */
    Bx *int64 `json:"bx,omitempty"`
    /* ai【{"length":"32"}】 */
    Ai *string `json:"ai,omitempty"`
    /* am【{"min":"-2147483648","max":"2147483647","step":"1"}】 */
    Am *int64 `json:"am,omitempty"`
}

type OtaServiceOutput struct {
    /* result【{"length":"10240"}】 */
    Result *string `json:"result,omitempty"`
}

type PostEventOutput struct {
    /* 电池电压【{"min":"0","max":"80000","step":"1","unit":"mV","unitName":"毫伏"}】 */
    BatteryVoltage *int64 `json:"BatteryVoltage,omitempty"`
    /* 设备信号强度【{"min":"0","max":"100","step":"1"}】 */
    NetState *int64 `json:"NetState,omitempty"`
    /* 设备固件版本【{"length":"64"}】 */
    FirmwareVersion *string `json:"FirmwareVersion,omitempty"`
    /* 设备状态 */
    RunningState *RunningState `json:"RunningState,omitempty"`
}

type Properties struct {
    /* 电池电压[accessMode:r](电池电压。单位为mv)【{"min":"0","max":"80000","step":"1","unit":"mV","unitName":"毫伏"}】 */
    BatteryVoltage *int64 `json:"BatteryVoltage,omitempty"`
    /* 设备信号强度[accessMode:r](设备4G信号强度，单位dBm)【{"min":"0","max":"100","step":"1"}】 */
    NetState *int64 `json:"NetState,omitempty"`
    /* 设备固件版本[accessMode:rw]【{"length":"64"}】 */
    FirmwareVersion *string `json:"FirmwareVersion,omitempty"`
    /* 设备状态[accessMode:rw](设备运行状态，2分钟上报1次
后续把 开关状态、烟雾、红外、网络、电压、蓝牙绑定等都放在这里) */
    RunningState *RunningState `json:"RunningState,omitempty"`
}

type RefreshPropertiesServiceInput struct {
    /* messageID【{"length":"10240"}】 */
    MessageID *string `json:"messageID,omitempty"`
}

type RefreshPropertiesServiceOutput struct {
    /* result【{"length":"10240"}】 */
    Result *string `json:"result,omitempty"`
}

type SetServiceInput struct {
    /* 设备固件版本【{"length":"64"}】 */
    FirmwareVersion *string `json:"FirmwareVersion,omitempty"`
    /* 设备状态 */
    RunningState *RunningState `json:"RunningState,omitempty"`
}

type SetServiceOutput struct {
}

type ErrorReport_error int64
var (
    // ErrorReport_error_0 恢复
    ErrorReport_error_0 ErrorReport_error = 0
    // ErrorReport_error_1 触发
    ErrorReport_error_1 ErrorReport_error = 1
)
var ErrorReport_errorDescMap = map[ErrorReport_error]string{
    ErrorReport_error_0: "恢复",
    ErrorReport_error_1: "触发",
    
}
func (e *ErrorReport_error) String() string {
    return ErrorReport_errorDescMap[*e]
}

type ErrorReport_type string
var (
    // ErrorReport_type_1 电机故障
    ErrorReport_type_1 ErrorReport_type = "1"
)
var ErrorReport_typeDescMap = map[ErrorReport_type]string{
    ErrorReport_type_1: "电机故障",
    
}
func (e *ErrorReport_type) String() string {
    return ErrorReport_typeDescMap[*e]
}

type LowBatteyEvent_trigger int64
var (
    // LowBatteyEvent_trigger_0 解除
    LowBatteyEvent_trigger_0 LowBatteyEvent_trigger = 0
    // LowBatteyEvent_trigger_1 触发
    LowBatteyEvent_trigger_1 LowBatteyEvent_trigger = 1
)
var LowBatteyEvent_triggerDescMap = map[LowBatteyEvent_trigger]string{
    LowBatteyEvent_trigger_0: "解除",
    LowBatteyEvent_trigger_1: "触发",
    
}
func (e *LowBatteyEvent_trigger) String() string {
    return LowBatteyEvent_triggerDescMap[*e]
}

type Properties_bluetooth int64
var (
    // Properties_bluetooth_0 未连接
    Properties_bluetooth_0 Properties_bluetooth = 0
    // Properties_bluetooth_1 已连接
    Properties_bluetooth_1 Properties_bluetooth = 1
)
var Properties_bluetoothDescMap = map[Properties_bluetooth]string{
    Properties_bluetooth_0: "未连接",
    Properties_bluetooth_1: "已连接",
    
}
func (e *Properties_bluetooth) String() string {
    return Properties_bluetoothDescMap[*e]
}
