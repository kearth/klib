package main

import (
	"fmt"

	"time"

	"github.com/kearth/kgolib/ktime"
)

func main() {
	//	kf := ktime.PHPLayout("Y-m-d H:i:s N")
	//	var t1 time.Time
	fmt.Printf("%+v\n", ktime.PHPFormat(time.Now().AddDate(0, 0, 0), "Y-m-d H:i:s  t"))
	//	fmt.Printf("%+v\n", ktime.PHPFormat(time.Now().AddDate(+9000, 0, 0), "Y-m-d H:i:s  X"))
	//	fmt.Printf("%+v\n", ktime.PHPFormat(time.Now().AddDate(-2000, 0, 0), "Y-m-d H:i:s  X"))
	//	a, b := t1.ISOWeek()
	loc, _ := time.LoadLocation("Europe/Berlin")
	fmt.Printf("%+v", time.Now().In(loc).Location())
}
