# 构建阶段
FROM golang:1.21-alpine AS builder

# 接收构建参数
ARG VERSION
ARG GIT_COMMIT
ARG BUILD_TIME

# 设置工作目录
WORKDIR /app

# 设置 Go 环境变量
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# 复制 go.mod 和 go.sum
COPY go.mod go.sum ./

# 下载依赖
RUN go mod tidy

# 复制源代码
COPY . .

# 构建应用 - 注入版本信息
RUN go build -ldflags "-X main.Version=${VERSION} -X main.GitCommit=${GIT_COMMIT} -X main.BuildTime=${BUILD_TIME}" -o workflow-server cmd/main.go

# 运行阶段
FROM alpine:latest

# 设置工作目录
WORKDIR /app

# 安装基础工具和证书
RUN apk --no-cache add ca-certificates tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

# 从构建阶段复制二进制文件
COPY --from=builder /app/workflow-server .
# 复制配置文件
COPY --from=builder /app/config/config.yaml ./config/

# 暴露端口
EXPOSE 8080

# 设置启动命令
CMD ["./workflow-server"] 