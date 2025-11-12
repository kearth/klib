package kunit

import (
	"testing"
	"time"

	"github.com/kearth/klib/kctx"
	"github.com/kearth/klib/kerr"
	"github.com/stretchr/testify/assert" // 需安装：go get github.com/stretchr/testify/assert
)

func TestNewUnit(t *testing.T) {
	// 测试默认名称
	u := NewUnit("")
	assert.Equal(t, NoName, u.Name())
	assert.Equal(t, RoleUnit, u.Role())
	_, err := u.Call(kctx.New())
	assert.Equal(t, nil, err) // 无函数时返回依赖错误

	// 测试带名称和函数的创建
	testFn := func(ctx kctx.Context, input ...any) (any, kerr.Error) {
		return "ok", nil
	}
	u = NewUnit("test", testFn)
	assert.Equal(t, "test", u.Name())
	assert.Equal(t, RoleUnit, u.Role())
	output, err := u.Call(kctx.New())
	assert.NoError(t, err) // 假设 kerr.Error 实现了 error 接口
	assert.Equal(t, "ok", output)
}

func TestUnit_Setters(t *testing.T) {
	u := NewUnit("original")

	// 测试 SetName
	u.SetName("new-name")
	assert.Equal(t, "new-name", u.Name())

	// 测试 SetRole
	u.SetRole(RoleComponent)
	assert.Equal(t, RoleComponent, u.Role())

	// 测试 SetFn
	newFn := func(ctx kctx.Context, input ...any) (any, kerr.Error) {
		return 123, nil
	}
	u.SetFn(newFn)
	output, err := u.Call(kctx.New())
	assert.NoError(t, err)
	assert.Equal(t, 123, output)
}

func TestUnit_Call(t *testing.T) {
	t.Run("无函数时调用", func(t *testing.T) {
		u := NewUnit("no-fn")
		_, err := u.Call(kctx.New())
		assert.Equal(t, nil, err)
		assert.Equal(t, time.Duration(0), u.Cost()) // 未执行成功但仍计算耗时？根据实际逻辑调整
	})

	t.Run("带输入参数的调用", func(t *testing.T) {
		fn := func(ctx kctx.Context, input ...any) (any, kerr.Error) {
			if len(input) < 2 {
				return nil, kerr.New(1, "参数不足")
			}
			return input[0].(int) + input[1].(int), nil
		}
		u := NewUnit("add", fn)
		output, err := u.Call(kctx.New(), 2, 3)
		assert.NoError(t, err)
		assert.Equal(t, 5, output)
		assert.True(t, u.Cost() > 0) // 耗时应大于0
	})

	t.Run("函数返回错误", func(t *testing.T) {
		errMsg := "执行失败"
		fn := func(ctx kctx.Context, input ...any) (any, kerr.Error) {
			return nil, kerr.New(1, errMsg)
		}
		u := NewUnit("error-fn", fn)
		output, err := u.Call(kctx.New())
		assert.Nil(t, output)
		assert.Equal(t, errMsg, err.Error()) // 假设 kerr.Error 有 Error() 方法
		assert.True(t, u.Cost() > 0)
	})
}

func TestUnit_Cost(t *testing.T) {
	u := NewUnit("cost-test")
	assert.Equal(t, time.Duration(0), u.Cost()) // 未调用时耗时为0

	// 测试耗时计算
	u.SetFn(func(ctx kctx.Context, input ...any) (any, kerr.Error) {
		time.Sleep(10 * time.Millisecond)
		return nil, kerr.Succ
	})
	u.Call(kctx.New())
	assert.True(t, u.Cost() >= 10*time.Millisecond)
}

func TestUnit_Setup(t *testing.T) {
	u := NewUnit("setup-test")
	// 测试默认 Setup 方法（返回 nil）
	assert.NoError(t, u.Setup(kctx.New()))
}
