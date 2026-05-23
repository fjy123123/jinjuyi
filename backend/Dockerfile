# 构建阶段
FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装构建依赖
RUN apk --no-cache add gcc musl-dev git

# 设置 Go 模块代理
ENV GOPROXY=https://goproxy.cn,direct
ENV GO111MODULE=on

# 复制 go mod 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/zhixin-backend .

# 运行阶段
FROM alpine:3.19

# 设置时区
ENV TZ=Asia/Shanghai
RUN apk --no-cache add ca-certificates tzdata

# 创建应用用户
RUN addgroup -g 1000 zhixin && \
    adduser -D -u 1000 -G zhixin zhixin

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/zhixin-backend /app/zhixin-backend

# 创建必要的目录
RUN mkdir -p /app/logs /app/uploads && \
    chown -R zhixin:zhixin /app

# 切换到非 root 用户
USER zhixin

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/v1/system/config || exit 1

# 启动应用
CMD ["/app/zhixin-backend"]
