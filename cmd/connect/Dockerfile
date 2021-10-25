FROM golang:1.15 AS build
WORKDIR /go/src/baby-fried-rice
COPY . .
RUN go env -w GO111MODULE=off
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -tags=jsoniter -o connect cmd/connect/connect.go

FROM alpine:3.8
ENV TZ=Asia/Shanghai LANG=C.UTF-8 TERM=xterm
WORKDIR /go/src/baby-fried-rice/module/connect
RUN mkdir res && mkdir cert && mkdir -p /baby/logs/
COPY --from=build /go/src/baby-fried-rice/cmd/connect/res/config.yaml res/
COPY --from=build /go/src/baby-fried-rice/cmd/connect/cert cert/
COPY --from=build /go/src/baby-fried-rice/connect .

EXPOSE 8080
VOLUME ["/baby/logs"]
ENTRYPOINT ["./connect"]