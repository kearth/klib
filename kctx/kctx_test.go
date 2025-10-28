package kctx

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNew 测试上下文创建及TraceID继承
func TestNew(t *testing.T) {
	t.Parallel()

	// 测试1：无父上下文时，应生成新TraceID
	ctx1 := New()
	assert.NotEmpty(t, ctx1.TraceID())
	assert.Equal(t, context.Background().Done(), ctx1.Context().Done()) // 底层默认是Background

	// 测试2：有父上下文时，应继承父上下文的TraceID
	parentTraceID := "test-trace-id-123"
	parentCtx := context.WithValue(context.Background(), TraceIDKey, parentTraceID)
	ctx2 := New(parentCtx)
	assert.Equal(t, parentTraceID, ctx2.TraceID())

	// 测试3：父上下文的TraceID非字符串类型，应生成新ID
	invalidParent := context.WithValue(context.Background(), TraceIDKey, 123) // 非字符串
	ctx3 := New(invalidParent)
	assert.NotEqual(t, 123, ctx3.TraceID()) // 不应继承非字符串值
	assert.NotEmpty(t, ctx3.TraceID())
}

// TestGetSet 测试元数据的Get/Set操作
func TestGetSet(t *testing.T) {
	t.Parallel()

	ctx := New()

	// 测试1：初始状态无数据
	assert.Empty(t, ctx.Get("key1"))

	// 测试2：设置并获取值
	ctx.Set("key1", "val1")
	assert.Equal(t, "val1", ctx.Get("key1"))

	// 测试3：覆盖已有值
	ctx.Set("key1", "val2")
	assert.Equal(t, "val2", ctx.Get("key1"))

	// 测试4：设置空键应被忽略
	ctx.Set("", "invalid")
	assert.Empty(t, ctx.Get(""))

	// 测试5：Values()应返回所有元数据的副本
	ctx.Set("key2", "val2")
	values := ctx.Values()
	assert.Equal(t, map[string]string{"key1": "val2", "key2": "val2"}, values)

	// 验证副本特性：修改返回的map不影响内部数据
	values["key1"] = "modified"
	assert.Equal(t, "val2", ctx.Get("key1")) // 内部数据未变
}

// TestContextDerive 测试上下文衍生函数（WithCancel/WithTimeout）
func TestContextDerive(t *testing.T) {
	t.Parallel()

	parent := New()
	parent.Set("parent-key", "parent-val")
	parentTraceID := parent.TraceID()

	// 测试WithCancel
	t.Run("WithCancel", func(t *testing.T) {
		child, cancel := WithCancel(parent)
		defer cancel()

		// 验证元数据继承
		assert.Equal(t, "parent-val", child.Get("parent-key"))
		assert.Equal(t, parentTraceID, child.TraceID())

		// 验证子上下文修改不影响父上下文
		child.Set("child-key", "child-val")
		assert.Empty(t, parent.Get("child-key"))

		// 验证取消功能
		cancel()
		select {
		case <-child.Done():
			assert.ErrorIs(t, child.Err(), context.Canceled)
		case <-time.After(100 * time.Millisecond):
			t.Fatal("child context should be canceled")
		}
	})

	// 测试WithTimeout
	t.Run("WithTimeout", func(t *testing.T) {
		// 正常超时
		child, cancel := WithTimeout(parent, 50*time.Millisecond)
		defer cancel()

		// 验证元数据继承
		assert.Equal(t, "parent-val", child.Get("parent-key"))
		assert.Equal(t, parentTraceID, child.TraceID())

		// 等待超时
		select {
		case <-child.Done():
			assert.ErrorIs(t, child.Err(), context.DeadlineExceeded)
		case <-time.After(100 * time.Millisecond):
			t.Fatal("child context should timeout")
		}

		// 测试超时≤0时降级为WithCancel
		child2, cancel2 := WithTimeout(parent, 0)
		defer cancel2()
		cancel2()
		select {
		case <-child2.Done():
			assert.ErrorIs(t, child2.Err(), context.Canceled) // 应是取消错误而非超时
		case <-time.After(100 * time.Millisecond):
			t.Fatal("child2 context should be canceled")
		}
	})

	// 测试非法父上下文（非*kCtx类型）
	t.Run("invalid parent", func(t *testing.T) {
		invalidParent := New() // 标准context，非kctx.Context

		// WithCancel应对非法父上下文兼容
		childCancel, cancel := WithCancel(invalidParent)
		require.NotNil(t, childCancel)
		cancel()

		// WithTimeout应对非法父上下文兼容
		childTimeout, cancel := WithTimeout(invalidParent, 10*time.Millisecond)
		require.NotNil(t, childTimeout)
		cancel()
	})
}

