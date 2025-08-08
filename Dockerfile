# ---- Builder Stage ----
FROM golang:1.24 AS builder


# 安装编译依赖
RUN apt-get update && apt-get install -y build-essential

# 设置 Go 模块代理为国内源
ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o panshcms ./cmd/server

# ---- Final Stage ----
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y libc6

WORKDIR /app
COPY --from=builder /app/panshcms .
COPY config.yaml .
EXPOSE 8080
CMD ["./panshcms"]
