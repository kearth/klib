package web

import (
	"time"
)

const DATE_FORMAT = "2006-01-02 15:04:05"

// 毫秒 当前时间
func Now() int {
	return int(time.Now().Unix())
}

// 毫秒 时间戳转日期
func Timestamp2Date(i int) string {
	t := time.Unix(int64(i), 0)
	return t.Format(DATE_FORMAT)
}

// 毫秒 日期转时间戳
func Date2Timestamp(s string) (int, error) {
	// 东八区
	loc, _ := time.LoadLocation("Asia/Shanghai")
	t, err := time.ParseInLocation(DATE_FORMAT, s, loc)
	if err != nil {
		return 0, err
	}
	return int(t.Unix()), nil
}
