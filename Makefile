# 项目名称
PROJECT_NAME := klib
# 模块路径（替换为你的实际模块路径）
MODULE := github.com/kearth/klib
SHORT_MODULE := kearth/klib
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
	@CODE_VERSION=$$(grep -E 'return "' $(VERSION_FILE) 2>/dev/null | sed -E 's/.*return "(v?[0-9]+\.[0-9]+\.[0-9]+)".*/\1/'); \
	echo "当前代码版本: $$CODE_VERSION"; \
	if ! echo "$$CODE_VERSION" | grep -qE '^v[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9]+)?$$'; then \
		echo "❌ 版本号格式错误，需符合语义化版本（如 v0.1.0、v1.2.3-beta）"; \
		exit 1; \
	fi; \
	if ! git diff --quiet --exit-code; then \
		echo "❌ 工作区存在未提交的变更，请先提交或 stash"; \
		exit 1; \
	fi; \
	if git rev-parse "$$CODE_VERSION" >/dev/null 2>&1; then \
		echo "ℹ️  本地已存在 tag: $$CODE_VERSION，跳过创建"; \
	else \
		echo "🏷️  创建本地 tag: $$CODE_VERSION"; \
		git tag -a "$$CODE_VERSION" -m "Release $$CODE_VERSION"; \
	fi; \
	echo "📤 推送 tag 到远程..."; \
	git push origin "$$CODE_VERSION" || (echo "⚠️ 推送失败，请检查权限或远程状态" && exit 1); \
	echo "✅ Tag 操作完成：$$CODE_VERSION";

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

# --------------- GitHub Release 管理 ---------------
# 创建 GitHub Release（需先打 Tag，支持自动构建产物+上传）
# 用法：make release VERSION=v0.1.0 [BUILD_BIN=true]
release:
	@if ! command -v gh >/dev/null 2>&1; then \
		echo "❌ 未安装 GitHub CLI（gh），请先执行 'make install-gh' 安装"; \
		exit 1; \
	fi; \
	if ! gh auth status >/dev/null 2>&1; then \
		echo "❌ gh 未登录或无仓库权限，请执行 'gh auth login' 登录授权"; \
		exit 1; \
	fi; \
	if [ -z "$(VERSION)" ]; then \
		echo "❌ 请指定版本号，格式: make release VERSION=v0.1.0"; \
		exit 1; \
	fi; \
	if ! echo "$(VERSION)" | grep -qE '^v[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9]+)?$$'; then \
		echo "❌ 版本号格式错误，需符合语义化版本（如 v0.1.0、v1.2.3-beta）"; \
		exit 1; \
	fi; \
	if ! git rev-parse $(VERSION) >/dev/null 2>&1; then \
		echo "❌ 本地不存在 Tag $(VERSION)，请先执行 make tag VERSION=$(VERSION) 创建"; \
		exit 1; \
	fi; \
	if ! git ls-remote --tags origin $(VERSION) >/dev/null 2>&1; then \
		echo "❌ 远程不存在 Tag $(VERSION)，请先执行 make push-tag VERSION=$(VERSION) 推送"; \
		exit 1; \
	fi; \
	RELEASE_NOTES="" ; \
	if [ -f "CHANGELOG.md" ]; then \
		RELEASE_NOTES=$$(awk '/^## \['"$(VERSION)"'\]/{flag=1;next}/^## \[v/{flag=0}flag' CHANGELOG.md | sed '/^$$/d' | sed 's/^[[:space:]]*//'); \
		if [ -z "$$RELEASE_NOTES" ]; then \
			RELEASE_NOTES="Release $(VERSION)"; \
		fi; \
	fi; \
	echo "🚀 开始创建 GitHub Release: $(VERSION)" ;\
	echo "仓库: $(SHORT_MODULE)" ;\
	echo "标题: $(PROJECT_NAME) $(VERSION)" ;\
	echo "变更记录：$$RELEASE_NOTES" ;\
	TMP_NOTES=$$(mktemp) ;\
	echo "$$RELEASE_NOTES" > "$$TMP_NOTES" ;\
	gh release create $(VERSION) \
		--title "$(PROJECT_NAME) $(VERSION)" \
		--notes-file "$$TMP_NOTES" \
		--repo $(SHORT_MODULE);\
	rm -f "$$TMP_NOTES" ;\
	echo "🎉 GitHub Release 创建完成！" ;\
	echo "🔗 查看地址：https://github.com/kearth/klib/releases/tag/$(VERSION)"

