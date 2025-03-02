# 设置 Go 编译器和基本变量
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=workflow-server
MAIN_PATH=cmd/main.go

# 版本管理
VERSION ?= $(shell git describe --tags --always || echo "v0.0.1")
GIT_COMMIT ?= $(shell git rev-parse --short HEAD || echo "unknown")
BUILD_TIME ?= $(shell date -u '+%Y-%m-%d_%H:%M:%S')

# 镜像相关变量
IMAGE_NAME = workflow-server
IMAGE_TAG ?= $(VERSION)
DOCKER_IMAGE = $(IMAGE_NAME):$(IMAGE_TAG)

# 编译参数 - 注入版本信息到二进制文件
LDFLAGS += -X "main.Version=$(VERSION)"
LDFLAGS += -X "main.GitCommit=$(GIT_COMMIT)"
LDFLAGS += -X "main.BuildTime=$(BUILD_TIME)"

.PHONY: all build clean run test deps tidy help

# 默认目标
all: clean build

# 构建应用
build:
	@echo "Building..."
	$(GOBUILD) $(LDFLAGS) -o bin/$(BINARY_NAME) $(MAIN_PATH)

# 运行应用
run:
	@echo "Running..."
	$(GORUN) $(MAIN_PATH)

# 清理构建文件
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -f bin/$(BINARY_NAME)

# 运行测试
test:
	@echo "Testing..."
	$(GOTEST) -v ./...

# 更新依赖
deps:
	@echo "Updating dependencies..."
	$(GOGET) -u ./...

# 整理 go.mod
tidy:
	@echo "Tidying Go modules..."
	$(GOMOD) tidy

# 创建 Docker 镜像
docker-build:
	@echo "Building Docker image $(DOCKER_IMAGE)..."
	docker build -t $(DOCKER_IMAGE) \
		--build-arg VERSION=$(VERSION) \
		--build-arg GIT_COMMIT=$(GIT_COMMIT) \
		--build-arg BUILD_TIME=$(BUILD_TIME) \
		.

# 运行 Docker 容器
docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 $(DOCKER_IMAGE)

# 显示帮助信息
help:
	@echo "Make commands:"
	@echo "  make build       - Build the application"
	@echo "  make run        - Run the application"
	@echo "  make clean      - Clean build files"
	@echo "  make test       - Run tests"
	@echo "  make deps       - Update dependencies"
	@echo "  make tidy       - Tidy go.mod"
	@echo "  make docker-build - Build Docker image"
	@echo "  make docker-run   - Run Docker container"
	@echo "  make help       - Show this help message"

# 创建必要的目录
init:
	@echo "Initializing project structure..."
	mkdir -p bin
	mkdir -p pkg
	mkdir -p internal
	mkdir -p api
	mkdir -p config
	mkdir -p docs

# 生成 API 文档
docs:
	@echo "Generating API documentation..."
	swag init -g cmd/main.go -o docs/swagger

# 检查代码格式
fmt:
	@echo "Formatting code..."
	gofmt -s -w .

# 运行代码检查
lint:
	@echo "Running linter..."
	golangci-lint run

# 开发模式（使用 air 进行热重载）
dev:
	@echo "Running in development mode..."
	air 

# 添加 docker 标签命令
docker-tag:
	@echo "Tagging Docker image..."
	docker tag $(DOCKER_IMAGE) $(DOCKER_REGISTRY)/$(IMAGE_NAME):latest

# 添加推送镜像命令
docker-push:
	@echo "Pushing Docker image..."
	docker push $(DOCKER_IMAGE)
	docker push $(DOCKER_REGISTRY)/$(IMAGE_NAME):latest

# 添加版本标签命令
tag:
	@echo "Creating git tag $(VERSION)..."
	git tag $(VERSION)
	git push origin $(VERSION)

# 显示版本信息
version:
	@echo "Version: $(VERSION)"
	@echo "Git Commit: $(GIT_COMMIT)"
	@echo "Build Time: $(BUILD_TIME)" 