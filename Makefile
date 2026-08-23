# Makefile for mybbs
# Windows 用户：可用 MinGW-w64 make、GNU make (choco install make)、WSL。
# 没有 make 也没关系，直接复制粘贴对应命令运行即可。

GO              ?= go
APP             := mybbs
BIN_DIR         := ./bin
BINARY          := $(BIN_DIR)/$(APP).exe
MAIN_PKG        := ./cmd/main.go
CONFIG_FILE     := config/config.yaml

# 数据库配置（migrate-sql 目标使用）：如果环境变量没设就用默认
DB_HOST         ?= 127.0.0.1
DB_PORT         ?= 3306
DB_USER         ?= root
DB_PASS         ?=
DB_NAME         ?= mybbs
MYSQL_BIN       ?= mysql

MIGRATIONS_DIR  := migrations
SCHEMA_SQL      := $(MIGRATIONS_DIR)/tables.sql

GOPATH          ?= $(shell $(GO) env GOPATH)

.PHONY: all build run test vet lint fmt tidy migrate-sql clean help

all: vet test build  ## 默认：lint-like 检查 + 测试 + 构建

build:  ## 编译二进制到 ./bin/mybbs.exe
	@echo "[build] compiling $(MAIN_PKG) -> $(BINARY)"
	@mkdir -p $(BIN_DIR)
	$(GO) build -trimpath -ldflags "-s -w" -o $(BINARY) $(MAIN_PKG)

run:  ## 本地直接运行（必须在项目根目录，依赖相对路径 config/config.yaml）
	@echo "[run] go run $(MAIN_PKG)"
	@test -f $(CONFIG_FILE) || (echo "error: $(CONFIG_FILE) not found, please run from project root" >&2; exit 1)
	$(GO) run $(MAIN_PKG)

test:  ## 跑全部测试（含 migrations/tables_test.go 的 SQL<->model 一致性）
	@echo "[test] go test ./..."
	$(GO) test -v ./...

vet:  ## go vet 静态检查
	@echo "[vet] go vet ./..."
	$(GO) vet ./...

fmt:  ## gofmt 检查格式（只报告不修改）
	@echo "[fmt] gofmt -l ."
	@out=$$(gofmt -l . 2>/dev/null); \
	 if [ -n "$$out" ]; then \
	   echo "files need reformatting:"; echo "$$out"; exit 1; \
	 else echo "format OK"; fi

lint: vet fmt  ## 组合检查：vet + fmt

tidy:  ## go mod tidy + verify
	$(GO) mod tidy
	$(GO) mod verify

# 注意：此目标要求系统 PATH 中可执行到 mysql 客户端（MySQL Server 自带或单独装）
# 若 mysql 不是默认路径，设 MYSQL_BIN=<绝对路径> make migrate-sql
migrate-sql:  ## 用 MySQL 客户端执行 migrations/tables.sql（DROP + 重建 4 张表）
	@echo "[migrate-sql] applying $(SCHEMA_SQL) to $(DB_NAME) via $(MYSQL_BIN)"
	@test -f $(SCHEMA_SQL) || (echo "error: $(SCHEMA_SQL) not found" >&2; exit 1)
	@if [ -n "$(DB_PASS)" ]; then \
	  PASSARG="-p$(DB_PASS)"; \
	fi; \
	$(MYSQL_BIN) -h$(DB_HOST) -P$(DB_PORT) -u$(DB_USER) $$PASSARG \
	  -e "CREATE DATABASE IF NOT EXISTS \`$(DB_NAME)\` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;" && \
	$(MYSQL_BIN) -h$(DB_HOST) -P$(DB_PORT) -u$(DB_USER) $$PASSARG $(DB_NAME) < $(SCHEMA_SQL)
	@echo "[migrate-sql] done"

clean:  ## 删除编译产物
	@echo "[clean] rm -rf $(BIN_DIR)"
	@rm -rf $(BIN_DIR)

help:  ## 显示这份帮助
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