# 查看已发布的 Release
list-releases:
	@echo "📋 已发布的 GitHub Release："
	gh release list --repo $(MODULE)


# --------------- 依赖安装（新增）---------------
# 安装 GitHub CLI（gh）：自动检测系统，无则安装
install-gh:
	@echo "🔍 检查是否已安装 GitHub CLI（gh）..."
	@if ! command -v gh >/dev/null 2>&1; then \
		echo "❌ 未找到 gh，开始安装..."; \
		UNAME_S=$$(uname -s); \
		if [ "$$UNAME_S" = "Darwin" ]; then \
			if command -v brew >/dev/null 2>&1; then \
				brew install gh; \
			else \
				echo "❌ 未找到 Homebrew，请先安装 Homebrew：/bin/bash -c \"$$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)\""; \
				exit 1; \
			fi; \
		elif [ "$$UNAME_S" = "Linux" ]; then \
			if command -v apt >/dev/null 2>&1; then \
				sudo apt update && sudo apt install -y gh; \
			elif command -v dnf >/dev/null 2>&1; then \
				sudo dnf install -y gh; \
			elif command -v yum >/dev/null 2>&1; then \
				sudo yum install -y gh; \
			elif command -v pacman >/dev/null 2>&1; then \
				sudo pacman -S --noconfirm gh; \
			else \
				echo "❌ 不支持的 Linux 包管理器，请手动安装 gh：https://cli.github.com/manual/installation"; \
				exit 1; \
			fi; \
		elif [ "$$UNAME_S" = "Windows_NT" ]; then \
			echo "ℹ️ Windows 系统请通过 Chocolatey 安装：choco install gh"; \
			echo "或手动下载：https://github.com/cli/cli/releases/latest/download/gh_windows_amd64.msi"; \
			exit 1; \
		fi; \
		echo "✅ gh 安装完成！请执行 'gh auth login' 登录授权"; \
	else \
		echo "✅ gh 已安装（版本：$$(gh --version | grep -E 'gh version' | awk '{print $$3}')）"; \
	fi

publish:
	@make changelog
	@make tag
	@make release
# 帮助信息
help:
	@echo "可用命令:"
	@echo "  版本管理（推荐自动升级）:"
	@echo "    make version          显示当前代码版本和最新Tag"
	@echo "    make patch            升级补丁版本（vX.Y.Z → vX.Y.Z+1）"
	@echo "    make minor            升级次版本（vX.Y.Z → vX.Y+1.0）"
	@echo "    make major            升级主版本（vX.Y.Z → vX+1.0.0）"
	@echo "    make tag 创建本地版本Tag,推送Tag至远程, 格式: vX.Y.Z"
	@echo "  GitHub Release 管理（需先执行 make install-gh + gh auth login）:"
	@echo "    make install-gh       安装GitHub CLI（gh）工具"
	@echo "    make release VERSION=vX.Y.Z "
	@echo "    make list-releases    查看所有已发布的GitHub Release"
	@echo "  文档相关:"
	@echo "    make install-doc-tool  安装文档生成工具"
	@echo "    make gen-docs         生成所有模块文档"
	@echo "    make view-docs        打开文档目录预览"
	@echo "  其他:"
	@echo "    make test             运行测试（含race检测）"
	@echo "    make clean            清理文档和临时文件"
	@echo "    make help             显示帮助信息"
	@echo "    make publish          发布新的版本（自动升级、创建Tag、发布Release）"
	@echo "  快速 Commit 命令（简化+规范提交）:"
	@echo "    make commit-<类型> MSG=\"描述\"  快速提交（如：make commit-feat MSG=\"新增功能\"）"
	@echo "    make commit-help          查看快速 Commit 命令说明"
	@echo "  Git Commit 规范（强制提交格式）:"
	@echo "    make install-commit-hooks  安装提交规范钩子（自动校验格式）"
	@echo "    make uninstall-commit-hooks  卸载提交规范钩子"

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
	sed -i '' -E "s/return \"v?[0-9]+\.[0-9]+\.[0-9]+\"/return \"$$NEW_TAG\"/" $(VERSION_FILE) ; \
	git add $(VERSION_FILE) ; \
	git commit -m "$$NEW_TAG" ; 
	echo "✅ 已更新版本：v$$CURRENT_VERSION → $$NEW_TAG" ; 
