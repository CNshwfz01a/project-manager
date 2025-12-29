# 使用官方 Go 镜像作为构建环境
FROM golang:1.25.3 AS builder

# 设置工作目录
WORKDIR /app

# 复制 go.mod 和 go.sum 并下载依赖
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -o project-manager .

# 使用轻量级镜像作为运行环境
FROM alpine:latest

WORKDIR /root/

# 从构建阶段复制二进制文件
COPY --from=builder /app/project-manager .
COPY --from=builder /app/conf ./conf

# 暴露端口
EXPOSE 8086

# 运行应用
CMD ["./project-manager"]
