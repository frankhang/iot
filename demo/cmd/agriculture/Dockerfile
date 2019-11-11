FROM golang:1.13 AS builder
ENV GOPROXY="https://goproxy.cn,direct"
ENV PRJ_DIR="github.com/frankhang/iot/"
ENV APP_DIR="demo/cmd/agriculture/"
ENV APP_NAME="agriculture"
WORKDIR /go/src/${PRJ_DIR}
COPY . .
WORKDIR /go/src/${PRJ_DIR}${APP_DIR}
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ${APP_NAME} .

FROM alpine:latest
ENV PRJ_DIR="github.com/frankhang/iot/"
ENV APP_DIR="demo/cmd/agriculture/"
ENV APP_NAME="agriculture"
WORKDIR /root/
COPY --from=builder /go/src/${PRJ_DIR}${APP_DIR}${APP_NAME} .
CMD [ ":8080" ]
ENTRYPOINT ["./agriculture", "--listen-address"]