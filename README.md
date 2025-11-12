# Klib

[![GitHub License](https://img.shields.io/github/license/kearth/klib.svg?style=flat-square)](https://github.com/kearth/klib/blob/master/LICENSE)
[![Release](https://img.shields.io/github/release/kearth/klib.svg?style=flat-square)](https://github.com/kearth/klib/releases)
[![GitHub Tag](https://img.shields.io/github/v/tag/kearth/klib?label=version&style=flat-square)](https://github.com/kearth/klib/tags)
[![Go Version](https://img.shields.io/github/go-mod/go-version/kearth/tea.svg?style=flat-square)](https://github.com/kearth/tea)
[![Go Reference](https://pkg.go.dev/badge/github.com/kearth/klib.svg)](https://pkg.go.dev/github.com/kearth/klib?tab=doc)
[![Last Commit](https://img.shields.io/github/last-commit/kearth/klib.svg?style=flat-square)](https://github.com/kearth/klib/commits/master)



Klib 是一个专为个人学习目的而创建的Go语言库。它的主要目标是提供一些更加方便的用法，以便于您更轻松地学习和使用Go编程语言。

## 特性

- 提供一组简单而实用的Go函数和工具。
- 优化了Go语言的一些常见用法，以提高开发效率。
- 帮助您更深入地了解Go语言的工作原理和最佳实践。

## 使用示例

```go
import (
    "github.com/kearth/klib/kctx"
	"github.com/kearth/klib/klog"
	"github.com/kearth/klib/kutil"
)

func main() {
	// 初始化日志
	ctx := kctx.New()
	klog.Init()

	// 测试条件
	condition := 5
	fn := kutil.If[func()](condition == 5, func() {
		klog.ColorPrint(ctx, klog.Cyan, "条件为5")
	}, func() {
		klog.ColorPrint(ctx, klog.Cyan, "条件不为5")
	})

	// 执行函数
	fn()
}
```
结果：
```
条件为5
```

## 安装

您可以使用以下命令来安装GoLib 学习库：
```go
go get -u github.com/kearth/klib
```

希望你们喜欢这个库！如果您有任何问题或建议，请随时提出。

## 详细文档
- [kctx 上下文管理](docs/kctx.md)
- [kerr 错误处理](docs/kerr.md)
- [kutil 工具函数](docs/kutil.md)
- [kunit 单元容器](docs/kunit.md)
- [klog 日志管理](docs/klog.md)
