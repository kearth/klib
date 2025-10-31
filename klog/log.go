package klog

import (
	"context"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/kearth/klib/kctx"
	"github.com/kearth/klib/kutil"
)

// 颜色映射
var colorMaps = map[int]int{
	glog.LEVEL_DEBU: glog.COLOR_YELLOW,
	glog.LEVEL_INFO: glog.COLOR_GREEN,
	glog.LEVEL_NOTI: glog.COLOR_CYAN,
	glog.LEVEL_WARN: glog.COLOR_MAGENTA,
	glog.LEVEL_ERRO: glog.COLOR_RED,
	glog.LEVEL_CRIT: glog.COLOR_HI_RED,
	glog.LEVEL_PANI: glog.COLOR_HI_RED,
	glog.LEVEL_FATA: glog.COLOR_HI_RED,
}

// DefaultHandler 默认日志处理
func DefaultHandler(ctx context.Context, in *glog.HandlerInput) {
	var newCtx kctx.Context
	if casted, ok := ctx.(kctx.Context); ok {
		newCtx = casted
	} else {
		newCtx = kctx.New(ctx)
	}
	in.Buffer.WriteString((&Log{
		Time:     in.Time.Format("2006-01-02 15:04:05 Z07:00"),
		Level:    in.LevelFormat,
		LevelInt: in.Level,
		TraceID:  newCtx.TraceID(),
		Body:     in.Values,
		Add:      newCtx.Values(),
	}).String())
	in.Buffer.WriteString("\n")
	in.Next(ctx)
}

// 初始化日志
func Init() {
	glog.SetDefaultHandler(DefaultHandler)
}

// Log 日志结构体
type Log struct {
	Time     string
	Level    string
	LevelInt int
	TraceID  string
	Body     []any
	Add      map[string]string
}

// Logger 获取指定名称的日志实例，若名称为空则返回默认实例
//
//	name - 日志实例名称（对应配置中的日志名称）
//	返回：*glog.Logger 日志实例
func Logger(name ...string) *glog.Logger {
	if len(name) > 0 {
		return g.Log(name[0])
	}
	return g.Log()
}

// formatBody 格式化日志主体
func formatBody(body []any, add map[string]string) string {
	var b strings.Builder
	if len(add) > 0 {
		b.WriteString(" [")
		for k, v := range add {
			b.WriteString(fmt.Sprintf(" %s=%s,", k, v))
		}
		b.WriteString(" ]")
	}
	if len(body) > 0 {
		b.WriteString(" {")
		for _, v := range body {
			b.WriteString(gconv.String(v))
		}
		b.WriteString(" }")
	}
	return b.String()
}

// String 格式化日志
func (l *Log) String() string {
	if l.Level == "" {
		return formatBody(l.Body, l.Add)
	}
	return fmt.Sprintf(
		"%s %s %s %s",
		l.Time,
		color.New(color.Attribute(colorMaps[l.LevelInt])).Sprint("["+l.Level+"]"),
		kutil.If[string](l.TraceID == "", "-", l.TraceID),
		formatBody(l.Body, l.Add))
}

// AddToCtx 追加参数
func AddToCtx(ctx context.Context, key string, val string) context.Context {
	return AddMapToCtx(ctx, map[string]string{key: val})
}

// AddMapToCtx 追加参数 - map
func AddMapToCtx(ctx context.Context, kv map[string]string) context.Context {
	var newCtx kctx.Context
	var ok bool
	if newCtx, ok = ctx.(kctx.Context); ok {
		for k, v := range kv {
			newCtx.Set(k, v)
		}
	} else {
		ctx = context.WithValue(ctx, kctx.MetaMapKey, kv)
		newCtx = kctx.New(ctx)
	}
	return newCtx
}

// Info 打印日志
func Info(ctx context.Context, v ...any) {
	Logger().Info(ctx, v...)
}

// Debug 打印日志
func Debug(ctx context.Context, v ...any) {
	Logger().Debug(ctx, v...)
}

// Notice 打印日志
func Notice(ctx context.Context, v ...any) {
	Logger().Notice(ctx, v...)
}

// Warn 打印警告级日志（原 Warning 重命名，对齐 glog 命名）
func Warn(ctx context.Context, v ...any) {
	Logger().Warning(ctx, v...)
}

// Error 打印错误级日志
func Error(ctx context.Context, v ...any) {
	Logger().Error(ctx, v...)
}

// Panic 打印日志
func Panic(ctx context.Context, v ...any) {
	Logger().Panic(ctx, v...)
}
