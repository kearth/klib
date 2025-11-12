# 项目名称
PROJECT_NAME := klib
# 模块路径（替换为你的实际模块路径）
MODULE := github.com/kearth/klib
# 文档生成目录
DOCS_DIR := docs
# 默认日志包文档生成路径（可根据模块扩展）
KLOG_DOC := $(DOCS_DIR)/klog.md
KCTX_DOC := $(DOCS_DIR)/kctx.md
KERR_DOC := $(DOCS_DIR)/kerr.md
KUTIL_DOC := $(DOCS_DIR)/kutil.md
KUNIT_DOC := $(DOCS_DIR)/kunit.md

# 版本管理核心配置
VERSION_FILE := version.go  # 版本文件路径
DEFAULT_BRANCH := master      # 仓库默认分支（根据实际调整）

# --------------- 版本管理（打 Tag）---------------
# 显示当前版本（从 version.go 提取代码版本 + 从Git提取最新Tag）
# 逻辑：1. 提取版本文件中的版本 2. 容错处理 3. 读取Git Tag 4. 格式化输出
version:
	@CODE_VERSION=$$(grep -E 'return "' $(VERSION_FILE) 2>/dev/null | sed -E 's/.*return "(v?[0-9]+\.[0-9]+\.[0-9]+)".*/\1/'); \
	if [ -z "$$CODE_VERSION" ]; then \
		CODE_VERSION="未知（版本文件异常）"; \
	fi ; \
	TAG_VERSION=$$(git describe --abbrev=0 --tags 2>/dev/null || echo "无版本Tag") ; \
	FORMATTED_VERSION=$$(echo "$$CODE_VERSION" | sed 's/^v//') ; \
	echo "==================== 版本信息 ===================="; \
	echo "当前代码版本: $$CODE_VERSION"; \
	echo "当前最新Tag:  $$TAG_VERSION"; \
	echo "==================================================";


# 新增：语义化版本升级（补丁版本：修复bug，vX.Y.Z → vX.Y.Z+1）
patch:
	@$(call upgrade_version,patch)

# 新增：语义化版本升级（次版本：新增兼容功能，vX.Y.Z → vX.Y+1.0）
minor:
	@$(call upgrade_version,minor)

# 新增：语义化版本升级（主版本：不兼容变更，vX.Y.Z → vX+1.0.0）
major:
	@$(call upgrade_version,major)

# 打新 Tag（示例：make tag VERSION=v0.1.0）
# 支持语义化版本（如 v0.1.0、v1.2.3-beta）
tag:
	@if [ -z "$(VERSION)" ]; then \
		echo "请指定版本号，格式: make tag VERSION=v0.1.0"; \
		exit 1; \
	fi
	@if ! echo "$(VERSION)" | grep -qE '^v[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9]+)?$$'; then \
		echo "版本号格式错误，需符合语义化版本（如 v0.1.0）"; \
		exit 1; \
	fi
	@git tag -a $(VERSION) -m "Release $(VERSION)"
	@echo "已创建本地Tag: $(VERSION)"
	@echo "提示: 推送至远程仓库执行: make push-tag VERSION=$(VERSION)"

# 推送 Tag 到远程仓库
push-tag:
	@if [ -z "$(VERSION)" ]; then \
		echo "请指定版本号，格式: make push-tag VERSION=v0.1.0"; \
		exit 1; \
	fi
	@git push origin $(VERSION)
	@echo "已推送Tag $(VERSION)至远程仓库"

# --------------- 文档生成与更新 ---------------
# 安装文档生成工具（gomarkdoc）
install-doc-tool:
	@echo "安装文档生成工具 gomarkdoc..."
	go install github.com/princjef/gomarkdoc/cmd/gomarkdoc@latest

# 生成所有模块文档（依赖 gomarkdoc）
gen-docs:
	@echo "生成文档至 $(DOCS_DIR) 目录..."
	mkdir -p $(DOCS_DIR)
	# 生成 klog 模块文档
	gomarkdoc -o $(KLOG_DOC) $(MODULE)/klog
	# 生成 kctx 模块文档
	gomarkdoc -o $(KCTX_DOC) $(MODULE)/kctx
	# 生成 kerr 模块文档
	gomarkdoc -o $(KERR_DOC) $(MODULE)/kerr
	# 生成 kutil 模块文档
	gomarkdoc -o $(KUTIL_DOC) $(MODULE)/kutil
	# 生成 kunit 模块文档
	gomarkdoc -o $(KUNIT_DOC) $(MODULE)/kunit
	@echo "文档生成完成"

