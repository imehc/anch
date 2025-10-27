.PHONY: gen-apis build run clean help

GENTOOL := openapi-generator-cli

# 默认目标
.DEFAULT_GOAL := help

## help: 显示帮助信息
help:
	@echo "可用的 Make 命令："
	@echo ""
	@echo "  make gen-apis     - 根据 OpenAPI 规范生成代码"
	@echo "  make build        - 编译项目"
	@echo "  make run          - 运行服务器"
	@echo "  make clean        - 清理编译文件"
	@echo ""
	@echo "数据库部署："
	@echo "  cd ../deploy && ./init_db.sh"
	@echo ""

## gen-apis: 生成 API 代码
gen-apis:
	$(GENTOOL) generate -c openapi-generator.config.yaml

## build: 编译项目
build:
	@echo "正在编译..."
	go build -o bin/server ./cmd/main.go
	@echo "编译完成: bin/server"

## run: 运行服务器
run: build
	@echo "启动服务器..."
	./bin/server

## clean: 清理编译文件
clean:
	@echo "清理编译文件..."
	rm -rf bin/
	@echo "清理完成"
