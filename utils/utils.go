package utils

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

const (
	YearMonthDay     = "2006-01-02"
	HourMinuteSecond = "15:04:05"
	DefaultLayout    = YearMonthDay + " " + HourMinuteSecond
)

/**
 * 格式化数据
 */
func FormatDatetime(time time.Time) string {
	return time.Format("2006-01-02 03:04:05")
}

// GetDate 获取当前时间
func GetDate() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func GetDay() string {
	template := "20060102"
	return time.Now().Format(template)
}

// GetUnixNano 获取当前时间 纳秒
func GetUnixNano() int64 {
	return time.Now().UnixNano()
}

// GetUnix 获取当前时间戳
func GetUnix() int64 {
	return time.Now().Unix()
}

// TimeStrToTimeDefault 默认格式日期字符串转time
func TimeStrToTimeDefault(str string) time.Time {
	parseTime, _ := time.ParseInLocation(DefaultLayout, str, time.Local)
	return parseTime
}

// IsEmpty 判断是否未空  空false  有值为true
func IsEmpty(con interface{}) bool {
	switch con.(type) {
	case string:
		if len(con.(string)) == 0 {
			return false
		} else {
			return true
		}
	case int:
		if con.(int) == 0 {
			return false
		} else {
			return true
		}
	default:
		return false
	}
}

// RandString 生成随机字符串
func RandString(n int) string {
	result := make([]byte, n/2)
	rand.Read(result)
	return hex.EncodeToString(result)
}
