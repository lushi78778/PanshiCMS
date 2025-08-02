# File: Dockerfile

# ---- Builder Stage ----
# 使用官方的 Go 镜像作为构建环境
FROM golang:1.24-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制 go mod 和 go sum 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制所有源代码
COPY . .

# 构建应用，-o 指定输出文件名，CGO_ENABLED=0 保证静态编译
RUN CGO_ENABLED=0 GOOS=linux go build -o panshcms ./cmd/server

# ---- Final Stage ----
# 使用一个非常小的基础镜像
FROM alpine:latest

# 设置工作目录
WORKDIR /app

# 从 builder 阶段复制编译好的二进制文件
COPY --from=builder /app/panshcms .

# 复制配置文件
COPY config.yaml .

# 暴露端口
EXPOSE 8080

# 容器启动时执行的命令
CMD ["./panshcms"]