endef


# --------------- CHANGELOG 自动管理（优化：自动补中文前缀）---------------
changelog:
	@export LC_ALL=en_US.UTF-8; export LANG=en_US.UTF-8; \
	VERSION_FILE=version.go; \
	echo "🔍 读取 $$VERSION_FILE 中的最新版本号..."; \
	NEW_VERSION=$$(grep -E 'return "' $$VERSION_FILE | sed -E 's/.*return "(v?[0-9]+\.[0-9]+\.[0-9]+)".*/\1/' | tr -d '"'); \
	echo "✅ 最新版本：$$NEW_VERSION"; \
	if [ -z "$$NEW_VERSION" ] || ! echo "$$NEW_VERSION" | grep -qE '^v[0-9]+\.[0-9]+\.[0-9]+'; then echo "❌ version 解析失败"; exit 1; fi; \
	CURRENT_DATE=$$(date +%Y-%m-%d); \
	echo "📅 当前日期：$$CURRENT_DATE"; \
	LAST_TAG=$$(git describe --tags --abbrev=0 HEAD^ 2>/dev/null); COMMIT_RANGE=$$LAST_TAG..HEAD; \
	extract_commits() { type=$$1; case $$type in feat) prefix="- 新增";; fix) prefix="- 修复";; chore) prefix="- 优化";; refactor) prefix="- 重构/移除";; docs) prefix="- 更新";; test) prefix="- 完善";; security) prefix="- 加固";; deprecated) prefix="- 标记弃用";; esac; git log $$COMMIT_RANGE --pretty=format:"%s" --grep="^$$type:" 2>/dev/null | sed "s#^$$type: ##" | sort -u | grep -v '^$$' | sed "s#^#$$prefix #"; }; \
	ADDED=$$(extract_commits feat); \
	echo "📝 新增功能：$$ADDED"; \
	CHANGED=$$(printf "%s\n%s\n%s" "$$(extract_commits chore)" "$$(extract_commits docs)" "$$(extract_commits test)" | sort -u | grep -v '^$$'); \
	echo "🔧 优化/更新：$$CHANGED"; \
	FIXED=$$(extract_commits fix); \
	echo "🔧 修复问题：$$FIXED"; \
	REMOVED=$$(extract_commits refactor); \
	echo "🔧 重构/移除：$$REMOVED"; \
	SECURITY=$$(extract_commits security); \
	echo "🔧 加固：$$SECURITY"; \
	DEPRECATED=$$(extract_commits deprecated); \
	echo "🔧 标记弃用：$$DEPRECATED"; \
	SECTIONS=""; \
	if [ -n "$$ADDED" ]; then SECTIONS="$$SECTIONS\n### Added\n$$ADDED\n"; fi; \
	if [ -n "$$CHANGED" ]; then SECTIONS="$$SECTIONS\n### Changed\n$$CHANGED\n"; fi; \
	if [ -n "$$FIXED" ]; then SECTIONS="$$SECTIONS\n### Fixed\n$$FIXED\n"; fi; \
	if [ -n "$$REMOVED" ]; then SECTIONS="$$SECTIONS\n### Removed\n$$REMOVED\n"; fi; \
	if [ -n "$$SECURITY" ]; then SECTIONS="$$SECTIONS\n### Security\n$$SECURITY\n"; fi; \
	if [ -n "$$DEPRECATED" ]; then SECTIONS="$$SECTIONS\n### Deprecated\n$$DEPRECATED\n"; fi; \
	if [ -z "$$SECTIONS" ]; then \
		echo "ℹ️ $$NEW_VERSION 无显著变更，不生成版本块"; \
		exit 0; \
	fi; \
	NEW_VERSION_BLOCK=$$(printf "## [%s] - %s%s" "$$NEW_VERSION" "$$CURRENT_DATE" "$$SECTIONS");\
	echo "📝 新变更记录：$$NEW_VERSION_BLOCK"; \
	if [ ! -f CHANGELOG.md ]; then echo -e "# CHANGELOG\n所有显著的变更都会记录在本文件中。\n\n---\n" > CHANGELOG.md; fi; \
	echo "🔍 检查 $$NEW_VERSION 是否已存在于 CHANGELOG.md..."; \
	if [ -f CHANGELOG.md ] && grep -q "$$NEW_VERSION" CHANGELOG.md; then \
		echo "⚠️ $$NEW_VERSION 的变更记录已存在，无需重复生成"; \
		exit 0; \
	fi; \
	if grep -q "## [$$NEW_VERSION]" CHANGELOG.md; then echo "⚠️ $$NEW_VERSION 已存在"; exit 0; fi; \
	if ! grep -q "^---" CHANGELOG.md; then echo "---" >> CHANGELOG.md; fi; \
	printf "%b" "/^---/a\n$$NEW_VERSION_BLOCK\n.\nw\nq\n" | ed -s CHANGELOG.md >/dev/null; \
	echo "✅ CHANGELOG 更新成功：$$NEW_VERSION"; head -n 10 CHANGELOG.md | grep -E '##|\- ' | sed 's/^/ /'; \
	git add -A ; \
	git commit -m "Update CHANGELOG.md" ; 

	

