package ktime

import (
	"fmt"
	"time"
)

// PHP 预定义时间格式
const (
	PHP_ATOM             = "2006-01-02T15:04:05-07:00"        // 2022-12-28T10:50:29+00:00
	PHP_COOKIE           = "Monday, 02-Jan-2006 15:04:05 MST" // Wednesday, 28-Dec-2022 10:50:29 UTC
	PHP_ISO8601          = "2006-01-02T15:04:05-0700"         // 2022-12-28T10:50:29+0000
	PHP_RFC822           = "Mon, 02 Jan 06 15:04:05 -0700"    // Wed, 28 Dec 22 10:50:29 +0000
	PHP_RFC850           = "Monday, 02-Jan-06 15:04:05 MST"   // Wednesday, 28-Dec-22 10:50:29 UTC
	PHP_RFC1036          = "Mon, 02 Jan 06 15:04:05 -0700"    // Wed, 28 Dec 22 10:50:29 +0000
	PHP_RFC1123          = "Mon, 02 Jan 2006 15:04:05 -0700"  // Wed, 28 Dec 2022 10:50:29 +0000
	PHP_RFC7231          = "Mon, 02 Jan 2006 15:04:05 GMT"    // Wed, 28 Dec 2022 10:50:29 GMT
	PHP_RFC2822          = "Mon, 02 Jan 2006 15:04:05 -0700"  // Wed, 28 Dec 2022 10:50:29 +0000
	PHP_RFC3339          = "2006-01-02T15:04:05-07:00"        // 2022-12-28T10:50:29+00:00
	PHP_RFC3339_EXTENDED = "2006-01-02T15:04:05.000-07:00"    // 2022-12-28T10:50:29.000+00:00
	PHP_RSS              = "Mon, 02 Jan 2006 15:04:05 -0700"  // Wed, 28 Dec 2022 10:50:29 +0000
	PHP_W3C              = "2006-01-02T15:04:05-07:00"        // 2022-12-28T10:50:29+00:00
)

// PHP 格式映射
var (
	phpFormatMap = map[byte]string{
		'd': "02",                // 01 to 31 // Day of the month, 2 digits with leading zeros
		'D': "Mon",               // Mon through Sun // A textual representation of a day, three letters
		'j': "2",                 // 1 to 31 // Day of the month without leading zeros
		'l': "Monday",            // Sunday through Saturday // A full textual representation of the day of the week
		'N': formatFuncWeekday,   // 1 (for Monday) through 7 (for Sunday) // ISO 8601 numeric representation of the day of the week
		'S': formatFuncUnsupport, // st, nd, rd or th. Works well with j // English ordinal suffix for the day of the month, 2 characters
		'w': formatFuncWeekday,   // 0 (for Sunday) through 6 (for Saturday) // Numeric representation of the day of the week
		'z': formatFuncYearDay,   // 0 through 365 // The day of the year (starting from 0)
		'W': formatFuncISOWeek,   // Example: 42 (the 42nd week in the year) // ISO 8601 week number of year, weeks starting on Monday
		'F': "January",           // January through December // A full textual representation of a month, such as January or March
		'm': "01",                // 01 through 12 // Numeric representation of a month, with leading zeros
		'M': "Jan",               // Jan through Dec // A short textual representation of a month, three letters
		'n': "1",                 // 1 through 12 // Numeric representation of a month, without leading zeros
		't': formatFuncDate,      // 28 through 31 // Number of days in the given month
		'L': formatFuncYear,      // 1 if it is a leap year, 0 otherwise. // Whether it's a leap year
		'o': formatFuncISOWeek,   // Examples: 1999 or 2003 // ISO 8601 week-numbering year. This has the same value as Y, except that if the ISO week number (W) belongs to the previous or next year, that year is used instead.
		'X': formatFuncYear,      // Examples: -0055, +0787, +1999, +10191 // An expanded full numeric representation of a year, at least 4 digits, with - for years BCE, and + for years CE.
		'x': formatFuncYear,      // Examples: -0055, 0787, 1999, +10191 // An expanded full numeric representation if requried, or a standard full numeral representation if possible (like Y). At least four digits. Years BCE are prefixed with a -. Years beyond (and including) 10000 are prefixed by a +.
		'Y': "2006",              // Examples: -0055, 0787, 1999, 2003, 10191 // A full numeric representation of a year, at least 4 digits, with - for years BCE.
		'y': "06",                // Examples: 99 or 03 // A two digit representation of a year
		'a': "pm",                // am or pm // Lowercase Ante meridiem and Post meridiem
		'A': "PM",                // AM or PM // Uppercase Ante meridiem and Post meridiem
		'B': formatFuncClock,     // 000 through 999 // Swatch Internet time
		'g': "3",                 // 1 through 12 // 12-hour format of an hour without leading zeros
		'G': formatFuncHour,      // 0 through 23 // 24-hour format of an hour without leading zeros
		'h': "03",                // 01 through 12 // 12-hour format of an hour with leading zeros
		'H': "15",                // 00 through 23 // 24-hour format of an hour with leading zeros
		'i': "04",                // 00 to 59 // Minutes with leading zeros
		's': "05",                // 00 through 59 // Seconds with leading zeros
		'u': formatFuncUnsupport, // Example: 654321 // Microseconds. Note that date() will always generate 000000 since it takes an int parameter, whereas DateTime::format() does support microseconds if DateTime was created with microseconds.
		'v': "000",               // Example: 654 // Milliseconds. Same note applies as for u.
		'e': formatFuncLocation,  // Examples: UTC, GMT, Atlantic/Azores // Timezone identifier
		'I': formatFuncIsDST,     // 1 if Daylight Saving Time, 0 otherwise. // Whether or not the date is in daylight saving time
		'O': "-0700",             // Example: +0200 // Difference to Greenwich time (GMT) without colon between hours and minutes
		'P': "-07:00",            // Example: +02:00 // Difference to Greenwich time (GMT) with colon between hours and minutes
		'T': "MST",               // Examples: EST, MDT, +05 // Timezone abbreviation, if known; otherwise the GMT offset
		'Z': formatFuncUnsupport, // -43200 through 50400 // Timezone offset in seconds. The offset for timezones west of UTC is always negative, and for those east of UTC is always positive.
		'c': formatFuncFormat,    // 2004-02-12T15:19:21+00:00 // ISO 8601 date
		'r': formatFuncFormat,    // Example: Thu, 21 Dec 2000 16:01:07 +0200 // RFC 2822 or RFC 5322 formatted date
		'U': formatFuncUnix,      // See also time() // Seconds since the Unix Epoch (January 1 1970 00:00:00 GMT)
	}
)

