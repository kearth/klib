package ktime

import (
	"strconv"
	"strings"
	"time"
)

// 格式化映射 - 基于 standard time 包
//
// const (
// 	_                        = iota
// 	stdLongMonth             = iota + stdNeedDate  // "January"
// 	stdMonth                                       // "Jan"
// 	stdNumMonth                                    // "1"
// 	stdZeroMonth                                   // "01"
// 	stdLongWeekDay                                 // "Monday"
// 	stdWeekDay                                     // "Mon"
// 	stdDay                                         // "2"
// 	stdUnderDay                                    // "_2"
// 	stdZeroDay                                     // "02"
// 	stdUnderYearDay                                // "__2"
// 	stdZeroYearDay                                 // "002"
// 	stdHour                  = iota + stdNeedClock // "15"
// 	stdHour12                                      // "3"
// 	stdZeroHour12                                  // "03"
// 	stdMinute                                      // "4"
// 	stdZeroMinute                                  // "04"
// 	stdSecond                                      // "5"
// 	stdZeroSecond                                  // "05"
// 	stdLongYear              = iota + stdNeedDate  // "2006"
// 	stdYear                                        // "06"
// 	stdPM                    = iota + stdNeedClock // "PM"
// 	stdpm                                          // "pm"
// 	stdTZ                    = iota                // "MST"
// 	stdISO8601TZ                                   // "Z0700"  // prints Z for UTC
// 	stdISO8601SecondsTZ                            // "Z070000"
// 	stdISO8601ShortTZ                              // "Z07"
// 	stdISO8601ColonTZ                              // "Z07:00" // prints Z for UTC
// 	stdISO8601ColonSecondsTZ                       // "Z07:00:00"
// 	stdNumTZ                                       // "-0700"  // always numeric
// 	stdNumSecondsTz                                // "-070000"
// 	stdNumShortTZ                                  // "-07"    // always numeric
// 	stdNumColonTZ                                  // "-07:00" // always numeric
// 	stdNumColonSecondsTZ                           // "-07:00:00"
// 	stdFracSecond0                                 // ".0", ".00", ... , trailing zeros included
// 	stdFracSecond9                                 // ".9", ".99", ..., trailing zeros omitted
//
// 	stdNeedDate       = 1 << 8             // need month, day, year
// 	stdNeedClock      = 2 << 8             // need hour, minute, second
// 	stdArgShift       = 16                 // extra argument in high bits, above low stdArgShift
// 	stdSeparatorShift = 28                 // extra argument in high 4 bits for fractional second separators
// 	stdMask           = 1<<stdArgShift - 1 // mask out argument
// )

// 格式化操作符
const (
	formatReplace       = "%"
	formatFuncPrefix    = "formatFunc"
	formatFuncUnsupport = ""
	formatFuncWeekday   = "formatFuncWeekday"
	formatFuncYearDay   = "formatFuncYearDay"
	formatFuncISOWeek   = "formatFuncISOWeek"
	formatFuncYear      = "formatFuncYear"
	formatFuncClock     = "formatFuncClock"
	formatFuncHour      = "formatFuncHour"
	formatFuncFormat    = "formatFuncFormat"
	formatFuncLocation  = "formatFuncLocation"
	formatFuncIsDST     = "formatFuncIsDST"
	formatFuncUnix      = "formatFuncUnix"
	formatFuncDate      = "formatFuncDate"
)

// FormatFunc
type FormatFunc string

// contains
func (f FormatFunc) contains() string {
	if strings.HasPrefix(string(f), formatFuncPrefix) {
		return string(f)
	}
	return ""
}

// KTime
type KTime interface {
	getMap() map[byte]string
	parseFormatFunc(t time.Time, fk byte, fc FormatFunc) string
}

// format
func format(k KTime, t time.Time, s string) string {
	sLen := len(s)
	snl := make([]string, 0, sLen)
	fl := make([]string, 0, sLen)
	m := k.getMap()
	for i := 0; i < sLen; i++ {
		rs := m[s[i]]
		switch rs {
		case "":
			snl = append(snl, string(s[i]))
		case FormatFunc(rs).contains():
			fl = append(fl, k.parseFormatFunc(t, s[i], FormatFunc(rs)))
			snl = append(snl, formatReplace)
		default:
			snl = append(snl, rs)
		}
	}
	sn := t.Format(strings.Join(snl, ""))
	j, p := 0, 0
	for j < len(sn) {
		if string(sn[j]) == formatReplace {
			sn = sn[0:j] + fl[p] + sn[j+1:]
			j = j + len(fl[p]) - 1
			p++
		}
		j++
	}
	return sn
}

// layout
func layout(k KTime, s string) string {
	sLen := len(s)
	snl := make([]string, 0, sLen)
	m := k.getMap()
	for i := 0; i < sLen; i++ {
		rs := m[s[i]]
		switch rs {
		case "":
			snl = append(snl, string(s[i]))
		case FormatFunc(rs).contains():
			snl = append(snl, string(s[i]))
		default:
			snl = append(snl, rs)
		}
	}
	return strings.Join(snl, "")
}

// strToInt
func strToInt(i int) string {
	return strconv.Itoa(i)
}

// IsLeap
func IsLeap(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

// MonthDays
func MonthDays(year int, month time.Month) int {
	switch month {
	case time.January, time.March, time.May, time.July, time.August, time.October, time.December:
		return 31
	case time.April, time.June, time.September, time.November:
		return 30
	}
	if IsLeap(year) {
		return 29
	}
	return 28
}
