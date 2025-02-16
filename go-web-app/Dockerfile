FROM golang:1.23-alpine

WORKDIR /go/src/gin
COPY . /go/src/gin

RUN go env -w GOPROXY=https://goproxy.io,direct

WORKDIR /go/src/gin
RUN go mod tidy

RUN mkdir -p /go/src/gin/logs

EXPOSE 8080
CMD ["go", "run", "main.go"]
