package kerr

import (
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"strings"
)

// Error 定义了错误接口
type Error interface {
	error
	Code() int
	Display() string
	Wrap(err error) Error
	Unwrap() error
	Stack() string
	Is(target error) bool
	WithStack(skip ...int) Error
	WithDisplay(display string) Error
	ToJSON() string
}

// New 创建一个新的错误实例
func New(code int, msg string) *KError {
	return &KError{
		code:    code,
		msg:     msg,
		display: msg, // 默认 display = msg
	}
}

// KError 是 Error 的默认实现
type KError struct {
	code    int
	msg     string
	display string
	cause   error
	stack   []uintptr
}

// Error 返回错误信息
func (e *KError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s: %v", e.msg, e.cause)
	}
	return e.msg
}

// Code 返回错误码
func (e *KError) Code() int {
	return e.code
}

// Display 返回错误显示信息
func (e *KError) Display() string {
	if e.display != "" {
		return e.display
	}
	// 若 display 为空，兜底使用 msg（或错误码默认信息）
	return e.msg
}

// WithDisplay 设置错误显示信息
func (e *KError) WithDisplay(display string) Error {
	// 复制原实例的所有字段，仅修改 display
	return &KError{
		code:    e.code,
		msg:     e.msg,
		display: display, // 新值
		cause:   e.cause,
		stack:   e.stack,
	}
}

// Wrap 包装错误，返回新实例
func (e *KError) Wrap(err error) Error {
	if err == nil {
		return e
	}
	return &KError{
		code:    e.code,
		msg:     e.msg,
		cause:   err,
		display: e.display,
		stack:   e.stack,
	}
}

// Unwrap 返回被包装的错误
func (e *KError) Unwrap() error {
	return e.cause
}

// Stack 返回调用栈
func (e *KError) Stack() string {
	if e.stack == nil {
		return ""
	}
	var b strings.Builder
	frames := runtime.CallersFrames(e.stack)
	for {
		frame, more := frames.Next()
		b.WriteString(fmt.Sprintf("%s:%d %s\n", frame.File, frame.Line, frame.Function))
		if !more {
			break
		}
	}
	return b.String()
}

// Is 判断相同错误码的错误
func (e *KError) Is(target error) bool {
	if target == nil {
		return false
	}
	if t, ok := target.(*KError); ok && e.code == t.code {
		return true
	}
	if e.cause != nil {
		return errors.Is(e.cause, target)
	}
	return false
}

// WithStack 收集调用栈，返回新实例
func (e *KError) WithStack(skip ...int) Error {
	if e.stack != nil {
		return e // 已存在堆栈，直接返回原实例（不可变，安全）
	}
	// 新建堆栈，返回新实例
	stack := make([]uintptr, 32)
	skipCount := 2 // 默认跳过 WithStack 和直接调用者
	if len(skip) > 0 {
		skipCount += skip[0]
	}
	n := runtime.Callers(skipCount, stack)
	return &KError{
		code:    e.code,
		msg:     e.msg,
		display: e.display,
		cause:   e.cause,
		stack:   stack[:n], // 新堆栈仅在新实例中
	}
}

// Format 支持 fmt.Printf("%+v", err) 打印堆栈
func (e *KError) Format(f fmt.State, c rune) {
	switch c {
	case 'v':
		if f.Flag('+') {
			fmt.Fprintf(f, "%s (code=%d, display=%q)\n", e.msg, e.code, e.display)
			if e.stack != nil {
				fmt.Fprint(f, e.Stack())
			}
			if e.cause != nil {
				fmt.Fprintf(f, "Caused by: %+v\n", e.cause)
			}
			return
		}
		fallthrough
	case 's':
		fmt.Fprint(f, e.Error())
	default:
		fmt.Fprint(f, e.Error())
	}
}

// MarshalJSON 支持 JSON 序列化
func (e *KError) MarshalJSON() ([]byte, error) {
	// 递归处理底层错误的 JSON 序列化
	var causeJSON interface{}
	if e.cause != nil {
		if kerr, ok := e.cause.(*KError); ok {
			// 若底层是 KError，直接复用其序列化结果
			causeJSON, _ = kerr.MarshalJSON()
		} else {
			// 其他错误类型，序列化其错误信息
			causeJSON = e.cause.Error()
		}
	}

	// 堆栈处理保持不变
	stackTrace := []string{}
	for _, pc := range e.stack {
		if fn := runtime.FuncForPC(pc); fn != nil {
			file, line := fn.FileLine(pc)
			stackTrace = append(stackTrace, fmt.Sprintf("%s:%d", file, line))
		}
	}

	return json.Marshal(struct {
		Code    int         `json:"code"`
		Msg     string      `json:"message"`
		Display string      `json:"display,omitempty"`
		Cause   interface{} `json:"cause,omitempty"` // 支持嵌套序列化
		Stack   []string    `json:"stack,omitempty"`
	}{
		Code:    e.code,
		Msg:     e.msg,
		Display: e.display,
		Cause:   causeJSON,
		Stack:   stackTrace,
	})
}

func (e *KError) ToJSON() string {
	b, err := json.Marshal(e)
	if err != nil {
		return fmt.Sprintf(`{"code":%d,"message":"%s"}`, e.code, e.msg)
	}
	return string(b)
}