# --------------- 快速 Commit 命令（简化提交操作）---------------
# 定义通用 Commit 函数（内部使用，无需手动调用）
# 注意：函数内部命令前加 @，抑制 Makefile 回显
# 定义通用 Commit 函数（内部使用，无需手动调用）

define commit_func
	@COMMIT_MSG="$(1): $(filter-out $@,$(MAKECMDGOALS))"; \
	if [ -z "$$COMMIT_MSG" ]; then \
		echo "❌ 请提供提交描述，例如：make commit-$(1) 新增模块"; \
		exit 1; \
	fi; \
	MSG_LEN=$$(echo -n "$$COMMIT_MSG" | wc -m); \
	if [ $$MSG_LEN -lt 10 ]; then \
		echo "❌ 提交描述过短！至少 10 个字符（当前：$$MSG_LEN 个）"; \
		exit 1; \
	fi; \
	if git diff --cached --quiet && git diff --quiet; then \
		echo "⚠️  无文件变更，将跳过提交"; \
		exit 0; \
	fi; \
	git add -A; \
	echo "📤 提交信息：$$COMMIT_MSG"; \
	if git commit -m "$$COMMIT_MSG"; then \
		echo "✅ 提交成功！"; \
	else \
		echo "❌ 提交失败，请检查错误信息"; \
		exit 1; \
	fi;
endef

# -----------------------------
# 🧩 具体提交类型命令
# -----------------------------
commit-feat:      ## 新功能提交
	@$(call commit_func,feat)

commit-fix:       ## 修复问题提交
	@$(call commit_func,fix)

commit-chore:     ## 杂项提交（构建/依赖/配置）
	@$(call commit_func,chore)

commit-refactor:  ## 代码重构
	@$(call commit_func,refactor)

commit-docs:      ## 文档更新
	@$(call commit_func,docs)

commit-test:      ## 测试相关
	@$(call commit_func,test)

commit-security:  ## 安全修复
	@$(call commit_func,security)

commit-deprecated:## 废弃/移除功能
	@$(call commit_func,deprecated)

# 快速提交帮助
commit-help:
	@echo "📋 快速 Commit 命令使用说明"
	@echo "=========================="
	@echo "格式：make commit-<类型> \"描述信息\""
	@echo "支持的类型及含义："
	@echo "  commit-feat      新增功能（对应 CHANGELOG Added）"
	@echo "  commit-fix       修复 Bug（对应 CHANGELOG Fixed）"
	@echo "  commit-chore     功能优化/构建配置变更（对应 CHANGELOG Changed）"
	@echo "  commit-refactor  代码重构/移除功能（对应 CHANGELOG Removed）"
	@echo "  commit-docs      文档更新（对应 CHANGELOG Changed）"
	@echo "  commit-test      测试相关（新增/修改测试用例）"
	@echo "  commit-security  安全相关修复（对应 CHANGELOG Security）"
	@echo "  commit-deprecated 标记弃用功能（对应 CHANGELOG Deprecated）"
	@echo "=========================="
	@echo "示例："
	@echo "  make commit-feat \"跨平台二进制构建功能\""
	@echo "  make commit-fix \"gh 登录授权检测失败问题\""
	@echo "  make commit-docs \"CHANGELOG.md 格式说明\""
	@echo "=========================="