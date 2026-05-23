# 构建阶段
FROM node:18-alpine AS builder

# 设置工作目录
WORKDIR /app

# 设置 npm 镜像
RUN npm config set registry https://registry.npmmirror.com

# 复制 package 文件
COPY package*.json ./

# 安装依赖
RUN npm ci --ignore-scripts

# 复制源代码
COPY . .

# 构建应用
RUN npm run build

# 运行阶段 - 使用 Nginx 提供静态文件服务
FROM nginx:alpine

# 设置时区
ENV TZ=Asia/Shanghai

# 复制构建产物
COPY --from=builder /app/dist /usr/share/nginx/html

# 复制 Nginx 配置
COPY nginx.conf /etc/nginx/conf.d/default.conf

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost/ || exit 1

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
