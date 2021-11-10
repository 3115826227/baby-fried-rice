FROM golang:1.15 AS build
WORKDIR /go/src/baby-fried-rice
COPY . .
RUN go env -w GO111MODULE=off
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -tags=jsoniter -o sms_dao cmd/smsDao/sms_dao.go

FROM alpine:3.8
ENV TZ=Asia/Shanghai LANG=C.UTF-8 TERM=xterm
WORKDIR /go/src/baby-fried-rice/module/smsDao
RUN mkdir res && mkdir -p /baby/logs/
COPY --from=build /go/src/baby-fried-rice/cmd/smsDao/res/config.yaml res/
COPY --from=build /go/src/baby-fried-rice/cmd/smsDao/cert cert/
COPY --from=build /go/src/baby-fried-rice/sms_dao .
EXPOSE 8080
VOLUME ["/baby/logs"]
ENTRYPOINT ["./sms_dao"]
