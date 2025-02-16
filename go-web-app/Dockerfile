FROM golang:1.23-alpine

WORKDIR /go/src/gin
COPY . /go/src/gin

RUN go env -w GOPROXY=https://goproxy.io,direct

WORKDIR /go/src/gin
RUN go mod tidy

# # 创建非 root 用户并切换
# RUN addgroup -S appgroup && adduser -S appuser -G appgroup
# USER appuser


# 创建日志目录
RUN mkdir -p /go/src/gin/logs

EXPOSE 8080
CMD ["go", "run", "main.go"]
