package utils

import (
	"fmt"
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

func GetUnixNano() int64 {
	return time.Now().UnixNano()
}

func GetUnix() int64 {
	return time.Now().Unix()
}

// 默认格式日期字符串转time
func TimeStrToTimeDefault(str string) time.Time {
	parseTime, _ := time.ParseInLocation(DefaultLayout, str, time.Local)
	return parseTime
}

// 判断是否未空  空false  有值为true
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
		fmt.Println("执行false")
		return false
	}
}
