package util

import "github.com/dromara/carbon/v2"

// Now 获取当前时间
func Now() *carbon.Carbon {
	return carbon.Now()
}

// Yesterday 获取昨天
func Yesterday() *carbon.Carbon {
	return carbon.Yesterday()
}

// Tomorrow 获取明天
func Tomorrow() *carbon.Carbon {
	return carbon.Tomorrow()
}

// Parse 解析时间字符串（使用默认格式）
func Parse(value string) *carbon.Carbon {
	return carbon.Parse(value)
}

// ParseByFormat 通过格式解析时间
func ParseByFormat(value, format string) *carbon.Carbon {
	return carbon.ParseByFormat(value, format)
}

// ParseByLayout 通过布局解析时间
func ParseByLayout(value, layout string) *carbon.Carbon {
	return carbon.ParseByLayout(value, layout)
}

// FromTimestamp 从时间戳创建（秒级）
func FromTimestamp(timestamp int64) *carbon.Carbon {
	return carbon.CreateFromTimestamp(timestamp)
}

// FromTimestampMilli 从时间戳创建（毫秒级）
func FromTimestampMilli(timestamp int64) *carbon.Carbon {
	return carbon.CreateFromTimestampMilli(timestamp)
}

// FromDateTime 从日期时间创建
func FromDateTime(year, month, day, hour, minute, second int) *carbon.Carbon {
	return carbon.CreateFromDateTime(year, month, day, hour, minute, second)
}

// FromDate 从日期创建
func FromDate(year, month, day int) *carbon.Carbon {
	return carbon.CreateFromDate(year, month, day)
}

// Format 格式化时间
func Format(carb *carbon.Carbon, format string) string {
	return carb.Format(format)
}

// ToDateTimeString 格式化为日期时间字符串：Y-m-d H:i:s
func ToDateTimeString(carb *carbon.Carbon) string {
	return carb.ToDateTimeString()
}

// ToDateString 格式化为日期字符串：Y-m-d
func ToDateString(carb *carbon.Carbon) string {
	return carb.ToDateString()
}

// ToTimeString 格式化为时间字符串：H:i:s
func ToTimeString(carb *carbon.Carbon) string {
	return carb.ToTimeString()
}

// AddSeconds 增加秒数
func AddSeconds(carb *carbon.Carbon, seconds int) *carbon.Carbon {
	return carb.AddSeconds(seconds)
}

// AddMinutes 增加分钟数
func AddMinutes(carb *carbon.Carbon, minutes int) *carbon.Carbon {
	return carb.AddMinutes(minutes)
}

// AddHours 增加小时数
func AddHours(carb *carbon.Carbon, hours int) *carbon.Carbon {
	return carb.AddHours(hours)
}

// AddDays 增加天数
func AddDays(carb *carbon.Carbon, days int) *carbon.Carbon {
	return carb.AddDays(days)
}

// AddMonths 增加月数
func AddMonths(carb *carbon.Carbon, months int) *carbon.Carbon {
	return carb.AddMonths(months)
}

// AddYears 增加年数
func AddYears(carb *carbon.Carbon, years int) *carbon.Carbon {
	return carb.AddYears(years)
}

// SubSeconds 减少秒数
func SubSeconds(carb *carbon.Carbon, seconds int) *carbon.Carbon {
	return carb.SubSeconds(seconds)
}

// SubMinutes 减少分钟数
func SubMinutes(carb *carbon.Carbon, minutes int) *carbon.Carbon {
	return carb.SubMinutes(minutes)
}

// SubHours 减少小时数
func SubHours(carb *carbon.Carbon, hours int) *carbon.Carbon {
	return carb.SubHours(hours)
}

// SubDays 减少天数
func SubDays(carb *carbon.Carbon, days int) *carbon.Carbon {
	return carb.SubDays(days)
}

// DiffInSeconds 计算秒数差异
func DiffInSeconds(carb1, carb2 *carbon.Carbon) int64 {
	return carb1.DiffInSeconds(carb2)
}

// DiffInMinutes 计算分钟数差异
func DiffInMinutes(carb1, carb2 *carbon.Carbon) int64 {
	return carb1.DiffInMinutes(carb2)
}

// DiffInHours 计算小时数差异
func DiffInHours(carb1, carb2 *carbon.Carbon) int64 {
	return carb1.DiffInHours(carb2)
}

// DiffInDays 计算天数差异
func DiffInDays(carb1, carb2 *carbon.Carbon) int64 {
	return carb1.DiffInDays(carb2)
}

// IsToday 判断是否是今天
func IsToday(carb *carbon.Carbon) bool {
	return carb.IsToday()
}

// IsYesterday 判断是否是昨天
func IsYesterday(carb *carbon.Carbon) bool {
	return carb.IsYesterday()
}

// IsFuture 判断是否是未来时间
func IsFuture(carb *carbon.Carbon) bool {
	return carb.IsFuture()
}

// IsPast 判断是否是过去时间
func IsPast(carb *carbon.Carbon) bool {
	return carb.IsPast()
}

// StartOfDay 获取当天开始时间
func StartOfDay(carb *carbon.Carbon) *carbon.Carbon {
	return carb.StartOfDay()
}

// EndOfDay 获取当天结束时间
func EndOfDay(carb *carbon.Carbon) *carbon.Carbon {
	return carb.EndOfDay()
}

// StartOfMonth 获取当月开始时间
func StartOfMonth(carb *carbon.Carbon) *carbon.Carbon {
	return carb.StartOfMonth()
}

// EndOfMonth 获取当月结束时间
func EndOfMonth(carb *carbon.Carbon) *carbon.Carbon {
	return carb.EndOfMonth()
}

// StartOfYear 获取当年开始时间
func StartOfYear(carb *carbon.Carbon) *carbon.Carbon {
	return carb.StartOfYear()
}

// EndOfYear 获取当年结束时间
func EndOfYear(carb *carbon.Carbon) *carbon.Carbon {
	return carb.EndOfYear()
}