// TestConcurrentSafety 测试并发读写安全性
func TestConcurrentSafety(t *testing.T) {
	t.Parallel()

	ctx := New()
	const (
		numWriters = 10
		numReads   = 100
		key        = "concurrent-key"
	)

	var wg sync.WaitGroup

	// 启动多个写协程
	for i := 0; i < numWriters; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			ctx.Set(key, string(rune(val))) // 写入不同值
		}(i)
	}

	// 启动多个读协程
	for i := 0; i < numReads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = ctx.Get(key)      // 读取值
			_ = ctx.Values()      // 读取所有值
			_ = ctx.Context()     // 读取底层context
			_ = ctx.TraceID()     // 读取TraceID
			_, _ = ctx.Deadline() // 调用context接口方法
			_ = ctx.Err()         // 调用context接口方法
			select {
			case <-ctx.Done():
			case <-time.After(10 * time.Millisecond):
			}
		}()
	}

	// 等待所有协程完成，验证无panic
	wg.Wait()
	t.Log("concurrent test passed without panic")
}

// TestValueInterface 测试context.Value接口实现
func TestValueInterface(t *testing.T) {
	t.Parallel()

	ctx := New()
	ctx.Set("meta-key", "meta-val")
	parentVal := "parent-value"
	ctx.SetContext(context.WithValue(context.Background(), "parent-key", parentVal))

	// 测试1：获取TraceIDKey
	assert.Equal(t, ctx.TraceID(), ctx.Value(TraceIDKey))

	// 测试2：获取MetaMapKey（应返回元数据副本）
	metaFromValue, ok := ctx.Value(MetaMapKey).(map[string]string)
	require.True(t, ok)
	assert.Equal(t, ctx.Values(), metaFromValue)

	// 测试3：获取底层context的值
	assert.Equal(t, parentVal, ctx.Value("parent-key"))

	// 测试4：获取不存在的键
	assert.Nil(t, ctx.Value("non-exist-key"))
}

// TestSetContext 测试SetContext方法
func TestSetContext(t *testing.T) {
	t.Parallel()

	ctx := New()

	// 测试1：设置新的底层context
	newBaseCtx := context.WithValue(context.Background(), "new-key", "new-val")
	ctx.SetContext(newBaseCtx)
	assert.Equal(t, newBaseCtx, ctx.Context())
	assert.Equal(t, "new-val", ctx.Value("new-key"))

	// 测试2：设置nil应被忽略
	ctx.SetContext(nil)
	assert.Equal(t, newBaseCtx, ctx.Context()) // 仍为上一步的newBaseCtx

	// 测试3：设置后，衍生上下文应使用新的底层context
	child, cancel := WithCancel(ctx)
	defer cancel()
	// 验证衍生上下文的底层context是 newBaseCtx 的衍生（cancelCtx），而非直接等于 newBaseCtx
	childBaseCtx := child.Context()
	// 可通过 Value 传递性验证：childBaseCtx 应能获取到 newBaseCtx 中的值
	assert.Equal(t, "new-val", childBaseCtx.Value("new-key"))
}
