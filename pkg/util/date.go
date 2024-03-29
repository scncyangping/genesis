// @Author: YangPing
// @Create: 2023/10/23
// @Description: 日期工具类

package util

import (
	"time"
)

// NowDateFormat 当前时间年月日字符串
func NowDateFormat() string {
	t := time.Now()
	return t.Format("2006-01-02")
}

// NowDateTimeFormat 当前时间字符串
func NowDateTimeFormat() string {
	t := time.Now()
	return t.Format("2006-01-02 15:04:05")
}

// DateTimeFormatCustom 当前时间自定义格式化字符串
func DateTimeFormatCustom(format string) string {
	t := time.Now()
	return t.Format(format)
}

// NowTimestamp 当前时间戳
func NowTimestamp() int64 {
	return time.Now().Unix()
}

// TimeFormat 时间戳转年月日时分秒字符串
func TimeFormat(ts int64) string {
	t := time.Unix(ts, 0)
	return t.Format("2006-01-02 15:04:05")
}

// TimeFormatCustom 时间戳转自定义字符串
func TimeFormatCustom(ts int64, format string) string {
	t := time.Unix(ts, 0)
	return t.Format(format)
}
