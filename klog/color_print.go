package klog

import (
	"context"

	"github.com/gogf/gf/v2/os/glog"
)

// Color 颜色
type Color int

var (
	Yellow  Color = Color(glog.COLOR_YELLOW)
	Green   Color = Color(glog.COLOR_GREEN)
	Cyan    Color = Color(glog.COLOR_CYAN)
	Magenta Color = Color(glog.COLOR_MAGENTA)
	Red     Color = Color(glog.COLOR_RED)
	HiRed   Color = Color(glog.COLOR_HI_RED)
)

// ColorPrint  彩色打印日志
func ColorPrint(ctx context.Context, color Color, v ...any) {
	nv := append([]any{color}, v...)
	Logger().Print(ctx, nv...)
}