// PHPKTime
type PHPKTime struct{}

// parseformatFunc
func (p *PHPKTime) parseFormatFunc(t time.Time, fk byte, fc FormatFunc) string {
	switch string(fc) {
	case formatFuncWeekday:
		if fk == 'w' {
			return strToInt(int(t.Weekday()))
		}
		return t.Weekday().String()
	case formatFuncYearDay:
		return strToInt(t.YearDay())
	case formatFuncISOWeek:
		y, w := t.ISOWeek()
		if fk == 'o' {
			return strToInt(y)
		}
		return strToInt(w)
	case formatFuncYear:
		if fk == 'L' {
			if IsLeap(t.Year()) {
				return "1"
			}
			return "0"
		}
		if fk == 'X' {
			return fmt.Sprintf("%+05d", t.Year())
		}
		if t.Year() > 10000 || t.Year() < 0 {
			return fmt.Sprintf("%+05d", t.Year())
		}
		return fmt.Sprintf("%04d", t.Year())
	case formatFuncClock:
		loc := time.FixedZone("UTC+1", 60*60)
		h, m, _ := t.In(loc).Clock()
		return fmt.Sprintf("%d", int(float64(h*60+m)/1.44))
	case formatFuncHour:
		return strToInt(t.Hour())
	case formatFuncFormat:
		if fk == 'r' {
			return t.UTC().Format(PHP_RFC2822)
		}
		return t.UTC().Format("2006-01-02T15:04:05-07:00")
	case formatFuncLocation:
		return t.Location().String()
	case formatFuncIsDST:
		if t.IsDST() {
			return "1"
		}
		return "0"
	case formatFuncUnix:
		return fmt.Sprintf("%d", t.Unix())
	case formatFuncDate:
		y, m, _ := t.Date()
		return strToInt(MonthDays(y, m))
	}
	return string(fc)
}

// getMap
func (p *PHPKTime) getMap() map[byte]string {
	return phpFormatMap
}

// PHPFormat
func PHPFormat(t time.Time, s string) string {
	return format(new(PHPKTime), t, s)
}

// PHPLayout
func PHPLayout(s string) string {
	return layout(new(PHPKTime), s)
}
