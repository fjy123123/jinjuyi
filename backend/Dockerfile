# 构建阶段
FROM golang:1.21-alpine AS builder

WORKDIR /build

# 安装依赖
RUN apk add --no-cache git make

# 复制源码
COPY backend/ ./

# 下载依赖
RUN go mod download
RUN go mod verify

# 编译
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o chat-server .

# 运行阶段
FROM alpine:latest

WORKDIR /app

# 安装CA证书和时区数据
RUN apk --no-cache add ca-certificates tzdata

# 创建目录
RUN mkdir -p /app/uploads /app/logs

# 从构建阶段复制二进制文件
COPY --from=builder /build/chat-server .
COPY --from=builder /build/config.yaml .

# 复制上传目录
COPY --from=builder /build/uploads /app/uploads

# 设置时区
ENV TZ=Asia/Shanghai

# 暴露端口
EXPOSE 8080

# 启动命令
CMD ["./chat-server"]
