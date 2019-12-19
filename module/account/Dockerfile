#
# Build stage
#
FROM golang:1.11.1 AS build

WORKDIR /go/src/github.com/3115826227/baby-fried-rice/module/account

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -tags=jsoniter -o main main.go

#
# Production stage
#
FROM alpine:3.8

ENV TZ=Asia/Shanghai LANG=C.UTF-8 TERM=xterm \
    LOG_PATH="/baby/logs/runtime.log" \
    ACCESS_LOG_PATH="/baby/logs/access.log" \
    TASK_LOG=true

WORKDIR /go/src/github.com/3115826227/baby-fried-rice/module/account

RUN touch /etc/config.yaml && mkdir -p /baby/logs/

COPY --from=build /go/src/github.com/3115826227/baby-fried-rice/module/account/main .

EXPOSE 8080

VOLUME ["/baby/logs"]

ENTRYPOINT ["./main"]
