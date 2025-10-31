package kutil

import "github.com/gogf/gf/v2/util/gconv"

// If 实现两种功能：
// 1. 三元运算：当传入 trueVal 和 falseVal 时，返回对应值
// 2. 条件执行：当传入 f（或 f、g）时，执行对应函数
// 注：两种功能通过参数类型区分，互斥使用
// 注意：Go 语言中，函数参数会先求值再传递，因此如果传入的函数有副作用，需要注意顺序问题。
func If[T any](condition bool, trueVal any, falseVal ...any) T {
	// 场景1：三元运算（trueVal 和 falseVal[0] 为同类型值）
	if len(falseVal) == 1 {
		if condition {
			return trueVal.(T)
		}
		return falseVal[0].(T)
	}

	// 场景2：条件执行（trueVal 为函数，falseVal 可选为函数）
	if condition {
		trueVal.(func())()
	} else if len(falseVal) > 0 {
		falseVal[0].(func())()
	}

	// 执行函数时无返回值，返回 T 的零值（不影响使用）
	var zero T
	return zero
}

// As 尝试将任意类型 v 转换为 T 类型，内部基于 gconv 实现，支持：
// 1. 基本类型互转（如 string ↔ int、bool ↔ float 等）；
// 2. 复合类型转换（如 slice ↔ array、map ↔ struct 等，需符合 gconv 规则）。
// 返回转换后的值和是否成功（true 表示转换有效）。
func As[T any](v any) (T, bool) {
	var result T
	if err := gconv.Scan(v, &result); err != nil {
		return result, false
	}
	return result, true
}
