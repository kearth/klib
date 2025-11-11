package klog

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/gogf/gf/v2/os/glog"
	"github.com/kearth/klib/kctx"
	"github.com/stretchr/testify/assert"
)

// TestLogString 测试日志格式化输出
func TestLogString(t *testing.T) {
	t.Parallel()

	// 测试基础日志格式
	t.Run("basic format", func(t *testing.T) {
		log := &Log{
			Time:     "2024-01-01 12:00:00 UTC",
			Level:    "INFO",
			LevelInt: glog.LEVEL_INFO,
			TraceID:  "test-trace-123",
			Body:     []any{"user", "login"},
			Add:      map[string]string{"user_id": "1001"},
		}
		result := log.String()
		assert.Contains(t, result, "2024-01-01 12:00:00 UTC")
		assert.Contains(t, result, "[INFO]")
		assert.Contains(t, result, "test-trace-123")
		assert.Contains(t, result, "user_id=1001")
		assert.Contains(t, result, "userlogin")
	})

	// 测试无TraceID和附加信息的情况（修复核心错误）
	t.Run("no traceid and add", func(t *testing.T) {
		log := &Log{
			Time:     "2024-01-01 12:00:00 UTC",
			Level:    "DEBUG",
			LevelInt: glog.LEVEL_DEBU,
			TraceID:  "",
			Body:     []any{"system", "start"},
			Add:      map[string]string{}, // 无附加信息
		}
		result := log.String()
		assert.Contains(t, result, "-")       // 无TraceID显示"-"
		assert.Contains(t, result, "[DEBUG]") // 级别正常显示[]
		assert.NotContains(t, result, " [ ")  // 无附加信息，不显示元数据的[]块（关键验证）
		assert.NotContains(t, result, " ]")   // 无附加信息，不显示元数据的结尾]
	})

	// 测试空级别情况
	t.Run("empty level", func(t *testing.T) {
		log := &Log{
			Body: []any{"plain", "message"},
			Add:  map[string]string{"key": "val"},
		}
		result := log.String()
		assert.Contains(t, result, "key=val")
		assert.Contains(t, result, "plainmessage")
	})
}

// TestFormatBody 测试日志内容格式化
func TestFormatBody(t *testing.T) {
	t.Parallel()

	// 测试附加信息和主体内容都存在
	t.Run("with add and body", func(t *testing.T) {
		body := []any{"action", "submit"}
		add := map[string]string{"page": "home", "method": "post"}
		result := formatBody(body, add)
		assert.Contains(t, result, "page=home")
		assert.Contains(t, result, "method=post")
		assert.Contains(t, result, "actionsubmit")
	})

	// 测试只有附加信息
	t.Run("only add", func(t *testing.T) {
		body := []any{}
		add := map[string]string{"status": "success"}
		result := formatBody(body, add)
		assert.Contains(t, result, "status=success")
		assert.NotContains(t, result, "{") // 无主体内容时不显示{}
	})

	// 测试只有主体内容
	t.Run("only body", func(t *testing.T) {
		body := []any{"error", "occurred"}
		add := map[string]string{}
		result := formatBody(body, add)
		assert.Contains(t, result, "erroroccurred")
		assert.NotContains(t, result, "[") // 无附加信息时不显示[]
	})
}

// TestContextMetaExtraction 测试从上下文中提取元数据
func TestContextMetaExtraction(t *testing.T) {
	t.Parallel()

	// 测试kctx.Context类型的上下文
	t.Run("kctx context", func(t *testing.T) {
		ctx := kctx.New()
		ctx.Set("user_id", "1001")
		ctx.Set("role", "admin")
		traceID := ctx.TraceID()

		// 模拟日志处理器中的逻辑
		var newCtx kctx.Context
		if casted, ok := ctx.(kctx.Context); ok {
			newCtx = casted
		} else {
			newCtx = kctx.New(ctx)
		}

		assert.Equal(t, traceID, newCtx.TraceID())
		assert.Equal(t, "1001", newCtx.Get("user_id"))
		assert.Equal(t, "admin", newCtx.Get("role"))
	})

	// 测试标准context.Context类型
	t.Run("standard context", func(t *testing.T) {
		ctx := context.Background()
		// 模拟从标准context创建kctx
		newCtx := kctx.New(ctx)
		assert.NotEmpty(t, newCtx.TraceID()) // 应自动生成TraceID
	})
}

// TestLogFunctions 测试日志输出函数（验证无panic）
func TestLogFunctions(t *testing.T) {
	t.Parallel()

	ctx := kctx.New()
	testCases := []struct {
		name     string
		function func(ctx context.Context, v ...any)
	}{
		{"Info", Info},
		{"Debug", Debug},
		{"Notice", Notice},
		{"Warn", Warn},
		{"Error", Error},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 捕获panic，确保日志函数可正常调用
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("%s函数调用 panic: %v", tc.name, r)
				}
			}()
			tc.function(ctx, "test", tc.name, "log")
		})
	}
}

// TestAddToCtx 测试上下文添加元数据
func TestAddToCtx(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	key, val := "test_key", "test_val"

	newCtx := AddToCtx(ctx, key, val)
	// 从新上下文重建kctx验证
	kctxInstance := kctx.New(newCtx)
	assert.Equal(t, val, kctxInstance.Get(key))
}

// TestAddMapToCtx 测试批量添加元数据到上下文
func TestAddMapToCtx(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	kv := map[string]string{
		"k1": "v1",
		"k2": "v2",
	}

	newCtx := AddMapToCtx(ctx, kv)
	kctxInstance := kctx.New(newCtx)
	assert.Equal(t, "v1", kctxInstance.Get("k1"))
	assert.Equal(t, "v2", kctxInstance.Get("k2"))
}

// TestDefaultHandler 测试默认日志处理器（基础功能验证）
func TestDefaultHandler(t *testing.T) {
	t.Parallel()

	ctx := kctx.New()
	ctx.Set("handler", "test")

	// 构造测试输入
	input := &glog.HandlerInput{
		Level:       glog.LEVEL_INFO,
		LevelFormat: "INFO",
		Time:        time.Now(),
		Values:      []any{"handler", "test"},
		Buffer:      &bytes.Buffer{},
	}

	// 执行处理器
	DefaultHandler(ctx, input)
	result := input.Buffer.String()

	// 验证输出包含关键信息
	assert.Contains(t, result, ctx.TraceID())
	assert.Contains(t, result, "handler=test")
	assert.Contains(t, result, "[INFO]")
	assert.Contains(t, result, "handlertest")
}

func TestColorPrint(t *testing.T) {
	t.Parallel()
	ctx := kctx.New()
	Init()
	Print(ctx, "plain ", "log")
	ColorPrint(ctx, Yellow, "yellow ", "log")
	Notice(ctx, "green ", "log")
}
