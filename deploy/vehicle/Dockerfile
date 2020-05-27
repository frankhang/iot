FROM golang:1.14.1 AS builder
ENV GOPROXY="https://goproxy.cn,direct"
ENV PRJ_DIR="github.com/frankhang/iot/"
ENV APP_DIR="gateway/vehicle/"
ENV APP_NAME="vehicle"
WORKDIR /go/src/${PRJ_DIR}
COPY . .
WORKDIR /go/src/${PRJ_DIR}${APP_DIR}
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ${APP_NAME} .

FROM alpine:latest
#FROM ubuntu
ENV PRJ_DIR="github.com/frankhang/iot/"
ENV APP_DIR="gateway/vehicle/"
ENV APP_NAME="vehicle"
WORKDIR /root/
COPY --from=builder /go/src/${PRJ_DIR}${APP_DIR}${APP_NAME} .
#ADD https://raw.githubusercontent.com/frankhang/doppler/master/deploy/doppler/vehicle.toml .
COPY deploy/vehicle/vehicle.toml .
ENTRYPOINT ["./vehicle", "-log-file=", "-config=vehicle.toml"]