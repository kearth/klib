# CHANGELOG
所有显著的变更都会记录在本文件中。

## 格式规范
- 版本号按语义化版本（vX.Y.Z）排序，最新版本在最上方
- 每个版本包含以下分类（根据实际情况取舍）：
  - `Added`：新增功能
  - `Changed`：现有功能变更
  - `Fixed`：bug 修复
  - `Removed`：移除的功能
  - `Deprecated`：即将移除的功能（弃用提示）
  - `Security`：安全相关修复

---
## [v0.2.2] - 2025-11-13
### Fixed
- 修复 Makefile 的命令


## [v0.2.1] - 2025-11-13
### Fixed
- 修复 相关`release`命令变更

## [v0.2.0] - 2025-11-13
### Fixed
- 修复 相关命令变更

## [v0.1.2] - 2025-11-13
### Added
- 新增 CHANGELOG.md 文件，记录版本变更历史
- 新增 Makefile 的 changelog 相关能力和 release 版本发布能力




## [v0.1.0] - 2025-11-15
### Added
- 新增 `make release` 命令，支持自动创建 GitHub Release
- 新增 `install-gh` 命令，自动安装 GitHub CLI 工具
### Changed
- 优化 `version` 命令输出格式，兼容带 v 前缀的版本号
- 修复 `upgrade_version` 内部注释干扰 shell 执行的问题
### Fixed
- 解决 `make minor/patch/major` 时 Git Tag 重复校验失败的问题

## [v0.0.6] - 2025-11-10
### Added
- 初始发布 klib 库，包含 klog/kctx/kerr/kutil/kunit 模块
- 新增文档生成命令 `gen-docs`（依赖 gomarkdoc）
- 新增测试命令 `make test`（含 race 检测）
### Changed
- 优化模块目录结构，统一导出接口格式

## [v0.0.5] - 2025-11-05
### Added
- 新增 kunit 模块基础功能（任务调度、并发控制）
### Fixed
- 修复 kctx 模块上下文传递时的内存泄漏问题
