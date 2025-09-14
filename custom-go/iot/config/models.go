package config



type Bluetooth struct {
    /* 蓝牙mac地址【{"length":"10240"}】 */
    Mac *string `json:"mac,omitempty"`
}

type GetServiceInput struct {
}

type GetServiceOutput struct {
    /* 电压触发阀值【{"min":"-2147483648","max":"2147483647","step":"1"}】 */
    ThresholdOfBatteryVoltage *int64 `json:"thresholdOfBatteryVoltage,omitempty"`
    /* 电源类型 */
    PowerType *Get_PowerType `json:"PowerType,omitempty"`
    /* 休眠模式 */
    SleepMode *Get_sleepMode `json:"sleepMode,omitempty"`
    /* 休眠时间段 */
    SleepTimePeriod []string `json:"sleepTimePeriod,omitempty"`
    /* 蓝牙相关 */
    Bluetooth *Bluetooth `json:"bluetooth,omitempty"`
}

type PostEventOutput struct {
    /* 电压触发阀值【{"min":"-2147483648","max":"2147483647","step":"1"}】 */
    ThresholdOfBatteryVoltage *int64 `json:"thresholdOfBatteryVoltage,omitempty"`
    /* 电源类型 */
    PowerType *Post_PowerType `json:"PowerType,omitempty"`
    /* 休眠模式 */
    SleepMode *Post_sleepMode `json:"sleepMode,omitempty"`
    /* 休眠时间段 */
    SleepTimePeriod []string `json:"sleepTimePeriod,omitempty"`
    /* 蓝牙相关 */
    Bluetooth *Bluetooth `json:"bluetooth,omitempty"`
}

type Properties struct {
    /* 电压触发阀值[accessMode:rw]【{"min":"-2147483648","max":"2147483647","step":"1"}】 */
    ThresholdOfBatteryVoltage *int64 `json:"thresholdOfBatteryVoltage,omitempty"`
    /* 电源类型[accessMode:rw] */
    PowerType *Properties_PowerType `json:"PowerType,omitempty"`
    /* 休眠模式[accessMode:rw] */
    SleepMode *Properties_sleepMode `json:"sleepMode,omitempty"`
    /* 休眠时间段[accessMode:rw](["00:00-07:00","09:00-19:00","21:00-24:00"]) */
    SleepTimePeriod []string `json:"sleepTimePeriod,omitempty"`
    /* 蓝牙相关[accessMode:rw](后续蓝牙的开关，也可以加在这里) */
    Bluetooth *Bluetooth `json:"bluetooth,omitempty"`
}

type SetServiceInput struct {
    /* 电压触发阀值【{"min":"-2147483648","max":"2147483647","step":"1"}】 */
    ThresholdOfBatteryVoltage *int64 `json:"thresholdOfBatteryVoltage,omitempty"`
    /* 电源类型 */
    PowerType *Set_PowerType `json:"PowerType,omitempty"`
    /* 休眠模式 */
    SleepMode *Set_sleepMode `json:"sleepMode,omitempty"`
    /* 休眠时间段 */
    SleepTimePeriod []string `json:"sleepTimePeriod,omitempty"`
    /* 蓝牙相关 */
    Bluetooth *Bluetooth `json:"bluetooth,omitempty"`
}

type SetServiceOutput struct {
}

type Get_PowerType string
var (
    // Get_PowerType_1 电池
    Get_PowerType_1 Get_PowerType = "1"
    // Get_PowerType_2 电源插电
    Get_PowerType_2 Get_PowerType = "2"
)
var Get_PowerTypeDescMap = map[Get_PowerType]string{
    Get_PowerType_1: "电池",
    Get_PowerType_2: "电源插电",
    
}
func (e *Get_PowerType) String() string {
    return Get_PowerTypeDescMap[*e]
}

type Get_sleepMode string
var (
    // Get_sleepMode_1 浅层休眠
    Get_sleepMode_1 Get_sleepMode = "1"
    // Get_sleepMode_2 深度休眠
    Get_sleepMode_2 Get_sleepMode = "2"
)
var Get_sleepModeDescMap = map[Get_sleepMode]string{
    Get_sleepMode_1: "浅层休眠",
    Get_sleepMode_2: "深度休眠",
    
}
func (e *Get_sleepMode) String() string {
    return Get_sleepModeDescMap[*e]
}

type Post_PowerType string
var (
    // Post_PowerType_1 电池
    Post_PowerType_1 Post_PowerType = "1"
    // Post_PowerType_2 电源插电
    Post_PowerType_2 Post_PowerType = "2"
)
var Post_PowerTypeDescMap = map[Post_PowerType]string{
    Post_PowerType_1: "电池",
    Post_PowerType_2: "电源插电",
    
}
func (e *Post_PowerType) String() string {
    return Post_PowerTypeDescMap[*e]
}

type Post_sleepMode string
var (
    // Post_sleepMode_1 浅层休眠
    Post_sleepMode_1 Post_sleepMode = "1"
    // Post_sleepMode_2 深度休眠
    Post_sleepMode_2 Post_sleepMode = "2"
)
var Post_sleepModeDescMap = map[Post_sleepMode]string{
    Post_sleepMode_1: "浅层休眠",
    Post_sleepMode_2: "深度休眠",
    
}
func (e *Post_sleepMode) String() string {
    return Post_sleepModeDescMap[*e]
}

type Properties_PowerType string
var (
    // Properties_PowerType_1 电池
    Properties_PowerType_1 Properties_PowerType = "1"
    // Properties_PowerType_2 电源插电
    Properties_PowerType_2 Properties_PowerType = "2"
)
var Properties_PowerTypeDescMap = map[Properties_PowerType]string{
    Properties_PowerType_1: "电池",
    Properties_PowerType_2: "电源插电",
    
}
func (e *Properties_PowerType) String() string {
    return Properties_PowerTypeDescMap[*e]
}

type Properties_sleepMode string
var (
    // Properties_sleepMode_1 浅层休眠
    Properties_sleepMode_1 Properties_sleepMode = "1"
    // Properties_sleepMode_2 深度休眠
    Properties_sleepMode_2 Properties_sleepMode = "2"
)
var Properties_sleepModeDescMap = map[Properties_sleepMode]string{
    Properties_sleepMode_1: "浅层休眠",
    Properties_sleepMode_2: "深度休眠",
    
}
func (e *Properties_sleepMode) String() string {
    return Properties_sleepModeDescMap[*e]
}

type Set_PowerType string
var (
    // Set_PowerType_1 电池
    Set_PowerType_1 Set_PowerType = "1"
    // Set_PowerType_2 电源插电
    Set_PowerType_2 Set_PowerType = "2"
)
var Set_PowerTypeDescMap = map[Set_PowerType]string{
    Set_PowerType_1: "电池",
    Set_PowerType_2: "电源插电",
    
}
func (e *Set_PowerType) String() string {
    return Set_PowerTypeDescMap[*e]
}

type Set_sleepMode string
var (
    // Set_sleepMode_1 浅层休眠
    Set_sleepMode_1 Set_sleepMode = "1"
    // Set_sleepMode_2 深度休眠
    Set_sleepMode_2 Set_sleepMode = "2"
)
var Set_sleepModeDescMap = map[Set_sleepMode]string{
    Set_sleepMode_1: "浅层休眠",
    Set_sleepMode_2: "深度休眠",
    
}
func (e *Set_sleepMode) String() string {
    return Set_sleepModeDescMap[*e]
}
