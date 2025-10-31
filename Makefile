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

# --------------- 版本管理（打 Tag）---------------
# 显示当前版本（最后一个 Tag）
version:
	@echo "当前最新版本: $(shell git describe --abbrev=0 --tags 2>/dev/null || echo '无版本Tag')"

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
	@echo "  make version          显示当前最新版本"
	@echo "  make tag VERSION=vX.Y.Z  创建本地版本Tag"
	@echo "  make push-tag VERSION=vX.Y.Z  推送Tag至远程"
	@echo "  make install-doc-tool  安装文档生成工具"
	@echo "  make gen-docs         生成所有模块文档"
	@echo "  make view-docs        打开文档目录预览"
	@echo "  make test             运行测试（含race检测）"
	@echo "  make clean            清理文档和临时文件"
	@echo "  make help             显示帮助信息"

# 默认命令：显示帮助
.DEFAULT_GOAL := help