# 查看文档（本地预览）
view-docs:
	@echo "打开文档目录: $(DOCS_DIR)"
	# 不同系统打开命令（根据需要注释/启用）
	open $(DOCS_DIR)  # MacOS
	# xdg-open $(DOCS_DIR)  # Linux
	# start $(DOCS_DIR)    # Windows

# --------------- 辅助命令 ---------------
# 运行测试（含 race 检测）
test:
	go test -race ./... -v

# 清理生成的文档和临时文件
clean:
	rm -rf $(DOCS_DIR)/*.md
	go clean

# 帮助信息
help:
	@echo "可用命令:"
	@echo "  版本管理（推荐自动升级）:"
	@echo "    make version          显示当前代码版本和最新Tag"
	@echo "    make patch            升级补丁版本（vX.Y.Z → vX.Y.Z+1）"
	@echo "    make minor            升级次版本（vX.Y.Z → vX.Y+1.0）"
	@echo "    make major            升级主版本（vX.Y.Z → vX+1.0.0）"
	@echo "  版本管理（手动打Tag，兼容旧用法）:"
	@echo "    make tag VERSION=vX.Y.Z  创建本地版本Tag"
	@echo "    make push-tag VERSION=vX.Y.Z  推送Tag至远程"
	@echo "  文档相关:"
	@echo "    make install-doc-tool  安装文档生成工具"
	@echo "    make gen-docs         生成所有模块文档"
	@echo "    make view-docs        打开文档目录预览"
	@echo "  其他:"
	@echo "    make test             运行测试（含race检测）"
	@echo "    make clean            清理文档和临时文件"
	@echo "    make help             显示帮助信息"

# 默认命令：显示帮助
.DEFAULT_GOAL := help

# --------------- 内部函数：版本升级核心逻辑（无需修改）---------------
# --------------- 内部函数：版本升级核心逻辑（无shell注释，避免干扰）---------------
define upgrade_version
	CURRENT_VERSION=$$(grep -E 'return "' $(VERSION_FILE) | sed -E 's/.*return "(v?[0-9]+\.[0-9]+\.[0-9]+)".*/\1/' | sed 's/^v//') ; \
	if [ -z "$$CURRENT_VERSION" ]; then \
		echo "错误：未在 $(VERSION_FILE) 中找到有效版本号"; \
		exit 1; \
	fi ; \
	IFS='.' read -r MAJOR MINOR PATCH <<< "$$CURRENT_VERSION" ; \
	case "$1" in \
		major) NEW_MAJOR=$$((MAJOR+1)); NEW_MINOR=0; NEW_PATCH=0 ;; \
		minor) NEW_MAJOR=$$MAJOR; NEW_MINOR=$$((MINOR+1)); NEW_PATCH=0 ;; \
		patch) NEW_MAJOR=$$MAJOR; NEW_MINOR=$$MINOR; NEW_PATCH=$$((PATCH+1)) ;; \
	esac ; \
	NEW_VERSION="$$NEW_MAJOR.$$NEW_MINOR.$$NEW_PATCH" ; \
	NEW_TAG="v$$NEW_VERSION" ; \
	if ! git diff --quiet --exit-code; then \
		echo "错误：工作区存在未提交的变更，请先提交或 stash"; \
		exit 1; \
	fi ; \
	sed -i '' -E "s/return \"v?[0-9]+\.[0-9]+\.[0-9]+\"/return \"$$NEW_TAG\"/" $(VERSION_FILE) ; \
	echo "✅ 已更新版本：v$$CURRENT_VERSION → $$NEW_TAG" ; \
	git add $(VERSION_FILE) ; \
	git commit -m "$$NEW_TAG" ; \
	git tag -a $$NEW_TAG -m "Release $$NEW_TAG" ; \
	git push origin $(DEFAULT_BRANCH) ; \
	git push origin $$NEW_TAG ; \
	echo "✅ 已完成：提交代码 + 推送Tag $$NEW_TAG 至远程仓库" ;
endef